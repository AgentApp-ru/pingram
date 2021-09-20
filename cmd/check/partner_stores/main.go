package main

import (
	"pingram/internal/healthcheck"
	"pingram/internal/healthcheck/partner_store"
	"time"
)

func main() {
	domains := []string{
		"localhost:8000",
	}

	healthcheckers := make([]healthcheck.HealthChecker, len(domains))

	for i, domain := range domains {
		healthcheckers[i] = partner_store.New(domain)
	}

	healthcheck.NewCheckWithTimeout(
		healthcheckers,
		time.Duration(15)*time.Second,
	).Start()
}
