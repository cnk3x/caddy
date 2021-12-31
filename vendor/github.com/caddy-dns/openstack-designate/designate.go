package openstack

import (
	designate "github.com/libdns/openstack-designate"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *designate.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "dns.providers.openstack-designate",
		New: func() caddy.Module {
			return &Provider{new(designate.Provider)}
		},
	}
}

// Provision implements the Provisioner interface to initialize the OpenStack Designate client
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.AuthOpenStack.RegionName = repl.ReplaceAll(p.Provider.AuthOpenStack.RegionName, "")
	p.Provider.AuthOpenStack.TenantID = repl.ReplaceAll(p.Provider.AuthOpenStack.TenantID, "")
	p.Provider.AuthOpenStack.IdentityApiVersion = repl.ReplaceAll(p.Provider.AuthOpenStack.IdentityApiVersion, "")
	p.Provider.AuthOpenStack.Password = repl.ReplaceAll(p.Provider.AuthOpenStack.Password, "")
	p.Provider.AuthOpenStack.AuthURL = repl.ReplaceAll(p.Provider.AuthOpenStack.AuthURL, "")
	p.Provider.AuthOpenStack.Username = repl.ReplaceAll(p.Provider.AuthOpenStack.Username, "")
	p.Provider.AuthOpenStack.TenantName = repl.ReplaceAll(p.Provider.AuthOpenStack.TenantName, "")
	p.Provider.AuthOpenStack.EndpointType = repl.ReplaceAll(p.Provider.AuthOpenStack.EndpointType, "")

	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// openstack-designate {
//     region_name <string>
//     tenant_id <string>
//     identity_api_version <string>
//     password <string>
//     username <string>
//     tenant_name <string>
//     endpoint_type <string>
//     auth_url <string>
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "region_name":
				if d.NextArg() {
					p.Provider.AuthOpenStack.RegionName = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "tenant_id":
				if d.NextArg() {
					p.Provider.AuthOpenStack.TenantID = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "identity_api_version":
				if d.NextArg() {
					p.Provider.AuthOpenStack.IdentityApiVersion = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "password":
				if d.NextArg() {
					p.Provider.AuthOpenStack.Password = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "username":
				if d.NextArg() {
					p.Provider.AuthOpenStack.Username = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "tenant_name":
				if d.NextArg() {
					p.Provider.AuthOpenStack.TenantName = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "endpoint_type":
				if d.NextArg() {
					p.Provider.AuthOpenStack.EndpointType = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "auth_url":
				if d.NextArg() {
					p.Provider.AuthOpenStack.AuthURL = d.Val()
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
