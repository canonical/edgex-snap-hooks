package snapctl_test

import (
	"testing"

	"github.com/canonical/edgex-snap-hooks/v3/snapctl"
	"github.com/stretchr/testify/require"
)

func TestIsConnected(t *testing.T) {
	t.Run("snapctl is-connected", func(t *testing.T) {

		connected, err := snapctl.IsConnected("test-plug").Run()
		require.NoError(t, err, "Error checking plug status.")
		require.False(t, connected)

	})

}
