[![Go](https://github.com/muety/caddy-pirsch-plugin/workflows/Go/badge.svg)](https://github.com/muety/caddy-pirsch-plugin/actions)
![Coding Time](https://img.shields.io/endpoint?url=https://wakapi.dev/api/compat/shields/v1/n1try/interval:any/project:caddy-pirsch-plugin&color=blue&label=coding%20time)

# caddy-pirsch-plugin

A Caddy v2 plugin to track requests in [Pirsch Analytics](https://pirsch.io).

## Usage
```
pirsch [<matcher>] {
    client_id <pirsch-client-id>
    client_secret <pirsch-client-secret>
    host_name <pirsch-host-name>
    base_url <alternative-api-url>
}
```

You can obtain these parameters from the [Settings](https://dashboard.pirsch.io/settings) section of your Pirsch dashboard.

Because this directive does not come standard with Caddy, you need to [put the directive in order](https://caddyserver.com/docs/caddyfile/options). The correct place is up to you, but usually putting it near the end works if no other terminal directives match the same requests. It's common to pair a Pirsch handler with a `file_server`, so ordering it just before is often a good choice:

```
{
	order pirsch before file_server
}
```

Alternatively, you may use `route` to order it the way you want. For example:

```
localhost
root * /srv
route {
	pirsch * {
		[...]
	}
	file_server
}
```

### Example
Track all requests to HTML pages in Pirsch. You might want to extend the matcher regexp to also include `/` or, alternatively, match everything but assets (like `.css`, `.js`, ...) since usually you wouldn't want to track those.

```
{
    order pirsch before file_server
}

http://localhost:8080 {
    @html path_regexp .*\.html$

    pirsch @html {
        client_id cCfoZttXzRH5AyOpiu97wqXH3j5lYXcg
        client_secret olshVxS73jWQFhXJE86DdoR4McPBh02OendvyLtajX2EA3aasfywb3q3uZio9tDL
        host_name mysite.example.org
    }

    file_server
}
```

## License
Apache 2.0