# Olaf's Declarative Configuration

## Declarative Configuration

### Overview

Olaf's declarative configuration is inspired by Kong, and the configuration must be written in YAML.

For the core idea of the declarative configuration, see [Kong's Declarative Configuration](https://docs.konghq.com/2.2.x/db-less-and-declarative-config/#what-is-declarative-configuration).

### Entities

While following the same idea as Kong, Olaf's declarative configuration and the entities it contains are different in details.

The top-level entries:

| Entry | Required | Description |
| --- | --- | --- |
| `services` | √ | A list of Services. Similar to Kong's [Service Object](https://docs.konghq.com/2.2.x/admin-api/#service-object). |
| `plugins` | | A list of global Plugins. Default: `[]`. Similar to Kong's [Plugin Object](https://docs.konghq.com/2.2.x/admin-api/#plugin-object). |

The Service entity:

| Attribute | Required | Description |
| --- | --- | --- |
| `name` | | The name of this Service. Default: `"service_<i>"` (`<i>` is the index of this service in the array). |
| `upstream` | √ | The Upstream associated to this Service. Similar to Kong's [Upstream Object](https://docs.konghq.com/gateway-oss/2.2.x/admin-api/#upstream-object). |
| `routes` | √ | A list of Routes associated to this Service. Similar to Kong's [Route Object](https://docs.konghq.com/2.2.x/admin-api/#route-object). |
| `plugins` | | A list of Plugins applied to this Service. Default: `[]`. Similar to Kong's [Plugin Object](https://docs.konghq.com/2.2.x/admin-api/#plugin-object). |

The Upstream entity:

| Attribute | Required | Description |
| --- | --- | --- |
| `backends` | √ | See descriptions of [reverse_proxy.upstreams](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#upstreams). |
| `max_requests` | | See descriptions of [reverse_proxy.lb_policy](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#lb_policy). |
| `dial_timeout` | | The [duration string](https://caddyserver.com/docs/json/apps/http/servers/routes/handle/reverse_proxy/transport/http/dial_timeout/), which indicates how long to wait before timing out trying to connect to this Service. Default: `""` (no timeout). |
| `lb_policy` | | See descriptions of [reverse_proxy.lb_policy](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#lb_policy). |
| `lb_try_duration` | | See descriptions of [reverse_proxy.lb_try_duration](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#lb_try_duration). |
| `lb_try_interval` | | See descriptions of [reverse_proxy.lb_try_interval](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#lb_try_interval). |
| `health_uri` | | See descriptions of [reverse_proxy.health_uri](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#health_uri). |
| `health_port` | | See descriptions of [reverse_proxy.health_port](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#health_port). |
| `health_interval` | | See descriptions of [reverse_proxy.health_interval](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#health_interval). |
| `health_timeout` | | See descriptions of [reverse_proxy.health_timeout](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#health_timeout). |
| `health_status` | | See descriptions of [reverse_proxy.health_status](https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#health_status). |
| `header_up` | | Set, add or remove header fields in a request going upstream to the backend (see [docs](https://caddyserver.com/docs/json/apps/http/servers/routes/handle/reverse_proxy/headers/request/)). Default: `{}` (no header manipulation). |
| `header_down` | | Set, add or remove header fields in a response coming downstream from the backend (see [docs](https://caddyserver.com/docs/json/apps/http/servers/routes/handle/reverse_proxy/headers/response/)). Default: `{}` (no header manipulation). |

The Route entity:

| Attribute | Required | Description |
| --- | --- | --- |
| `name` | | The name of this Route. Default: `"<service_name>_route_<i>"` (`<i>` is the index of this route in the array). |
| `protocol` | | The request [protocol](https://caddyserver.com/docs/caddyfile/matchers#protocol) that matches this Route. Default: `""` (any protocol). |
| `methods` | | A list of [HTTP methods](https://caddyserver.com/docs/caddyfile/matchers#method) that match this Route. Default: `[]` (any HTTP method). |
| `hosts` | | A list of [hosts](https://caddyserver.com/docs/caddyfile/matchers#host) that match this Route. Default: `[]` (any host). |
| `paths` | √ | A list of [URI paths](https://caddyserver.com/docs/caddyfile/matchers#path) that match this Route. A special prefix `~:` means a [regexp path](https://caddyserver.com/docs/caddyfile/matchers#path-regexp). |
| `headers` | | A list of [headers](https://caddyserver.com/docs/caddyfile/matchers#header) that match this Route. Default: `[]` (any header). |
| `strip_prefix` | | The [prefix](https://caddyserver.com/docs/json/apps/http/servers/routes/handle/rewrite/strip_path_prefix/) that needs to be stripped from the request path. Default: `""` (no stripping). |
| `strip_suffix` | | The [suffix](https://caddyserver.com/docs/json/apps/http/servers/routes/handle/rewrite/strip_path_suffix/) that needs to be stripped from the request path. Default: `""` (no stripping). |
| `target_path` | | The final path when the request is proxied to the target service (using `$` as a placeholder for the request path, which may have been stripped). Default: `""` (leave the request path as is, i.e. `"$"`). |
| `add_prefix` | | The prefix that needs to be added to the final path. Default: `""` (no adding). |
| `priority` | | The priority of this Route. Default: `0`. All the services' routes will be matched from highest priority to lowest. |
| `plugins` | | A list of Plugins applied to this Route. Default: `[]`. Similar to Kong's [Plugin Object](https://docs.konghq.com/2.2.x/admin-api/#plugin-object). |
| `response` | | The static response (see `StaticResponse`) for this Route, which indicates that the request will not be proxied to the target service. Default: `{}` (no static response). |

The [StaticResponse](https://caddyserver.com/docs/json/apps/http/servers/routes/handle/static_response/) entity:

| Attribute | Required | Description |
| --- | --- | --- |
| `status_code` | | The HTTP status code to respond with. Default: `200`. |
| `headers` | | The header fields to set on the response. Default: `{}` (no extra header fields). |
| `body` | | The response body to respond with. Default: `""` (no response body). |
| `close` | | Whether to close the client's connection to the server after writing the response. Default: `false`. |

The Plugin entity:

| Attribute | Required | Description |
| --- | --- | --- |
| `disabled` | | Whether this Plugin is disabled. Default: `false`. |
| `name` | | The name of this Plugin. Default: `"plugin_<i>"` for global plugins, `"<service_name>_plugin_<i>"` for service plugins, or `"<route_name>_plugin_<i>"` for route plugins (`<i>` is the index of this plugin in the array). |
| `type` | √ | The type of this Plugin. Available plugin types: `"canary"` (built-in), or `"request_body_var"` (requires the [caddy-ext/requestbodyvar](https://github.com/RussellLuo/caddy-ext/tree/master/requestbodyvar) extension), or `"rate_limit"` (requires the [caddy-ext/ratelimit](https://github.com/RussellLuo/caddy-ext/tree/master/ratelimit) extension). |
| `order_after` | | The order of this Plugin. Default: `""` (the `type` of the previous Plugin, if any, in the Plugin array). |
| `config` | | The configuration of this Plugin. |

The Config of the Canary Plugin:

| Attribute | Required | Description |
| --- | --- | --- |
| `upstream` | √ | The name of the upstream service for this Plugin. |
| `key` | √ | The variable used to differentiate one client from another. Currently supported variables: `"{path.*}"`, `"{query.*}"`, `"{header.*}"`, `"{cookie.*}"` or `"{body.*}"` (requires the [caddy-ext/requestbodyvar](https://github.com/RussellLuo/caddy-ext/tree/master/requestbodyvar) extension). |
| `type` | | The type of key. Default: `""` (string). |
| `whitelist` | √ | The whitelist defined in a [CEL expression](https://caddyserver.com/docs/caddyfile/matchers#expression) (using `$` as a placeholder for the value of key). If the key value is in the whitelist, the corresponding request will be routed to the service specified by `upstream`. |
| `matcher` | | The advanced matcher, which can consist of various [Caddy matchers](https://caddyserver.com/docs/json/apps/http/servers/routes/match/) or your own ones. **NOTE**: `matcher` and (`key`, `type`, `whitelist`) are mutually exclusive. |
| `strip_prefix` | | The [prefix](https://caddyserver.com/docs/json/apps/http/servers/routes/handle/rewrite/strip_path_prefix/) that needs to be stripped from the request path. Default: `""` (no stripping). |
| `strip_suffix` | | The [suffix](https://caddyserver.com/docs/json/apps/http/servers/routes/handle/rewrite/strip_path_suffix/) that needs to be stripped from the request path. Default: `""` (no stripping). |
| `target_path` | | The final path when the request is proxied to the upstream service (using `$` as a placeholder for the request path, which may have been stripped). Default: `""` (leave the request path as is, i.e. `"$"`). |
| `add_prefix` | | The prefix that needs to be added to the final path. Default: `""` (no adding). |

### Example

See [apis.yaml](apis.yaml).


## Embedding Olaf in Caddyfile

### Serving your APIs

```
{
    order olaf last
}

example.com {
    olaf apis.yaml
}
```

### Serving both a website and your APIs

```
{
    order olaf last
}

example.com {
    route /* {
        file_server
    }

    route /api/* {
        uri strip_prefix /api
        olaf apis.yaml
    }
}
```

## Usage

### Build Caddy

Install xcaddy:

```bash
$ go get -u github.com/caddyserver/xcaddy/cmd/xcaddy
```

Build Caddy:

```bash
$ xcaddy build \
    --with github.com/RussellLuo/olaf/caddyconfig/adapter \
    --with github.com/RussellLuo/olaf/caddymodule \
    --with github.com/RussellLuo/caddy-ext/requestbodyvar \
    --with github.com/RussellLuo/caddy-ext/ratelimit
```

### Run Caddy

```bash
$ ./caddy run --config Caddyfile --adapter olaf
```

### Reload Config

Don't run, just test the configuration:

```bash
$ ./caddy adapt --config Caddyfile --adapter olaf --validate > /dev/null
```

Reload the configuration:

```bash
$ ./caddy reload --config Caddyfile --adapter olaf
```
