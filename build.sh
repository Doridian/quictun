#!/bin/sh
set -ex

go build -o quictun -trimpath -ldflags="-w -s -X main.VERSION=$(git rev-parse HEAD)"
