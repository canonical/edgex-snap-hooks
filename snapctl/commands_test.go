package snapctl

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// utility testing functions
 
 func setConfigValue(t *testing.T, key, value string) {
	 err := exec.Command("snapctl", "set", fmt.Sprintf("%s=%s", key, value)).Run()
	 require.NoError(t, err, "Error setting config value via snapctl.")
 }
 
 func getConfigValue(t *testing.T, key string) string {
	 out, err := exec.Command("snapctl", "get", key).Output()
	 require.NoError(t, err, "Error getting config value via snapctl.")
	 return strings.TrimSpace(string(out))
 }
 