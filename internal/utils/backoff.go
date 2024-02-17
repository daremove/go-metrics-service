package utils

import "time"

type Backoff struct {
	step     int
	strategy []time.Duration
}

func NewBackoff() Backoff {
	return Backoff{
		step:     -1,
		strategy: []time.Duration{time.Second, 2 * time.Second, 5 * time.Second},
	}
}

func (b *Backoff) Duration() (time.Duration, bool) {
	b.step++

	if b.step > len(b.strategy)-1 {
		return 0, false
	}

	return b.strategy[b.step], true
}

func (b *Backoff) Reset() {
	b.step = -1
}
