package main

import (
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/web"
)

func main() {
	configuration.Setup()

	web.Serve()
}
