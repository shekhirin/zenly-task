package zenly

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shekhirin/zenly-task/zenly/enricher"
	"github.com/shekhirin/zenly-task/zenly/metrics"
	"github.com/shekhirin/zenly-task/zenly/pb"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	EnricherTimeout = 100 * time.Millisecond

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

				enrich := enricher.Enrich(payload)
				enricherElapsed := time.Since(enricherStart)

				metrics.EnricherTimeMS.With(prometheus.Labels{
					"enricher": enricher.String(),
					"timeout": fmt.Sprintf("%t", enricherElapsed > EnricherTimeout),
				}).Observe(float64(enricherElapsed.Milliseconds()))

				log.WithFields(log.Fields{
					"enricher":   enricher.String(),
					"timeout": fmt.Sprintf("%t", enricherElapsed > EnricherTimeout),
					"elapsed_ms": enricherElapsed.Milliseconds(),
				}).Debug("finish enricher")

				// We don't give control of the context to enricher because of the possibility of forgetting
				// to check timeout before setting the submessage inside the enricher.
				// So there we enrich only if context is not timeout-ed
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
