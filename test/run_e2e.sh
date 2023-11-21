#!/bin/bash
set -e

if ! docker network ls --format '{{.Name}}' | grep e2e > /dev/null;then
	echo "[*] creating a e2e docker network"
	docker network create e2e
fi

# cleanup on prev. run
rm -f e2e/testdata/debug.yml
echo "peers: []" > e2e/testdata/debug.yml
docker rm e2e-swn-provider e2e-nats-server

# run

rm -f e2e/testdata/debug.yml
docker rm cwn1
docker network rm e2e