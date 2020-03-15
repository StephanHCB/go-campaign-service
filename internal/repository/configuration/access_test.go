package configuration

import (
	"github.com/StephanHCB/go-autumn-config"
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

func TestContains_yes_first(t *testing.T) {
	actual := contains([]string{"development", "squirrel", "local"}, "development")
	require.True(t, actual)
}

func TestContains_yes_last(t *testing.T) {
	actual := contains([]string{"development", "squirrel", "local"}, "local")
	require.True(t, actual)
}

func TestContains_yes_only(t *testing.T) {
	actual := contains([]string{"local"}, "local")
	require.True(t, actual)
}

func TestContains_no_empty(t *testing.T) {
	actual := contains([]string{}, "local")
	require.False(t, actual)
}

func TestContains_no(t *testing.T) {
	actual := contains([]string{"development", "local"}, "cat")
	require.False(t, actual)
}
