#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

PACKAGE_PATH='github.com/pgillich/sample-blog'

BUILD_TAG=$(git describe --tags)
BUILD_COMMIT=$(git log -1 --pretty=format:%h)
BUILD_BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIME=$(date --rfc-3339=seconds | tr ' ' 'T')

go mod verify
go vet
golangci-lint run
go test -v ./...

echo "BUILD_TAG=$BUILD_TAG BUILD_COMMIT=$BUILD_COMMIT BUILD_BRANCH=$BUILD_BRANCH BUILD_TIME=$BUILD_TIME"
go build -ldflags "-X $PACKAGE_PATH/configs.BuildTag=$BUILD_TAG -X $PACKAGE_PATH/configs.BuildCommit=$BUILD_COMMIT -X $PACKAGE_PATH/configs.BuildBranch=$BUILD_BRANCH -X $PACKAGE_PATH/configs.BuildTime=$BUILD_TIME"

echo "Building image ${PACKAGE_PATH}:$BUILD_TAG"
docker build --tag "${PACKAGE_PATH}:$BUILD_TAG" --build-arg BUILD_TAG=$BUILD_TAG --build-arg BUILD_COMMIT=$BUILD_COMMIT --build-arg BUILD_BRANCH=$BUILD_BRANCH --build-arg BUILD_TIME=$BUILD_TIME .
