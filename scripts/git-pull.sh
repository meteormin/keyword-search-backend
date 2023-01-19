#!/usr/bin/env bash

branch_name=$1
repo_name="origin"

echo "[git pull]"

if ! git --version >/dev/null 2>&1; then
  echo "git is not installed"
  exit 1
fi

echo "git pull $repo_name $branch_name"
git pull $repo_name "$branch_name"

if [ ! -f ".env" ]; then
  echo "generate .env file"
  cp .env.example .env
fi

if [ ! -f ".env.local" ]; then
  echo "generate .env.local file"
  cp .env .env.local
fi
