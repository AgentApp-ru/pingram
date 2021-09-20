package check_domain

import (
	"fmt"
	"pingram/pkg/config"
	"pingram/pkg/log"
	"pingram/pkg/models"
	"pingram/pkg/redis"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

func contains(array []string, elem string) bool {
	for _, x := range array {
		if elem == x {
			return true
		}
	}
	return false
}

func GetNewUp(allDomain, failedDomains []string) []*models.DomainWithTimeout {
	var (
		upDomains         []string
		extendedUpDomains []*models.DomainWithTimeout
		err               error
		unixTime          int64
	)

	for _, domain := range allDomain {
		if !contains(failedDomains, domain) {
			upDomains = append(upDomains, domain)
		}
	}

	for _, domain := range upDomains {
		startDownTimeDomain := fmt.Sprintf("start.%s", domain)
		unixTime, err = redis.GetInt64(startDownTimeDomain)
		if err != nil {
			if err != redigo.ErrNil {
				log.Logger.Error("redis.get error: %v", err)
			}
			continue
		}
		if err = redis.Del(startDownTimeDomain); err != nil {
			if err != redigo.ErrNil {
				log.Logger.Error("redis.del error: %v", err)
			}
		}
		PSDownDomain := fmt.Sprintf("ps-down.%s", domain)
		if err = redis.Del(PSDownDomain); err != nil {
			if err != redigo.ErrNil {
				log.Logger.Error("redis.del error: %v", err)
			}
		}

		totalDownTime := time.Now().Unix() - unixTime
		if totalDownTime >= int64(config.Settings.DowntimeWithoutAlertsInSeconds) {
			extendedUpDomains = append(extendedUpDomains, &models.DomainWithTimeout{
				Domain:   domain,
				Downtime: totalDownTime,
			})
		}

		if err = redis.Del(domain); err != nil {
			log.Logger.Error("redis.get error: %v", err)
		}
	}

	return extendedUpDomains
}
