package enricher

import (
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/pb/enrichers"
	"math/rand"
)

const maxTransportType = int(enrichers.Transport_TRANSPORT_PLANE)

type transport struct{}

func NewTransport() Enricher {
	return &transport{}
}

func (e transport) Enrich(payload Payload) SetFunc {
	var transport enrichers.Transport

	transport.Type = enrichers.Transport_Type(rand.Intn(maxTransportType))

	return func(gle *pb.GeoLocationEnriched) {
		gle.Transport = &transport
	}
}
