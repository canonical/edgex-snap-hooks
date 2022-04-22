package log

import (
	"bytes"
	"os"
	"os/exec"
)

var (
	debug           bool
	snapInstanceKey string // used as default syslog tag and tag prefix
	tag             string // syslog tag and stderr prefix
)

func init() {
	Init()
}

func Init() {
	value, err := exec.Command("snapctl", "get", "debug").CombinedOutput()
	if err != nil {
		stderr(err)
		os.Exit(1)
	}
	debug = (string(bytes.TrimSpace(value)) == "true")

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
