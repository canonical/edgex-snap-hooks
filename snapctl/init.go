package snapctl

import (
	"os"
)

const snapNameEnv = "SNAP_NAME"

var snapName string

func init() {
	snapName = os.Getenv(snapNameEnv)
	if snapName == "" {
		panic("SNAP_NAME is not set")
	}
}
