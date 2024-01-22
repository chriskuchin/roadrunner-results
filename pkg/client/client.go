package client

import (
	"strings"

	"github.com/chriskuchin/roadrunner-results/pkg/client/athletic_net"
	"github.com/chriskuchin/roadrunner-results/pkg/client/meettrax"
	"github.com/chriskuchin/roadrunner-results/pkg/client/mile_split"
	"github.com/chriskuchin/roadrunner-results/pkg/client/model"
)

func GetEventInformation(eventURL string) []model.Event {
	if strings.Contains(eventURL, "athletic.net") {
		return athletic_net.GetEventInformation(eventURL)
	} else if strings.Contains(eventURL, "milesplit.com") && strings.Contains(eventURL, "type=raw") {
		return mile_split.GetEventInformation(eventURL)
	} else if strings.Contains(eventURL, "meettrax.com") {
		return meettrax.GetEventInformation(eventURL)
	}

	return nil
}
