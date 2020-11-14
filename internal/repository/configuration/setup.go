package configuration

import (
	"github.com/StephanHCB/go-autumn-config"
	"github.com/StephanHCB/go-autumn-config-api"
	aulogging "github.com/StephanHCB/go-autumn-logging"
)

// initialize configuration with full setup - you need to call this
func Setup() {
	auconfig.Setup(configItems, fail, warn)
	auconfig.Load()
}

// use this in unit tests
func SetupForUnitTestDefaultsOnlyNoErrors() {
	auconfig.SetupDefaultsOnly(configItems, func(err error) {}, func(message string) {})
}

// use this in integration tests
func SetupForIntegrationTest(failFunc auconfigapi.ConfigFailFunc, warnFunc auconfigapi.ConfigWarnFunc, configPath string, secretsPath string) {
	auconfig.ResetForTesting()
	auconfig.SetupWithOverriddenConfigPath(configItems, failFunc, warnFunc, configPath, secretsPath)
	auconfig.Load()
}

func fail(err error) {
	// this will os.exit 1
	aulogging.Logger.NoCtx().Fatal().WithErr(err).Print(err.Error())
}

func warn(message string) {
	aulogging.Logger.NoCtx().Warn().Print(message)
}
