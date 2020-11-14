package logging

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
)

func Setup() {
	// assume JSON logging at first, until configuration is loaded, and we do not know the service id yet
	auzerolog.SetupJsonLogging("")
}

func PostConfigSetup() {
	if configuration.IsProfileActive("local") {
		auzerolog.SetupPlaintextLogging()
		aulogging.Logger.NoCtx().Info().Print("switching to developer friendly console log because profile 'local' is active")
	} else {
		// stay with JSON logging and add ECS service.id field
		auzerolog.SetupJsonLogging(configuration.ServiceName())
	}
}
