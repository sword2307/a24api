#!/bin/bash

mkdir -p ./build

for platform in 386 amd64 arm arm64; do
    for os in linux darwin windows; do

        if [ "${os}" == "windows" -o "${os}" == "darwin" ] && [ "${platform}" == "arm" -o "${platform}" == "arm64" ]; then continue; fi

        BINARYNAME=a24api-${os}-${platform}
        if [ "${os}" == "windows" ]; then BINARYNAME="${BINARYNAME}.exe"; fi

        echo "################################################################################"
        echo "## Building for ${os}-${platform}"
        echo "################################################################################"
        env GOOS=${os} GOARCH=${platform} go build -ldflags="-s -w" -o ./build/${BINARYNAME} && /usr/bin/upx --brute ./build/${BINARYNAME}

    done
done
