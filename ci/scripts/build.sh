#!/bin/bash -eux

pushd dp-frontend-release-calendar
  make build
  cp build/dp-frontend-release-calendar Dockerfile.concourse ../build
popd
