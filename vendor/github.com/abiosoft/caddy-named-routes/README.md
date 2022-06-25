# caddy-named-routes

named routes is a Caddy v2 module for creating reusable named http routes.

## Installation

```
xcaddy build v2.4.1 \
    --with github.com/abiosoft/caddy-named-routes
```

## Usage

Named route is only currently limited to Caddy API and not available for the Caddyfile.

```jsonc
{
  "app": {
    // define named routes
    "named_routes": {
      // the name for the route
      "<route name>": [
        // ... list of http routes
      ]
    },

    // use the routes in the http handler
    "http": {
      "servers": {
        "srv0": {
          "routes": [
            {
              "handle": {
                "handler": "named_route",
                "name": "<route name>"
              }
            }
          ]
        }
      }
    }
  }
}
```

## Why?

Caddy's API is simple enough. What is this unnecessary complexity?

### The Problem

Caddy's configuration API is flexible and gives more control than the Caddyfile.
However, composing http routes can get easily messy and unclear if you have many of them.

Consider the following YAML config snippet for the configuration of two reverse proxies and a static file server. The average use case cannot get any simpler.

**Note:** even though YAML is used here, Caddy's native configuration language is JSON. You need [config adapters](https://caddyserver.com/docs/config-adapters#known-config-adapters) for other than JSON.

```yaml
http:
  servers:
    default:
      routes:
        - match:
            - host: [localhost] # other global conditions
          handle:
            - handler: subroute
              routes:
                # reverse proxy to API
                - match:
                    - path: [/api/*]
                  handle:
                    # strip prefix before reverse proxy
                    - handler: rewrite
                      strip_path_prefix: /api
                    - handler: subroute
                      routes:
                        # API v2
                        - match:
                            - header:
                                X-API-Version: [v2]
                          handle:
                            - handler: reverse_proxy
                              upstreams:
                                - dial: localhost:8080
                        # API legacy
                        - handle:
                            - handler: reverse_proxy
                              upstreams:
                                - dial: localhost:8888
            # blog
            - handler: file_server
              root: /home/blog/static
            # fallback handler
            - handler: static_response
              status_code: "404"
```

This is still relatively readable thanks to the comments and the simple use case. Now imagine a more complex structure and even worse; imagine the complex structure as a JSON config while bearing in mind pure JSON does not support comments.

### The Alternative

The following is a composition of same YAML config with named routes.

```yaml
http:
  servers:
    default:
      routes:
        - match:
            - host: [localhost] # other global conditions
          handle:
            # api
            - handler: named_route
              name: api
            # blog
            - handler: file_server
              root: /home/blog/static
            # fallback handler
            - handler: static_response
              status_code: "404"

named_routes:
  api:
    - match:
        - path: [/api/*]
      handle:
        # strip prefix before reverse proxy
        - handler: rewrite
          strip_path_prefix: /api
        # attempt v2
        - handler: named_route
          name: api.v2
        # otherwise fall back to legacy
        - handler: named_route
          name: api.legacy

  api.v2:
    - match:
        header:
          X-API-Version: [v2]
      handle:
        - handler: reverse_proxy
          upstreams:
            - dial: localhost:8080

  api.legacy:
    - handle:
        - handler: reverse_proxy
          upstreams:
            - dial: localhost:8888
```

Even though the latter config is 10 lines longer, it looks cleaner and more readable.

## License

Apache 2
