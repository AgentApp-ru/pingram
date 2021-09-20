package views

import (
	"pingram/pkg/models"
)

func GetAllDomains() []*models.DomainWithSRV {
	return models.Domains
}
