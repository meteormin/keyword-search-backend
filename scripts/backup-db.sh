#!/usr/bin/env bash
BASEDIR="$(dirname "$0")/../"
DATA_DIR="$BASEDIR/backup/db"
ENV_FILE="$BASEDIR/.env"
CONTAINER_NAME="$1"

if [ "$CONTAINER_NAME" == "" ]; then
  echo "first parameter is required. input your db container name"
  return 1
fi

if [ ! -d "$DATA_DIR" ]; then
  mkdir "$DATA_DIR"
fi

export "$(grep -v "^#" "$ENV_FILE" | xargs -d "\n")"

docker exec "$CONTAINER_NAME" pg_dump -h localhost -U "$DB_USERNAME" -d "$DB_DATABASE" -F t > "${DATA_DIR}/pgsql_$(date "+%Y-%m-%d").tar"