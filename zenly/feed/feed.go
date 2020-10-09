package feed

import "github.com/shekhirin/zenly-task/zenly/pb"

type Feed interface {
	Publish(message *pb.FeedMessage) error
}
