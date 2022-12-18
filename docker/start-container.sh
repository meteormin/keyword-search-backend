#!/usr/bin/env bash

if [ ! -z "$GO_USER" ]; then
    usermod -u $GO_USER gofiber
fi

export GOPATH=/usr/local/go/bin
export PATH=$PATH:$GOPATH

/usr/local/go/bin/go mod download
make build

exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf