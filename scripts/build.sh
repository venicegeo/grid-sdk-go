#!/bin/bash

go test */**_test.go

if [[ $TRAVIS_TAG =~ ^v[0-9\.]+ ]]; then
  go get github.com/mitchellh/gox
  gox -os="linux" -arch="amd64" -build-toolchain
  gox -os="linux" -arch="amd64" -output="grid_{{.OS}}_x86_64"
fi
