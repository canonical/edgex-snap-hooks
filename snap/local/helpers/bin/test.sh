#!/bin/sh

set -e

echo "Running tests from $0"

echo "Change directory to $SNAP"
cd $SNAP
go test -v $1 --cover

# go vet ./...