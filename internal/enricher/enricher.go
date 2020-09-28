package enricher

import (
	"context"
	"github.com/shekhirin/zenly-task/internal/pb"
)

type Enricher interface {
	Enrich(ctx context.Context, gle *pb.GeoLocationEnriched)
}
