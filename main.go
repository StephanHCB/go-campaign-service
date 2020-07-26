package main

import (
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database"
	"github.com/StephanHCB/go-campaign-service/internal/repository/logging"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice/mailserviceclient"
	"github.com/StephanHCB/go-campaign-service/internal/service/campaignsrv"
	"github.com/StephanHCB/go-campaign-service/web"
)

func main() {
	logging.Setup()
	configuration.Setup()
	logging.PostConfigSetup()

	database.Open()
	defer database.Close()
	database.MigrateIfEnabled()

	web.StartMetricsServer()

	server := web.Create(campaignsrv.Create(mailserviceclient.Create()))
	web.Serve(server)
}
