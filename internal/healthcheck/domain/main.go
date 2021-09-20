package domain

import (
	"pingram/internal/checkers/domain_checkers"
	"pingram/internal/healthcheck"
	"sync"
)

type DomainHealthChecker struct {
	mainDomain string
	domains    []string
	checkers   []healthcheck.Checker
}

func New(domains []string) *DomainHealthChecker {
	return &DomainHealthChecker{
		mainDomain: domains[0],
		domains:    domains,
		checkers:   make([]healthcheck.Checker, 0),
	}
}

func (d *DomainHealthChecker) Update() {
	d.checkers = make([]healthcheck.Checker, 0)

	d.AddChecker(domain_checkers.NewSSLChecker(d.domains))
}

func (d *DomainHealthChecker) AddChecker(checker healthcheck.Checker) {
	d.checkers = append(d.checkers, checker)
}

func (d *DomainHealthChecker) CheckAll() {
	wg := new(sync.WaitGroup)

	for _, checker := range d.checkers {
		wg.Add(1)
		go healthcheck.WaitForCheck(checker.Check, wg)
	}

	wg.Wait()
}
