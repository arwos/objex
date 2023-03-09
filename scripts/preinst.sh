#!/bin/bash


if ! [ -d /var/lib/artifactory/ ]; then
    mkdir /var/lib/artifactory
fi

if [ -f "/etc/systemd/system/artifactory.service" ]; then
    systemctl stop artifactory
    systemctl disable artifactory
    systemctl daemon-reload
fi
