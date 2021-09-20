package models

import (
	"bufio"
	"fmt"
	"os"
	"pingram/pkg/config"
	"pingram/pkg/log"
	"strings"
)

type DomainWithSRV struct {
	Server      string
	Name        string
	IsImportant bool
	Domain      string
}

type DomainWithTimeout struct {
	Domain   string `json:"domain"`
	Downtime int64  `json:"downtime"`
}

type HttpErrorDomain struct {
	DomainWithSRV
	NginxReferer      string
	FiveMinutesErrors float64
	DayErrors         float64
	WeekErrors        float64
}

func (d *DomainWithSRV) GetNewHttpErrorDomain() *HttpErrorDomain {
	return &HttpErrorDomain{
		DomainWithSRV: *d,
		NginxReferer:  fmt.Sprintf("%s.host", d.Name),
		DayErrors:     0,
		WeekErrors:    0,
	}
}

type FailedDomain struct {
	DomainWithSRV
	Error               string
	PingStatus          bool
	Errors5xxStatus     bool
	DatabasesStatus     bool
	DatabasesDetails    string
	CachesStatus        bool
	CeleryStatus        bool
	CeleryDetails       map[string]bool
	ExternalAppsStatus  bool
	ExternalAppsDetails map[string]bool
}

func (d *DomainWithSRV) GetNewFailedDomain() *FailedDomain {
	return &FailedDomain{
		DomainWithSRV:   *d,
		PingStatus:      true,
		Errors5xxStatus: true,
		DatabasesStatus: true,
		CachesStatus:    true,
		CeleryStatus:    true,
	}
}

var Domains []*DomainWithSRV

func GetAllDomains() []*DomainWithSRV {
	file, err := os.Open(config.Settings.FileName)
	if err != nil {
		log.Logger.Fatal(err)
	}
	defer file.Close()

	var allDomains []*DomainWithSRV
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ps := strings.Split(scanner.Text(), ";")

		allDomains = append(allDomains, &DomainWithSRV{
			Server:      ps[0],
			Name:        ps[1],
			IsImportant: ps[2] == "True",
			Domain:      fmt.Sprintf("%s", ps[1]),
		})
	}
	if err := scanner.Err(); err != nil {
		log.Logger.Fatal(err)
	}

	return allDomains
}

func init() {
	Domains = GetAllDomains()
}
