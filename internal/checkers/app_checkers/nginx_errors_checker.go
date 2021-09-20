package app_checkers

import (
	"encoding/json"
)

type (
	errorsStats struct {
		Total      uint    `json:"total"`
		Errors     uint    `json:"errors"`
		Percentage float64 `json:"percentage"`
		Status     bool    `json:"status"`
	}

	nginxStats struct {
		FiveMinErrors *errorsStats `json:"5m"`
		OneDayErrors  *errorsStats `json:"1d"`
		OneWeekError  *errorsStats `json:"1w"`
	}

	NginxErrorsChecker struct {
		health struct{}
		stats  *nginxStats
	}
)

func NewNginxErrorsChecker(healthRaw, statsRaw json.RawMessage) *NginxErrorsChecker {
	var (
		stats *nginxStats
	)
	json.Unmarshal(statsRaw, &stats)

	return &NginxErrorsChecker{
		health: struct{}{},
		stats:  stats,
	}
}

func (c *NginxErrorsChecker) Check() bool {
	if c.stats == nil {
		println(6, 1)
		return false
	}

	if !c.stats.FiveMinErrors.Status {
		println(6, 2, c.stats.FiveMinErrors.Status)
		return false
	}

	return true
}
