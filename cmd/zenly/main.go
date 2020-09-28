package main

import (
	"flag"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/nats-io/nats.go"
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/zenly"
	"github.com/shekhirin/zenly-task/internal/zenly/bus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var (
	env        = flag.String("env", "debug", "App environment")
	addr       = flag.String("addr", ":8080", "Server addr")
	natsAddr   = flag.String("nats-addr", ":4222", "NATS addr")
	busSubject = flag.String("bus-subject", "zenly", "Bus subject")
)

func main() {
	flag.Parse()

	if *env == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.WithError(err).Fatal("listen")
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(
			grpcRecovery.StreamServerInterceptor(),
		),
	)

	natsConn, err := nats.Connect(*natsAddr)
	if err != nil {
		log.WithError(err).Fatal("connect to NATS")
	}

	natsBus := bus.NewNats(natsConn, *busSubject)

	zenlyServer := zenly.NewServer(natsBus, zenly.DefaultEnrichers)

	pb.RegisterZenlyService(grpcServer, zenlyServer.Service())

	log.WithError(grpcServer.Serve(lis)).Fatal("serve grpc server")
}
