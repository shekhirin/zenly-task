package enricher

import (
	"github.com/shekhirin/zenly-task/internal/pb"
)

type Enricher interface {
	Enrich(payload Payload) SetFunc // Pass Payload by value to prevent modification
}

// Additional struct to pass only required fields to enricher
type Payload struct {
	UserId int32
	Lat    float64
	Lng    float64
}

type SetFunc func(gle *pb.GeoLocationEnriched)

func EnrichChannel(enricher Enricher, payload Payload) <-chan SetFunc {
	ch := make(chan SetFunc, 1)

	go func() {
		ch <- enricher.Enrich(payload)
	}()

	return ch
}
