package main

import (
	"pingram/internal/healthcheck"
	"pingram/internal/healthcheck/domain"
	"time"
)

func main() {
	domainGroups := [][]string{
		{"localhost:8000"},
	}

	healthcheckers := make([]healthcheck.HealthChecker, len(domainGroups))

	for i, domains := range domainGroups {
		healthcheckers[i] = domain.New(domains)
	}

	healthcheck.NewCheckWithTimeout(
		healthcheckers,
		time.Duration(24)*time.Hour,
	).Start()
}
