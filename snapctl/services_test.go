package snapctl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServices(t *testing.T) {
	mockService := snapName + ".mock-service"
	mockServiceDisabled := snapName + ".mock-service-disabled"

	t.Run("snapctl services", func(t *testing.T) {
		t.Run("one", func(t *testing.T) {
			services, err := Services(mockService).Run()
			require.NoError(t, err, "Error getting services.")
			require.Len(t, services, 1)
			for k, v := range services {
				require.Equal(t, mockService, k)
				require.True(t, v.Enabled, "Service not enabled")
				require.True(t, v.Active, "Service not active")
			}
		})

		t.Run("all", func(t *testing.T) {
			services, err := Services().Run()
			require.NoError(t, err, "Error getting services.")
			require.Len(t, services, 2)
			for k := range services {
				require.Contains(t, []string{mockService, mockServiceDisabled}, k)
			}
		})

		t.Run("service name invalid", func(t *testing.T) {
			_, err := Services("non-existed").Run()
			require.Error(t, err)
		})

		t.Run("service name with space", func(t *testing.T) {
			_, err := Services("bad name").Run()
			require.Error(t, err)
		})

	})
}
