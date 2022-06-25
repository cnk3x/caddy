# caddy-tls-format

This packages contains two log field filters to log TLS version and cipher suites in a more readable
form.

## Installation

```sh
xcaddy build --with github.com/ueffel/caddy-tls-format
```

## Usage

See [caddy log filter documentation](https://caddyserver.com/docs/caddyfile/directives/log#filter).
There will be two new filters to use:

### tls_version

```caddy-d
<field> tls_version [prefix]
```

* **field** Probably the only sensible field to use here is: `request>tls>version`
* **prefix** string that is added before the TLS version string.

### tls_cipher

```caddy-d
<field> tls_cipher
```

* **field** Probably the only sensible field to use here is: `request>tls>cipher_suite`

## Example configuration

The following example configuration uses the [Formatted Log
Encoder](https://github.com/caddyserver/format-encoder)

```caddy-d
format filter {
    wrap formatted "\"{request>method} {request>uri} {request>proto}\" {request>tls>version}/{request>tls>cipher_suite}"
    fields {
        request>tls>version tls_version TLSv
        request>tls>cipher_suite tls_cipher
    }
}
```

Log output (with and without HTTPS):

```plain
"GET / HTTP/2.0" TLSv1.3/TLS_AES_128_GCM_SHA256
"GET / HTTP/1.1" -/-
```

> For reference the configuration and output without filters:
>
> ```caddy-d
> format formatted "\"{request>method} {request>uri} {request>proto}\" {request>tls>version}/{request>tls>cipher_suite}"
> ```
>
> Log output:
>
> ```plain
> "GET / HTTP/2.0" 772/4865
> "GET / HTTP/1.1" -/-
> ```
