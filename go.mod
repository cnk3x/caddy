module caddy

go 1.13

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.43.0
	github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.15
	github.com/lucaslorentz/caddy-supervisor => github.com/shuxs/caddy-supervisor v0.1.1-0.20190806004419-0c50fbdb9f42
	github.com/miquella/caddy-awses => github.com/whalehub/caddy-awses v0.0.0-20190709150835-656ad4af91bb
	github.com/payintech/caddy-datadog => github.com/whalehub/caddy-datadog v0.0.0-20190709094746-9b0422a27b18
	github.com/simia-tech/caddy-locale => github.com/shuxs/caddy-locale v0.0.0-20190806002554-a01b17fc5d14
	github.com/zikes/gopkg => github.com/fawick/gopkg v1.0.2-0.20190706112402-6c2f2452db80
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190731235908-ec7cb31e5a56
	golang.org/x/image => github.com/golang/image v0.0.0-20190802002840-cff245a6509b
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190719004257-d2bd2a29d028
	golang.org/x/net => github.com/golang/net v0.0.0-20190724013045-ca1201d0de80
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190804053845-51ab0e2deafa
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190805222050-c5a2fd39b72a
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.7.0
	google.golang.org/appengine => github.com/golang/appengine v1.6.1
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190801165951-fa694d86fc64
	google.golang.org/grpc => github.com/grpc/grpc-go v1.22.1
)

require (
	blitznote.com/src/http.upload v1.9.9
	github.com/BTBurke/caddy-jwt v3.7.1+incompatible
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.20.0+incompatible // indirect
	github.com/SchumacherFM/mailout v1.3.0
	github.com/Xumeiquer/nobots v0.1.1
	github.com/aablinov/caddy-geoip v0.0.0-20190710083220-6705babb56ca
	github.com/abiosoft/caddy-git v0.0.0-20190703061829-f8cc2f20c9e7
	github.com/caddyserver/caddy v1.0.1
	github.com/caddyserver/dnsproviders v0.3.0
	github.com/caddyserver/forwardproxy v0.0.0-20190707023537-05540a763b63
	github.com/captncraig/caddy-realip v0.0.0-20190710144553-6df827e22ab8
	github.com/captncraig/cors v0.0.0-20190703115713-e80254a89df1
	github.com/casbin/caddy-authz v1.0.2
	github.com/casbin/casbin v1.9.1 // indirect
	github.com/coopernurse/caddy-awslambda v1.0.0
	github.com/dhaavi/caddy-permission v0.6.0
	github.com/echocat/caddy-filter v0.14.0
	github.com/epicagency/caddy-expires v1.1.1
	github.com/freman/caddy-reauth v0.0.0-20190703021030-0863eef919a2
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/gorilla/handlers v1.4.2 // indirect
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/hacdias/caddy-minify v1.0.2
	github.com/hacdias/caddy-service v1.0.1
	github.com/hacdias/caddy-webdav v1.0.1
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/improbable-eng/grpc-web v0.10.0 // indirect
	github.com/jung-kurt/caddy-cgi v1.11.4
	github.com/jung-kurt/caddy-pubsub v0.5.6
	github.com/linkonoid/caddy-dyndns v0.0.0-20190718171622-2414d6236b0f
	github.com/lucaslorentz/caddy-docker-proxy v0.3.1-0.20190709175318-d258b595b82a
	github.com/lucaslorentz/caddy-supervisor v0.0.0-00010101000000-000000000000
	github.com/mastercactapus/caddy-proxyprotocol v0.0.3
	github.com/miekg/caddy-prometheus v0.0.0-20190709133612-1fe4cb19becd
	github.com/miolini/datacounter v0.0.0-20190724021726-aa48df3a02c1 // indirect
	github.com/miquella/caddy-awses v0.0.0-00010101000000-000000000000
	github.com/mmcloughlin/geohash v0.9.0 // indirect
	github.com/nicolasazrak/caddy-cache v0.3.4
	github.com/payintech/caddy-datadog v0.0.0-00010101000000-000000000000
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pieterlouw/caddy-grpc v0.1.0
	github.com/pieterlouw/caddy-net v0.1.5
	github.com/pteich/caddy-tlsconsul v0.0.0-20190709201921-ebc221e392e1
	github.com/pyed/ipfilter v1.1.4
	github.com/quasoft/memstore v0.0.0-20180925164028-84a050167438 // indirect
	github.com/restic/caddy v0.2.2-0.20190709151628-d755491f9a25
	github.com/restic/rest-server v0.9.8 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/simia-tech/caddy-locale v0.0.0-00010101000000-000000000000
	github.com/steambap/captcha v1.3.0 // indirect
	github.com/tarent/loginsrv v1.3.1
	github.com/techknowlogick/caddy-s3browser v0.0.0-20190710044735-655ba503d3ea
	github.com/tinylib/msgp v1.1.0 // indirect
	github.com/xuqingfeng/caddy-rate-limit v1.6.4
	github.com/zikes/gopkg v0.0.0-00010101000000-000000000000
	go.okkur.org/torproxy v0.2.0
	goji.io v2.0.2+incompatible // indirect
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.16.1 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
)
