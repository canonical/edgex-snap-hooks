package snapctl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartStop(t *testing.T) {
	// mockService := snapName + ".mock-service"
	mockServiceDisabled := snapName + ".mock-service-disabled"

	t.Run("snapctl start", func(t *testing.T) {
		err := Start(mockServiceDisabled).Run()
		require.NoError(t, err)
		_, active := getServiceStatus(t, mockServiceDisabled)
		require.True(t, active, "active")
	})

	t.Run("snapctl stop", func(t *testing.T) {
		err := Stop(mockServiceDisabled).Run()
		require.NoError(t, err)
		_, active := getServiceStatus(t, mockServiceDisabled)
		require.False(t, active, "active")
	})

	t.Run("snapctl start --enable", func(t *testing.T) {
		err := Start(mockServiceDisabled).Enable().Run()
		require.NoError(t, err)
		enabled, active := getServiceStatus(t, mockServiceDisabled)
		require.True(t, enabled, "enabled")
		require.True(t, active, "active")
	})

	t.Run("snapctl stop --disable", func(t *testing.T) {
		err := Stop(mockServiceDisabled).Disable().Run()
		require.NoError(t, err)
		enabled, active := getServiceStatus(t, mockServiceDisabled)
		require.False(t, enabled, "enabled")
		require.False(t, active, "active")
	})

	t.Run("reject name with space", func(t *testing.T) {
		err := Start("bad name").Run()
		require.Error(t, err)
	})

}
