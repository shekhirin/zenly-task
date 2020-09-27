package main

import (
	"flag"
	"fmt"
	"github.com/shekhirin/zenly-task/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

var (
	addr = flag.String("addr", ":8080", "")
)

type zenlyServer struct {
}

func newServer() *zenlyServer {
	s := &zenlyServer{}
	return s
}

func (s *zenlyServer) Service() *pb.ZenlyService {
	return &pb.ZenlyService{
		Publish:   s.Publish,
		Subscribe: s.Subscribe,
	}
}

func (s *zenlyServer) Publish(stream pb.Zenly_PublishServer) error {
	for {
		publishRequest, err := stream.Recv()
		switch err {
		case nil:
			break
		case io.EOF:
			return stream.SendAndClose(&pb.PublishResponse{
				Success: true,
			})
		default:
			return err
		}

		fmt.Printf("request: %+v\n", publishRequest)
	}
}

func (s *zenlyServer) Subscribe(request *pb.SubscribeRequest, stream pb.Zenly_SubscribeServer) error {
	for {
		subscribeResponse := &pb.SubscribeResponse{
			UserId: rand.Int31n(int32(len(request.UserId))),
			GeoLocation: &pb.GeoLocation{
				Lat: float64(rand.Int31n(180 * 1_000_000)) / 1_000_000 - 90,
				Lng: float64(rand.Int31n(360 * 1_000_000)) / 1_000_000 - 180,
				CreatedAt: timestamppb.Now(),
			},
		}

		if err := stream.Send(subscribeResponse); err != nil {
			return err
		}

		fmt.Printf("response: %+v\n", subscribeResponse)

		time.Sleep(time.Duration(rand.Int31n(500) + 500) * time.Millisecond)
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterZenlyService(grpcServer, newServer().Service())
	panic(grpcServer.Serve(lis))
}
