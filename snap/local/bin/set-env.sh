#!/bin/sh

set -e

echo "Setting environment from $0"

export PATH=$PATH:$SNAP/go/bin
export CGO_ENABLED=0

exec "$@"