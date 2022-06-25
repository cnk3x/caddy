# Request Quantity Limiter for Caddy
The `caddy-quantity-limiter` module for Caddy allows to limit the number of requests for a specified token on resources
and only allow requests after the token has been set.

This pattern makes it possible to enable requests from an external event and automatically block again after a
configurable amount of requests.

## Configuration
```
quantity_limiter [<matcher>] {
  parameterNamePrefix <prefix>
  quantity <quantity>
}
```
All options are optional:
- `matcher` according to [matcher](https://caddyserver.com/docs/caddyfile/concepts#matchers)
- `parameterNamePrefix` prefix for token parameter (default: ql_)
- `quantity` number of allowed requests for a token after set (default: 1)

The module is unordered by default and needs to be ordered using the global option
```
{
    order quantity_limiter before file_server
}
```
or be used inside a [route](https://caddyserver.com/docs/caddyfile/directives/route) block.

## Example Caddyfile
```
{
  order quantity_limiter before file_server
}

:8080 {
  log
  file_server
  quantity_limiter
}
```

## Usage
Without defining any of the modules parameters (default: `ql_set` and `ql_get`) there is no special behaviour of
the module - it will pass through all requests.

### Setting
To set a token's counter to `quantity` a *GET* request can be issued on any resource matched by the module
(all by default). For example `/does/not/need/to/exist?ql_set=token`. This will set the counter `token` to the number
of the `quantity` option and return status `202`.

### Getting
To limit the requests the *GET* parameter `ql_get` (default) is used. For example `/existing/file.ext?ql_get=token`.
This will request the file `/existing/file.ext` but will only be returned, if the `token` has prior been set
and requests are left on the counter (initialized with `quantity`). After the amount of remaining requests is 0, the
token is deleted and further requests will return the status code `404`.

## Building
See [xcaddy](https://github.com/caddyserver/xcaddy), short:
```
xcaddy build \
  --with github.com/cubic3d/caddy-quantity-limiter
```
or download from https://caddyserver.com/download by selecting this module.