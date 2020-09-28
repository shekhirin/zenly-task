package weather

import "golang.org/x/exp/rand"

type Service interface {
	Temperature(lat, lng float64) float64
}

type service struct{}

func New() Service {
	return &service{}
}

func (s service) Temperature(lat, lng float64) float64 {
	return float64(rand.Intn(50*100)-25*100) / 100
}
