package snapctl

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	key, value := "mykey", "myvalue"

	setConfigValue(t, key, value)

	t.Run("snapctl get", func(t *testing.T) {
		retrievedValue, err := Get().Keys(key).Run()
		require.NoError(t, err, "Error getting config.")
		require.Equal(t, value, retrievedValue)
	})

	t.Run("snapctl get -d", func(t *testing.T) {
		retrievedValue, err := Get().Keys(key).Document().Run()
		require.NoError(t, err, "Error getting config as document.")
		compact := new(bytes.Buffer)
		err = json.Compact(compact, []byte(retrievedValue))
		require.NoError(t, err, "Error parsing response as JSON.")
		require.Equal(t, `{"mykey":"myvalue"}`, compact.String())
	})
}
