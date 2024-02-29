#!/bin/bash


if [ -f "/etc/systemd/system/objex.service" ]; then
    systemctl stop objex
    systemctl disable objex
    systemctl daemon-reload
fi
