#!/bin/bash

set -e

docker-compose -f e2e/docker-compose.yml up -d swn1 swn2

# TODO: run cwn as docker container
sleep 2

cp e2e/testdata/debug.yml e2e/cwn/
cd e2e/cwn/
go run main.go
cd -

docker-compose -f e2e/docker-compose.yml stop swn1 swn2
rm -f e2e/testdata/debug.yml