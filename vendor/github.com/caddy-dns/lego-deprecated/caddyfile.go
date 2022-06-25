package legodeprecated

import "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//     lego_deprecated <provider>
//
func (ld *LegoDeprecated) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if !d.Next() { // consume module name
		return d.Err("expected tokens")
	}
	if !d.NextArg() { // get provider name
		return d.ArgErr()
	}
	ld.ProviderName = d.Val()
	if d.NextArg() {
		return d.ArgErr()
	}
	return nil
}
