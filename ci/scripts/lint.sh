#!/bin/bash -eux

pushd dp-frontend-release-calendar
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
  make lint
popd
