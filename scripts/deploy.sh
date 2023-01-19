#!/usr/bin/env bash
BASEDIR=$(dirname "$0")

"$BASEDIR"/git-pull.sh
"$BASEDIR"/install.sh
"$BASEDIR"/docker-compose.sh
