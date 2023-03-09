#!/bin/bash


if [ -f "/etc/systemd/system/artifactory.service" ]; then
    systemctl start artifactory
    systemctl enable artifactory
    systemctl daemon-reload
fi
