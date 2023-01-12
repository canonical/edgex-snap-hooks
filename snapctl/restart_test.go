package snapctl_test

import (
	"testing"

	"github.com/canonical/edgex-snap-hooks/v3/snapctl"
	"github.com/stretchr/testify/require"
)

func TestRestart(t *testing.T) {
	// make sure services are stopped
	stopAndDisableAllServices(t)

	t.Run("snapctl restart", func(t *testing.T) {
		t.Run("one", func(t *testing.T) {
			t.Cleanup(func() { stopAndDisableAllServices(t) })

			err := snapctl.Restart(mockService).Run()
			require.NoError(t, err)
			_, active := getServiceStatus(t, mockService)
			require.True(t, active, "active")
		})

		t.Run("multiple", func(t *testing.T) {
			t.Cleanup(func() { stopAndDisableAllServices(t) })

			err := snapctl.Restart(mockService, mockService2).Run()
			require.NoError(t, err)
			_, active := getServiceStatus(t, mockService)
			require.True(t, active, "active")
			_, active = getServiceStatus(t, mockService2)
			require.True(t, active, "active")
		})
	})

	t.Run("snapctl restart --reload", func(t *testing.T) {
		t.Skip("TODO")
	})

	t.Run("reject name with space", func(t *testing.T) {
		err := snapctl.Restart("bad name").Run()
		require.Error(t, err)
	})
}
