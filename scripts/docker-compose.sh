#!/usr/bin/env bash

echo "[docker compose up]"

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