#!/bin/sh -e
cd "$(dirname "$0")"

rm -f ./test-no-escalate
CGO_ENABLED=0 go build
