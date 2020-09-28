package bus

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/shekhirin/zenly-task/internal/pb"
)

type natsBus struct {
	nats    *nats.Conn
	subject string
}

func NewNats(natsConn *nats.Conn, subject string) Bus {
	return &natsBus{
		nats:    natsConn,
		subject: subject,
	}
}

func (bus natsBus) Publish(message *pb.BusMessage) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	return bus.nats.Publish(bus.subject, data)
}

func (bus natsBus) Subscribe(userIds []int32) (<-chan *pb.BusMessage, context.CancelFunc, error) {
	var ch = make(chan *pb.BusMessage)
	var natsCh = make(chan *nats.Msg)

	var userIdsMapping = make(map[int32]bool)
	for _, userId := range userIds {
		userIdsMapping[userId] = true
	}

	sub, err := bus.nats.ChanSubscribe(bus.subject, natsCh)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = sub.Unsubscribe()
				return
			case msg := <-natsCh:
				var message pb.BusMessage
				if err := proto.Unmarshal(msg.Data, &message); err != nil {
					continue
				}

				if !userIdsMapping[message.UserId] {
					continue
				}

				ch <- &message
			}
		}
	}()

	return ch, cancel, nil
}
