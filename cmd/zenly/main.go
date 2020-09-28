package main

import (
	"flag"
	"github.com/nats-io/nats.go"
	"github.com/shekhirin/zenly-task/internal/pb"
	"github.com/shekhirin/zenly-task/internal/zenly"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	addr     = flag.String("addr", ":8080", "Server addr")
	natsAddr = flag.String("nats-addr", ":4222", "NATS addr")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	natsConn, err := nats.Connect(*natsAddr)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}

	zenlyServer := zenly.NewServer(natsConn, zenly.DefaultEnrichers)

	pb.RegisterZenlyService(grpcServer, zenlyServer.Service())

	log.Fatal(grpcServer.Serve(lis))
}
