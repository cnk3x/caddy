package tlsformat

import (
	"crypto/tls"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/logging"
	"go.uber.org/zap/zapcore"
)

// TLSVersionFilter is a Caddy log field filter that replaces the numeric TLS version with the
// string version and optionally adds a prefix.
type TLSVersionFilter struct {
	// Prefix is a constant string that will be added before the replaced version string.
	Prefix string `json:"prefix,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (TLSVersionFilter) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.logging.encoders.filter.tls_version",
		New: func() caddy.Module { return new(TLSVersionFilter) },
	}
}

// Filter replaces the input field containing the numeric TLS version with the string version and
// adds the prefix.
func (f *TLSVersionFilter) Filter(in zapcore.Field) zapcore.Field {
	in.Type = zapcore.StringType
	in.String = f.Prefix
	switch in.Integer {
	case tls.VersionTLS10:
		in.String += "1.0"
	case tls.VersionTLS11:
		in.String += "1.1"
	case tls.VersionTLS12:
		in.String += "1.2"
	case tls.VersionTLS13:
		in.String += "1.3"
	default:
		in.String = strconv.FormatInt(in.Integer, 16)
	}
	return in
}

// UnmarshalCaddyfile sets up the module from Caddyfile tokens.
func (f *TLSVersionFilter) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			f.Prefix = d.Val()
		}
	}
	return nil
}

// TLSCipherFilter is Caddy log field filter that replaces the numeric TLS cipher_suite value with
// the string representation.
type TLSCipherFilter struct{}

// CaddyModule returns the Caddy module information.
func (TLSCipherFilter) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.logging.encoders.filter.tls_cipher",
		New: func() caddy.Module { return new(TLSCipherFilter) },
	}
}

// Filter replaces the input field containing numeric TLS cipher_suite with the corresponding string
// representation.
func (f *TLSCipherFilter) Filter(in zapcore.Field) zapcore.Field {
	in.Type = zapcore.StringType
	in.String = tls.CipherSuiteName(uint16(in.Integer))
	return in
}

// UnmarshalCaddyfile sets up the module from Caddyfile tokens.
func (f *TLSCipherFilter) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

func init() {
	caddy.RegisterModule(TLSVersionFilter{})
	caddy.RegisterModule(TLSCipherFilter{})
}

// Interface guards.
var (
	_ caddy.Module           = (*TLSVersionFilter)(nil)
	_ logging.LogFieldFilter = (*TLSVersionFilter)(nil)
	_ caddyfile.Unmarshaler  = (*TLSVersionFilter)(nil)
	_ caddy.Module           = (*TLSCipherFilter)(nil)
	_ logging.LogFieldFilter = (*TLSCipherFilter)(nil)
	_ caddyfile.Unmarshaler  = (*TLSCipherFilter)(nil)
)
