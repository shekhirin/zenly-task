package load

import (
	"context"
	"github.com/shekhirin/zenly-task/load/stats"
	"github.com/shekhirin/zenly-task/zenly/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

type Loader struct {
	grpcAddr       string
	rps            int
	duration       time.Duration
	userIdsNum     int
	publishStats   stats.Publish
	subscribeStats stats.Subscribe
}

func NewLoader(grpcAddr string, rps int, duration time.Duration, userIdsNum int) Loader {
	return Loader{
		grpcAddr:       grpcAddr,
		rps:            rps,
		duration:       duration,
		userIdsNum:     userIdsNum,
		publishStats:   stats.NewPublish(),
		subscribeStats: stats.NewSubscribe(),
	}
}

func (l *Loader) Load() {
	var ctx context.Context
	var cancel context.CancelFunc
	if l.duration > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), l.duration)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	conn, err := grpc.Dial(l.grpcAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("dial grpc on %s", l.grpcAddr)
	}

	client := pb.NewZenlyClient(conn)

	publishWaitCh := l.loadPublish(ctx, client)
	subscribeWaitCh := l.loadSubscribe(ctx, client)

	var signalCh = make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)

	select {
	case <-signalCh:
	case <-publishWaitCh:
	case <-subscribeWaitCh:
	}
}

func (l *Loader) loadPublish(ctx context.Context, client pb.ZenlyClient) <-chan struct{} {
	publishClient, err := client.Publish(ctx)
	if err != nil {
		log.WithError(err).Fatal("init publish client")
	}

	tick := time.Second / time.Duration(l.rps)
	timer := time.NewTicker(tick)

	durationStr := l.duration.String()
	if l.duration == 0 {
		durationStr = "infinite"
	}
	log.WithFields(log.Fields{
		"duration":     durationStr,
		"rps":          l.rps,
		"tick_every":   tick,
		"user_ids_num": l.userIdsNum,
	}).Info("start publish")

	var waitCh = make(chan struct{})

	go func() {
		defer close(waitCh)
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				go func() {
					start := time.Now()
					err := publishClient.Send(&pb.PublishRequest{
						UserId: rand.Int31n(int32(l.userIdsNum)),
						GeoLocation: &pb.GeoLocation{
							Lat:       rand.Float64(),
							Lng:       rand.Float64(),
							CreatedAt: timestamppb.Now(),
						},
					})
					finish := time.Now()

					if err != nil {
						log.WithError(err).Fatal("send publish request")
					}

					l.publishStats.Observe(start, finish)
				}()
			}
		}
	}()

	return waitCh
}

func (l *Loader) loadSubscribe(ctx context.Context, client pb.ZenlyClient) <-chan struct{} {
	var subscribeUserIds []int32
	for i := 0; i < l.userIdsNum; i++ {
		subscribeUserIds = append(subscribeUserIds, int32(i))
	}

	subscribeClient, err := client.Subscribe(ctx, &pb.SubscribeRequest{
		UserId: subscribeUserIds,
	})
	if err != nil {
		log.WithError(err).Fatal("init subscribe client")
	}

	log.WithFields(log.Fields{
		"user_ids": subscribeUserIds,
	}).Info("start subscribe")

	var waitCh = make(chan struct{})

	go func() {
		defer close(waitCh)
		for {
			select {
			case <-ctx.Done():
				err := subscribeClient.CloseSend()
				if err != nil {
					log.WithError(err).Fatal("close subscribe stream send direction")
				}
			default:
				message, err := subscribeClient.Recv()
				received := time.Now()
				switch err {
				case nil:
					break
				case io.EOF:
					return
				default:
					log.WithError(err).Fatal("receive from subscribe stream")
				}

				l.subscribeStats.Observe(received, message)
			}
		}
	}()

	return waitCh
}

func (l *Loader) PrintStats() {
	log.WithFields(log.Fields{
		"actual_rps":  l.publishStats.RPS(),
		"requests":    l.publishStats.Requests(),
		"elapsed_min": l.publishStats.ElapsedMin,
		"elapsed_max": l.publishStats.ElapsedMax,
		"elapsed_avg": l.publishStats.ElapsedAverage(),
	}).Info("publish complete")

	log.WithFields(log.Fields{
		"messages_per_second": l.subscribeStats.MPS(),
		"messages":            l.subscribeStats.Messages(),
		"lag_avg":             l.subscribeStats.LagAverage(),
	}).Info("subscribe complete")
}
