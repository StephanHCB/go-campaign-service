package configuration

import (
	"fmt"
	"github.com/spf13/viper"
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

// helper functions

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}