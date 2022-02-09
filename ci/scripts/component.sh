#!/bin/bash -eux

pushd dp-frontend-release-calendar
  make test-component
popd
