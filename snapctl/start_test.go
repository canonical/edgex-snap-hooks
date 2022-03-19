package snapctl_test

import (
	"testing"

	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	// make sure services are stopped
	stopAndDisableAllServices(t)

	t.Run("snapctl start", func(t *testing.T) {
		t.Run("one", func(t *testing.T) {
			t.Cleanup(func() { stopAndDisableAllServices(t) })

			err := snapctl.Start(mockService).Run()
			require.NoError(t, err)
			_, active := getServiceStatus(t, mockService)
			require.True(t, active, "active")
		})

		t.Run("multiple", func(t *testing.T) {
			t.Cleanup(func() { stopAndDisableAllServices(t) })

			err := snapctl.Start(mockService, mockService2).Run()
			require.NoError(t, err)
			_, active := getServiceStatus(t, mockService)
			require.True(t, active, "active")
			_, active = getServiceStatus(t, mockService2)
			require.True(t, active, "active")
		})
	})

	t.Run("snapctl start --enable", func(t *testing.T) {
		t.Cleanup(func() { stopAndDisableAllServices(t) })

		err := snapctl.Start(mockService).Enable().Run()
		require.NoError(t, err)
		enabled, active := getServiceStatus(t, mockService)
		require.True(t, enabled, "enabled")
		require.True(t, active, "active")
	})

	t.Run("reject name with space", func(t *testing.T) {
		err := snapctl.Start("bad name").Run()
		require.Error(t, err)
	})
}
