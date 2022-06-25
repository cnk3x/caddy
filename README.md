# 全插件编译 caddy 离线源码版本

官网要编译一个大全插件版本根本不可能，依赖地狱太痛苦了，所以搞了一个项目来编译 caddy, 尽可能的集成官网列出来的插件

```cmd
set CGO_ENABLED=0
go build -mod vendor -trimpath -ldflags '-s -w -extldflags -static' -v
```

```powershell
$env:CGO_ENABLED = 0
go build -mod vendor -trimpath -ldflags '-s -w -extldflags -static' -v
```

```bash
CGO_ENABLED=0 go build -mod vendor -trimpath -ldflags '-s -w -extldflags -static' -v
```

```text
下列插件多年没更新，已经不适合最新的代码了，暂时过滤掉

firecow/caddy-forward-auth  //forward_auth: 在2.5.1中已经内置了
github.com/francislavoie/caddy-hcl
github.com/techknowlogick/certmagic-s3
github.com/hslatman/caddy-openapi-validator
github.com/mohammed90/caddy-ssh
github.com/dunglas/vulcain/caddy
github.com/dunglas/mercure/caddy
github.com/RussellLuo/caddy-ext/flagr
github.com/caddyserver/cache-handler
```
