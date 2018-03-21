#!/usr/bin/env bash

set -o errexit
set -o pipefail

VENDOR=informaticslab
declare -a DRIVERS=("pysssix-flex-volume-driver" "goofys-flex-volume-driver")

for DRIVER in "${DRIVERS[@]}"
do
    echo "Installing $DRIVER"
    driver_dir=$VENDOR${VENDOR:+"~"}${DRIVER}
    if [ ! -d "/flexmnt/$driver_dir" ]; then
        mkdir "/flexmnt/$driver_dir"
    fi

    cp "/$DRIVER" "/flexmnt/$driver_dir/.$DRIVER"
    mv -f "/flexmnt/$driver_dir/.$DRIVER" "/flexmnt/$driver_dir/$DRIVER"
done

echo "Listing installed drivers:"
ls /flexmnt
