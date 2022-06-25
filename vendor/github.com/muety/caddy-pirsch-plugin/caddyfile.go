package caddy_pirsch_plugin

import (
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"regexp"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("pirsch", parseCaddyfile)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	p := new(PirschPlugin)

	for h.Next() {
		// configuration should be in a block
		for h.NextBlock(0) {
			switch h.Val() {
			case "client_id":
				var clientId string
				if !h.AllArgs(&clientId) {
					return nil, h.ArgErr()
				}
				p.ClientId = clientId
			case "client_secret":
				var clientSecret string
				if !h.AllArgs(&clientSecret) {
					return nil, h.ArgErr()
				}
				p.ClientSecret = clientSecret
			case "host_name":
				var hostName string
				if !h.AllArgs(&hostName) {
					return nil, h.ArgErr()
				}
				p.HostName = hostName
			case "base_url":
				var baseUrl string
				if !h.AllArgs(&baseUrl) {
					return nil, h.ArgErr()
				}
				urlRegex := regexp.MustCompile(`https?://(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
				if !urlRegex.MatchString(baseUrl) {
					return nil, h.Errf("'%s' is not a valid url", baseUrl)
				}
				p.BaseURL = baseUrl
			default:
				return nil, h.Errf("unrecognized option '%s'", h.Val())
			}
		}
	}

	if p.ClientId == "" || p.ClientSecret == "" || p.HostName == "" {
		return nil, fmt.Errorf("missing configuration option (one of 'client_id', 'client_secret', 'host_name')")
	}

	return p, nil
}
