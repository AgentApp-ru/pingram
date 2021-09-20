package partner_store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pingram/internal/checkers/app_checkers"
	"pingram/internal/healthcheck"
	"sync"
	"time"
)

type PSHealthChecker struct {
	mainDomain string
	checkers   []healthcheck.Checker
}

func New(domain string) *PSHealthChecker {
	return &PSHealthChecker{
		mainDomain: domain,
		checkers:   make([]healthcheck.Checker, 0),
	}
}

func (ps *PSHealthChecker) Update() {
	ps.checkers = make([]healthcheck.Checker, 0)

	health, stats := ps.getHealthAndStats()
	ps.AddChecker(app_checkers.NewDatabaseChecker(health["databases"], stats["postgres"]))
	ps.AddChecker(app_checkers.NewCeleryChecker(health["celery"], stats["celery"]))
	ps.AddChecker(app_checkers.NewCachesChecker(health["caches"], nil))
	ps.AddChecker(app_checkers.NewExtAppsChecker(nil, stats["ext_apps"]))
	ps.AddChecker(app_checkers.NewNginxErrorsChecker(nil, stats["50x_errors"]))
}

func (ps *PSHealthChecker) AddChecker(checker healthcheck.Checker) {
	ps.checkers = append(ps.checkers, checker)
}

func (ps *PSHealthChecker) CheckAll() {
	wg := new(sync.WaitGroup)

	for _, checker := range ps.checkers {
		wg.Add(1)
		go healthcheck.WaitForCheck(checker.Check, wg)
	}

	wg.Wait()
}

func (ps *PSHealthChecker) getHealthAndStats() (map[string]json.RawMessage, map[string]json.RawMessage) {
	client := http.Client{Timeout: 50 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://%s/health/", ps.mainDomain))
	if err != nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	healthContainer := new(health)
	if err := json.Unmarshal(body, healthContainer); err != nil {
		return nil, nil
	}

	return healthContainer.Health, healthContainer.Stats
}
