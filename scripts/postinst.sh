#!/bin/bash

set -e

#------- Create DB ------------------------------------------
if grep -q 'new-artifactory-user-passwd' /etc/artifactory/config.yaml; then

  DBPASSWD="$(openssl rand -base64 20)"
  DBPASSWD=${DBPASSWD//[^0-9a-zA-Z]/X}
  DBUSER="artifactory"

  RESULT_VARIABLE="$(mysql -u root -sse "SELECT EXISTS(SELECT 1 FROM mysql.user WHERE user = '${DBUSER}')")"
  if [ "$RESULT_VARIABLE" = 1 ]; then
    mysql -u root -e "ALTER USER '${DBUSER}'@'localhost' IDENTIFIED BY '${DBPASSWD}';"
  else
    mysql -u root -e "CREATE DATABASE '${DBUSER}' /*\!40100 DEFAULT CHARACTER SET utf8 */;"
    mysql -u root -e "CREATE USER '${DBUSER}'@'localhost' IDENTIFIED BY '${DBPASSWD}';"
    mysql -u root -e "GRANT ALL PRIVILEGES ON '${DBUSER}'.* TO '${DBUSER}'@'localhost';"
  fi
  mysql -u root -e "FLUSH PRIVILEGES;"

  sed -i "s/new-artifactory-user-passwd/${DBPASSWD}/g" /etc/artifactory/config.yaml

  CURTIME="$(date +%s)"
  cp /etc/artifactory/config.yaml "/etc/artifactory/config.yaml.${CURTIME}"
fi

#------- Service ------------------------------------------
set +e

if [ -f "/etc/systemd/system/artifactory.service" ]; then
  systemctl start artifactory
  systemctl enable artifactory
  systemctl daemon-reload
fi
