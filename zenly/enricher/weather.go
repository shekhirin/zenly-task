package enricher

import (
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/pb/enricher"
	weatherService "github.com/shekhirin/zenly-task/zenly/service/weather"
)

type weather struct {
	service weatherService.Service
}

func NewWeather(service weatherService.Service) Enricher {
	return &weather{service: service}
}

func (e weather) Enrich(payload Payload) SetFunc {
	var weather enricher.Weather

	weather.Temperature = e.service.Temperature(payload.Lat, payload.Lng)

	simulateIO()

	return func(gle *pb.GeoLocationEnriched) {
		gle.Weather = &weather
	}
}

func (e weather) String() string {
	return "weather"
}
