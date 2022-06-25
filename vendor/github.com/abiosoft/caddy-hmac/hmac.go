package hmac

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// Interface guards
var (
	_ caddy.Provisioner           = (*HMAC)(nil)
	_ caddy.Validator             = (*HMAC)(nil)
	_ caddyhttp.MiddlewareHandler = (*HMAC)(nil)
	_ caddyfile.Unmarshaler       = (*HMAC)(nil)
)

func init() {
	caddy.RegisterModule(HMAC{})
	httpcaddyfile.RegisterHandlerDirective("hmac", parseCaddyfile)
}

// HMAC implements an HTTP handler that
// validates request body with hmac.
type HMAC struct {
	Algorithm string `json:"algorithm,omitempty"`
	Secret    string `json:"secret,omitempty"`
	Name      string `json:"name,omitempty"`

	hasher func() hash.Hash
}

// CaddyModule returns the Caddy module information.
func (HMAC) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.hmac",
		New: func() caddy.Module { return new(HMAC) },
	}
}

// Provision implements caddy.Provisioner.
func (m *HMAC) Provision(ctx caddy.Context) error {
	switch hashAlgorithm(m.Algorithm) {
	case algSha1:
		m.hasher = sha1.New
	case algSha256:
		m.hasher = sha256.New
	case algMd5:
		m.hasher = md5.New
	}
	return nil
}

// Validate implements caddy.Validator.
func (m HMAC) Validate() error {
	if !hashAlgorithm(m.Algorithm).valid() {
		return fmt.Errorf("unsupported hash type '%s'", m.Algorithm)
	}
	if m.hasher == nil {
		// this will never happen
		return fmt.Errorf("hasher is null")
	}
	return nil
}

func (m HMAC) replacerKey() string {
	if m.Name != "" {
		return fmt.Sprintf("hmac.%s.signature", m.Name)
	}
	return "hmac.signature"
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
//    hmac [<name>] <algorithm> <secret>
//
func (m *HMAC) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		args := d.RemainingArgs()

		switch len(args) {
		case 2:
			m.Algorithm, m.Secret = args[0], args[1]
		case 3:
			m.Name, m.Algorithm, m.Secret = args[0], args[1], args[2]
		default:
			return d.Err("unexpected number of arguments")
		}
	}

	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m HMAC
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}
