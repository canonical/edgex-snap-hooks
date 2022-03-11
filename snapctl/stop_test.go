package snapctl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStop(t *testing.T) {
	// make sure services are started
	stopAndEnableAllServices(t)

	t.Run("snapctl stop", func(t *testing.T) {
		t.Run("one", func(t *testing.T) {
			t.Cleanup(func() { stopAndEnableAllServices(t) })

			err := Stop(mockService).Run()
			require.NoError(t, err)
			_, active := getServiceStatus(t, mockService)
			require.False(t, active, "active")
		})

		t.Run("multiple", func(t *testing.T) {
			t.Cleanup(func() { stopAndEnableAllServices(t) })

			err := Stop(mockService, mockService2).Run()
			require.NoError(t, err)
			_, active := getServiceStatus(t, mockService)
			require.False(t, active, "active")
			_, active = getServiceStatus(t, mockService2)
			require.False(t, active, "active")
		})
	})

	t.Run("snapctl stop --disable", func(t *testing.T) {
		t.Cleanup(func() { stopAndEnableAllServices(t) })

		err := Stop(mockService).Disable().Run()
		require.NoError(t, err)
		enabled, active := getServiceStatus(t, mockService)
		require.False(t, enabled, "enabled")
		require.False(t, active, "active")
	})

	t.Run("reject name with space", func(t *testing.T) {
		err := Start("bad name").Run()
		require.Error(t, err)
	})
}
