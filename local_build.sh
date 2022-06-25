#!/usr/bin/env sh
set -eu
export XCADDY_GO_BUILD_FLAGS="-ldflags '-s -w -extldflags -static'"
export GOPROXY=https://proxy.golang.com.cn,direct
export GOBIN=/mnt/d/devlops/caddy
export XCADDY_SKIP_CLEANUP=1
go install -v github.com/caddyserver/xcaddy/cmd/xcaddy@master

./xcaddy build \
    --with github.com/Elegant996/scgi-transport \
    --with github.com/HeavenVolkoff/caddy-authelia/plugin \
    --with github.com/RussellLuo/caddy-ext/requestbodyvar \
    --with github.com/RussellLuo/olaf/caddyconfig/adapter \
    --with github.com/WingLim/caddy-webhook \
    --with github.com/abiosoft/caddy-exec \
    --with github.com/abiosoft/caddy-hmac \
    --with github.com/abiosoft/caddy-json-parse \
    --with github.com/abiosoft/caddy-json-schema \
    --with github.com/abiosoft/caddy-named-routes \
    --with github.com/abiosoft/caddy-yaml \
    --with github.com/aksdb/caddy-cgi/v2 \
    --with github.com/baldinof/caddy-supervisor \
    --with github.com/caddy-dns/alidns \
    --with github.com/caddy-dns/azure \
    --with github.com/caddy-dns/cloudflare \
    --with github.com/caddy-dns/digitalocean \
    --with github.com/caddy-dns/dnspod \
    --with github.com/caddy-dns/duckdns \
    --with github.com/caddy-dns/gandi \
    --with github.com/caddy-dns/godaddy \
    --with github.com/caddy-dns/googleclouddns \
    --with github.com/caddy-dns/hetzner \
    --with github.com/caddy-dns/lego-deprecated \
    --with github.com/caddy-dns/metaname \
    --with github.com/caddy-dns/netcup \
    --with github.com/caddy-dns/netlify \
    --with github.com/caddy-dns/openstack-designate \
    --with github.com/caddy-dns/route53 \
    --with github.com/caddy-dns/vultr \
    --with github.com/caddyserver/jsonc-adapter \
    --with github.com/caddyserver/nginx-adapter@master \
    --with github.com/caddyserver/ntlm-transport \
    --with github.com/caddyserver/replace-response \
    --with github.com/caddyserver/transform-encoder \
    --with github.com/casbin/caddy-authz/v2 \
    --with github.com/chukmunnlee/caddy-openapi \
    --with github.com/circa10a/caddy-geofence \
    --with github.com/cubic3d/caddy-ct \
    --with github.com/cubic3d/caddy-quantity-limiter \
    --with github.com/darkweak/souin/plugins/caddy \
    --with github.com/gamalan/caddy-tlsredis \
    --with github.com/gbox-proxy/gbox \
    --with github.com/ggicci/caddy-jwt \
    --with github.com/git001/caddyv2-upload \
    --with github.com/greenpau/caddy-git \
    --with github.com/greenpau/caddy-security \
    --with github.com/greenpau/caddy-trace \
    --with github.com/hairyhenderson/caddy-teapot-module \
    --with github.com/hslatman/caddy-crowdsec-bouncer \
    --with github.com/imgk/caddy-trojan \
    --with github.com/kirsch33/realip \
    --with github.com/lindenlab/caddy-s3-proxy \
    --with github.com/lolPants/caddy-requestid \
    --with github.com/lucaslorentz/caddy-docker-proxy/v2 \
    --with github.com/mastercactapus/caddy2-proxyprotocol \
    --with github.com/mholt/caddy-dynamicdns \
    --with github.com/mholt/caddy-l4 \
    --with github.com/mholt/caddy-ratelimit \
    --with github.com/mholt/caddy-webdav \
    --with github.com/mpilhlt/caddy-conneg \
    --with github.com/muety/caddy-pirsch-plugin \
    --with github.com/muety/caddy-remote-host \
    --with github.com/porech/caddy-maxmind-geolocation \
    --with github.com/pteich/caddy-tlsconsul \
    --with github.com/shift72/caddy-geo-ip \
    --with github.com/sillygod/cdp-cache \
    --with github.com/sjtug/caddy2-filter \
    --with github.com/techknowlogick/caddy-s3browser \
    --with github.com/tosie/caddy-dns-linode \
    --with github.com/ueffel/caddy-basic-auth-filter \
    --with github.com/ueffel/caddy-brotli \
    --with github.com/ueffel/caddy-imagefilter/defaults \
    --with github.com/ueffel/caddy-tls-format \
    --with magnax.ca/caddy/gopkg \
    --with github.com/darkweak/souin@v1.6.6
