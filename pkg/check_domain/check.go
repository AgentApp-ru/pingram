package check_domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pingram/pkg/log"
	"pingram/pkg/models"
	"pingram/pkg/redis"
	"strings"
	"sync"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var PATH = "health/"

type Health struct {
	Ping      map[string]bool `json:"ping"`
	Databases map[string]bool `json:"databases"`
	Caches    map[string]bool `json:"caches"`
	Celery    map[string]bool `json:"celery"`
}

type Postgres struct {
	MaxConnections     uint8
	CurrentConnections uint8
	Percentage         string
	Status             bool
}

type Stats struct {
	Postgres Postgres        `json:"postgres"`
	Celery   map[string]int  `json:"celery"`
	IcDicts  map[string]bool `json:"ic_dicts"`
}

type PartnerStoreHealth struct {
	Health Health `json:"health"`
	Stats  Stats  `json:"stats"`
}

const THRESHOLD = 0.02

func getErrorStatus(domainName string) bool {
	key := fmt.Sprintf("http-errors.%s", domainName)

	var data []byte
	data, err := redis.GetBytes(key)
	if err != nil {
		if err != redigo.ErrNil {
			log.Logger.Error("redis.get error: %v", err)
		}
		return false
	}
	loadedDomain := new(models.HttpErrorDomain)
	if err := json.Unmarshal(data, loadedDomain); err != nil {
		log.Logger.Error("json decode error: %v", err)
		return false
	}

	return loadedDomain.FiveMinutesErrors > THRESHOLD
}

func checkDomain(domain *models.DomainWithSRV, fails chan<- *models.FailedDomain, wg *sync.WaitGroup) {
	defer wg.Done()

	failedDomain := domain.GetNewFailedDomain()

	errorsStatus := getErrorStatus(failedDomain.Name)
	if errorsStatus {
		failedDomain.Error = "Too many 5xx errors"
		failedDomain.Errors5xxStatus = false
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/%s", domain.Domain, PATH))
	if err != nil {
		var respErr string
		if strings.Contains(err.Error(), "x509") {
			respErr = "Certificate"
		} else if strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
			respErr = "Timeout"
		} else {
			respErr = err.Error()
		}
		failedDomain.PingStatus = false
		failedDomain.Error = respErr

		fails <- failedDomain
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		failedDomain.PingStatus = false
		failedDomain.Error = fmt.Sprintf("Response error: %s", err)

		fails <- failedDomain
		return
	}

	if resp.StatusCode != http.StatusOK {
		failedDomain.PingStatus = false
		failedDomain.Error = fmt.Sprintf("Status: %d", resp.StatusCode)

		fails <- failedDomain
		return
	}

	j := new(PartnerStoreHealth)
	if err := json.Unmarshal(body, j); err != nil {
		failedDomain.PingStatus = false
		failedDomain.Error = fmt.Sprintf("Json decode error: %s", err)

		fails <- failedDomain
		return
	}

	for _, v := range j.Health.Ping {
		if !v {
			failedDomain.Error = "Ping error"
			failedDomain.PingStatus = false
		}
	}

	for _, v := range j.Health.Databases {
		if !v {
			failedDomain.Error = "Databases error"
			failedDomain.DatabasesStatus = false
		}
	}
	postgresStats := j.Stats.Postgres
	if postgresStats.MaxConnections != 0 && !postgresStats.Status {
		msg := fmt.Sprintf("Database connections: %s", postgresStats.Percentage)
		failedDomain.Error = msg
		failedDomain.DatabasesDetails = msg
		failedDomain.DatabasesStatus = false
	}
	for _, v := range j.Health.Caches {
		if !v {
			failedDomain.Error = "Cached error"
			failedDomain.CachesStatus = false
		}
	}
	failedDomain.CeleryDetails = map[string]bool{}
	for k, v := range j.Stats.Celery {
		if v >= 50 {
			failedDomain.CeleryDetails[k] = false
			failedDomain.Error = fmt.Sprintf("Celery error, overloaded queue: %s", k)
			failedDomain.CeleryStatus = false
		}
	}

	failedDomain.ExternalAppsDetails = map[string]bool{}
	for k, v := range j.Stats.IcDicts {
		if !v {
			failedDomain.ExternalAppsDetails[k] = false
			failedDomain.Error = fmt.Sprintf("External apps error: %s", k)
			failedDomain.ExternalAppsStatus = false
		}
	}

	if failedDomain.Error != "" {
		log.Logger.Error(fmt.Sprintf("ERROR 5: %s - %s", failedDomain.Name, failedDomain.Error))
		fails <- failedDomain
	}
}

func CheckDomains(domains []*models.DomainWithSRV) []*models.FailedDomain {
	fails := make(chan *models.FailedDomain)

	var failedDomains []*models.FailedDomain
	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)
		go checkDomain(domain, fails, &wg)
	}

	go func() {
		wg.Wait()
		close(fails)
	}()

	for failDomain := range fails {
		failedDomains = append(failedDomains, failDomain)
	}

	return failedDomains
}
