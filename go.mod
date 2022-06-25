module caddy

go 1.18

require (
	github.com/Elegant996/scgi-transport v0.6.0
	github.com/HeavenVolkoff/caddy-authelia/plugin v0.0.0-20220412230450-9b5a2af7f160
	github.com/RussellLuo/caddy-ext/requestbodyvar v0.1.0
	github.com/RussellLuo/olaf v0.0.0-20220424065813-bebb59a8494e
	github.com/WingLim/caddy-webhook v1.0.8
	github.com/abiosoft/caddy-exec v0.0.0-20210526181020-06d4f7218eb8
	github.com/abiosoft/caddy-hmac v0.0.0-20210522205451-976ca0a419ef
	github.com/abiosoft/caddy-json-parse v0.0.0-20210522205405-c57039f26567
	github.com/abiosoft/caddy-json-schema v0.0.0-20220621031927-c4d6e132f3af
	github.com/abiosoft/caddy-named-routes v0.0.0-20210526091612-80ad81d5e162
	github.com/abiosoft/caddy-yaml v0.0.0-20210522210701-64fbdd07cf02
	github.com/aksdb/caddy-cgi/v2 v2.0.1
	github.com/baldinof/caddy-supervisor v0.6.0
	github.com/caddy-dns/alidns v1.0.23
	github.com/caddy-dns/azure v0.2.0
	github.com/caddy-dns/cloudflare v0.0.0-20210607183747-91cf700356a1
	github.com/caddy-dns/digitalocean v0.0.0-20220527005842-9c71e343246b
	github.com/caddy-dns/dnspod v0.0.4
	github.com/caddy-dns/duckdns v0.3.1
	github.com/caddy-dns/gandi v1.0.2
	github.com/caddy-dns/godaddy v1.0.2
	github.com/caddy-dns/googleclouddns v1.0.3
	github.com/caddy-dns/hetzner v0.0.1
	github.com/caddy-dns/lego-deprecated v0.0.0-20220510003557-83194f3f2958
	github.com/caddy-dns/metaname v0.2.0
	github.com/caddy-dns/netcup v0.1.0
	github.com/caddy-dns/netlify v1.0.1
	github.com/caddy-dns/openstack-designate v0.1.0
	github.com/caddy-dns/route53 v1.1.3
	github.com/caddy-dns/vultr v0.0.0-20211122185502-733392841379
	github.com/caddyserver/caddy/v2 v2.5.1
	github.com/caddyserver/jsonc-adapter v0.0.0-20200325004025-825ee096306c
	github.com/caddyserver/nginx-adapter v0.0.5-0.20220621222301-0f0bba94f5c6
	github.com/caddyserver/ntlm-transport v0.1.1
	github.com/caddyserver/replace-response v0.0.0-20211108214007-d32dc3ffff0c
	github.com/caddyserver/transform-encoder v0.0.0-20220319234440-17b694fb69eb
	github.com/casbin/caddy-authz/v2 v2.0.0
	github.com/chukmunnlee/caddy-openapi v0.7.0
	github.com/circa10a/caddy-geofence v0.5.3
	github.com/cubic3d/caddy-ct v1.0.1
	github.com/cubic3d/caddy-quantity-limiter v1.0.0
	github.com/darkweak/souin v1.6.10
	github.com/darkweak/souin/plugins/caddy v0.0.0-20220622063614-d1b8dd76e1b8
	github.com/gamalan/caddy-tlsredis v0.2.9
	github.com/gbox-proxy/gbox v1.0.6
	github.com/ggicci/caddy-jwt v0.7.1
	github.com/git001/caddyv2-upload v0.0.0-20220608225501-32f8c1dd6c31
	github.com/greenpau/caddy-git v1.0.7
	github.com/greenpau/caddy-security v1.1.14
	github.com/greenpau/caddy-trace v1.1.10
	github.com/hairyhenderson/caddy-teapot-module v0.0.2
	github.com/hslatman/caddy-crowdsec-bouncer v0.2.0
	github.com/imgk/caddy-trojan v0.0.0-20220608002302-81f36a3e396d
	github.com/kirsch33/realip v1.6.1
	github.com/lindenlab/caddy-s3-proxy v0.5.6
	github.com/lolPants/caddy-requestid v1.1.1
	github.com/lucaslorentz/caddy-docker-proxy/v2 v2.7.1
	github.com/mastercactapus/caddy2-proxyprotocol v0.0.2
	github.com/mholt/caddy-dynamicdns v0.0.0-20220312031409-f638ea80fe56
	github.com/mholt/caddy-l4 v0.0.0-20220503192553-2ecee94d269f
	github.com/mholt/caddy-ratelimit v0.0.0-20220428144044-9c011f665e5d
	github.com/mholt/caddy-webdav v0.0.0-20210914165325-f7b67f8ca1e6
	github.com/mpilhlt/caddy-conneg v0.1.4
	github.com/muety/caddy-pirsch-plugin v0.0.0-20220516213216-660c2aee8c2c
	github.com/muety/caddy-remote-host v0.0.0-20211013090634-b21775afa730
	github.com/porech/caddy-maxmind-geolocation v0.0.0-20210828161002-89d86498ab7d
	github.com/pteich/caddy-tlsconsul v1.4.1
	github.com/shift72/caddy-geo-ip v0.6.0
	github.com/sillygod/cdp-cache v0.4.6
	github.com/sjtug/caddy2-filter v0.0.0-20220427014017-828ee1cb05be
	github.com/techknowlogick/caddy-s3browser v1.0.0
	github.com/tosie/caddy-dns-linode v0.0.0-20210701121230-e16993235880
	github.com/ueffel/caddy-basic-auth-filter v1.0.0
	github.com/ueffel/caddy-brotli v1.2.0
	github.com/ueffel/caddy-imagefilter v1.2.0
	github.com/ueffel/caddy-tls-format v1.0.0
	magnax.ca/caddy/gopkg v1.2.0
)

require (
	cloud.google.com/go/compute v1.6.1 // indirect
	filippo.io/edwards25519 v1.0.0-rc.1 // indirect
	github.com/99designs/gqlgen v0.17.2 // indirect
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/Azure/azure-sdk-for-go v58.0.0+incompatible // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.20 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.15 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c // indirect
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/EpicStep/go-simple-geo/v2 v2.0.1 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/NebulousLabs/go-upnp v0.0.0-20181203152547-b32978b8ccbf // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/OpenDNS/vegadns2client v0.0.0-20180418235048-a3fa4a771d87 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/XiaoMi/pegasus-go-client v0.0.0-20210427083443-f3b6b08bc4c2 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/ajeddeloh/go-json v0.0.0-20200220154158-5ae607161559 // indirect
	github.com/ajeddeloh/yaml v0.0.0-20170912190910-6b94386aeefd // indirect
	github.com/akamai/AkamaiOPEN-edgegrid-golang v1.1.1 // indirect
	github.com/alecthomas/chroma v0.10.0 // indirect
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1183 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/antlr/antlr4 v0.0.0-20200503195918-621b933c7a7f // indirect
	github.com/armon/go-metrics v0.4.0 // indirect
	github.com/aryann/difflib v0.0.0-20210328193216-ff5ff6dc229b // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.41.14 // indirect
	github.com/basgys/goxml2json v1.1.0 // indirect
	github.com/beevik/etree v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20220106215444-fb4bf637b56d // indirect
	github.com/bsm/redislock v0.7.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/buraksezer/connpool v0.5.0 // indirect
	github.com/buraksezer/consistent v0.0.0-20191006190839-693edf70fd72 // indirect
	github.com/buraksezer/olric v0.4.5 // indirect
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/caddyserver/caddy v1.0.4 // indirect
	github.com/caddyserver/certmagic v0.16.1 // indirect
	github.com/casbin/casbin/v2 v2.8.6 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cheekybits/genny v1.0.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/circa10a/go-geofence v0.5.0 // indirect
	github.com/cloudflare/cloudflare-go v0.20.0 // indirect
	github.com/coocood/freecache v1.2.1 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/coreos/ignition v0.35.0 // indirect
	github.com/cpu/goacmedns v0.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/crewjam/httperr v0.2.0 // indirect
	github.com/crewjam/saml v0.4.6 // indirect
	github.com/crowdsecurity/crowdsec v1.0.2 // indirect
	github.com/crowdsecurity/go-cs-bouncer v0.0.0-20201130114000-e5b8016e5bf3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deepmap/oapi-codegen v1.6.1 // indirect
	github.com/dgraph-io/badger v1.6.2 // indirect
	github.com/dgraph-io/badger/v2 v2.2007.4 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.2 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/digitalocean/godo v1.41.0 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/dnsimple/dnsimple-go v0.70.1 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.16+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.1-0.20200219035652-afde56e7acac // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0 // indirect
	github.com/eko/gocache/v2 v2.3.0 // indirect
	github.com/elnormous/contenttype v1.0.3 // indirect
	github.com/emersion/go-sasl v0.0.0-20211008083017-0b9dcfb154ac // indirect
	github.com/emersion/go-smtp v0.15.0 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/emvi/null v1.3.1 // indirect
	github.com/exoscale/egoscale v0.67.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/getkin/kin-openapi v0.76.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-acme/lego/v3 v3.7.0 // indirect
	github.com/go-acme/lego/v4 v4.6.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.1 // indirect
	github.com/go-chi/chi v4.1.2+incompatible // indirect
	github.com/go-chi/stampede v0.5.1 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/go-ldap/ldap/v3 v3.4.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.19.16 // indirect
	github.com/go-openapi/errors v0.19.9 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/loads v0.20.0 // indirect
	github.com/go-openapi/runtime v0.19.24 // indirect
	github.com/go-openapi/spec v0.20.0 // indirect
	github.com/go-openapi/strfmt v0.19.11 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-openapi/validate v0.20.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gobwas/httphead v0.0.0-20180130184737-2c6c146eadee // indirect
	github.com/gobwas/pool v0.2.0 // indirect
	github.com/gobwas/ws v1.0.4 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.2.0 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/cel-go v0.7.3 // indirect
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gophercloud/gophercloud v0.16.0 // indirect
	github.com/gophercloud/utils v0.0.0-20210216074907-f6de111f2eae // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/greenpau/go-authcrunch v1.0.35 // indirect
	github.com/greenpau/versioned v1.0.27 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/hairyhenderson/go-which v0.2.0 // indirect
	github.com/hashicorp/consul/api v1.13.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.2.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/hashicorp/golang-lru v0.5.5-0.20200511160909-eb529947af53 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/memberlist v0.3.0 // indirect
	github.com/hashicorp/serf v0.9.8 // indirect
	github.com/hslatman/cidranger v1.0.3-0.20210102151717-b2292da972c3 // indirect
	github.com/hslatman/ipstore v0.0.0-20210131120430-64b55d649887 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/icholy/replace v0.4.0 // indirect
	github.com/iij/doapi v0.0.0-20190504054126-0bbf12d6d7df // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/imgk/memory-go v0.0.0-20220328012817-37cdd311f1a3 // indirect
	github.com/infobloxopen/infoblox-go-client v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.9.0 // indirect
	github.com/jackc/pgx/v4 v4.14.0 // indirect
	github.com/jarcoal/httpmock v1.0.8 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jensneuse/abstractlogger v0.0.4 // indirect
	github.com/jensneuse/byte-template v0.0.0-20200214152254-4f3cf06e5c68 // indirect
	github.com/jensneuse/graphql-go-tools v1.51.0 // indirect
	github.com/jensneuse/pipeline v0.0.0-20200117120358-9fb4de085cd6 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k0kubun/go-ansi v0.0.0-20180517002512-3bf9e2903213 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/kinvolk/container-linux-config-transpiler v0.9.1 // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	github.com/kolo/xmlrpc v0.0.0-20200310150728-e0350524596b // indirect
	github.com/labbsr0x/bindman-dns-webhook v1.0.2 // indirect
	github.com/labbsr0x/goh v1.0.1 // indirect
	github.com/libdns/alidns v1.0.2-x2 // indirect
	github.com/libdns/azure v0.2.0 // indirect
	github.com/libdns/cloudflare v0.1.0 // indirect
	github.com/libdns/digitalocean v0.0.0-20220518195853-a541bc8aa80f // indirect
	github.com/libdns/dnspod v0.0.3 // indirect
	github.com/libdns/duckdns v0.1.1 // indirect
	github.com/libdns/gandi v1.0.2 // indirect
	github.com/libdns/godaddy v0.0.0-20220126161229-bb81e9eae213 // indirect
	github.com/libdns/googleclouddns v1.0.2 // indirect
	github.com/libdns/hetzner v0.0.1 // indirect
	github.com/libdns/libdns v0.2.1 // indirect
	github.com/libdns/metaname v0.3.0 // indirect
	github.com/libdns/netcup v0.1.0 // indirect
	github.com/libdns/netlify v1.0.1 // indirect
	github.com/libdns/openstack-designate v0.1.0 // indirect
	github.com/libdns/route53 v1.1.2 // indirect
	github.com/libdns/vultr v0.0.0-20211122184636-cd4cb5c12e51 // indirect
	github.com/linode/linodego v1.0.0 // indirect
	github.com/liquidweb/go-lwApi v0.0.5 // indirect
	github.com/liquidweb/liquidweb-cli v0.6.9 // indirect
	github.com/liquidweb/liquidweb-go v1.6.3 // indirect
	github.com/lucas-clemente/quic-go v0.26.0 // indirect
	github.com/mailgun/groupcache/v2 v2.2.1 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/marten-seemann/qpack v0.2.1 // indirect
	github.com/marten-seemann/qtls-go1-16 v0.1.5 // indirect
	github.com/marten-seemann/qtls-go1-17 v0.1.1 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.1 // indirect
	github.com/mastercactapus/proxyprotocol v0.0.3 // indirect
	github.com/matoous/go-nanoid/v2 v2.0.0 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/maxmind/geoipupdate/v4 v4.9.0 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mholt/acmez v1.0.2 // indirect
	github.com/mholt/certmagic v0.8.3 // indirect
	github.com/micromdm/scep/v2 v2.1.0 // indirect
	github.com/miekg/dns v1.1.49 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/minio/minio-go/v6 v6.0.48 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/muhammadmuzzammil1998/jsonc v0.0.0-20200303171503-1e787b591db7 // indirect
	github.com/namedotcom/go v0.0.0-20180403034216-08470befbe04 // indirect
	github.com/nats-io/jwt/v2 v2.2.0 // indirect
	github.com/nats-io/nats.go v1.11.1-0.20210623165838-4b75fc59ae30 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/netlify/open-api/v2 v2.9.0 // indirect
	github.com/nrdcg/auroradns v1.0.1 // indirect
	github.com/nrdcg/desec v0.6.0 // indirect
	github.com/nrdcg/dnspod-go v0.4.0 // indirect
	github.com/nrdcg/freemyip v0.2.0 // indirect
	github.com/nrdcg/goinwx v0.8.1 // indirect
	github.com/nrdcg/namesilo v0.2.1 // indirect
	github.com/nrdcg/porkbun v0.1.1 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/open-policy-agent/opa v0.41.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/oracle/oci-go-sdk v24.3.0+incompatible // indirect
	github.com/oschwald/maxminddb-golang v1.8.0 // indirect
	github.com/ovh/go-ovh v1.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pegasus-kv/thrift v0.13.0 // indirect
	github.com/pirsch-analytics/pirsch-go-sdk v1.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/pquerna/otp v1.3.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.34.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/pteich/errors v1.0.1 // indirect
	github.com/qri-io/jsonpointer v0.1.1 // indirect
	github.com/qri-io/jsonschema v0.2.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/xid v1.2.1 // indirect
	github.com/russellhaering/goxmldsig v1.1.1 // indirect
	github.com/russross/blackfriday v1.5.2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sacloud/libsacloud v1.36.2 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.7.0.20210127161313-bd30bebeac4f // indirect
	github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/slackhq/nebula v1.5.2 // indirect
	github.com/smallstep/certificates v0.19.0 // indirect
	github.com/smallstep/cli v0.18.0 // indirect
	github.com/smallstep/nosql v0.4.0 // indirect
	github.com/smallstep/truststore v0.11.0 // indirect
	github.com/smartystreets/go-aws-auth v0.0.0-20180515143844-0c1422d1fdb9 // indirect
	github.com/softlayer/softlayer-go v1.0.3 // indirect
	github.com/softlayer/xmlrpc v0.0.0-20200409220501-5f089df7cb7e // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/tailscale/tscert v0.0.0-20220125204807-4509a5fbaf74 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.287 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.287 // indirect
	github.com/tidwall/gjson v1.11.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/sjson v1.0.4 // indirect
	github.com/tosie/libdns-linode v0.1.0 // indirect
	github.com/transip/gotransip/v6 v6.6.1 // indirect
	github.com/urfave/cli v1.22.5 // indirect
	github.com/vektah/gqlparser/v2 v2.4.4 // indirect
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50 // indirect
	github.com/vinyldns/go-vinyldns v0.9.16 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vultr/govultr/v2 v2.11.0 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xujiajun/mmap-go v1.0.1 // indirect
	github.com/xujiajun/nutsdb v0.8.0 // indirect
	github.com/xujiajun/utils v0.0.0-20190123093513-8bf096c4f53b // indirect
	github.com/yashtewari/glob-intersection v0.1.0 // indirect
	github.com/yuin/goldmark v1.4.8 // indirect
	github.com/yuin/goldmark-highlighting v0.0.0-20220208100518-594be1970594 // indirect
	gitlab.com/NebulousLabs/fastrand v0.0.0-20181126182046-603482d69e40 // indirect
	gitlab.com/NebulousLabs/go-upnp v0.0.0-20181011194642-3a71999ed0d3 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.etcd.io/etcd/api/v3 v3.5.4 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.4 // indirect
	go.etcd.io/etcd/client/v3 v3.5.4 // indirect
	go.mongodb.org/mongo-driver v1.4.4 // indirect
	go.mozilla.org/pkcs7 v0.0.0-20210826202110-33d05740a352 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	go.opentelemetry.io/proto/otlp v0.16.0 // indirect
	go.step.sm/cli-utils v0.7.0 // indirect
	go.step.sm/crypto v0.16.1 // indirect
	go.step.sm/linkedca v0.15.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277 // indirect
	go.uber.org/zap v1.21.0 // indirect
	go4.org v0.0.0-20201209231011-d4a079459e60 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/exp v0.0.0-20220428152302-39d4317da171 // indirect
	golang.org/x/image v0.0.0-20220413100746-70e8d0d3baa9 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/net v0.0.0-20220607020251-c690dde0001d // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	golang.org/x/text v0.3.8-0.20211004125949-5bd84dd9b33b // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df // indirect
	google.golang.org/api v0.80.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220524023933-508584e28198 // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/ns1/ns1-go.v2 v2.6.2 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	howett.net/plist v1.0.0 // indirect
	k8s.io/api v0.22.5 // indirect
	k8s.io/apimachinery v0.23.5 // indirect
	k8s.io/client-go v0.22.5 // indirect
	k8s.io/klog/v2 v2.30.0 // indirect
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65 // indirect
	k8s.io/utils v0.0.0-20211116205334-6203023598ed // indirect
	nhooyr.io/websocket v1.8.7 // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
