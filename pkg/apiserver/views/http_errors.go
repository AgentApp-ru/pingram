package views

import (
	"encoding/json"
	"pingram/pkg/log"
	"pingram/pkg/models"
	"pingram/pkg/redis"

	redigo "github.com/gomodule/redigo/redis"
)

func GetInfoHttpErrors() []*models.HttpErrorDomain {
	var (
		err        error
		domainKeys []string
	)

	domainKeys, err = redis.GetKeys("http-errors.*")
	if err != nil && err != redigo.ErrNil {
		log.Logger.Error("redis.get error: %v", err)
	}

	var (
		data    []byte
		domain  *models.HttpErrorDomain
		domains []*models.HttpErrorDomain
	)

	domains = []*models.HttpErrorDomain{}
	for _, key := range domainKeys {
		data, err = redis.GetBytes(key)
		if err != nil {
			if err != redigo.ErrNil {
				log.Logger.Error("redis.get error: %v", err)
			}
			continue
		}
		domain = new(models.HttpErrorDomain)
		if err := json.Unmarshal(data, domain); err != nil {
			log.Logger.Error("json decode error: %v", err)
			continue
		}
		domains = append(domains, domain)
	}

	return domains
}
