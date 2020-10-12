package zenly

import (
	"github.com/golang/mock/gomock"
	"github.com/shekhirin/zenly-task/internal/pb"
	pbEnricher "github.com/shekhirin/zenly-task/internal/pb/enricher"
	busMocks "github.com/shekhirin/zenly-task/zenly/bus/mocks"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	enricherMocks "github.com/shekhirin/zenly-task/zenly/enricher/mocks"
	feedMocks "github.com/shekhirin/zenly-task/zenly/feed/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestZenly_Enrich(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWeatherEnricher := enricherMocks.NewMockEnricher(ctrl)
	mockWeatherEnricher.EXPECT().
		String().
		Return("weather").
		AnyTimes()
	mockWeatherEnricher.EXPECT().
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

	mockTransportEnricher := enricherMocks.NewMockEnricher(ctrl)
	mockTransportEnricher.EXPECT().
		String().
		Return("transport").
		AnyTimes()
	mockTransportEnricher.EXPECT().
		Enrich(enricher.Payload{
			UserId: 1,
			Lat:    2.2,
			Lng:    3.3,
		}).
		DoAndReturn(func(payload enricher.Payload) enricher.SetFunc {
			time.Sleep(100 * time.Millisecond)

			return func(gle *pb.GeoLocationEnriched) {
				gle.Transport = &pbEnricher.Transport{
					Type: pbEnricher.Transport_CAR,
				}
			}
		})

	enrichers := []enricher.Enricher{mockWeatherEnricher, mockTransportEnricher}

	zenly := New(busMocks.NewMockBus(ctrl), feedMocks.NewMockFeed(ctrl), enrichers)

	payload := enricher.Payload{
		UserId: 1,
		Lat:    2.2,
		Lng:    3.3,
	}

	var geoLocationEnriched pb.GeoLocationEnriched

	zenly.Enrich(payload, &geoLocationEnriched)

	assert.NotNil(t, geoLocationEnriched.Weather)
	assert.Nil(t, geoLocationEnriched.Transport)
}
