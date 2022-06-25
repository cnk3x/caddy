# caddy-basic-auth-filter

This packages contains a log field filter to extract the user from a basic Authorization
HTTP-Header.

## Installation

```sh
xcaddy build --with github.com/ueffel/caddy-basic-auth-filter
```

## Usage

See [caddy log filter documentation](https://caddyserver.com/docs/caddyfile/directives/log#filter).
There will be a new filters to use:

```caddy-d
<field> basic_auth_user
```

* **field** Probably the only sensible field to use here is: `request>headers>Authorization`

## Example configuration

The following example configuration uses the [Formatted Log
Encoder](https://github.com/caddyserver/format-encoder)

```caddy-d
format filter {
    wrap formatted "{request>host} {request>headers>Authorization} [{ts}] \"{request>method} {request>uri} {request>proto}\""
    fields {
        request>headers>Authorization basic_auth_user
    }
}
```

```plain
localhost admin [1620840157.514536] "GET /some/path HTTP/2.0" 
```

> For reference the configuration and output without filters:
>
> ```caddy-d
> format formatted "{request>host} {request>headers>Authorization} [{ts}] \"{request>method} {request>uri} {request>proto}\""
> ```
>
> Log output:
>
> ```plain
> localhost ["Basic YWRtaW46YWRtaW4="] [1638732239.578346] "GET /some/path HTTP/2.0"
> ```
