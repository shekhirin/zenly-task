package multisub

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/multierr"
)

type SubscribeUnsubscribe interface {
	Subscribe(subject string, cb nats.MsgHandler) error
	SubscribeMany(subjects []string, cb nats.MsgHandler) error
	Unsubscribe() error
}

type multiSub struct {
	nats        *nats.Conn
	rootSubject string
	subs        []*nats.Subscription
}

func New(nats *nats.Conn, rootSubject string) SubscribeUnsubscribe {
	return &multiSub{
		nats:        nats,
		rootSubject: rootSubject,
	}
}

func (ms *multiSub) Subscribe(subject string, cb nats.MsgHandler) error {
	sub, err := ms.nats.Subscribe(fmt.Sprintf("%s.%s", ms.rootSubject, subject), cb)
	if err != nil {
		return err
	}

	ms.subs = append(ms.subs, sub)

	return nil
}

func (ms *multiSub) SubscribeMany(subjects []string, cb nats.MsgHandler) error {
	for _, subj := range subjects {
		err := ms.Subscribe(subj, cb)
		if err != nil {
			_ = ms.Unsubscribe()
			return err
		}
	}

	return nil
}

func (ms *multiSub) Unsubscribe() error {
	var err error

	var subs = make([]*nats.Subscription, 0)

	for _, sub := range ms.subs {
		subErr := sub.Unsubscribe()
		if subErr != nil {
			err = multierr.Append(err, fmt.Errorf("%s: %w", sub.Subject, subErr))
			subs = append(subs, sub)
		}
	}

	ms.subs = subs

	return err
}
