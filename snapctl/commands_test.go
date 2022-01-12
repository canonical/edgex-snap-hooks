package snapctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	key, value := "mykey", "myvalue"

	setConfigValue(t, key, value)

	t.Run("snapctl get", func(t *testing.T) {
		retrievedValue, err := Get(key).Run()
		require.NoError(t, err, "Error getting config.")
		require.Equal(t, value, retrievedValue)
	})

	t.Run("snapctl get -d", func(t *testing.T) {
		retrievedValue, err := Get(key).Document().Run()
		require.NoError(t, err, "Error getting config as document.")
		compact := new(bytes.Buffer)
		err = json.Compact(compact, []byte(retrievedValue))
		require.NoError(t, err, "Error parsing response as JSON.")
		require.Equal(t, `{"mykey":"myvalue"}`, compact.String())
	})
}

//  func TestSetConfig(t *testing.T) {
// 	 key, value := "mykey", "myvalue"
 
// 	 cli := NewSnapCtl()
// 	 err := cli.SetConfig(key, value)
// 	 require.NoError(t, err, "Error setting config.")
 
// 	 // check using snapctl
// 	 require.Equal(t, value, getConfigValue(t, key))
//  }
 
//  func TestUnsetConfig(t *testing.T) {
// 	 key, value := "mykey2", "myvalue"
 
// 	 // make sure this isn't already set
// 	 require.Equal(t, "", getConfigValue(t, key))
 
// 	 // set using snapctl
// 	 setConfigValue(t, key, value)
 
// 	 // check using snapctl
// 	 require.Equal(t, value, getConfigValue(t, key))
 
// 	 // set using the library
// 	 cli := NewSnapCtl()
// 	 err := cli.UnsetConfig(key)
// 	 require.NoError(t, err, "Error un-setting config.")
 
// 	 // make sure it has been unset
// 	 require.Equal(t, "", getConfigValue(t, key))
//  }
 
//  func TestStartMultiple(t *testing.T) {
// 	 t.Skipf("need to run in an active snap")
// 	 cli := NewSnapCtl()
// 	 err := cli.StartMultiple(false, "edgexfoundry.consul")
// 	 require.NoError(t, err, "Error getting services.")
//  }
 
//  func TestEnabledServices(t *testing.T) {
// 	 t.Skipf("need to run in an active snap")
// 	 cli := NewSnapCtl()
// 	 services, err := cli.EnabledServices()
// 	 require.NoError(t, err, "Error getting enabled services.")
// 	 t.Logf("services: %v", services)
//  }
 
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
 