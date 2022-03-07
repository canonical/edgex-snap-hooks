package snapctl

import (
	"reflect"
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
			serviceName := reflect.ValueOf(services).MapKeys()[0].String()
			require.Equal(t, mockService, serviceName)
		})

		t.Run("all", func(t *testing.T) {
			services, err := Services().Run()
			require.NoError(t, err, "Error getting services.")
			require.Len(t, services, 2)
			for k := range services {
				require.Contains(t, []string{mockService, mockServiceDisabled}, k)
			}
		})

		t.Run("enabled and active", func(t *testing.T) {
			services, err := Services(mockService).Run()
			require.NoError(t, err, "Error getting services.")
			require.Len(t, services, 1)
			for k, v := range services {
				require.Equal(t, mockService, k)
				require.True(t, v.Enabled, "Service not enabled")
				require.True(t, v.Active, "Service not active")
			}
		})

		t.Run("disabled and inactive", func(t *testing.T) {
			services, err := Services(mockServiceDisabled).Run()
			require.NoError(t, err, "Error getting services.")
			require.Len(t, services, 1)
			for k, v := range services {
				require.Equal(t, mockServiceDisabled, k)
				require.False(t, v.Enabled, "Service not disabled")
				require.False(t, v.Active, "Service not inactive")
			}
		})

		t.Run("service not found", func(t *testing.T) {
			_, err := Services("non-existed").Run()
			require.Error(t, err)
		})

		t.Run("service name with space", func(t *testing.T) {
			_, err := Services("bad name").Run()
			require.Error(t, err)
		})

	})
}
