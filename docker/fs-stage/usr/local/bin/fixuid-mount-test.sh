#!/bin/sh

expected_uid=$1
expected_gid=$2

rc=0

files="/home/docker/tmp/test-dir /home/docker/tmp/test-dir/test-file /home/docker/tmp/test-file"
for file in $files
do
    file_uid=$(stat -c "%u" $file)
    if [ "$file_uid" != "$expected_uid" ]
    then
        >&2 echo "$file expected owning uid: $expected_uid, actual owning uid: $file_uid"
        rc=1
    fi

    file_gid=$(stat -c "%g" $file)
    if [ "$file_gid" != "$expected_gid" ]
    then
        >&2 echo "$file expected owning gid: $expected_gid, actual owning gid: $file_gid"
        rc=1
    fi

done

>&2 echo "mount test complete, RC=$rc"
exit $rc
