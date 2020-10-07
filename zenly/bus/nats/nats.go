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

func (bus natsBus) Subscribe(userIds []int32) (<-chan *pb.BusMessage, context.CancelFunc, error) {
	ch := make(chan *pb.BusMessage)
	ctx, cancel := context.WithCancel(context.Background())

	var userIdsMapping = make(map[int32]bool)
	for _, userId := range userIds {
		userIdsMapping[userId] = true
	}

	sub, err := bus.nats.Subscribe(bus.subject, bus.MessageHandler(userIdsMapping, ch))
	if err != nil {
		return nil, nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			_ = sub.Unsubscribe()
		}
	}()

	return ch, cancel, nil
}

func (bus natsBus) MessageHandler(
	userIdsMapping map[int32]bool,
	ch chan<- *pb.BusMessage,
) nats.MsgHandler {
	return func(msg *nats.Msg) {
		var message pb.BusMessage
		if err := proto.Unmarshal(msg.Data, &message); err != nil {
			return
		}

		// TODO: implement more convenient user id routing
		if !userIdsMapping[message.UserId] {
			return
		}

		ch <- &message
	}
}
