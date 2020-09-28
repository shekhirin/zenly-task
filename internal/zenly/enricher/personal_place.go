package enricher

import (
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/pb/enrichers"
	"math/rand"
)

const maxPersonalPlaceType = int(enrichers.PersonalPlace_PERSONAL_PLACE_SCHOOL)

type personalPlace struct{}

func NewPersonalPlace() Enricher {
	return &personalPlace{}
}

func (e personalPlace) Enrich(payload Payload) SetFunc {
	var personalPlace enrichers.PersonalPlace

	personalPlace.Type = enrichers.PersonalPlace_Type(rand.Intn(maxPersonalPlaceType))

	simulateIO()

	return func(gle *pb.GeoLocationEnriched) {
		gle.PersonalPlace = &personalPlace
	}
}

func (e personalPlace) String() string {
	return "personal_place"
}
