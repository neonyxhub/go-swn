#!/bin/bash

set -e

swn() {
	docker-compose -f e2e/docker-compose.yml up swn1 swn2
}

echo "peers: []" > e2e/testdata/debug.yml
swn &

docker-compose -f e2e/docker-compose.yml up cwn1

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