package loyalty

import "time"

type limiter struct {
	C        chan struct{}
	ticker   *time.Ticker
	interval time.Duration
}

func newLimiter(interval time.Duration) *limiter {
	limiter := limiter{
		interval: interval,
		C:        make(chan struct{}),
	}

	go func() {
		limiter.C <- struct{}{}

		limiter.ticker = time.NewTicker(limiter.interval)

		for range limiter.ticker.C {
			limiter.C <- struct{}{}
		}
	}()

	return &limiter
}

func (l *limiter) Close() {
	l.ticker.Stop()
}
