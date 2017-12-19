#!/bin/sh

set -o errexit
set -o pipefail

VENDOR=informaticslab
DRIVER=s3-fuse-flex-volume
PLUGIN_DIR=/rootfs/usr/libexec/kubernetes/kubelet-plugins/volume/exec/

## Install driver
# Assuming the single driver file is located at /$DRIVER inside the DaemonSet image.

driver_dir=$VENDOR${VENDOR:+"~"}${DRIVER}
if [ ! -d "$PLUGIN_DIR/$driver_dir" ]; then
    mkdir "$PLUGIN_DIR/$driver_dir"
fi

cp "/$DRIVER" "$PLUGIN_DIR/$driver_dir/.$DRIVER"
mv -f "$PLUGIN_DIR/$driver_dir/.$DRIVER" "$PLUGIN_DIR/$driver_dir/$DRIVER"
