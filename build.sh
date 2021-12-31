#!/usr/bin/env sh

set -eu

cd $(dirname $0)
ROOT=$(pwd)

export GOPROXY="https://goproxy.cn,direct"
export OUTPUT=${ROOT}/release
export GOWORK=off
export CGO_ENABLED=0
export GOROOT=$(go env GOROOT)
export GOOS=$(go env GOOS)
export GOARCH=$(go env GOARCH)

echo "build caddy ${GOOS} ${GOARCH}"
archive="${OUTPUT}/caddy.tar.gz"
binary="caddy"
if [ "${GOOS}" == "windows" ]; then
    binary=${binary}.exe
fi

# -workfile=off
${GOROOT}/bin/go build -mod=vendor -ldflags '-extldflags "-static"' -o ${binary} .
tar zcvf ${archive} ${binary}

${root}/${binary} list-modules
${root}/${binary} version

rm ${binary}

echo "build file at ${archive}"
