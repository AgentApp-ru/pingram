package views

import (
	"encoding/json"
	"pingram/pkg/log"
	"pingram/pkg/models"
	"pingram/pkg/redis"

	redigo "github.com/gomodule/redigo/redis"
)

func GetDowned() []*models.FailedDomain {
	var (
		err              error
		failedDomainKeys []string
	)

	failedDomainKeys, err = redis.GetKeys("ps-down.*")
	if err != nil && err != redigo.ErrNil {
		log.Logger.Error("redis.get error: %v", err)
	}

	var (
		data          []byte
		failedDomain  *models.FailedDomain
		failedDomains []*models.FailedDomain
	)

	failedDomains = []*models.FailedDomain{}
	for _, domain := range failedDomainKeys {
		data, err = redis.GetBytes(domain)
		if err != nil {
			if err != redigo.ErrNil {
				log.Logger.Error("redis.get error: %v", err)
			}
			continue
		}
		failedDomain = new(models.FailedDomain)
		if err := json.Unmarshal(data, failedDomain); err != nil {
			log.Logger.Error("json decode error: %v", err)
			continue
		}
		failedDomains = append(failedDomains, failedDomain)
	}

	return failedDomains
}
