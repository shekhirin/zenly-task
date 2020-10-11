package multisub

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/multierr"
	"sync"
)

type MultiSub interface {
	Subscribe(cb nats.MsgHandler, subject ...string) error
	UnsubscribeAll() error
}

type multiSub struct {
	mu          sync.Mutex
	nats        *nats.Conn
	rootSubject string
	queue       string
	subs        []*nats.Subscription
}

func New(nats *nats.Conn, rootSubject string, queue string) MultiSub {
	return &multiSub{
		nats:        nats,
		rootSubject: rootSubject,
		queue:       queue,
	}
}

func (ms *multiSub) Subscribe(cb nats.MsgHandler, subjects ...string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, subject := range subjects {
		sub, err := ms.nats.QueueSubscribe(fmt.Sprintf("%s.%s", ms.rootSubject, subject), ms.queue, cb)
		if err != nil {
			_ = ms.UnsubscribeAll()
			return err
		}

		ms.subs = append(ms.subs, sub)
	}

	return nil
}

func (ms *multiSub) UnsubscribeAll() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

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
