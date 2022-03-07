package snapctl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnset(t *testing.T) {
	testKey, testValue := "test-key", "test-value"

	// set via snapctl as starting point
	setConfigValue(t, testKey, testValue)

	t.Run("snapctl unset", func(t *testing.T) {
		t.Run("one", func(t *testing.T) {
			// set via snapctl as starting point
			setConfigValue(t, testKey, testValue)
			require.Equal(t, `"test-value"`, getConfigStrictValue(t, testKey))

			err := Unset(testKey).Run()
			require.NoError(t, err)
			require.Equal(t, `null`, getConfigStrictValue(t, testKey))
		})

		t.Run("multiple", func(t *testing.T) {
			testKey2 := "test-key2"
			// set via snapctl as starting point
			setConfigValue(t, testKey, testValue)
			require.Equal(t, `"test-value"`, getConfigStrictValue(t, testKey))
			setConfigValue(t, testKey2, testValue)
			require.Equal(t, `"test-value"`, getConfigStrictValue(t, testKey2))

			err := Unset(testKey, testKey2).Run()
			require.NoError(t, err)
			require.Equal(t, `null`, getConfigStrictValue(t, testKey))
			require.Equal(t, `null`, getConfigStrictValue(t, testKey2))
		})

		t.Run("reject key with space", func(t *testing.T) {
			err := Unset("bad key").Run()
			require.Error(t, err)
		})
	})
}
