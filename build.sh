#!/usr/bin/env sh

set -eu

cd $(dirname $0)
ROOT=$(pwd)

os=$(uname | tr 'A-Z' 'a-z')
arch=$(uname -m)
case $arch in
x86_64)
    arch=amd64
    ;;
armv8 | arm8 | arm64)
    arch=arm64
    ;;
armv*)
    arch=arm
    ;;
*)
    echo "not support arch ${arch}"
    ;;
esac

echo "build caddy ${os} ${arch}"

docker run --rm \
    -v $(pwd):/build \
    -e CGO_ENABLED=0 \
    -e GOOS=${os} \
    -e GOARCH=${arch} \
    -e GOPROXY=https://goproxy.cn \
    caddy:builder \
    xcaddy build --output /build/caddy \
    --with github.com/caddyserver/caddy/v2/modules/standard \
    --with github.com/caddy-dns/alidns \
    --with github.com/caddy-dns/azure \
    --with github.com/caddy-dns/cloudflare \
    --with github.com/caddy-dns/digitalocean \
    --with github.com/caddy-dns/dnspod \
    --with github.com/caddy-dns/duckdns \
    --with github.com/caddy-dns/gandi \
    --with github.com/caddy-dns/hetzner \
    --with github.com/caddy-dns/openstack-designate \
    --with github.com/caddy-dns/route53 \
    --with github.com/caddy-dns/vultr \
    --with github.com/mholt/caddy-dynamicdns \
    --with github.com/abiosoft/caddy-exec \
    --with github.com/hslatman/caddy-crowdsec-bouncer/crowdsec \
    --with github.com/ss098/certmagic-s3 \
    --with github.com/gamalan/caddy-tlsredis \
    --with github.com/silinternational/certmagic-storage-dynamodb/v2 \
    --with github.com/pteich/caddy-tlsconsul \
    --with github.com/caddyserver/format-encoder \
    --with github.com/mastercactapus/caddy2-proxyprotocol \
    --with github.com/ggicci/caddy-jwt \
    --with github.com/ueffel/caddy-brotli \
    --with github.com/HeavenVolkoff/caddy-authelia/plugin \
    --with github.com/greenpau/caddy-auth-portal \
    --with github.com/casbin/caddy-authz/v2 \
    --with github.com/aksdb/caddy-cgi/v2 \
    --with github.com/hslatman/caddy-crowdsec-bouncer/http \
    --with github.com/cubic3d/caddy-ct \
    --with github.com/dunglas/mercure/caddy \
    --with github.com/abiosoft/caddy-json-parse \
    --with github.com/abiosoft/caddy-hmac \
    --with magnax.ca/caddy/gopkg \
    --with github.com/sjtug/caddy2-filter \
    --with github.com/caddyserver/replace-response \
    --with github.com/kirsch33/realip \
    --with github.com/mholt/caddy-ratelimit \
    --with github.com/cubic3d/caddy-quantity-limiter \
    --with github.com/lindenlab/caddy-s3-proxy \
    --with github.com/lolPants/caddy-requestid \
    --with github.com/caddyserver/ntlm-transport \
    --with github.com/porech/caddy-maxmind-geolocation \
    --with github.com/WingLim/caddy-webhook \
    --with github.com/mholt/caddy-webdav \
    --with github.com/dunglas/vulcain/caddy \
    --with github.com/greenpau/caddy-trace \
    --with github.com/mholt/caddy-l4/layer4 \
    --with github.com/mholt/caddy-l4/modules/l4echo \
    --with github.com/mholt/caddy-l4/modules/l4proxy \
    --with github.com/mholt/caddy-l4/modules/l4tee \
    --with github.com/mholt/caddy-l4/modules/l4tls \
    --with github.com/mholt/caddy-l4/modules/l4ssh \
    --with github.com/mholt/caddy-l4/modules/l4http \
    --with github.com/hslatman/caddy-crowdsec-bouncer/layer4 \
    --with github.com/baldinof/caddy-supervisor

./caddy list-modules
./caddy version
