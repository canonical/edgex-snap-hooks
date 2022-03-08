#!/bin/sh

set -e

echo "Running tests from $0"

# setup the environment
export PATH=$PATH:$SNAP/go/bin
export CGO_ENABLED=0
echo "Change directory to $SNAP"
cd $SNAP

go test --cover "$@"
echo "✅ go test"

go vet "$@"
echo "✅ go vet"