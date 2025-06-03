#!/bin/sh

set -e

source /app/app.env
echo "run db migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
echo "Executing: $@"
exec "$@"