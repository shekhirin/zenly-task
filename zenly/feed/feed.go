package feed

import "github.com/shekhirin/zenly-task/internal/pb"

type Feed interface {
	Publish(message *pb.FeedMessage) error
}
