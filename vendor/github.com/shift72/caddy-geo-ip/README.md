# GeoIP

Provides middleware for resolving a users IP address against the Maxmind Geo IP Database.

Manages Downloading and Refreshing the Maxmind Database via https://github.com/maxmind/geoipupdate

## Examples

```

{
  http_port     8080
  https_port    8443
  order geo_ip first
}

localhost:8080 {

  geo_ip {
    db_path GeoLite2-Country.mmdb
    trust_header X-Real-IP
  }

  respond / 200 {
    body "Hello from {geoip.country_code}"
  }
}

```

## Configuration

`account` - this is the Maxmind Account, required for all download operations. Can be supplied as an ENV Variable `GEO_ACCOUNT_ID`.

`db_path` - is the path to load the database from. The filename is used to determine the edition of the file to download
     Valid values tested with are GeoIP2-Country | GeoLite2-Country

`trust_header` - this is used to determine the header to load the users ip address from, if empty it will use the requests `RemoteAddr`

`api_key` - this is a Maxmind API Key. If blank no attempt will be made to download the database.  Can be supplied as an ENV Variable `GEO_API_KEY`.

`download_frequency` - this is how often to download the database from the Maxmind Server (requires an APIKey)

`reload_frequency` - this is how often to check for updated versions of the database on disk. This can be used when an external process is responsible for downloading the database. If the database is being managed via the `api_key` and `download_frequency` then there is no need
to specify the `reload_frequency`

`override_country_code` - this setting allows for the country code to be configured for local development. It will attempt to use the database if there is one, and use the value if no record could be found.

## Release Process

```
git tag v0.4.0
git push origin v0.4.0
```

## Builds on the good work by

https://github.com/porech/caddy-maxmind-geolocation



