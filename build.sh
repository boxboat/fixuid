#!/bin/sh
cd $(dirname $0)

rm -f ./fixuid
CGO_ENABLED=0 go build
