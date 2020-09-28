package zenly

import (
	"context"
	"fmt"
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/zenly/bus"
	"github.com/shekhirin/zenly-task/internal/zenly/enricher"
	weatherService "github.com/shekhirin/zenly-task/internal/zenly/service/weather"
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

type Server struct {
	bus       bus.Bus
	enrichers []enricher.Enricher
}

func NewServer(bus bus.Bus, enrichers []enricher.Enricher) *Server {
	return &Server{
		bus:       bus,
		enrichers: enrichers,
	}
}

func (s *Server) Service() *pb.ZenlyService {
	return &pb.ZenlyService{
		Publish:   s.Publish,
		Subscribe: s.Subscribe,
	}
}

func (s *Server) Publish(stream pb.Zenly_PublishServer) error {
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
		wg.Add(len(s.enrichers))
		waitCh := make(chan struct{})

		ctx, _ := context.WithTimeout(context.Background(), EnricherTimeout)

		go func() {
			for _, serverEnricher := range s.enrichers {
				go func(enricher enricher.Enricher) {
					defer wg.Done()

					start := time.Now()

					// Don't give control of the context to enricher because of the possibility of forgetting
					// to check timeout before setting the submessage inside the enricher
					enrich := enricher.Enrich(payload)
					elapsed := time.Since(start)

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

		if err := s.bus.Publish(message); err != nil {
			return err
		}
	}
}

func (s *Server) Subscribe(request *pb.SubscribeRequest, stream pb.Zenly_SubscribeServer) error {
	ch, cancel, err := s.bus.Subscribe(request.UserId)
	if err != nil {
		return err
	}
	defer cancel()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
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
