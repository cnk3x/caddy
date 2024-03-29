#!/usr/bin/env sh

set -eu

cd $(dirname $0)
ROOT=$(pwd)

tag=${1:-ghcr.io/cnk3x/caddy:latest}
docker build --tag ${tag} .
docker run --rm ${tag} list-modules
docker run --rm ${tag} version
echo "you can use \`docker push ${tag}\` to publish this repo"
