package legodeprecated

import (
	"context"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/mholt/acmez"
	"github.com/mholt/acmez/acme"
)

func init() {
	caddy.RegisterModule(LegoDeprecated{})
}

// LegoDeprecated is a shim module that allows any and all of the
// DNS providers in go-acme/lego to be used with Caddy. They must
// be configured via environment variables, they do not support
// cancellation in the case of frequent config changes.
//
// Even though this module is in the dns.providers namespace, it
// is only a special case for solving ACME challenges, intended to
// replace the modules that used to be in the now-defunct tls.dns
// namespace. Using it in other places of the Caddy config will
// result in errors.
//
// This module will eventually go away in favor of the modules that
// make use of the libdns APIs: https://github.com/libdns
type LegoDeprecated struct {
	ProviderName string `json:"provider_name,omitempty"`

	prov challenge.Provider
}

// CaddyModule returns the Caddy module information.
func (LegoDeprecated) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.lego_deprecated",
		New: func() caddy.Module { return new(LegoDeprecated) },
	}
}

// Provision initializes the underlying DNS provider.
func (ld *LegoDeprecated) Provision(ctx caddy.Context) error {
	prov, err := dns.NewDNSChallengeProviderByName(ld.ProviderName)
	if err != nil {
		return err
	}
	ld.prov = prov
	return nil
}

// Present wraps the go-acme/lego/v4/challenge.Provider interface
// with the certmagic.ACMEDNSProvider interface. Normally, DNS providers
// in the caddy-dns repositories would implement the libdns interfaces
// (https://github.com/libdns/libdns) instead, but this module is a
// special case to give time for more DNS providers to be ported over
// to the libdns interfaces from the deprecated lego interface.
func (ld LegoDeprecated) Present(_ context.Context, challenge acme.Challenge) error {
	return ld.prov.Present(challenge.Identifier.Value, challenge.Token, challenge.KeyAuthorization)
}

// Wait waits just a few seconds before proceeding. We don't have a clean way of
// doing true propagation polling from this layer of abstraction, unfortunately.
// If there is a way to do that with lego v4, then I don't know what it is.
func (LegoDeprecated) Wait(ctx context.Context, challenge acme.Challenge) error {
	select {
	case <-time.After(10 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// CleanUp wraps the go-acme/lego/v4/challenge.Provider interface
// with the acmez.Solver interface. Normally, DNS providers
// in the caddy-dns repositories would implement the libdns interfaces
// (https://github.com/libdns/libdns) instead, but this module is a
// special case to give time for more DNS providers to be ported over
// to the libdns interfaces from the deprecated lego interface.
func (ld LegoDeprecated) CleanUp(_ context.Context, challenge acme.Challenge) error {
	return ld.prov.CleanUp(challenge.Identifier.Value, challenge.Token, challenge.KeyAuthorization)
}

// Interface guard
var _ acmez.Solver = (*LegoDeprecated)(nil)
