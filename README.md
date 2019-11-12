# all plugins caddy 1

```shell
https://github.com/shuxs/caddy-builder.git
cd caddy-builder
go build
./caddy --plugins
```

## main.go

```go

package main

import (
    "github.com/caddyserver/caddy/caddy/caddymain"

    //Plugins

    //Caddyfile Loaders
    _ "github.com/lucaslorentz/caddy-docker-proxy/plugin" //docker

    //DNS Providers
    _ "github.com/caddyserver/dnsproviders/acmedns"
    _ "github.com/caddyserver/dnsproviders/alidns"
    _ "github.com/caddyserver/dnsproviders/auroradns"
    _ "github.com/caddyserver/dnsproviders/azure"
    _ "github.com/caddyserver/dnsproviders/cloudflare"
    _ "github.com/caddyserver/dnsproviders/cloudxns"
    _ "github.com/caddyserver/dnsproviders/conoha"
    _ "github.com/caddyserver/dnsproviders/digitalocean"
    _ "github.com/caddyserver/dnsproviders/dnsimple"
    _ "github.com/caddyserver/dnsproviders/dnsmadeeasy"
    _ "github.com/caddyserver/dnsproviders/dnspod"
    _ "github.com/caddyserver/dnsproviders/duckdns"
    _ "github.com/caddyserver/dnsproviders/dyn"
    _ "github.com/caddyserver/dnsproviders/exoscale"
    _ "github.com/caddyserver/dnsproviders/fastdns"
    _ "github.com/caddyserver/dnsproviders/gandi"
    _ "github.com/caddyserver/dnsproviders/gandiv5"
    _ "github.com/caddyserver/dnsproviders/generic"
    _ "github.com/caddyserver/dnsproviders/glesys"
    _ "github.com/caddyserver/dnsproviders/godaddy"
    _ "github.com/caddyserver/dnsproviders/googlecloud"
    _ "github.com/caddyserver/dnsproviders/httpreq"
    _ "github.com/caddyserver/dnsproviders/inwx"
    _ "github.com/caddyserver/dnsproviders/lightsail"
    _ "github.com/caddyserver/dnsproviders/linode"
    _ "github.com/caddyserver/dnsproviders/linodev4"
    _ "github.com/caddyserver/dnsproviders/namecheap"
    _ "github.com/caddyserver/dnsproviders/namedotcom"
    _ "github.com/caddyserver/dnsproviders/namesilo"
    _ "github.com/caddyserver/dnsproviders/nifcloud"
    _ "github.com/caddyserver/dnsproviders/ns1"
    _ "github.com/caddyserver/dnsproviders/otc"
    _ "github.com/caddyserver/dnsproviders/ovh"
    _ "github.com/caddyserver/dnsproviders/pdns"
    _ "github.com/caddyserver/dnsproviders/rackspace"
    _ "github.com/caddyserver/dnsproviders/rfc2136"
    _ "github.com/caddyserver/dnsproviders/route53"
    _ "github.com/caddyserver/dnsproviders/selectel"
    _ "github.com/caddyserver/dnsproviders/stackpath"
    _ "github.com/caddyserver/dnsproviders/transip"
    _ "github.com/caddyserver/dnsproviders/vscale"
    _ "github.com/caddyserver/dnsproviders/vultr"

    //Directives/Middleware
    _ "github.com/BTBurke/caddy-jwt"                        //http.jwt
    _ "github.com/SchumacherFM/mailout"                     //http.mailout
    _ "github.com/Xumeiquer/nobots"                         //http.nobots
    _ "github.com/aablinov/caddy-geoip"                     //http.geoip
    _ "github.com/abiosoft/caddy-git"                       //http.git
    _ "github.com/caddyserver/forwardproxy"                 //http.forwardproxy
    _ "github.com/captncraig/caddy-realip"                  //http.realip
    _ "github.com/captncraig/cors/caddy"                    //http.cors
    _ "github.com/casbin/caddy-authz"                       //http.authz
    _ "github.com/coopernurse/caddy-awslambda"              //http.awslambda
    _ "github.com/echocat/caddy-filter"                     //http.filter
    _ "github.com/epicagency/caddy-expires"                 //http.expires
    _ "github.com/freman/caddy-reauth"                      //http.reauth
    _ "github.com/hacdias/caddy-minify"                     //http.minify
    _ "github.com/hacdias/caddy-webdav"                     //http.webdav
    _ "github.com/jung-kurt/caddy-cgi"                      //http.cgi
    _ "github.com/jung-kurt/caddy-pubsub"                   //http.pubsub
    _ "github.com/linkonoid/caddy-dyndns"                   //http.dyndns
    _ "github.com/lucaslorentz/caddy-supervisor/httpplugin" //http.supervisor
    _ "github.com/mastercactapus/caddy-proxyprotocol"       //http.proxyprotocol
    _ "github.com/miekg/caddy-prometheus"                   //http.prometheus
    _ "github.com/miquella/caddy-awses"                     //http.awses
    _ "github.com/nicolasazrak/caddy-cache"                 //http.cache
    _ "github.com/payintech/caddy-datadog"                  //http.datadog
    _ "github.com/pieterlouw/caddy-grpc"                    //http.grpc
    _ "github.com/pyed/ipfilter"                            //http.ipfilter

    _ "github.com/restic/caddy"                   //http.restic
    _ "github.com/shuxs/gopkgr"                   //http.gopkgr
    _ "github.com/simia-tech/caddy-locale"        //http.locale
    _ "github.com/tarent/loginsrv/caddy"          //http.login
    _ "github.com/techknowlogick/caddy-s3browser" //http.s3browser
    _ "github.com/xuqingfeng/caddy-rate-limit"    //http.ratelimit
    _ "github.com/zikes/gopkg"                    //http.gopkg

    // _ "go.okkur.org/torproxy" //http.torproxy: registered dev directive

    //More Directives/Middleware
    _ "blitznote.com/src/http.upload"      //http.upload
    _ "github.com/dhaavi/caddy-permission" //http.permission

    //Event Hooks
    _ "github.com/hacdias/caddy-service" //hook.service issus: panic: close of closed channel

    //Server Types
    _ "github.com/lucaslorentz/caddy-supervisor/servertype" //supervisor
    _ "github.com/pieterlouw/caddy-net/caddynet"            //net

    //TLS Clustering
    _ "github.com/pteich/caddy-tlsconsul" //consul
)

func main() {
    // optional: disable telemetry
    caddymain.EnableTelemetry = false
    caddymain.Run()
}

```

## go.mod replacer

```go
module caddy

go 1.13

replace (
    github.com/h2non/gock => gopkg.in/h2non/gock.v1 latest
    github.com/lucaslorentz/caddy-supervisor => github.com/shuxs/caddy-supervisor master
    github.com/miquella/caddy-awses => github.com/leelynne/caddy-awses master
    github.com/restic/caddy => github.com/restic/caddy master
    github.com/simia-tech/caddy-locale => github.com/shuxs/caddy-locale master
    github.com/zikes/gopkg => github.com/fawick/gopkg master
)

```
