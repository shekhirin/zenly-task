package bus

import (
	"context"
	"github.com/shekhirin/zenly-task/zenly/pb"
)

type MessageFunc func(message *pb.BusMessage) error

type Bus interface {
	Publish(message *pb.BusMessage) error
	Subscribe(userIds []int32, messageFunc MessageFunc) (context.CancelFunc, error)
}
