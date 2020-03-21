package acceptance

import (
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database"
	"github.com/StephanHCB/go-campaign-service/internal/repository/logging"
	"github.com/StephanHCB/go-campaign-service/web"
	"net/http/httptest"
)

// placing these here because they are package global

var (
	ts *httptest.Server
	failures []error
	warnings []string
)

const tstValidConfigurationPath =  "resources/validconfig"
const tstInvalidConfigurationPath = "resources/invalidconfig"

func tstSetup(configAndSecretsPath string) {
	tstSetupConfig(configAndSecretsPath, configAndSecretsPath)
	if !tstHadFailures() {
		tstSetupDatabase()
		tstSetupHttpTestServer()
	}
}

func tstFail(err error) {
	failures = append(failures, err)
}

func tstWarn(msg string) {
	warnings = append(warnings, msg)
}

func tstSetupConfig(configPath string, secretsPath string) {
	failures = []error{}
	warnings = []string{}
	logging.SetupForTesting()
	configuration.SetupForIntegrationTest(tstFail, tstWarn, configPath, secretsPath)
	if !tstHadFailures() {
		logging.PostConfigSetup()
	}
}

func tstHadFailures() bool {
	return len(failures) > 0
}

func tstSetupHttpTestServer() {
	router := web.Create()
	ts = httptest.NewServer(router)
}

func tstSetupDatabase() {
	database.Open()
	database.MigrateIfEnabled()
}

func tstShutdown() {
	if !tstHadFailures() {
		ts.Close()
		database.Close()
	}
}
