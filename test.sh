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
docker compose build

echo "\nalpine default user/group cmd"
docker run --rm fixuid-alpine fixuid-test.sh docker docker
echo "\ncentos default user/group cmd"
docker run --rm fixuid-centos fixuid-test.sh docker docker
echo "\ndebian default user/group cmd"
docker run --rm fixuid-debian fixuid-test.sh docker docker
echo "\nalpine default user/group entrypoint"
docker run --rm --entrypoint fixuid fixuid-alpine fixuid-test.sh docker docker
echo "\ncentos default user/group entrypoint"
docker run --rm --entrypoint fixuid fixuid-centos fixuid-test.sh docker docker
echo "\ndebian default user/group entrypoint"
docker run --rm --entrypoint fixuid fixuid-debian fixuid-test.sh docker docker

echo "\nalpine 1001:1001 cmd"
docker run --rm -u 1001:1001 fixuid-alpine fixuid-test.sh docker docker
echo "\ncentos 1001:1001 cmd"
docker run --rm -u 1001:1001 fixuid-centos fixuid-test.sh docker docker
echo "\ndebian 1001:1001 cmd"
docker run --rm -u 1001:1001 fixuid-debian fixuid-test.sh docker docker
echo "\nalpine 1001:1001 entrypoint"
docker run --rm -u 1001:1001 --entrypoint fixuid fixuid-alpine fixuid-test.sh docker docker
echo "\ncentos 1001:1001 entrypoint"
docker run --rm -u 1001:1001 --entrypoint fixuid fixuid-centos fixuid-test.sh docker docker
echo "\ndebian 1001:1001 entrypoint"
docker run --rm -u 1001:1001 --entrypoint fixuid fixuid-debian fixuid-test.sh docker docker "docker users"

echo "\nalpine 0:0 cmd"
docker run --rm -u 0:0 fixuid-alpine fixuid-test.sh root root
echo "\ncentos 0:0 cmd"
docker run --rm -u 0:0 fixuid-centos fixuid-test.sh root root
echo "\ndebian 0:0 cmd"
docker run --rm -u 0:0 fixuid-debian fixuid-test.sh root root
echo "\nalpine 0:0 entrypoint"
docker run --rm -u 0:0 --entrypoint fixuid fixuid-alpine fixuid-test.sh root root
echo "\ncentos 0:0 entrypoint"
docker run --rm -u 0:0 --entrypoint fixuid fixuid-centos fixuid-test.sh root root
echo "\ndebian 0:0 entrypoint"
docker run --rm -u 0:0 --entrypoint fixuid fixuid-debian fixuid-test.sh root root

echo "\nalpine 0:1001 cmd"
docker run --rm -u 0:1001 fixuid-alpine fixuid-test.sh root docker
echo "\ncentos 0:1001 cmd"
docker run --rm -u 0:1001 fixuid-centos fixuid-test.sh root docker
echo "\ndebian 0:1001 cmd"
docker run --rm -u 0:1001 fixuid-debian fixuid-test.sh root docker
echo "\nalpine 0:1001 entrypoint"
docker run --rm -u 0:1001 --entrypoint fixuid fixuid-alpine fixuid-test.sh root docker
echo "\ncentos 0:1001 entrypoint"
docker run --rm -u 0:1001 --entrypoint fixuid fixuid-centos fixuid-test.sh root docker
echo "\ndebian 0:1001 entrypoint"
docker run --rm -u 0:1001 --entrypoint fixuid fixuid-debian fixuid-test.sh root docker "docker users"

echo "\nalpine 1001:0 cmd"
docker run --rm -u 1001:0 fixuid-alpine fixuid-test.sh docker root
echo "\ncentos 1001:0 cmd"
docker run --rm -u 1001:0 fixuid-centos fixuid-test.sh docker root
echo "\ndebian 1001:0 cmd"
docker run --rm -u 1001:0 fixuid-debian fixuid-test.sh docker root
echo "\nalpine 1001:0 entrypoint"
docker run --rm -u 1001:0 --entrypoint fixuid fixuid-alpine fixuid-test.sh docker root "root docker"
echo "\ncentos 1001:0 entrypoint"
docker run --rm -u 1001:0 --entrypoint fixuid fixuid-centos fixuid-test.sh docker root "root docker"
echo "\ndebian 1001:0 entrypoint"
docker run --rm -u 1001:0 --entrypoint fixuid fixuid-debian fixuid-test.sh docker root "root users docker"

echo "\nalpine run twice cmd"
docker run --rm fixuid-alpine sh -c "fixuid-test.sh docker docker && fixuid fixuid-test.sh docker docker"
echo "\ncentos run twice cmd"
docker run --rm fixuid-centos sh -c "fixuid-test.sh docker docker && fixuid fixuid-test.sh docker docker"
echo "\ndebian run twice cmd"
docker run --rm fixuid-debian sh -c "fixuid-test.sh docker docker && fixuid fixuid-test.sh docker docker"
echo "\nalpine run twice entrypoint"
docker run --rm --entrypoint fixuid fixuid-alpine sh -c "fixuid-test.sh docker docker && fixuid fixuid-test.sh docker docker"
echo "\ncentos run twice entrypoint"
docker run --rm --entrypoint fixuid fixuid-centos sh -c "fixuid-test.sh docker docker && fixuid fixuid-test.sh docker docker"
echo "\ndebian run twice entrypoint"
docker run --rm --entrypoint fixuid fixuid-debian sh -c "fixuid-test.sh docker docker && fixuid fixuid-test.sh docker docker"

echo "\nalpine should not chown mount"
docker run --rm -v $(pwd)/docker/fs-stage/tmp:/home/docker/mnt-dir -v $(pwd)/docker/fs-stage/tmp/test-file:/home/docker/mnt-file -u 1234:1234 fixuid-alpine sh -c "fixuid-test.sh docker docker && fixuid-mount-test.sh $(id -u) $(id -g)"
echo "\ncentos should not chown mount"
docker run --rm -v $(pwd)/docker/fs-stage/tmp:/home/docker/mnt-dir -v $(pwd)/docker/fs-stage/tmp/test-file:/home/docker/mnt-file -u 1234:1234 fixuid-centos sh -c "fixuid-test.sh docker docker && fixuid-mount-test.sh $(id -u) $(id -g)"
echo "\ndebian should not chown mount"
docker run --rm -v $(pwd)/docker/fs-stage/tmp:/home/docker/mnt-dir -v $(pwd)/docker/fs-stage/tmp/test-file:/home/docker/mnt-file -u 1234:1234 fixuid-debian sh -c "fixuid-test.sh docker docker && fixuid-mount-test.sh $(id -u) $(id -g)"

echo "\nalpine quiet cmd"
docker run --rm -e "FIXUID_FLAGS=-q" fixuid-alpine fixuid-test.sh docker docker
echo "\ncentos quiet cmd"
docker run --rm -e "FIXUID_FLAGS=-q" fixuid-centos fixuid-test.sh docker docker
echo "\ndebian quiet cmd"
docker run --rm -e "FIXUID_FLAGS=-q" fixuid-debian fixuid-test.sh docker docker
echo "\nalpine quiet entrypoint"
docker run --rm --entrypoint fixuid fixuid-alpine -q fixuid-test.sh docker docker
echo "\ncentos quiet entrypoint"
docker run --rm --entrypoint fixuid fixuid-centos -q fixuid-test.sh docker docker
echo "\ndebian quiet entrypoint"
docker run --rm --entrypoint fixuid fixuid-debian -q fixuid-test.sh docker docker

printf "\npaths:\n  - /\n  - /home/docker\n  - /tmp/space dir\n  - /does/not/exist" >> docker/alpine/stage/etc/fixuid/config.yml
printf "\npaths:\n  - /\n  - /home/docker\n  - /tmp/space dir\n  - /does/not/exist" >> docker/centos/stage/etc/fixuid/config.yml
printf "\npaths:\n  - /\n  - /home/docker\n  - /tmp/space dir\n  - /does/not/exist" >> docker/debian/stage/etc/fixuid/config.yml
docker compose build

echo "\nalpine 1001:1001 cmd"
docker run --rm -u 1001:1001 -v /home/docker -v "/tmp/space dir" fixuid-alpine fixuid-test.sh docker docker
echo "\ncentos 1001:1001 cmd"
docker run --rm -u 1001:1001 -v /home/docker -v "/tmp/space dir" fixuid-centos fixuid-test.sh docker docker
echo "\ndebian 1001:1001 cmd"
docker run --rm -u 1001:1001 -v /home/docker -v "/tmp/space dir" fixuid-debian fixuid-test.sh docker docker
echo "\nalpine 1001:1001 entrypoint"
docker run --rm -u 1001:1001 -v /home/docker -v "/tmp/space dir" --entrypoint fixuid fixuid-alpine fixuid-test.sh docker docker
echo "\ncentos 1001:1001 entrypoint"
docker run --rm -u 1001:1001 -v /home/docker -v "/tmp/space dir" --entrypoint fixuid fixuid-centos fixuid-test.sh docker docker
echo "\ndebian 1001:1001 entrypoint"
docker run --rm -u 1001:1001 -v /home/docker -v "/tmp/space dir" --entrypoint fixuid fixuid-debian fixuid-test.sh docker docker "docker users"
