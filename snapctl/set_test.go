package snapctl_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	testKey, testValue := "test-key", "test-value"
	testBoolValue := "true"
	testJSONValue := `{"nested-key":"nested-value"}`

	t.Run("snapctl set", func(t *testing.T) {
		t.Run("one", func(t *testing.T) {
			err := snapctl.Set(testKey, testValue).Run()
			require.NoError(t, err)
			require.Equal(t, `"test-value"`, getConfigStrictValue(t, testKey))
		})

		t.Run("bool", func(t *testing.T) {
			err := snapctl.Set(testKey, testBoolValue).Run()
			require.NoError(t, err)
			// should not have double quotations
			require.Equal(t, `true`, getConfigStrictValue(t, testKey))
		})

		testKey2, testValue2 := "test-key2", "test-value2"
		t.Run("multiple", func(t *testing.T) {
			err := snapctl.Set(
				testKey, testValue,
				testKey2, testValue2,
			).Run()
			require.NoError(t, err)

			assert.Equal(t, `"test-value"`, getConfigStrictValue(t, testKey))
			assert.Equal(t, `"test-value2"`, getConfigStrictValue(t, testKey2))
		})

		t.Run("reject bad pair", func(t *testing.T) {
			err := snapctl.Set(testKey, testValue, testKey2).Run()
			require.Error(t, err)
		})

		t.Run("reject key with space", func(t *testing.T) {
			_, err := snapctl.Get("bad key").Run()
			require.Error(t, err)
		})
	})

	t.Run("snapctl set -t", func(t *testing.T) {
		err := snapctl.Set(testKey, testJSONValue).Document().Run()
		require.NoError(t, err)

		// read the compacted set value
		compact := new(bytes.Buffer)
		err = json.Compact(compact, []byte(getConfigStrictValue(t, testKey)))
		require.NoError(t, err, "Error parsing response as JSON.")

		// should NOT be escaped
		require.Equal(t, `{"nested-key":"nested-value"}`, compact.String())
	})

	t.Run("snapctl set -s", func(t *testing.T) {
		t.Run("bool as string", func(t *testing.T) {
			err := snapctl.Set(testKey, testBoolValue).String().Run()
			require.NoError(t, err)

			require.Equal(t, `"true"`, getConfigStrictValue(t, testKey))
		})

		t.Run("json as string", func(t *testing.T) {
			err := snapctl.Set(testKey, testJSONValue).String().Run()
			require.NoError(t, err)

			// should be escaped JSON
			require.Equal(t, `"{\"nested-key\":\"nested-value\"}"`, getConfigStrictValue(t, testKey))
		})

	})

	t.Run("snapctl set :interface", func(t *testing.T) {
		// interface attributes can only be read during the execution of interface hooks
		t.Skip("TODO: test interface hooks")

		// t.Run("reject colon prefix", func(t *testing.T) {
		// 	err := Set(testKey, testValue).Interface(":test-plug").Run()
		// 	require.Error(t, err)
		// })
	})
}
