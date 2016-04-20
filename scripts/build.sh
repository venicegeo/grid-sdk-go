#!/bin/bash

gox -osarch="darwin/amd64 linux/amd64 windows/amd64" \
    -output "dist/{{.Dir}}-{{.OS}}-{{.Arch}}"
