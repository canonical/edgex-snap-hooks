# edgex-snap-hooks
Snap hooks library use by Go EdgeX service snaps.  

### What is this repository for? ###

This module provides support for Go-based configure an install hooks...

### Installation ###

* Make sure you have modules enabled, i.e. have an initialized  go.mod file
* If your code is in your GOPATH then make sure ```GO111MODULE=on``` is set
* Run ```go get github.com/canonical/edgex-snap-hooks```
  * This will add the edgex-snap-hooks to the go.mod file and download it into the module cache

### How to Use ###

TBA

### Testing
The tests need to run in a snap environment:

```bash
snapcraft build
```