package configuration

import (
	"github.com/StephanHCB/go-autumn-config"
	"github.com/StephanHCB/go-autumn-config-api"
)

const configKeyServerAddress = "server.address"
const configKeyServerPort = "server.port"
const configKeyServiceName = "service.name"
const configKeyDatabaseUse = "database.use"
const configKeyDatabaseMigrate = "database.migrate"
const configKeyDatabaseMysqlUsername = "database.mysql.username"
const configKeyDatabaseMysqlPassword = "database.mysql.password"
const configKeyDatabaseMysqlDatabase = "database.mysql.database"
const configKeyDatabaseMysqlParameters = "database.mysql.parameters"
const configKeySecuritySecret = "security.secret"
const configKeyDownstreamMailerserviceUrl = "downstream.mailerservice.url"
const configKeyDownstreamMailerserviceTimeoutMs = "downstream.mailerservice.timeoutMs"

var configAllowedDatabases = []string{"mysql", "inmemory"}

var configItems = []auconfigapi.ConfigItem{
	auconfig.ConfigItemProfile,
	{
		Key:         configKeyServerAddress,
		Default:     "",
		Description: "ip address or hostname to listen on, can be left blank for localhost",
		Validate:    func(key string) error { return checkLength(0, 255, key) },
	}, {
		Key:         configKeyServerPort,
		Default:     uint(8080),
		Description: "port to listen on, defaults to 8080 if not set",
		Validate:    checkValidPortNumber,
	}, {
		Key:         configKeyServiceName,
		Default:     "unnamed-service",
		Description: "name of service, used for logging",
		Validate:    func(key string) error { return checkLength(1, 255, key) },
	},
	// database configuration
	{
		Key:         configKeyDatabaseUse,
		Default:     "inmemory",
		Description: "choice of database, either 'mysql' or 'inmemory'",
		Validate:    func(key string) error { return checkEnum(key, configAllowedDatabases) },
	}, {
		Key:         configKeyDatabaseMigrate,
		Default:     false,
		Description: "automatically migrate the database, boolean, defaults to false",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:         configKeyDatabaseMysqlUsername,
		Default:     "",
		Description: "username to connect to database",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:         configKeyDatabaseMysqlPassword,
		Default:     "",
		Description: "password to connect to database",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:         configKeyDatabaseMysqlDatabase,
		Default:     "tcp(localhost:3306)/dbname",
		Description: "database connect string",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:         configKeyDatabaseMysqlParameters,
		Default:     []string{},
		Description: "list of database connection parameters",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	// security configuration
	{
		Key:         configKeySecuritySecret,
		Default:     "",
		Description: "secret used for signing jwt tokens",
		Validate:    func(key string) error { return checkLength(1, 255, key) },
	},
	// downstream services configuration
	{
		Key:         configKeyDownstreamMailerserviceUrl,
		Default:     "",
		Description: "url to access the mailer-service",
		Validate:    func(key string) error { return checkLength(1, 255, key) },
	}, {
		Key:         configKeyDownstreamMailerserviceTimeoutMs,
		Default:     uint(2000),
		Description: "timeout in milliseconds until calls are aborted",
		Validate:    func(key string) error { return checkUintRange(100, 10000, key) },
	},
}
