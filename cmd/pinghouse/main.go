package main

import (
	"pingram/pkg/chat_bot"
	"pingram/pkg/check_domain"
	"pingram/pkg/config"
	"pingram/pkg/log"
	"pingram/pkg/models"
	"time"
)

func main() {
	allDomains := models.Domains

	for {
		var plainAllDomains []string
		for _, domain := range allDomains {
			plainAllDomains = append(plainAllDomains, domain.Domain)
		}

		failedDomains := check_domain.CheckDomains(allDomains)
		var plainFailedDomains []string
		for _, domain := range failedDomains {
			plainFailedDomains = append(plainFailedDomains, domain.Domain)
		}
		log.Logger.Info(plainFailedDomains)

		newFailedDomains := check_domain.UpdateFailed(failedDomains)
		for _, failedDomain := range newFailedDomains {
			chat_bot.SendFailedDomainToChat(failedDomain)
		}

		newUpDomains := check_domain.GetNewUp(plainAllDomains, plainFailedDomains)
		for _, upDomain := range newUpDomains {
			chat_bot.SendUpDomainToChat(upDomain)
		}

		time.Sleep(time.Duration(config.Settings.TimeoutBetweenChecks) * time.Second)
	}
}
