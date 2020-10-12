package stats

import (
	"github.com/shekhirin/zenly-task/internal/pb"
	"sync"
	"time"
)

type Subscribe struct {
	messagesPerSecond map[time.Time]int
	lagSum            time.Duration
	mu                sync.Mutex
}

func NewSubscribe() Subscribe {
	return Subscribe{}
}

func (s *Subscribe) Observe(received time.Time, message *pb.SubscribeResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.messagesPerSecond == nil {
		s.messagesPerSecond = make(map[time.Time]int)
	}

	s.messagesPerSecond[received.Truncate(time.Second)] += 1
	s.lagSum += received.Sub(message.GeoLocation.GeoLocation.CreatedAt.AsTime())
}

func (s *Subscribe) LagAverage() time.Duration {
	if s.Messages() == 0 {
		return time.Duration(-1)
	}

	return s.lagSum / time.Duration(s.Messages())
}

func (s *Subscribe) Messages() int {
	var messages int
	for _, count := range s.messagesPerSecond {
		messages += count
	}
	return messages
}

func (s *Subscribe) MPS() float64 {
	var sum, total float64
	for _, count := range s.messagesPerSecond {
		sum += float64(count)
		total += 1
	}

	if total == 0 {
		return -1
	}

	return sum / total
}

