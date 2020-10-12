package zenly

import (
	"context"
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	log "github.com/sirupsen/logrus"
	"io"
)

func (z *Zenly) Service() *pb.ZenlyService {
	return &pb.ZenlyService{
		Publish:   z.Publish,
		Subscribe: z.Subscribe,
	}
}

func (z *Zenly) Publish(stream pb.Zenly_PublishServer) error {
	for {
		message, err := stream.Recv()
		switch err {
		case nil:
			break
		case io.EOF:
			return stream.SendAndClose(&pb.PublishResponse{
				Success: true,
			})
		default:
			log.WithError(err).Error("receive from publish stream")
			return err
		}

		geoLocationEnriched := &pb.GeoLocationEnriched{
			GeoLocation: message.GeoLocation,
		}

		payload := enricher.Payload{
			UserId: message.UserId,
			Lat:    message.GeoLocation.Lat,
			Lng:    message.GeoLocation.Lng,
		}

		z.Enrich(payload, geoLocationEnriched)

		busMessage := &pb.BusMessage{
			UserId:      message.UserId,
			GeoLocation: geoLocationEnriched,
		}

		busErr := z.bus.Publish(busMessage)
		if busErr != nil {
			log.WithError(busErr).Error("publish bus message")
		}

		feedMessage := &pb.FeedMessage{
			UserId:       message.UserId,
			GeoLocation:  geoLocationEnriched,
			BusPublished: busErr == nil,
		}

		if err := z.feed.Publish(feedMessage); err != nil {
			log.WithError(err).Error("publish feed message")
		}

		if busErr != nil {
			return busErr
		}
	}
}

func (z *Zenly) Subscribe(request *pb.SubscribeRequest, stream pb.Zenly_SubscribeServer) error {
	ch, cancel, err := z.bus.Subscribe(request.UserId)
	if err != nil {
		log.WithError(err).WithField("user_ids", request.UserId).Error("subscribe to bus")
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
				log.WithError(err).Error("stream context done")
				return err
			}
		case message := <-ch:
			subscribeResponse := &pb.SubscribeResponse{
				UserId:      message.UserId,
				GeoLocation: message.GeoLocation,
			}

			if err := stream.Send(subscribeResponse); err != nil {
				log.WithError(err).Error("send to subscribe stream")
				return err
			}
		}
	}
}
