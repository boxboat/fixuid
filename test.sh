#!/bin/sh
cd $(dirname $0)
set -e

./build.sh
mv fixuid docker/fs-stage/usr/local/bin
rm -rf docker/alpine/stage
cp -r docker/fs-stage docker/alpine/stage
rm -rf docker/centos/stage
cp -r docker/fs-stage docker/centos/stage
rm -rf docker/debian/stage
cp -r docker/fs-stage docker/debian/stage
docker-compose build
docker run --rm fixuid-alpine
docker run --rm fixuid-centos
docker run --rm fixuid-debian
docker run --rm -u 1001:1001 fixuid-alpine
docker run --rm -u 1001:1001 fixuid-centos
docker run --rm -u 1001:1001 fixuid-debian
