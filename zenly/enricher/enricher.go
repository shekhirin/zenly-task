package enricher

import (
	"fmt"
	"github.com/shekhirin/zenly-task/internal/pb"
	"math/rand"
	"time"
)

const minIOSimulation = 1 * time.Millisecond
const maxIOSimulation = 100 * time.Millisecond

type Enricher interface {
	fmt.Stringer
	Enrich(payload Payload) SetFunc // Pass Payload by value to prevent modification
}

// Additional struct to pass only required fields to enricher
type Payload struct {
	UserId int32
	Lat    float64
	Lng    float64
}

type SetFunc func(gle *pb.GeoLocationEnriched)

func simulateIO() {
	time.Sleep(time.Duration(rand.Intn(int(minIOSimulation+maxIOSimulation)) + int(minIOSimulation)))
}
