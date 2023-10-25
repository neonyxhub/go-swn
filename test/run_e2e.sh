#!/bin/bash

set -e

run_with_pwn=${1:-cwn}
DOCKER_ARGS=${2:-}

export RUN_WITH_PWN=$run_with_pwn

if ! docker network ls --format '{{.Name}}' | grep e2e > /dev/null;then
	echo "[*] creating a e2e docker network"
	docker network create e2e
fi

echo "[*] running swn1 swn2"
echo "peers: []" > e2e/testdata/debug.yml

echo "[*] running cwn1"
docker-compose -f e2e/docker-compose.yml up ${DOCKER_ARGS} cwn1

while true; do
	status=$(docker-compose -f e2e/docker-compose.yml ps cwn1 | grep 'Exit' || true)
	if [[ ! -z "$status" ]]; then
		echo "cwn1 has exited, stopping other services..."
		docker-compose -f e2e/docker-compose.yml stop swn1 swn2
		break
	fi
	echo "[*] waiting for cwn1 to exit..."
	sleep 1
done

rm -f e2e/testdata/debug.yml
docker rm cwn1 swn1 swn2
docker network rm e2e