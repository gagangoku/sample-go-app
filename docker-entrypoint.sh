#!/bin/bash

set -x
set -e

PORT="${PORT:-4001}"
PRETTY_LOG="${PRETTY_LOG:-true}"

./sample-go-app.out --port="$PORT" --prettyLog="$PRETTY_LOG"
