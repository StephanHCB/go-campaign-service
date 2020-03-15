package main

import (
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/logging"
	"github.com/StephanHCB/go-campaign-service/web"
)

func main() {
	logging.Setup()
	configuration.Setup()
	logging.PostConfigSetup()

	web.Serve()
}
