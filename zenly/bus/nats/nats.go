package nats

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/shekhirin/zenly-task/zenly/bus"
	"github.com/shekhirin/zenly-task/zenly/bus/nats/multisub"
	"github.com/shekhirin/zenly-task/zenly/pb"
	log "github.com/sirupsen/logrus"
	"strconv"
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
		return fmt.Errorf("marshal message to proto: %w", err)
	}

	return bus.nats.Publish(fmt.Sprintf("%s.%d", bus.subject, message.UserId), data)
}

func (bus natsBus) Subscribe(userIds []int32) (<-chan *pb.BusMessage, context.CancelFunc, error) {
	ch := make(chan *pb.BusMessage)
	ctx, cancel := context.WithCancel(context.Background())

	sub := multisub.New(bus.nats, bus.subject)
	go func() {
		select {
		case <-ctx.Done():
			_ = sub.UnsubscribeAll()
		}
	}()

	var subjects []string
	for _, userId := range userIds {
		subjects = append(subjects, strconv.Itoa(int(userId)))
	}

	if err := sub.Subscribe(bus.MessageHandler(ch), subjects...); err != nil {
		return nil, nil, err
	}

	return ch, cancel, nil
}

func (bus natsBus) MessageHandler(
	ch chan<- *pb.BusMessage,
) nats.MsgHandler {
	return func(msg *nats.Msg) {
		var message pb.BusMessage
		if err := proto.Unmarshal(msg.Data, &message); err != nil {
			log.WithError(err).Error("unmarshal proto to bus message")
			return
		}

		ch <- &message
	}
}
