#!/bin/bash

TESTDOMAIN=$1
BINARYNAME="a24api"

if [ "${TESTDOMAIN}" == "" ]; then echo "Usage: $0 domain"; exit 1; fi

mkdir -p ./build
go build -ldflags="-s -w" -o ./build/${BINARYNAME} && /usr/bin/upx --brute ./build/${BINARYNAME}

