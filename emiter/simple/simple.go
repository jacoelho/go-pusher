package simple

import (
	"time"
)

type Simple struct {
	interval time.Duration
	events   chan []byte
}

type Option func(*Simple) error

func SetInterval(t time.Duration) Option {
	return func(s *Simple) error {
		s.interval = t
		return nil
	}
}

func newWithDefaults() *Simple {
	return &Simple{
		interval: time.Second * 5,
		events:   make(chan []byte, 256),
	}
}

func New(opts ...Option) (*Simple, error) {
	s := newWithDefaults()

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return s, err
		}
	}

	return s, nil
}

func (s *Simple) Watch(channel string) {
	ticker := time.NewTicker(s.interval)

	go func() {

		for {
			select {
			case <-ticker.C:
				s.events <- []byte(channel)
			}
		}
		ticker.Stop()
		close(s.events)
	}()
}

func (s *Simple) Events() <-chan []byte {
	return s.events
}
