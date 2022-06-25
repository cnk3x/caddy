package azure

import (
	"github.com/libdns/azure"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *azure.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "dns.providers.azure",
		New: func() caddy.Module {
			return &Provider{new(azure.Provider)}
		},
	}
}

// Provision implements the Provisioner interface to initialize the Azure client
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.TenantId = repl.ReplaceAll(p.Provider.TenantId, "")
	p.Provider.ClientId = repl.ReplaceAll(p.Provider.ClientId, "")
	p.Provider.ClientSecret = repl.ReplaceAll(p.Provider.ClientSecret, "")
	p.Provider.SubscriptionId = repl.ReplaceAll(p.Provider.SubscriptionId, "")
	p.Provider.ResourceGroupName = repl.ReplaceAll(p.Provider.ResourceGroupName, "")

	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// azure {
//     tenant_id <string>
//     client_id <string>
//     client_secret <string>
//     subscription_id <string>
//     resource_group_name <string>
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "tenant_id":
				if d.NextArg() {
					p.Provider.TenantId = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "client_id":
				if d.NextArg() {
					p.Provider.ClientId = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "client_secret":
				if d.NextArg() {
					p.Provider.ClientSecret = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "subscription_id":
				if d.NextArg() {
					p.Provider.SubscriptionId = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "resource_group_name":
				if d.NextArg() {
					p.Provider.ResourceGroupName = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}

	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
