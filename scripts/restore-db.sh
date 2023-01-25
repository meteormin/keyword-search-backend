#!/usr/bin/env bash
BASEDIR="$(dirname "$0")/../"
DATA_DIR="$BASEDIR/backup/db"
ENV_FILE="$BASEDIR/.env"
CONTAINER_NAME="$1"
BACKUP_FILENAME="$2"

if [ "$CONTAINER_NAME" == "" ]; then
  echo "first parameter is required. input your db container name"
  return 1
fi

if [ ! -f "$DATA_DIR/$BACKUP_FILENAME" ]; then
    echo "file not exist"
    return 1
fi

export "$(grep -v "^#" "$ENV_FILE" | xargs -d "\n")"

docker cp "$DATA_DIR/$BACKUP_FILENAME" "go-pgsql:/tmp/"
docker exec -i "$CONTAINER_NAME" pg_restore -h localhost -U "$DB_USERNAME" -d "$DB_DATABASE" -c -F t -v "/tmp/$BACKUP_FILENAME"
docker exec "$CONTAINER_NAME" rm "/temp/$BACKUP_FILENAME"