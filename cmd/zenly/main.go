package main

import (
	"flag"
	"github.com/Shopify/sarama"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shekhirin/zenly-task/zenly"
	natsBus "github.com/shekhirin/zenly-task/zenly/bus/nats"
	kafkaFeed "github.com/shekhirin/zenly-task/zenly/feed/kafka"
	"github.com/shekhirin/zenly-task/zenly/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/benchmark/flags"
	"net"
	"net/http"
)

var (
	env             = flag.String("env", "debug", "App environment")
	addr            = flag.String("addr", ":8080", "Server addr")
	diagnosticsAddr = flag.String("diagnostics-addr", ":7777", "Diagnostics addr")
	natsAddr        = flag.String("nats-addr", ":4222", "NATS addr")
	natsSubject     = flag.String("nats-subject", "zenly", "NATS subject")
	kafkaBrokers    = flags.StringSlice("kafka-brokers", []string{":9092"}, "Kafka brokers (comma separated list)")
	kafkaTopic      = flag.String("kafka-topic", "feed", "Kafka topic")
)

var logEntry = log.NewEntry(log.StandardLogger())

func init() {
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

	bus := natsBus.New(natsConn, *natsSubject)

	kafkaProducer, err := sarama.NewAsyncProducer(*kafkaBrokers, sarama.NewConfig())
	if err != nil {
		log.WithError(err).Fatalf("connect to kafkaProducer on %+v", *kafkaBrokers)
	}

	feed := kafkaFeed.New(kafkaProducer, *kafkaTopic)

	zenlyService := zenly.New(bus, feed, zenly.DefaultEnrichers).Service()

	pb.RegisterZenlyService(grpcServer, zenlyService)

	grpcPrometheus.Register(grpcServer)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(*diagnosticsAddr, http.DefaultServeMux)
		log.WithError(err).Fatalf("serve diagnostics server on %s", *diagnosticsAddr)
	}()

	log.WithError(grpcServer.Serve(lis)).Fatal("serve grpc server")
}
