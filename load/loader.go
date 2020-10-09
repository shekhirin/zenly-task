package load

import (
	"context"
	"github.com/shekhirin/zenly-task/load/stats"
	"github.com/shekhirin/zenly-task/zenly/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

type Loader struct {
	grpcAddr   string
	rps        int
	duration   time.Duration
	userIdsNum int
	stats      stats.Stats
}

func NewLoader(grpcAddr string, rps int, duration time.Duration, userIdsNum int) Loader {
	return Loader{
		grpcAddr:   grpcAddr,
		rps:        rps,
		duration:   duration,
		userIdsNum: userIdsNum,
		stats:      stats.New(),
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

	publishClient, err := client.Publish(ctx)
	if err != nil {
		log.WithError(err).Fatal("init publish client")
	}

	tick := time.Second / time.Duration(l.rps)
	timer := time.NewTicker(tick)
	defer timer.Stop()

	durationStr := l.duration.String()
	if l.duration == 0 {
		durationStr = "infinite"
	}
	log.WithFields(log.Fields{
		"duration":     durationStr,
		"rps":          l.rps,
		"tick_every":   tick,
		"user_ids_num": l.userIdsNum,
	}).Info("start")

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
					err = publishClient.Send(&pb.PublishRequest{
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

					l.stats.Observe(start, finish)
				}()
			}
		}
	}()

	var signalCh = make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)

	select {
	case <-signalCh:
	case <-waitCh:
	}
}

func (l *Loader) PrintStats() {
	log.WithFields(log.Fields{
		"actual_rps":  l.stats.RPS(),
		"requests":    l.stats.Requests(),
		"elapsed_min": l.stats.ElapsedMin,
		"elapsed_max": l.stats.ElapsedMax,
		"elapsed_avg": l.stats.ElapsedAverage(),
	}).Info("complete")
}
