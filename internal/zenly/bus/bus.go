package bus

import (
	"context"
	"github.com/shekhirin/zenly-task/internal/pb"
)

type Bus interface {
	Publish(message *pb.BusMessage) error
	Subscribe(userIds []int32) (<-chan *pb.BusMessage, context.CancelFunc, error)
}
