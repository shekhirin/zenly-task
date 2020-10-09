package enricher

import (
	"github.com/shekhirin/zenly-task/zenly/pb"
	"github.com/shekhirin/zenly-task/zenly/pb/enricher"
	"math/rand"
)

const maxTransportType = int(enricher.Transport_PLANE)

type transport struct{}

func NewTransport() Enricher {
	return &transport{}
}

func (e transport) Enrich(payload Payload) SetFunc {
	var transport enricher.Transport

	transport.Type = enricher.Transport_Type(rand.Intn(maxTransportType))

	simulateIO()

	return func(gle *pb.GeoLocationEnriched) {
		gle.Transport = &transport
	}
}

func (e transport) String() string {
	return "transport"
}
