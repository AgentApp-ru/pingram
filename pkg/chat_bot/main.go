package chat_bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pingram/pkg/config"
	"pingram/pkg/log"
	"pingram/pkg/models"
)

var UptimeChatChannel = ""

func sendToChat(data map[string]string) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		log.Logger.Error(err)
	}
	resp, err := http.Post(UptimeChatChannel, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Logger.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Logger.Error("not OK from chat")
	}
}

func SendFailedDomainToChat(domain *models.FailedDomain) {
	data := map[string]string{
		"check_name":       fmt.Sprintf("[%s] %s", domain.Server, domain.Domain),
		"current_state":    "DOWN",
		"long_description": domain.Error,
	}
	if config.Settings.SendToChat && domain.IsImportant {
		sendToChat(data)
	} else {
		log.Logger.Info(data)
	}
}

func SendUpDomainToChat(upDomainWithTimeout *models.DomainWithTimeout) {
	data := map[string]string{
		"check_name":       upDomainWithTimeout.Domain,
		"current_state":    "UP",
		"long_description": fmt.Sprintf("downtime: %f min", float64(upDomainWithTimeout.Downtime)/60),
	}

	isImportantDomain := false
	for _, d := range models.Domains {
		if upDomainWithTimeout.Domain == d.Domain {
			isImportantDomain = d.IsImportant
			break
		}
	}

	if config.Settings.SendToChat && isImportantDomain {
		sendToChat(data)
	} else {
		log.Logger.Info(data)
	}
}
