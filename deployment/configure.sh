#!/bin/bash

DEPLOYMENT_DIR=./deployment
failed=0

echo "[*] Checking necessary tools.."

print_status() {
	local component=$1
	local status=$2
	printf "* %-40s %s\n" "$component" "$status"
}

tools=( golangci-lint docker go protoc protoc-gen-go protoc-gen-go-grpc )

for tool in "${tools[@]}"; do
	if ! which $tool > /dev/null; then
		print_status $tool FAIL
		((failed++))
	else
		print_status $tool OK
	fi
done

echo "[*] Checking pip3 tools.."
pip3 install -q -r ${DEPLOYMENT_DIR}/requirements.txt

if [ $failed -gt 0 ]; then
	echo "[-] Configuration failed"
	exit 1
fi