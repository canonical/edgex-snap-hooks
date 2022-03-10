package log

import (
	"os"

	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
)

var (
	debug           bool
	snapInstanceKey string // used as default syslog tag and tag prefix
	tag             string // syslog tag and stderr prefix
)

func init() {
	initialize()
}

func initialize() {
	value, err := snapctl.Get("debug").Run()
	if err != nil {
		stderr(err)
		os.Exit(1)
	}
	debug = (value == "true")

	snapInstanceKey = os.Getenv("SNAP_INSTANCE_NAME")
	if snapInstanceKey == "" {
		stderr("SNAP_INSTANCE_NAME environment variable not set.")
		os.Exit(1)
	}
	tag = snapInstanceKey

	if err := setupSyslogWriter(tag); err != nil {
		stderr(err)
		os.Exit(1)
	}
}
