#!/usr/bin/env sh

set -e

root=$(
    cd $(dirname $0)
    pwd
)

cd ${root}

modules=$(cat ${root}/package.json |
    grep '"path": "' |
    sed 's/"path": "//g; s/",//g; s/ //g' |
    grep -v 'firecow/caddy-forward-auth' |
    grep -v 'RussellLuo/caddy-ext/ratelimit' |
    grep -v 'francislavoie/caddy-hcl' |
    grep -v 'techknowlogick/certmagic-s3' |
    grep -v 'hslatman/caddy-openapi-validator' |
    grep -v 'mohammed90/caddy-ssh' |
    grep -v 'dunglas/vulcain/caddy' |
    grep -v 'dunglas/mercure/caddy' |
    grep -v 'silinternational/certmagic-storage-dynamodb' |
    # sed 's|silinternational/certmagic-storage-dynamodb/v2|silinternational/certmagic-storage-dynamodb/v3|g' |
    grep -v 'firecow/caddy-elastic-encoder' |
    grep -v 'RussellLuo/caddy-ext/flagr' |
    grep -v 'caddyserver/cache-handler' |
    sed 's|caddyserver/nginx-adapter|caddyserver/nginx-adapter@master|g' |
    sed 's|lucaslorentz/caddy-docker-proxy/plugin/v2|lucaslorentz/caddy-docker-proxy/v2|g' |
    sort -u)

shfn=local_build.sh
cat >${shfn} <<EOF
#!/usr/bin/env sh
set -eu
export XCADDY_GO_BUILD_FLAGS="-ldflags '-s -w -extldflags -static'"
export GOPROXY=https://proxy.golang.com.cn,direct
export GOBIN=/mnt/d/devlops/caddy
export XCADDY_SKIP_CLEANUP=1
go install -v github.com/caddyserver/xcaddy/cmd/xcaddy@master

./xcaddy build \\
EOF
for n in $modules; do
    echo "    --with ${n} \\" >>${shfn}
done
echo "    --with github.com/darkweak/souin@v1.6.6" >>${shfn}

psfn=local_build.ps1
cat >${psfn} <<EOF
\$env:XCADDY_GO_BUILD_FLAGS = "-ldflags '-s -w -extldflags -static'"
\$env:GOPROXY = "https://proxy.golang.com.cn,direct"
\$env:GOBIN = Get-Location
\$env:XCADDY_SKIP_CLEANUP = 1
go install -v github.com/caddyserver/xcaddy/cmd/xcaddy@master

./xcaddy build \`
EOF
for n in $modules; do
    echo "    --with ${n} \`" >>${psfn}
done
echo "    --with github.com/darkweak/souin@v1.6.6" >>${psfn}
