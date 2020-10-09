package zenly

import (
	"github.com/shekhirin/zenly-task/zenly/bus"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	"github.com/shekhirin/zenly-task/zenly/feed"
	weatherService "github.com/shekhirin/zenly-task/zenly/service/weather"
	"time"
)

const EnricherTimeout = 100 * time.Millisecond

var DefaultEnrichers = []enricher.Enricher{
	enricher.NewWeather(weatherService.New()),
	enricher.NewPersonalPlace(),
	enricher.NewTransport(),
}

type Zenly struct {
	bus       bus.Bus
	feed      feed.Feed
	enrichers []enricher.Enricher
}

func New(bus bus.Bus, feed feed.Feed, enrichers []enricher.Enricher) *Zenly {
	return &Zenly{
		bus:       bus,
		feed:      feed,
		enrichers: enrichers,
	}
}
