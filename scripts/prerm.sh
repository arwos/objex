#!/bin/bash


if [ -f "/etc/systemd/system/artifactory.service" ]; then
    systemctl stop artifactory
    systemctl disable artifactory
    systemctl daemon-reload
fi
