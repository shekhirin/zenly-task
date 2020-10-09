package zenly

import (
	"context"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	"github.com/shekhirin/zenly-task/zenly/pb"
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

		geoLocationEnriched := &pb.GeoLocationEnriched{
			GeoLocation: publishRequest.GeoLocation,
		}

		payload := enricher.Payload{
			UserId: publishRequest.UserId,
			Lat:    publishRequest.GeoLocation.Lat,
			Lng:    publishRequest.GeoLocation.Lng,
		}

		z.Enrich(payload, geoLocationEnriched)

		busMessage := &pb.BusMessage{
			UserId:      publishRequest.UserId,
			GeoLocation: geoLocationEnriched,
		}

		busErr := z.bus.Publish(busMessage)

		feedMessage := &pb.FeedMessage{
			UserId:       publishRequest.UserId,
			GeoLocation:  geoLocationEnriched,
			BusPublished: busErr == nil,
		}

		// TODO: publish to feed from bus in a separate worker
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
