# edgex-snap-hooks
[![Go Reference](https://pkg.go.dev/badge/github.com/canonical/edgex-snap-hooks.svg)](https://pkg.go.dev/github.com/canonical/edgex-snap-hooks/v2)

Snap hooks library used by Go EdgeX service snaps.  
It provides utilites to implement snap hooks, including some wrappers for `snapctl` commands.

### Installation ###
Make sure you have modules enabled, i.e. have an initialized `go.mod` file.

Run:
```
go get github.com/canonical/edgex-snap-hooks/v2
```
This will add the edgex-snap-hooks to the go.mod file and download it into the module cache.

### How to Use ###

TBA

### Testing (WIP)
The tests need to run in a snap environment:

```bash
snapcraft build
```
