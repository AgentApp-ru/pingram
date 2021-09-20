package views

import (
	"fmt"
	"pingram/pkg/log"
	"pingram/pkg/models"
	"pingram/pkg/redis"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

func getDownedAtMonth(monthlyDomainNamePrefix string) []*models.DomainWithTimeout {
	failedDomains, err := redis.GetKeys(fmt.Sprintf("%s*", monthlyDomainNamePrefix))
	if err != nil {
		if err != redigo.ErrNil {
			log.Logger.Error("redis.get error: %v", err)
		}
	}

	var detailedDomains []*models.DomainWithTimeout
	for _, domain := range failedDomains {
		durationInSeconds, err := redis.GetInt64(domain)
		if err != nil {
			log.Logger.Error("redis.get error: %v", err)
			continue
		}
		detailedDomains = append(detailedDomains, &models.DomainWithTimeout{
			Domain:   domain[len(monthlyDomainNamePrefix):],
			Downtime: durationInSeconds,
		})
	}
	if detailedDomains == nil {
		return []*models.DomainWithTimeout{}
	}

	return detailedDomains
}

func GetDownedAtCurrentMonth() []*models.DomainWithTimeout {
	now := time.Now()
	monthlyDomainNamePrefix := fmt.Sprintf("monthly.%d.%d.", now.Year(), int(now.Month()))

	return getDownedAtMonth(monthlyDomainNamePrefix)
}

func GetDownedAtPreviousMonth() []*models.DomainWithTimeout {
	now := time.Now()
	year, month := now.Year(), int(now.Month())
	if int(now.Month()) == 1 {
		year, month = year-1, 12
	} else {
		month = month - 1
	}
	monthlyDomainNamePrefix := fmt.Sprintf("monthly.%d.%d.", year, month)

	return getDownedAtMonth(monthlyDomainNamePrefix)
}
