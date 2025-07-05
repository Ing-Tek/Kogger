#!/bin/bash
set -e

echo "Generating kogger package"
protoc --go_out=koggerrpc --go_opt=paths=source_relative  --go-grpc_out=koggerrpc --go-grpc_opt=paths=source_relative kogger.proto

echo "cleaning modules"
go mod tidy

set +e
echo "Building package"
goBuildResult=$(go build -a -installsuffix cgo -o app . 2>&1)
if [[ "$?" != "0" ]]; then
    echo -e "\e[31m!!Building package: fail \e[0m"
    echo -e "$goBuildResult"
    exit 
fi