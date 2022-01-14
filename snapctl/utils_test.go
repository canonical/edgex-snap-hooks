// Utility testing functions

package snapctl

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setConfigValue(t *testing.T, key, value string) {
	output, err := exec.Command("snapctl", "set", fmt.Sprintf("%s=%s", key, value)).CombinedOutput()
	assert.NoError(t, err,
		"Error setting config value via snapctl: %s", output)
}

func getConfigValue(t *testing.T, key string) string {
	output, err := exec.Command("snapctl", "get", key).CombinedOutput()
	require.NoError(t, err,
		"Error getting config value via snapctl: %s", output)
	return strings.TrimSpace(string(output))
}
