package stats

import (
	"sync"
	"time"
)

type Publish struct {
	requestsPerSecond map[time.Time]int
	ElapsedMin        time.Duration
	ElapsedMax        time.Duration
	elapsedSum        time.Duration
	mu                sync.Mutex
}

func NewPublish() Publish {
	return Publish{
		requestsPerSecond: make(map[time.Time]int),
	}
}

func (s *Publish) Observe(start time.Time, finish time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	elapsed := finish.Sub(start)

	if s.requestsPerSecond == nil {
		s.requestsPerSecond = make(map[time.Time]int)
	}

	s.requestsPerSecond[finish.Truncate(time.Second)]++
	if s.ElapsedMin == 0 || elapsed < s.ElapsedMin {
		s.ElapsedMin = elapsed
	}
	if s.ElapsedMax == 0 || elapsed > s.ElapsedMax {
		s.ElapsedMax = elapsed
	}
	s.elapsedSum += elapsed
}

func (s *Publish) ElapsedAverage() time.Duration {
	if s.Requests() == 0 {
		return time.Duration(-1)
	}

	return s.elapsedSum / time.Duration(s.Requests())
}

func (s *Publish) Requests() int {
	var requests int
	for _, count := range s.requestsPerSecond {
		requests += count
	}
	return requests
}

func (s *Publish) RPS() float64 {
	var sum, total float64
	for _, count := range s.requestsPerSecond {
		sum += float64(count)
		total += 1
	}

	if total == 0 {
		return -1
	}

	return sum / total
}
