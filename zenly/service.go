package zenly

import (
	"context"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	"github.com/shekhirin/zenly-task/zenly/pb"
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

		message := &pb.BusMessage{
			UserId:      publishRequest.UserId,
			GeoLocation: geoLocationEnriched,
		}

		if err := z.bus.Publish(message); err != nil {
			return err
		}
	}
}

func (z *Zenly) Subscribe(request *pb.SubscribeRequest, stream pb.Zenly_SubscribeServer) error {
	cancel, err := z.bus.Subscribe(request.UserId, func(message *pb.BusMessage) error {
		subscribeResponse := &pb.SubscribeResponse{
			UserId:      message.UserId,
			GeoLocation: message.GeoLocation,
		}

		return stream.Send(subscribeResponse)
	})
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
		}
	}
}
