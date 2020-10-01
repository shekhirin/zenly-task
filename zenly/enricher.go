package zenly

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	"github.com/shekhirin/zenly-task/zenly/metrics"
	"github.com/shekhirin/zenly-task/zenly/pb"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	EnrichComplete = "complete"
	EnrichTimeout  = "timeout"
)

func (z *Zenly) Enrich(payload enricher.Payload, gle *pb.GeoLocationEnriched) {
	var wg sync.WaitGroup
	wg.Add(len(z.enrichers))
	waitCh := make(chan struct{})

	ctx, _ := context.WithTimeout(context.Background(), EnricherTimeout)

	reasonCh := make(chan string, 1)

	go func() {
		start := time.Now()

		for _, targetEnricher := range z.enrichers {
			go func(enricher enricher.Enricher) {
				defer wg.Done()

				enricherStart := time.Now()

				// Don't give control of the context to enricher because of the possibility of forgetting
				// to check timeout before setting the submessage inside the enricher
				enrich := enricher.Enrich(payload)
				enricherElapsed := time.Since(enricherStart)

				metrics.EnricherTimeMS.With(prometheus.Labels{
					"enricher": enricher.String(),
				}).Observe(float64(enricherElapsed.Milliseconds()))
				log.WithFields(log.Fields{
					"enricher":   enricher.String(),
					"elapsed_ms": enricherElapsed.Milliseconds(),
				}).Debug("finish enricher")

				if ctx.Err() == nil {
					enrich(gle)
				}
			}(targetEnricher)
		}

		wg.Wait()
		close(waitCh)

		elapsed := time.Since(start)

		reason := <-reasonCh
		metrics.EnrichFinishMS.With(prometheus.Labels{"reason": reason}).Observe(float64(elapsed.Milliseconds()))
		log.WithFields(log.Fields{"reason": reason, "elapsed_ms": elapsed}).Debug("finish enrich")
	}()

	select {
	case <-ctx.Done():
		reasonCh <- EnrichTimeout
	case <-waitCh:
		reasonCh <- EnrichComplete
	}
}
