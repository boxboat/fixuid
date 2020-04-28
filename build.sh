#!/bin/sh
cd $(dirname $0)

rm -f ./fixuid
GOOS=linux CGO_ENABLED=0 go build
