package app_checkers

import (
	"encoding/json"
)

type (
	CeleryChecker struct {
		health map[string]bool
		stats  map[string]uint
	}
)

func NewCeleryChecker(healthRaw, statsRaw json.RawMessage) *CeleryChecker {
	var (
		health map[string]bool
		stats  map[string]uint
	)
	json.Unmarshal(healthRaw, &health)
	json.Unmarshal(statsRaw, &stats)

	return &CeleryChecker{
		health: health,
		stats:  stats,
	}
}

func (c *CeleryChecker) Check() bool {
	if c.health == nil || c.stats == nil {
		println(3, 1)
		return false
	}

	for worker := range c.health {
		if !c.health[worker] {
			println(3, 2, worker, c.health[worker])
		}
	}

	for worker := range c.stats {
		if c.stats[worker] > 50 {

			println(3, 3, worker, c.stats[worker])
		}
	}

	return true
}
