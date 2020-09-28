package enricher

import (
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/pb/enrichers"
	weatherService "github.com/shekhirin/zenly-task/internal/service/weather"
)

type weather struct {
	service weatherService.Service
}

func NewWeather(service weatherService.Service) Enricher {
	return &weather{service: service}
}

func (e weather) Enrich(payload Payload) SetFunc {
	var weather enrichers.Weather

	weather.Temperature = e.service.Temperature(payload.Lat, payload.Lng)

	return func(gle *pb.GeoLocationEnriched) {
		gle.Weather = &weather
	}
}
