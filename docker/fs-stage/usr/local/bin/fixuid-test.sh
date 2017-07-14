#!/bin/sh

expected_user=$1
expected_group=$2

set -e
fixuid
set +e

rc=0

user=$(id -u -n)
if [ "$user" != "$expected_user" ]
then
    >&2 echo "expected user: $expected_user, actual user: $user"
    rc=1
fi

group=$(id -g -n)
if [ "$group" != "$expected_group" ]
then
    >&2 echo "expected group: $expected_group, actual group: $group"
    rc=1
fi

files="/tmp/test-dir /tmp/test-dir/test-file /tmp/test-file /home/docker"
for file in $files
do
    file_user=$(stat -c "%U" $file)
    if [ "$file_user" != "$expected_user" ]
    then
        >&2 echo "$file expected owning user: $expected_user, actual owning user: $file_user"
        rc=1
    fi

    file_group=$(stat -c "%G" $file)
    if [ "$file_group" != "$expected_group" ]
    then
        >&2 echo "$file expected owning group: $expected_group, actual owning group: $file_group"
        rc=1
    fi

done

exit $rc
