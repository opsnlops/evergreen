#!/usr/bin/env bash

docker build -t 10gen/evergreen:latest \
    --build-arg GIT_HASH="$(git rev-parse --short HEAD)" \
    --build-arg GOOS="linux" \
    --build-arg GOARCH="arm64" \
    --build-arg GOVERSION="1.20.4" \
    .
