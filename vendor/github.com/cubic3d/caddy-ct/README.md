# Container Linux Config Transpiler for Caddy
The `caddy-ct` module for Caddy allows to transpile YAML based configuration into a JSON `ignition` to be used with
[Flatcar](https://www.kinvolk.io/flatcar-container-linux/) or
[Fedora CoreOS](https://getfedora.org/en/coreos?stream=stable).

It targets to replace [Matchbox](https://matchbox.psdn.io/) with an open and flexible approach of templating,
matching and providing metadata using Caddy's `templates` for static configurations (no API for terraform).

## Configuration
```
ct [<matcher>] {
  strict
  mime <MIMEType> [<MIMEType...>]
  platform <platformName>
}
```
All options are optional:
- `matcher` according to [matcher](https://caddyserver.com/docs/caddyfile/concepts#matchers)
- `strict` fail on non critical errors (default: false)
- `mime` only transpile specific MIME types (default: all)
- `platform` only for
[dynamic data](https://kinvolk.io/docs/flatcar-container-linux/latest/provisioning/config-transpiler/dynamic-data/)
must be one of
[those](https://github.com/kinvolk/container-linux-config-transpiler/blob/flatcar-master/config/platform/platform.go#L17)
  (default: none)

The module is unordered by default and needs to be ordered using the global option
```
{
    order ct before templates
}
```
or be used inside a [route](https://caddyserver.com/docs/caddyfile/directives/route) block.

## Example Caddyfile
The following example allows files with specific MIME types to be templated and transpiled after to `ignition` config.
```
{
  order ct before templates
}

:8080 {
  log
  respond / "OK"
  file_server
  templates {
    mime text/html text/plain application/json application/x-yaml
    between [[ ]]
  }
  ct {
    strict
    mime application/x-yaml
  }
}
```

## Building
See [xcaddy](https://github.com/caddyserver/xcaddy), short:
```
xcaddy build \
  --with github.com/cubic3d/caddy-ct
```
or download from https://caddyserver.com/download by selecting this module.