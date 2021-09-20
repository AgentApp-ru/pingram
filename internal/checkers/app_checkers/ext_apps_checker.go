package app_checkers

import (
	"encoding/json"
)

type (
	ExtAppsChecker struct {
		health struct{}
		stats  map[string]bool
	}
)

func NewExtAppsChecker(healthRaw, statsRaw json.RawMessage) *ExtAppsChecker {
	var (
		stats map[string]bool
	)
	json.Unmarshal(statsRaw, &stats)

	return &ExtAppsChecker{
		health: struct{}{},
		stats:  stats,
	}
}

func (c *ExtAppsChecker) Check() bool {
	if c.stats == nil {
		println(4, 1)
		return false
	}

	for app := range c.stats {
		if !c.stats[app] {
			println(4, 2, app, c.stats[app])
		}
	}

	return true
}
