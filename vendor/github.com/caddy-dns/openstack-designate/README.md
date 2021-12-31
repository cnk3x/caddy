OpenStack Designate DNS module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records in OpenStack Designate DNS zones.

## Caddy module name

```
dns.providers.openstack-designate
```

## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/openstack-designate) for important information about credentials.

## Building

To compile this Caddy module, follow the steps describe at the [Caddy Build from Source](https://github.com/caddyserver/caddy#build-from-source) instructions and import the `github.com/caddy-dns/openstack-designate` plugin

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "openstack-designate",
        "region_name": "{env.OS_REGION_NAME}",
        "tenant_id": "{env.OS_TENANT_ID}",
        "identity_api_version": "{env.OS_IDENTITY_API_VERSION}",
        "password": "{env.OS_PASSWORD}",
        "username": "{env.OS_USERNAME}",
        "tenant_name": "{env.OS_TENANT_NAME}",
        "auth_url": "{env.OS_AUTH_URL}",
        "endpoint_type": "{env.OS_ENDPOINT_TYPE}"
      }
    }
  }
}
```

or with the Caddyfile:

```
tls {
  dns openstack-designate {
    region_name {$OS_REGION_NAME}
    tenant_id {$OS_TENANT_ID}
    identity_api_version {$OS_IDENTITY_API_VERSION}
    password {$OS_PASSWORD}
    username {$OS_USERNAME}
    tenant_name {$OS_TENANT_NAME}
    auth_url {$OS_AUTH_URL}
    endpoint_type {$OS_ENDPOINT_TYPE}
  }
}
```
