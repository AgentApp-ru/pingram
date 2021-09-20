package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pingram/internal/checkers/server_checkers"
	"pingram/internal/healthcheck"
	"sync"
	"time"
)

type ServerHealthChecker struct {
	address  string
	checkers []healthcheck.Checker
}

func New(address string) *ServerHealthChecker {
	return &ServerHealthChecker{
		address:  address,
		checkers: make([]healthcheck.Checker, 0),
	}
}

func (s *ServerHealthChecker) Update() {
	s.checkers = make([]healthcheck.Checker, 0)

	health := s.getStats()

	s.AddChecker(server_checkers.NewSpaceChecker(health.Space))
	// s.AddChecker(server_checkers.NewUptimeChecker(health.Uptime))
}

func (s *ServerHealthChecker) AddChecker(checker healthcheck.Checker) {
	s.checkers = append(s.checkers, checker)
}

func (s *ServerHealthChecker) CheckAll() {
	wg := new(sync.WaitGroup)

	for _, checker := range s.checkers {
		wg.Add(1)
		go healthcheck.WaitForCheck(checker.Check, wg)
	}

	wg.Wait()
}

func (s *ServerHealthChecker) getStats() *health {
	client := http.Client{Timeout: 50 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s/health", s.address))
	if err != nil {
		println(1, err.Error())
		return &health{}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(212, err.Error())
		return &health{}
	}

	if resp.StatusCode != http.StatusOK {
		println(3, resp.StatusCode)
		return &health{}
	}

	healthContainer := new(health)
	if err := json.Unmarshal(body, healthContainer); err != nil {
		println(4)
		return &health{}
	}

	return healthContainer
}
