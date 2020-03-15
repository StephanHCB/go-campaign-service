package main

import (
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database"
	"github.com/StephanHCB/go-campaign-service/internal/repository/logging"
	"github.com/StephanHCB/go-campaign-service/web"
)

func main() {
	logging.Setup()
	configuration.Setup()
	logging.PostConfigSetup()

	database.Open()
	defer database.Close()
	database.MigrateIfSwitchedOn()

	web.Serve()
}
