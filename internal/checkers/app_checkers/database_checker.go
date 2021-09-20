package app_checkers

import (
	"encoding/json"
	"fmt"
)

type (
	databaseHealth struct {
		Default bool `json:"default"`
	}

	databaseStats struct {
		MaxConnections     uint   `json:"max_connections"`
		CurrentConnections uint   `json:"current_connections"`
		Percentage         string `json:"percentage"`
		Status             bool   `json:"status"`
	}

	DatabaseChecker struct {
		health *databaseHealth
		stats  *databaseStats
	}
)

func NewDatabaseChecker(healthRaw, statsRaw json.RawMessage) *DatabaseChecker {
	var (
		health *databaseHealth
		stats  *databaseStats
	)
	json.Unmarshal(healthRaw, &health)
	json.Unmarshal(statsRaw, &stats)

	return &DatabaseChecker{
		health: health,
		stats:  stats,
	}
}

func (c *DatabaseChecker) Check() bool {
	if c.health == nil || c.stats == nil {
		println(2, 1)
		return false
	}

	println(2, 2, c.health.Default)
	if !c.health.Default {
		println(2, 2, c.health.Default)
		return false
	}

	println(2, 3, c.stats.Status)
	if !c.stats.Status {
		println(2, 3)
		return false
	}

	fmt.Printf("%d/%d\n", c.stats.CurrentConnections, c.stats.MaxConnections)

	return true
}
