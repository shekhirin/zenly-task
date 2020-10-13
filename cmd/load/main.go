package main

import (
	"flag"
	"github.com/shekhirin/zenly-task/load"
	"time"
)

var (
	grpcAddr   = flag.String("grpc-addr", "localhost:8080", "gRPC addr")
	rps        = flag.Int("rps", 10, "RPS")
	duration   = flag.Duration("duration", 1*time.Minute, "Duration (0 for infinite)")
	users = flag.Int("users", 10, "Users amount to publish and subscribe")
)

func main() {
	flag.Parse()

	loader := load.NewLoader(*grpcAddr, *rps, *duration, *users)

	loader.Load()

	loader.PrintStats()
}
