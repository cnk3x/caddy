# all plugins caddy 1

```shell
https://github.com/shuxs/caddy.git
cd caddy
go build
./caddy --plugins
```

```go
// 这四个玩意编译不过
// _ "github.com/leelynne/caddy-awses"    //http.awses
// _ "github.com/payintech/caddy-datadog" //http.datadog
// _ "github.com/restic/caddy"            //http.restic
// _ "go.okkur.org/gomods"                //http.gomods
```
