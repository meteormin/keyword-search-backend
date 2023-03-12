#!/usr/bin/env bash

echo "[DOCKER COMPOSE]"

docker_compose="docker compose"

if ! docker info >/dev/null 2>&1; then
  echo "Docker is not running." >&2
  exit 1
fi

if ! $docker_compose >/dev/null 2>&1; then
  echo "docker compose is not installed"
  exit 1
fi

if [ -f .env ]; then
  app_env=$(cat .env | grep 'APP_ENV=')
  IFS='=' read -ra app_env_arr <<<"$app_env"
  app_env=${app_env_arr[1]}

  if [ ! -f ".env.${app_env}" ]; then
    echo "copy .env .env.${app_env}"
    cp .env ".env.${app_env}"
  fi

fi

echo "$docker_compose up -d --build"
$docker_compose up -d --build

echo "[CLEAN UP IMAGES] '<none>'"
docker image rm $(docker image list -f 'dangling=true' -q --no-trunc)
echo ""

echo "[CLEAN UP BUILDER] until 24h"
echo "y" | docker builder prune --filter until=24h
echo ""