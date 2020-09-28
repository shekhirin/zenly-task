package zenly

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shekhirin/zenly-task/zenly/bus"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	"github.com/shekhirin/zenly-task/zenly/pb"
	weatherService "github.com/shekhirin/zenly-task/zenly/service/weather"
	log "github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"io"
	"sync"
	"time"
)

const EnricherTimeout = 100 * time.Millisecond

var DefaultEnrichers = []enricher.Enricher{
	enricher.NewWeather(weatherService.New()),
	enricher.NewPersonalPlace(),
	enricher.NewTransport(),
}

type Zenly struct {
	bus            bus.Bus
	enricherTimeMS *prometheus.HistogramVec
	enrichers      []enricher.Enricher
}

func New(bus bus.Bus, enricherTimeMS *prometheus.HistogramVec, enrichers []enricher.Enricher) *Zenly {
	return &Zenly{
		bus:            bus,
		enricherTimeMS: enricherTimeMS,
		enrichers:      enrichers,
	}
}

func (z *Zenly) Service() *pb.ZenlyService {
	return &pb.ZenlyService{
		Publish:   z.Publish,
		Subscribe: z.Subscribe,
	}
}

func (z *Zenly) Publish(stream pb.Zenly_PublishServer) error {
	for {
		publishRequest, err := stream.Recv()
		switch err {
		case nil:
			break
		case io.EOF:
			return stream.SendAndClose(&pb.PublishResponse{
				Success: true,
			})
		default:
			return err
		}

		var geoLocationEnriched = pb.GeoLocationEnriched{
			GeoLocation: publishRequest.GeoLocation,
		}

		payload := enricher.Payload{
			UserId: publishRequest.UserId,
			Lat:    publishRequest.GeoLocation.Lat,
			Lng:    publishRequest.GeoLocation.Lng,
		}

		var completed atomic.Int32

		var wg sync.WaitGroup
		wg.Add(len(z.enrichers))
		waitCh := make(chan struct{})

		ctx, _ := context.WithTimeout(context.Background(), EnricherTimeout)

		go func() {
			for _, serverEnricher := range z.enrichers {
				go func(enricher enricher.Enricher) {
					defer wg.Done()

					start := time.Now()

					// Don't give control of the context to enricher because of the possibility of forgetting
					// to check timeout before setting the submessage inside the enricher
					enrich := enricher.Enrich(payload)
					elapsed := time.Since(start)

					z.enricherTimeMS.With(prometheus.Labels{"enricher": enricher.String()}).Observe(float64(elapsed.Milliseconds()))

					debugString := fmt.Sprintf("%s took %s\n", enricher.String(), elapsed.String())

					if ctx.Err() != nil {
						log.WithField("enriched", false).Debug(debugString)
						return
					}

					log.WithField("enriched", true).Debug(debugString)
					enrich(&geoLocationEnriched)

					completed.Inc()
				}(serverEnricher)
			}

			wg.Wait()
			close(waitCh)
		}()

		select {
		case <-ctx.Done():
			log.Debugf("%d enrichers completed", completed.Load())
		case <-waitCh:
			log.Debugf("all enrichers completed")
		}

		message := &pb.BusMessage{
			UserId:      publishRequest.UserId,
			GeoLocation: &geoLocationEnriched,
		}

		if err := z.bus.Publish(message); err != nil {
			return err
		}
	}
}

func (z *Zenly) Subscribe(request *pb.SubscribeRequest, stream pb.Zenly_SubscribeServer) error {
	ch, cancel, err := z.bus.Subscribe(request.UserId)
	if err != nil {
		return err
	}
	defer cancel()

	for {
		select {
		case <-stream.Context().Done():
			switch stream.Context().Err() {
			case context.Canceled:
				return nil
			default:
				return err
			}
		case message := <-ch:
			subscribeResponse := &pb.SubscribeResponse{
				UserId:      message.UserId,
				GeoLocation: message.GeoLocation,
			}

			if err := stream.Send(subscribeResponse); err != nil {
				return err
			}
		}
	}
}
