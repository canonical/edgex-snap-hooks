#!/bin/sh

set -e

echo "Running tests from $0"

echo "Change directory to $SNAP"
cd $SNAP
go test -v  --cover
go test -v ./snapctl --cover

# go vet ./...