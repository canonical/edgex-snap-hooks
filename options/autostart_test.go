package options_test

import (
	"testing"

	"github.com/canonical/edgex-snap-hooks/v2/options"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
	"github.com/stretchr/testify/require"
)

const (
	mockApp      = "mock-service"
	mockApp2     = "mock-service-2"
	mockService  = "edgex-snap-hooks." + mockApp
	mockService2 = "edgex-snap-hooks." + mockApp2
)

func TestProcessAutostartGlobal(t *testing.T) {
	require.NoError(t, snapctl.Stop(mockService, mockService2).Disable().Run())
	t.Cleanup(func() {
		require.NoError(t, snapctl.Stop(mockService, mockService2).Disable().Run())
		require.NoError(t, snapctl.Unset("autostart").Run())
	})

	t.Run("true", func(t *testing.T) {
		require.NoError(t, snapctl.Set("autostart", "true").Run())
		require.NoError(t, options.ProcessAutostart(mockApp, mockApp2))

		services, err := snapctl.Services(mockService, mockService2).Run()
		require.NoError(t, err)
		for name, status := range services {
			require.True(t, status.Active, name+" active")
			require.True(t, status.Enabled, name+" enabled")
		}
	})

	t.Run("unset", func(t *testing.T) { // should have no effect
		require.NoError(t, snapctl.Unset("autostart").Run())
		require.NoError(t, options.ProcessAutostart(mockApp, mockApp2))

		services, err := snapctl.Services(mockService, mockService2).Run()
		require.NoError(t, err)
		for name, status := range services {
			require.True(t, status.Active, name+" active")
			require.True(t, status.Enabled, name+" enabled")
		}
	})

	t.Run("false", func(t *testing.T) {
		require.NoError(t, snapctl.Set("autostart", "false").Run())
		require.NoError(t, options.ProcessAutostart(mockApp, mockApp2))

		services, err := snapctl.Services(mockService, mockService2).Run()
		t.Log(services)
		require.NoError(t, err)
		for name, status := range services {
			require.False(t, status.Active, name+" active")
			require.False(t, status.Enabled, name+" enabled")
		}
	})
}

func TestProcessAutostartApp(t *testing.T) {
	require.NoError(t, snapctl.Stop(mockService, mockService2).Disable().Run())
	t.Cleanup(func() {
		require.NoError(t, snapctl.Stop(mockService, mockService2).Disable().Run())
		require.NoError(t, snapctl.Unset("autostart").Run())
		require.NoError(t, snapctl.Unset("apps").Run())
	})

	t.Run("true", func(t *testing.T) {
		require.NoError(t, snapctl.Set("apps."+mockApp+".autostart", "true").Run())
		require.NoError(t, options.ProcessAutostart(mockApp, mockApp2))

		services, err := snapctl.Services(mockService, mockService2).Run()
		require.NoError(t, err)
		require.True(t, services[mockService].Active, mockApp+" active")
		require.True(t, services[mockService].Enabled, mockApp+" enabled")
		require.False(t, services[mockService2].Active, mockApp2+" active")
		require.False(t, services[mockService2].Enabled, mockApp2+" enabled")
	})

	t.Run("unset", func(t *testing.T) { // should have no effect
		require.NoError(t, snapctl.Unset("apps."+mockApp+".autostart").Run())
		require.NoError(t, options.ProcessAutostart(mockApp, mockApp2))

		services, err := snapctl.Services(mockService, mockService2).Run()
		require.NoError(t, err)
		require.True(t, services[mockService].Active, mockApp+" active")
		require.True(t, services[mockService].Enabled, mockApp+" enabled")
		require.False(t, services[mockService2].Active, mockApp2+" active")
		require.False(t, services[mockService2].Enabled, mockApp2+" enabled")
	})

	t.Run("false", func(t *testing.T) {
		require.NoError(t, snapctl.Set("apps."+mockApp+".autostart", "false").Run())
		require.NoError(t, options.ProcessAutostart(mockApp, mockApp2))

		services, err := snapctl.Services(mockService, mockService2).Run()
		require.NoError(t, err)
		require.False(t, services[mockService].Active, mockApp+" active")
		require.False(t, services[mockService].Enabled, mockApp+" enabled")
		require.False(t, services[mockService2].Active, mockApp2+" active")
		require.False(t, services[mockService2].Enabled, mockApp2+" enabled")
	})
}

func TestProcessAutostartAppOverride(t *testing.T) {
	require.NoError(t, snapctl.Stop(mockService, mockService2).Disable().Run())
	t.Cleanup(func() {
		require.NoError(t, snapctl.Stop(mockService, mockService2).Disable().Run())
		require.NoError(t, snapctl.Unset("autostart").Run())
		require.NoError(t, snapctl.Unset("apps").Run())
	})

	// set globally and override it for one app
	require.NoError(t, snapctl.Set("autostart", "true").Run())
	require.NoError(t, snapctl.Set("apps."+mockApp2+".autostart", "false").Run())
	require.NoError(t, options.ProcessAutostart(mockApp, mockApp2))

	services, err := snapctl.Services(mockService, mockService2).Run()
	require.NoError(t, err)
	require.True(t, services[mockService].Active, mockApp+" active")
	require.True(t, services[mockService].Enabled, mockApp+" enabled")
	require.False(t, services[mockService2].Active, mockApp2+" active")
	require.False(t, services[mockService2].Enabled, mockApp2+" enabled")
}
