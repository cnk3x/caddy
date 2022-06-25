package metaname

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	metaname "github.com/libdns/metaname"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *metaname.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.metaname",
		New: func() caddy.Module { return &Provider{new(metaname.Provider)} },
	}
}

// TODO: This is just an example. Useful to allow env variable placeholders; update accordingly.
// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	replacer := caddy.NewReplacer()
	p.Provider.APIKey = replacer.ReplaceAll(p.Provider.APIKey, "")
	p.Provider.AccountReference = replacer.ReplaceAll(p.Provider.AccountReference, "")
	p.Provider.Endpoint = replacer.ReplaceAll(p.Provider.Endpoint, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// providername [<account_reference>] [<api_key>] {
//     account_reference <account_reference>
//     api_key <api_token>
//     endpoint <endpoint>
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.AccountReference = d.Val()
		}
		if d.NextArg() {
			p.Provider.APIKey = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_key":
				if p.Provider.APIKey != "" {
					return d.Err("API key already set")
				}
				p.Provider.APIKey = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "account_reference":
				if p.Provider.AccountReference != "" {
					return d.Err("Account reference already set")
				}
				p.Provider.AccountReference = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "endpoint":
				p.Provider.Endpoint = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.AccountReference == "" || p.Provider.APIKey == "" {
		return d.Err("missing API key or account reference")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
