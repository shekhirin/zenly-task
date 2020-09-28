package enricher

import (
	"context"
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/pb/enrichers"
	weatherService "github.com/shekhirin/zenly-task/internal/service/weather"
)

type weatherEnricher struct {
	service weatherService.Service
}

func NewWeather(service weatherService.Service) Enricher {
	return &weatherEnricher{service: service}
}

func (e weatherEnricher) Enrich(ctx context.Context, gle *pb.GeoLocationEnriched) {
	var weather enrichers.Weather

	weather.Temperature = e.service.Temperature(gle.GeoLocation.Lat, gle.GeoLocation.Lng)

	gle.Weather = &weather
}
