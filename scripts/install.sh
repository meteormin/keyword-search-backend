#!/usr/bin/env bash

echo "[GO MOD INSTALL & DOWNLOAD]"

if ! go version > /dev/null 2>&1; then
  echo "go is not installed." >&2
  echo ""
  exit 1
fi

echo "go mod download"
go mod download
echo ""

echo "go install swaggo"
go install github.com/swaggo/swag/cmd/swag@latest
echo ""

echo "go install godoc"
go install golang.org/x/tools/cmd/godoc@latest
echo ""
