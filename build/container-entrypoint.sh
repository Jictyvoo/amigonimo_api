#!/bin/sh

# URL-encode @ in the MySQL password
# (replace @ with %40)
MYSQL_PASSWORD_URLENCODED="$(printf %s "$MYSQL_PASSWORD" | sed 's/@/%40/g')"

# Run the Atlas migration
atlas migrate apply \
  --dir file://build/migrations \
  --url "mysql://${MYSQL_USER}:${MYSQL_PASSWORD_URLENCODED}@localhost:3306/${MYSQL_DATABASE}"

# Start the server
exec /bin/anonymigo_api
