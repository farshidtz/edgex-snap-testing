package test

import (
	"edgex-snap-testing/test/utils"
	"testing"
)

// Deprecated
func TestEnvConfig(t *testing.T) {
	utils.SetEnvConfig(t, deviceVirtualSnap, deviceVirtualApp, defaultServicePort)
}

func TestAppConfig(t *testing.T) {
	utils.SetAppConfig(t, deviceVirtualSnap, deviceVirtualApp, defaultServicePort)
}

func TestGlobalConfig(t *testing.T) {
	utils.SetGlobalConfig(t, deviceVirtualSnap, deviceVirtualApp, defaultServicePort)
}

func TestMixedConfig(t *testing.T) {
	utils.SetMixedConfig(t, deviceVirtualSnap, deviceVirtualApp, defaultServicePort)
}