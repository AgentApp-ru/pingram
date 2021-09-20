package app_checkers

import (
	"encoding/json"
)

type (
	cachesHealth struct {
		Default bool `json:"default"`
	}

	CachesChecker struct {
		health *cachesHealth
		stats  struct{}
	}
)

func NewCachesChecker(healthRaw, statsRaw json.RawMessage) *CachesChecker {
	var (
		health *cachesHealth
	)
	json.Unmarshal(healthRaw, &health)

	return &CachesChecker{
		health: health,
		stats:  struct{}{},
	}
}

func (c *CachesChecker) Check() bool {
	if c.health == nil {
		println(5, 1)
		return false
	}

	if !c.health.Default {
		println(5, 2, c.health.Default)
		return false
	}

	return true
}
