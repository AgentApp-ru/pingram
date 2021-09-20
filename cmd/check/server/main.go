package main

import (
	"pingram/internal/healthcheck"
	"pingram/internal/healthcheck/server"
	"time"
)

func main() {
	servers := []string{
	}

	healthcheckers := make([]healthcheck.HealthChecker, len(servers))

	for i, address := range servers {
		healthcheckers[i] = server.New(address)
	}

	healthcheck.NewCheckWithTimeout(
		healthcheckers,
		time.Duration(30)*time.Minute,
	).Start()
}
