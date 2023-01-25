#!/usr/bin/env bash

echo "[go mod download]"

if ! go version > /dev/null 2>&1; then
  echo "go is not installed." >&2
  exit 1
fi

echo "go mod download"
go mod download

echo "go install swaggo"
go install github.com/swaggo/swag/cmd/swag@latest

echo "go install godoc"
go install golang.org/x/tools/cmd/godoc@latest