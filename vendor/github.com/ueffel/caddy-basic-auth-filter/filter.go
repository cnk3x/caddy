package basic_auth_filter

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/logging"
	"go.uber.org/zap/zapcore"
)

// BasicAuthFilter is a Caddy log field filter that replaces the a base64 encoded authorization
// header with just the user name.
type BasicAuthFilter struct{}

// CaddyModule returns the Caddy module information.
func (BasicAuthFilter) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.logging.encoders.filter.basic_auth_user",
		New: func() caddy.Module { return new(BasicAuthFilter) },
	}
}

// Filter extracts the user name from the field, if it is a basic authorization, and returns it.
func (f *BasicAuthFilter) Filter(in zapcore.Field) zapcore.Field {
	authHeader, ok := in.Interface.(caddyhttp.LoggableStringArray)
	if !ok {
		return in
	}
	fakeReq := &http.Request{Header: http.Header{"Authorization": authHeader}}
	userID, _, ok := fakeReq.BasicAuth()
	in.Type = zapcore.StringType
	if ok {
		in.String = userID
	} else {
		in.String = ""
	}
	return in
}

// UnmarshalCaddyfile sets up the module from Caddyfile tokens.
func (*BasicAuthFilter) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

func init() {
	caddy.RegisterModule(BasicAuthFilter{})
}

// Interface guards.
var (
	_ logging.LogFieldFilter = (*BasicAuthFilter)(nil)
	_ caddyfile.Unmarshaler  = (*BasicAuthFilter)(nil)
)
