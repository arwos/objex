#!/bin/bash

# ------- Service -----------------------------------------------------
if [ -f "/etc/systemd/system/artifactory.service" ]; then
    systemctl stop artifactory
    systemctl disable artifactory
    systemctl daemon-reload
fi

# ------- User -----------------------------------------------------
set -e

if ! getent group artifactory >/dev/null; then
	addgroup --system artifactory >/dev/null
fi

if ! getent passwd artifactory >/dev/null; then
	adduser \
	  --system \
          --disabled-login \
	  --ingroup artifactory \
	  --no-create-home \
	  --home /nonexistent \
	  --gecos "Artifactory Server" \
	  --shell /bin/false \
	  artifactory  >/dev/null
fi

# ------- Store -----------------------------------------------------
if ! [ -d /var/lib/artifactory/ ]; then
    mkdir /var/lib/artifactory
fi

chown artifactory:artifactory -R /var/lib/artifactory

# ------- Log -----------------------------------------------------
if ! [ -f /var/log/artifactory.log ]; then
    touch /var/log/artifactory.log
fi

chown artifactory:artifactory /var/log/artifactory.log

exit 0
