package caddy

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/cnk3x/libdns-westcn"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *westcn.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.westcn",
		New: func() caddy.Module { return &Provider{new(westcn.Provider)} },
	}
}

// Provision Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *Provider) Provision(caddy.Context) error {
	p.Provider.Username = caddy.NewReplacer().ReplaceAll(p.Provider.Username, "")
	p.Provider.Password = caddy.NewReplacer().ReplaceAll(p.Provider.Password, "")
	p.Provider.Endpoint = caddy.NewReplacer().ReplaceAll(p.Provider.Endpoint, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// westcn <username> <password> {
//     username <username>
//     password <password>
//	   endpoint <endpoint>
// }
//
// Expansion of placeholders in the API token is left to the JSON config caddy.Provisioner (above).
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.Username = d.Val()
		}
		if d.NextArg() {
			p.Provider.Password = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "username":
				if p.Provider.Username != "" {
					return d.Err("username already set")
				}
				p.Provider.Username = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "password":
				if p.Provider.Password != "" {
					return d.Err("password already set")
				}
				p.Provider.Password = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "endpoint":
				if p.Provider.Endpoint != "" {
					return d.Err("endpoint already set")
				}
				p.Provider.Endpoint = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized sub directive '%s'", d.Val())
			}
		}
	}
	if p.Provider.Username == "" || p.Provider.Password == "" {
		return d.Err("missing username or password")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
