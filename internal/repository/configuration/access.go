package configuration

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// public functions for accessing all configuration values.

func ServerAddress() string {
	return fmt.Sprintf("%v:%d", viper.GetString(configKeyServerAddress), viper.GetUint(configKeyServerPort))
}

func ServiceName() string {
	return viper.GetString(configKeyServiceName)
}

func IsProfileActive(profileName string) bool {
	profiles := viper.GetStringSlice("profiles")
	return contains(profiles, profileName)
}

func DatabaseUse() string {
	return viper.GetString(configKeyDatabaseUse)
}

func DatabaseMysqlConnectString() string {
	username := viper.GetString(configKeyDatabaseMysqlUsername)
	password := viper.GetString(configKeyDatabaseMysqlPassword)
	database := viper.GetString(configKeyDatabaseMysqlDatabase)
	parameters := viper.GetStringSlice(configKeyDatabaseMysqlParameters)

	result := username + ":" + password + "@" + database
	if len(parameters) > 0 {
		result += "?" + strings.Join(parameters, "&")
	}
	return result
}

func MigrateDatabase() bool {
	return viper.GetBool(configKeyDatabaseMigrate)
}

func SecuritySecret() string {
	return viper.GetString(configKeySecuritySecret)
}
