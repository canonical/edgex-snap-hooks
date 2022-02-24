package snapctl

import (
	"os"
)

const (
	snapNameEnv = "SNAP_NAME"
	debugEnv    = "DEBUG"
)

var (
	snapName string
	debug    bool
)

func init() {
	snapName = os.Getenv(snapNameEnv)
	if snapName == "" {
		panic("SNAP_NAME is not set")
	}

	if os.Getenv(debugEnv) == "true" {
		debug = true
	}

}
