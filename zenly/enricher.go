package zenly

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	"github.com/shekhirin/zenly-task/zenly/metrics"
	"github.com/shekhirin/zenly-task/zenly/pb"
	"sync"
	"time"
)

func (z *Zenly) Enrich(payload enricher.Payload, gle *pb.GeoLocationEnriched) {
	var wg sync.WaitGroup
	wg.Add(len(z.enrichers))
	waitCh := make(chan struct{})

	ctx, _ := context.WithTimeout(context.Background(), EnricherTimeout)

	go func() {
		for _, serverEnricher := range z.enrichers {
			go func(enricher enricher.Enricher) {
				defer wg.Done()

				start := time.Now()

				// Don't give control of the context to enricher because of the possibility of forgetting
				// to check timeout before setting the submessage inside the enricher
				enrich := enricher.Enrich(payload)
				elapsed := time.Since(start)

				metrics.EnricherTimeMS.With(prometheus.Labels{"enricher": enricher.String()}).Observe(float64(elapsed.Milliseconds()))

				if ctx.Err() == nil {
					enrich(gle)
				}
			}(serverEnricher)
		}

		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-ctx.Done():
	case <-waitCh:
	}
}
