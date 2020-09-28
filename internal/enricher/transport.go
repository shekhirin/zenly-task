package enricher

import (
	"context"
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/pb/enrichers"
	"math/rand"
)

const maxTransportType = int(enrichers.Transport_TRANSPORT_PLANE)

type transportEnricher struct{}

func NewTransport() Enricher {
	return &transportEnricher{}
}

func (e transportEnricher) Enrich(ctx context.Context, gle *pb.GeoLocationEnriched) {
	var transport enrichers.Transport

	transport.Type = enrichers.Transport_Type(rand.Intn(maxTransportType))

	select {
	case <-ctx.Done():
		return
	default:
		gle.Transport = &transport
	}
}
