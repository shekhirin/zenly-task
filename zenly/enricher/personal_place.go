package enricher

import (
	"github.com/shekhirin/zenly-task/zenly/pb"
	"github.com/shekhirin/zenly-task/zenly/pb/enricher"
	"math/rand"
)

const maxPersonalPlaceType = int(enricher.PersonalPlace_SCHOOL)

type personalPlace struct{}

func NewPersonalPlace() Enricher {
	return &personalPlace{}
}

func (e personalPlace) Enrich(payload Payload) SetFunc {
	var personalPlace enricher.PersonalPlace

	personalPlace.Type = enricher.PersonalPlace_Type(rand.Intn(maxPersonalPlaceType))

	simulateIO()

	return func(gle *pb.GeoLocationEnriched) {
		gle.PersonalPlace = &personalPlace
	}
}

func (e personalPlace) String() string {
	return "personal_place"
}
