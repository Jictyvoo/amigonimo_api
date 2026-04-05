#!/bin/sh

# URL-encode @ in the MySQL password
# (replace @ with %40)
DB_PASSWORD_URLENCODED="$(printf %s "$DATABASE_PASSWORD" | sed 's/@/%40/g')"

# Run the Atlas migration
atlas migrate apply \
  --dir file://build/migrations \
  --url "mysql://${DATABASE_USER}:${DB_PASSWORD_URLENCODED}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}"

# Start the server
exec /bin/anonymigo_api
