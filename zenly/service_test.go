package zenly

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/shekhirin/zenly-task/internal/pb"
	pbEnricher "github.com/shekhirin/zenly-task/internal/pb/enricher"
	pbMocks "github.com/shekhirin/zenly-task/internal/pb/mocks"
	busMocks "github.com/shekhirin/zenly-task/zenly/bus/mocks"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	enricherMocks "github.com/shekhirin/zenly-task/zenly/enricher/mocks"
	feedMocks "github.com/shekhirin/zenly-task/zenly/feed/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestZenly_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stream := pbMocks.NewMockZenly_PublishServer(ctrl)
	stream.EXPECT().
		Recv().
		Return(&pb.PublishRequest{
			UserId: 1,
			GeoLocation: &pb.GeoLocation{
				Lat: 2.2,
				Lng: 3.3,
			},
		}, nil)
	stream.EXPECT().
		Recv().
		Return(nil, io.EOF)
	stream.EXPECT().
		SendAndClose(&pb.PublishResponse{Success: true}).
		Return(nil)

	mockEnricher := enricherMocks.NewMockEnricher(ctrl)
	mockEnricher.EXPECT().
		String().
		Return("mock").
		AnyTimes()
	mockEnricher.EXPECT().
		Enrich(enricher.Payload{
			UserId: 1,
			Lat:    2.2,
			Lng:    3.3,
		}).
		Return(func(gle *pb.GeoLocationEnriched) {
			gle.Weather = &pbEnricher.Weather{
				Temperature: 6.9,
			}
		})

	bus := busMocks.NewMockBus(ctrl)
	bus.EXPECT().
		Publish(&pb.BusMessage{
			UserId: 1,
			GeoLocation: &pb.GeoLocationEnriched{
				GeoLocation: &pb.GeoLocation{
					Lat: 2.2,
					Lng: 3.3,
				},
				Weather: &pbEnricher.Weather{
					Temperature: 6.9,
				},
			},
		})

	feed := feedMocks.NewMockFeed(ctrl)
	feed.EXPECT().
		Publish(&pb.FeedMessage{
			UserId: 1,
			GeoLocation: &pb.GeoLocationEnriched{
				GeoLocation: &pb.GeoLocation{
					Lat: 2.2,
					Lng: 3.3,
				},
				Weather: &pbEnricher.Weather{
					Temperature: 6.9,
				},
			},
			BusPublished: true,
		})

	zenly := New(bus, feed, []enricher.Enricher{mockEnricher})

	assert.NoError(t, zenly.Publish(stream))
}

func TestZenly_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	busCtx, busCtxCancel := context.WithCancel(context.Background())

	busCh := make(chan *pb.BusMessage, 1)
	busCh <- &pb.BusMessage{
		UserId: 1,
		GeoLocation: &pb.GeoLocationEnriched{
			GeoLocation: &pb.GeoLocation{
				Lat: 2.2,
				Lng: 3.3,
			},
			Weather: &pbEnricher.Weather{
				Temperature: 6.9,
			},
		},
	}

	bus := busMocks.NewMockBus(ctrl)
	bus.EXPECT().
		Subscribe([]int32{1}).
		Return(busCh, busCtxCancel, nil)

	streamCtx, streamCtxCancel := context.WithCancel(context.Background())

	stream := pbMocks.NewMockZenly_SubscribeServer(ctrl)
	stream.EXPECT().
		Send(&pb.SubscribeResponse{
			UserId: 1,
			GeoLocation: &pb.GeoLocationEnriched{
				GeoLocation: &pb.GeoLocation{
					Lat: 2.2,
					Lng: 3.3,
				},
				Weather: &pbEnricher.Weather{
					Temperature: 6.9,
				},
			},
		}).
		Return(nil).
		Do(func(*pb.SubscribeResponse) {
			streamCtxCancel()
		})
	stream.EXPECT().
		Context().
		Return(streamCtx).
		AnyTimes()

	feed := feedMocks.NewMockFeed(ctrl)

	request := &pb.SubscribeRequest{
		UserId: []int32{1},
	}

	zenly := New(bus, feed, []enricher.Enricher{})

	assert.NoError(t, zenly.Subscribe(request, stream))
	assert.Equal(t, context.Canceled, busCtx.Err())
}
