package zenly

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/shekhirin/zenly-task/internal/enricher"
	"github.com/shekhirin/zenly-task/internal/pb"
	weatherService "github.com/shekhirin/zenly-task/internal/service/weather"
	"io"
)

var DefaultEnrichers = []enricher.Enricher{
	enricher.NewWeather(weatherService.NewService()),
	enricher.NewPersonalPlace(),
	enricher.NewTransport(),
}

type busMessage struct {
	UserId              int32                   `json:"user_id"`
	GeoLocationEnriched *pb.GeoLocationEnriched `json:"geo_location_enriched"`
}

type Server struct {
	nats      *nats.Conn
	enrichers []enricher.Enricher
}

func NewServer(nats *nats.Conn, enrichers []enricher.Enricher) *Server {
	s := &Server{
		nats:      nats,
		enrichers: enrichers,
	}
	return s
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

		for _, serverEnricher := range s.enrichers {
			serverEnricher.Enrich(context.Background(), &geoLocationEnriched)
		}

		message := busMessage{
			UserId:              publishRequest.UserId,
			GeoLocationEnriched: &geoLocationEnriched,
		}

		data, err := json.Marshal(message)
		if err != nil {
			return err
		}

		if err := s.nats.Publish("zenly", data); err != nil {
			return err
		}
	}
}

func (s *Server) Subscribe(request *pb.SubscribeRequest, stream pb.Zenly_SubscribeServer) error {
	var userIds = make(map[int32]bool)
	for _, userId := range request.UserId {
		userIds[userId] = true
	}

	var ch = make(chan *nats.Msg)

	sub, err := s.nats.ChanSubscribe("zenly", ch)
	if err != nil {
		return err
	}

	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	for {
		select {
		case msg := <-ch:
			var message busMessage
			if err := json.Unmarshal(msg.Data, &message); err != nil {
				continue
			}

			if !userIds[message.UserId] {
				continue
			}

			subscribeResponse := &pb.SubscribeResponse{
				UserId:      message.UserId,
				GeoLocation: message.GeoLocationEnriched,
			}

			if err := stream.Send(subscribeResponse); err != nil {
				return err
			}
		}
	}
}
