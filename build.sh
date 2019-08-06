#!/bin/sh

export GOPROXY="https://goproxy.io"
export CGO_ENABLED=0

set -e

oss="linux darwin windows"
archs="amd64 386"

mkdir -p bin
for os in $oss; do
    for arch in $archs; do
        echo "build caddy-${os}-${arch}"
        GOOS=${os} GOARCH=${arch} go build -ldflags '-s -w' -v -o caddy-${os}-${arch}
        tar zcf caddy-${os}-${arch}.tar.gz caddy-${os}-${arch}
        mv caddy-${os}-${arch}* bin/
    done
done
