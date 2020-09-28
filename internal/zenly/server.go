package zenly

import (
	"context"
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/zenly/bus"
	"github.com/shekhirin/zenly-task/internal/zenly/enricher"
	weatherService "github.com/shekhirin/zenly-task/internal/zenly/service/weather"
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

		wg := sync.WaitGroup{}
		wg.Add(len(s.enrichers))

		for _, serverEnricher := range s.enrichers {
			go func(targetEnricher enricher.Enricher) {
				ctx, cancel := context.WithTimeout(context.Background(), EnricherTimeout)
				defer cancel()
				defer wg.Done()

				// Don't give control of the context to enricher because of the possibility of forgetting
				// to check timeout before setting the submessage inside the enricher
				select {
				case <-ctx.Done():
					return
				case enrich := <-enricher.EnrichChannel(targetEnricher, payload):
					enrich(&geoLocationEnriched)
				}
			}(serverEnricher)
		}

		wg.Wait()

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
