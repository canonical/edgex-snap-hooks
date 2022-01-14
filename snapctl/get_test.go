package snapctl

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	testKey, testValue := "test-key", "test-value"

	// set via snapctl as starting point
	setConfigValue(t, testKey, testValue)

	t.Run("snapctl get", func(t *testing.T) {
		t.Run("one", func(t *testing.T) {
			retrievedValue, err := Get(testKey).Run()
			require.NoError(t, err, "Error getting config.")
			require.Equal(t, testValue, retrievedValue)
		})

		t.Run("multiple", func(t *testing.T) {
			t.Skip("TODO")
		})
	})

	t.Run("snapctl get -d", func(t *testing.T) {
		retrievedValue, err := Get(testKey).Document().Run()
		require.NoError(t, err, "Error getting config as document.")
		compact := new(bytes.Buffer)
		err = json.Compact(compact, []byte(retrievedValue))
		require.NoError(t, err, "Error parsing response as JSON.")
		require.Equal(t, `{"test-key":"test-value"}`, compact.String())
	})

	t.Run("snapctl get -t", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			retrievedValue, err := Get(testKey).Strict().Run()
			require.NoError(t, err, "Error getting config.")
			require.Equal(t, `"test-value"`, retrievedValue)
		})

		t.Run("null", func(t *testing.T) {
			retrievedValue, err := Get("some-other-key").Strict().Run()
			require.NoError(t, err, "Error getting config.")
			require.Equal(t, "null", retrievedValue)
		})
	})

	t.Run("snapctl get :interface", func(t *testing.T) {
		t.Skip("TODO: test interface hooks")
		// interface attributes can only be read during the execution of interface hooks
	})

	t.Run("snapctl get :interface --slot", func(t *testing.T) {
		t.Skip("TODO: test interface hooks")
		// interface attributes can only be read during the execution of interface hooks
		// cannot use --plug or --slot without <snap>:<plug|slot> argument
	})

	t.Run("snapctl get :interface --plug", func(t *testing.T) {
		t.Skip("TODO: test interface hooks")
		// interface attributes can only be read during the execution of interface hooks
		// cannot use --plug or --slot without <snap>:<plug|slot> argument
	})
}

// TestGetValidate tests the validation done via this library
func TestGetValidate(t *testing.T) {
	t.Run("interface colon prefix", func(t *testing.T) {
		_, err := Get().Interface(":test-plug").Run()
		require.Error(t, err, "interface has colon as prefix")
	})
}
