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

echo "[*] running cwn1"
run &

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