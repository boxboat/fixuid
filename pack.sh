#!/bin/sh
cd $(dirname $0)

display_usage() {
    echo "Usage:\n$0 [version]"
}

# check whether user had supplied -h or --help . If yes display usage
if [ $# = "--help" ] || [ $# = "-h" ]
then
    display_usage
    exit 0
fi

# check number of arguments
if [ $# -ne 1 ]
then
    display_usage
    exit 1
fi

./build.sh
rm -f fixuid-*-linux-amd64.tar.gz
perm="$(id -u):$(id -g)"
sudo chown root:root fixuid
sudo chmod u+s fixuid
tar -cvzf fixuid-$1-linux-amd64.tar.gz fixuid
sudo chmod u-s fixuid
sudo chown $perm fixuid

