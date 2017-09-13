package recurrent

import (
	"time"

	"github.com/efritz/glock"
)

type (
	// Scheduler periodically executes the a target function. A scheduler can
	// be signaled to execute the function immediately as often as the user likes
	// via the Signal method (if the scheduler is configured to allow it).
	Scheduler interface {
		// Start the scheduler in a goroutine.
		Start()

		// Stop the scheduler. No additional signals are meaningful. This method
		// must not be called twice.
		Stop()

		// Signal the scheduler to execute the function immediately. This method is
		// always non-blocking, and may be ignored depending on if the scheduler is
		// throttling signals or not. This method must not be called after Stop.
		Signal()
	}

	scheduler struct {
		target   func()
		interval time.Duration
		clock    glock.Clock
		withChan func(f func(chan struct{}))
		quit     chan struct{}
		signal   chan struct{}
	}

	// ConfigFunc is a function used to initialize a new scheduler.
	ConfigFunc func(*scheduler)
)

// NewScheduler creates a new scheduler that will invoke the target function.
func NewScheduler(target func(), configs ...ConfigFunc) Scheduler {
	withChan := func(f func(chan struct{})) {
		quit := make(chan struct{})
		defer close(quit)

		f(hammer(quit))
	}

	scheduler := &scheduler{
		target:   target,
		interval: time.Second,
		clock:    glock.NewRealClock(),
		withChan: withChan,
		quit:     make(chan struct{}),
		signal:   make(chan struct{}, 1),
	}

	for _, config := range configs {
		config(scheduler)
	}

	return scheduler
}

// WithInterval sets the interval at which the scheduler will invoke the
// scheduled function (default is one second).
func WithInterval(interval time.Duration) ConfigFunc {
	return func(s *scheduler) {
		s.interval = interval
	}
}

// WithThrottle sets the minimum duration between two invocations of the
// scheduled function (there is no default minimum).
func WithThrottle(minInterval time.Duration) ConfigFunc {
	return func(s *scheduler) {
		s.withChan = func(f func(chan struct{})) {
			ticker := s.clock.NewTicker(minInterval)
			defer ticker.Stop()

			f(convert(ticker.Chan()))
		}
	}
}

func withClock(clock glock.Clock) ConfigFunc {
	return func(s *scheduler) { s.clock = clock }
}

func (s *scheduler) Start() {
	go func() {
		defer close(s.signal)

		s.withChan(func(c chan struct{}) {
			t := throttle(c, s.signal)

			for {
				select {
				case <-t:
					s.target()

				case <-s.clock.After(s.interval):
					s.Signal()

				case <-s.quit:
					return
				}
			}
		})
	}()
}

func (s *scheduler) Stop() {
	close(s.quit)
}

func (s *scheduler) Signal() {
	select {
	case s.signal <- struct{}{}:
	default:
	}
}

func hammer(quit <-chan struct{}) chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		for {
			select {
			case ch <- struct{}{}:
			case <-quit:
				return
			}
		}
	}()

	return ch
}

func convert(ch1 <-chan time.Time) chan struct{} {
	ch2 := make(chan struct{})

	go func() {
		defer close(ch2)

		for range ch1 {
			ch2 <- struct{}{}
		}
	}()

	return ch2
}

func throttle(ch1 chan struct{}, ch2 chan struct{}) chan struct{} {
	ch3 := make(chan struct{})

	go func() {
		defer close(ch3)

		for range ch1 {
			ch3 <- <-ch2
		}
	}()

	return ch3
}
