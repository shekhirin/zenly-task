package enricher

import (
	"context"
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/pb/enrichers"
	"math/rand"
)

const maxPersonalPlaceType = int(enrichers.PersonalPlace_PERSONAL_PLACE_SCHOOL)

type personalPlaceEnricher struct{}

func NewPersonalPlace() Enricher {
	return &personalPlaceEnricher{}
}

func (e personalPlaceEnricher) Enrich(ctx context.Context, gle *pb.GeoLocationEnriched) {
	var personalPlace enrichers.PersonalPlace

	personalPlace.Type = enrichers.PersonalPlace_Type(rand.Intn(maxPersonalPlaceType))

	gle.PersonalPlace = &personalPlace
}
