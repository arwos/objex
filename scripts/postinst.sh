#!/bin/bash


if [ -f "/etc/systemd/system/objex.service" ]; then
    systemctl start objex
    systemctl enable objex
    systemctl daemon-reload
fi
