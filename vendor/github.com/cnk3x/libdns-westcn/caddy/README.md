# West.cn 西部数码 dns provider for caddy

## Caddy module name
```
dns.providers.westcn
```

## Config examples

To use this module for the ACME DNS challenge, configure the ACME issuer in your Caddy JSON like so:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "westcn",
        "username": "{env.WESTCN_USERNAME}",
        "password": "{env.WESTCN_PASSWORD}",
        "endpoint": "{env.WESTCN_ENDPOINT}"
      }
    }
  }
}
```

or with the Caddyfile

```
tls {
	dns westcn {env.WESTCN_USERNAME} {env.WESTCN_PASSWORD}
}
```