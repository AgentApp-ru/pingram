package check_domain

import (
	"encoding/json"
	"fmt"
	"pingram/pkg/config"
	"pingram/pkg/log"
	"pingram/pkg/models"
	"pingram/pkg/redis"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

func CreateFailedDomain(domain *models.FailedDomain) (bool, error) {
	created := false
	key := fmt.Sprintf("ps-down.%s", domain.Domain)
	_, err := redis.GetBytes(key)

	if err != nil {
		if err == redigo.ErrNil {
			created = true
		} else {
			log.Logger.Error("redis.get error: %v", err)
			return created, err
		}
	}

	data, err := json.Marshal(domain)
	if err != nil {
		return created, err
	}
	if err = redis.SetBytes(key, data); err != nil {
		log.Logger.Error("redis.create error: %v", err)
	}
	return created, nil
}

func UpdateFailed(domains []*models.FailedDomain) []*models.FailedDomain {
	var newDomains []*models.FailedDomain

	for _, domain := range domains {
		startDownTimeDomain := fmt.Sprintf("start.%s", domain.Domain)
		unixDatetime, err := redis.GetInt64(startDownTimeDomain)
		if err != nil {
			if err == redigo.ErrNil {
				if err = redis.CreateOrUpdateShort(startDownTimeDomain); err != nil {
					log.Logger.Error("redis.create error: %v", err)
				}
			}
			continue
		}

		now := time.Now()
		if now.After(
			time.Unix(unixDatetime, 0).Add(
				time.Duration(config.Settings.DowntimeWithoutAlertsInSeconds) * time.Second),
		) {
			created, err := CreateFailedDomain(domain)
			if err != nil {
				continue
			}
			if created {
				newDomains = append(newDomains, domain)
			}
		}

		if !domain.PingStatus && config.Settings.WorkingHoursStart <= now.Hour() && now.Hour() < config.Settings.WorkingHoursEnd {
			monthlyDomainName := fmt.Sprintf("monthly.%d.%d.%s", now.Year(), int(now.Month()), domain.Domain)
			duration, err := redis.GetInt64(monthlyDomainName)
			if err != nil {
				if err == redigo.ErrNil {
					duration = 0
				} else {
					log.Logger.Error("redis.get error: %v", err)
				}
			}
			duration = duration + int64(config.Settings.TimeoutBetweenChecks)
			if err = redis.CreateOrUpdateFull(monthlyDomainName, duration); err != nil {
				log.Logger.Error("redis.get error: %v", err)
			}
		}
	}

	return newDomains
}
