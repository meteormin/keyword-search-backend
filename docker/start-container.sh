#!/usr/bin/env bash

if [ ! -z "$GO_USER" ]; then
    usermod -u $GO_USER gofiber
fi

export GOPATH=/usr/local/go/bin
export PATH=$PATH:$GOPATH

/usr/local/go/bin/go mod download

swag --version &> /dev/null

if [ $? != 0 ]; then
  go install github.com/swaggo/swag/cmd/swag@latest
  export SWAG_PATH=/usr/local/go/bin/bin
  export PATH=$PATH:$SWAG_PATH
fi

make build

exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf