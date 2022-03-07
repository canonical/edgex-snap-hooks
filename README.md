# edgex-snap-hooks
[![Go Reference](https://pkg.go.dev/badge/github.com/canonical/edgex-snap-hooks.svg)](https://pkg.go.dev/github.com/canonical/edgex-snap-hooks/v2)

Snap hooks library used by [EdgeX Foundry](https://docs.edgexfoundry.org/) Go service snaps.  
It provides utilites to implement snap hooks, including some wrappers for the [`snapctl`](https://snapcraft.io/docs/using-snapctl) commands.

### Usage
Download or upgrade to the latest version:
```
go get github.com/canonical/edgex-snap-hooks/v2
```
Please refer to [go get docs](https://pkg.go.dev/cmd/go#hdr-Add_dependencies_to_current_module_and_install_them) for details.

#### Example

```go
package main

import (
	"fmt"
	"os"

	hooks "github.com/canonical/edgex-snap-hooks/v2"
)

func main() {
	var err error

	if err = hooks.Init(false, "edgex-device-example"); err != nil {
		fmt.Printf("initialization failure: %s", err)
		os.Exit(1)
	}

	// copy file from $SNAP to $SNAP_DATA
	if err = hooks.CopyFile(hooks.Snap+"/config.json", hooks.SnapData+"config.json"); err != nil {
		hooks.Error(err.Error())
		os.Exit(1)
	}
  
	// read env var override configuration
	cli := hooks.NewSnapCtl()
	envJSON, err := cli.Config(hooks.EnvConfig)
	if err != nil {
		hooks.Error(fmt.Sprintf("Reading config 'env' failed: %v", err))
		os.Exit(1)
	}
	hooks.Debug(fmt.Sprintf("envJSON: %s", envJSON))
}

```

### Testing
The tests need to run in a snap environment:

Build and install:
```bash
snapcraft
sudo snap install --dangerous ./edgex-snap-hooks_test_amd64.snap
```

The tests files are read relative to project source inside the snap.
The `edgex-snap-hooks.test` command runs `go test -v --cover` internally and accepts
all other go test arguments.

Run top-level tests:
```bash
sudo edgex-snap-hooks.test
```

Run tests in one package, e.g. `snapctl`:
```bash
sudo edgex-snap-hooks.test ./snapctl
```

Run one unit test, e.g. `TestGet`:
```bash
sudo edgex-snap-hooks.test ./snapctl -run TestGet
```

#### Development
```
snapcraft try
snap try prime
sudo edgex-snap-hooks.test ./snapctl
```

You can now edit the files locally, copy them to prime directory, and re-run the
tests without rebuilding the project. E.g.:

```
cp -r snapctl prime/ && \
sudo edgex-snap-hooks.test ./snapctl
```