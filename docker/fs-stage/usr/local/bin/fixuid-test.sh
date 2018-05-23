#!/bin/sh

expected_user=$1
expected_group=$2

if [ ! -f /var/run/fixuid.ran ]
then
    set -e
    set_home=$( fixuid $FIXUID_FLAGS )
    set +e
    eval $set_home
fi

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

OLD_IFS="$IFS"
IFS="|"
files="/tmp/test-dir|/tmp/test-dir/test-file|/tmp/test-file|/home/docker|/home/docker/aaa|/home/docker/zzz|/tmp/space dir|/tmp/space dir/space file|/tmp/space file"
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
IFS="$OLD_IFS"

if [ "$user" = "root" ]
then
    if [ "$HOME" != "/root" ]
    then
        >&2 echo "expected home directory: /root, actual home directory: $HOME"
        rc=1
    fi
elif [ "$HOME" != "/home/$user" ]
then
    >&2 echo "expected home directory: /home/$user, actual home directory: $HOME"
    rc=1
fi

>&2 echo "test complete, RC=$rc"
exit $rc
