package nats

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/shekhirin/zenly-task/zenly/bus"
	"github.com/shekhirin/zenly-task/zenly/pb"
)

type natsBus struct {
	nats    *nats.Conn
	subject string
}

func New(natsConn *nats.Conn, subject string) bus.Bus {
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

func (bus natsBus) Subscribe(userIds []int32, messageFunc func(message *pb.BusMessage) error) (context.CancelFunc, error) {
	var ch = make(chan *nats.Msg)

	var userIdsMapping = make(map[int32]bool)
	for _, userId := range userIds {
		userIdsMapping[userId] = true
	}

	sub, err := bus.nats.ChanSubscribe(bus.subject, ch)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = sub.Unsubscribe()
				return
			case msg := <-ch:
				var message pb.BusMessage
				if err := proto.Unmarshal(msg.Data, &message); err != nil {
					continue
				}

				if !userIdsMapping[message.UserId] {
					continue
				}

				if err := messageFunc(&message); err != nil {
					cancel()
				}
			}
		}
	}()

	return cancel, nil
}
