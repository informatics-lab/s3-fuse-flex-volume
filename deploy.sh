#!/bin/sh

set -o errexit
set -o pipefail

VENDOR=informaticslab
DRIVER=s3-fuse-flex-volume

# Assuming the single driver file is located at /$DRIVER inside the DaemonSet image.

driver_dir=$VENDOR${VENDOR:+"~"}${DRIVER}
if [ ! -d "/flexmnt/$driver_dir" ]; then
    mkdir "/flexmnt/$driver_dir"
fi

cp "/$DRIVER" "/flexmnt/$driver_dir/.$DRIVER"
mv -f "/flexmnt/$driver_dir/.$DRIVER" "/flexmnt/$driver_dir/$DRIVER"

cp "/usr/local/src/s3fuse.py" "/hostbin/.s3fuse.py"
mv -f "/hostbin/.s3fuse.py" "/hostbin/s3fuse.py"

cp "/usr/local/src/s3fuseenv.pex" "/hostbin/.s3fuseenv.pex"
mv -f "/hostbin/.s3fuseenv.pex" "/hostbin/s3fuseenv.pex"

while : ; do
    sleep 3600
done
