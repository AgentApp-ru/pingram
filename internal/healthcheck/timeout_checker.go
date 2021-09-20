package healthcheck

import (
	"sync"
	"time"
)

type CheckWithTimeout struct {
	healthCheckers []HealthChecker
	timeout        time.Duration
}

func NewCheckWithTimeout(healthCheckers []HealthChecker, timeout time.Duration) *CheckWithTimeout {
	return &CheckWithTimeout{
		healthCheckers: healthCheckers,
		timeout:        timeout,
	}
}

func (c *CheckWithTimeout) Start() {
	var wg *sync.WaitGroup

	for {
		wg = new(sync.WaitGroup)
		for _, healthChecker := range c.healthCheckers {
			wg.Add(1)
			go check(healthChecker, wg)
		}

		wg.Wait()
		time.Sleep(c.timeout)
	}
}

func check(healthChecker HealthChecker, wg *sync.WaitGroup) {
	defer wg.Done()

	healthChecker.Update()
	healthChecker.CheckAll()
}
