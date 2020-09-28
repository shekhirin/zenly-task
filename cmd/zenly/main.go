package main

import (
	"flag"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shekhirin/zenly-task/zenly"
	"github.com/shekhirin/zenly-task/zenly/bus"
	"github.com/shekhirin/zenly-task/zenly/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

var (
	env             = flag.String("env", "debug", "App environment")
	addr            = flag.String("addr", ":8080", "Server addr")
	diagnosticsAddr = flag.String("diagnostics-addr", ":8081", "Diagnostics addr")
	natsAddr        = flag.String("nats-addr", ":4222", "NATS addr")
	busSubject      = flag.String("bus-subject", "zenly", "Bus subject")
)

var (
	enricherTimeMS = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "enricher_time_ms",
		Buckets: []float64{25, 50, 75, 100},
	}, []string{"enricher"})
)

var logEntry = log.NewEntry(log.StandardLogger())

func init() {
	prometheus.MustRegister(enricherTimeMS)

	grpcLogrus.ReplaceGrpcLogger(logEntry)
}

func main() {
	flag.Parse()

	if *env == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.WithError(err).Fatalf("listen tcp on %s", *addr)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcPrometheus.StreamServerInterceptor,
			grpcLogrus.StreamServerInterceptor(logEntry),
			grpcRecovery.StreamServerInterceptor(),
		)),
	)

	natsConn, err := nats.Connect(*natsAddr)
	if err != nil {
		log.WithError(err).Fatalf("connect to NATS on %s", *natsAddr)
	}

	natsBus := bus.NewNats(natsConn, *busSubject)

	zenlyService := zenly.New(natsBus, enricherTimeMS, zenly.DefaultEnrichers).Service()

	pb.RegisterZenlyService(grpcServer, zenlyService)

	grpcPrometheus.Register(grpcServer)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(*diagnosticsAddr, http.DefaultServeMux)
		log.WithError(err).Fatalf("serve diagnostics server on %s", *diagnosticsAddr)
	}()

	log.WithError(grpcServer.Serve(lis)).Fatal("serve grpc server")
}
