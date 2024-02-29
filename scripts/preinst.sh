#!/bin/bash


if ! [ -d /var/lib/objex/ ]; then
    mkdir /var/lib/objex
fi

if [ -f "/etc/systemd/system/objex.service" ]; then
    systemctl stop objex
    systemctl disable objex
    systemctl daemon-reload
fi
