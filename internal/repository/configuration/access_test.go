package configuration

import (
	"github.com/StephanHCB/go-autumn-config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServerAddress(t *testing.T) {
	auconfig.SetupDefaultsOnly(configItems, failFunction, warnFunction)

	expected := ":8080"
	actual := ServerAddress()
	require.Equal(t, expected, actual)
}

func TestServiceName(t *testing.T) {
	auconfig.SetupDefaultsOnly(configItems, failFunction, warnFunction)

	expected := "unnamed-service"
	actual := ServiceName()
	require.Equal(t, expected, actual)
}

func TestIsProfileActive(t *testing.T) {
	auconfig.SetupDefaultsOnly(configItems, failFunction, warnFunction)

	actual := IsProfileActive("production")
	require.False(t, actual)
}

func TestDatabaseMysqlConnectString_NoParams(t *testing.T) {
	auconfig.SetupDefaultsOnly(configItems, failFunction, warnFunction)
	viper.Set(configKeyDatabaseMysqlUsername, "usr")
	viper.Set(configKeyDatabaseMysqlPassword, "pwd")

	expected := "usr:pwd@tcp(localhost:3306)/dbname"
	actual := DatabaseMysqlConnectString()
	require.Equal(t, expected, actual)
}

func TestDatabaseMysqlConnectString_Params(t *testing.T) {
	auconfig.SetupDefaultsOnly(configItems, failFunction, warnFunction)
	viper.Set(configKeyDatabaseMysqlUsername, "usr")
	viper.Set(configKeyDatabaseMysqlPassword, "pwd")
	viper.Set(configKeyDatabaseMysqlParameters, []string{ "charset=utf8mb4", "cat=meow"})

	expected := "usr:pwd@tcp(localhost:3306)/dbname?charset=utf8mb4&cat=meow"
	actual := DatabaseMysqlConnectString()
	require.Equal(t, expected, actual)
}

func TestDatabaseUse(t *testing.T) {
	auconfig.SetupDefaultsOnly(configItems, failFunction, warnFunction)

	actual := DatabaseUse()
	expected := "inmemory"
	require.Equal(t, expected, actual)
}
