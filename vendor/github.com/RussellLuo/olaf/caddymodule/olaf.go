package caddymodule

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Olaf{})
	httpcaddyfile.RegisterHandlerDirective("olaf", parseCaddyfile)
}

// Olaf implements a handler that embeds Olaf's declarative configuration, which
// will be expanded later by a config adapter named `olaf`.
type Olaf struct {
	// The config source type.
	Type string `json:"type,omitempty"`

	// The path to the config.
	//
	//    Type: TypeFile => Path: filename
	//    Type: TypeHTTP => Path: url
	Path string `json:"path,omitempty"`

	// Maximum time allowed for a complete connection and request. This
	// option is useful only if Type is TypeHTTP.
	Timeout caddy.Duration `json:"timeout,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (Olaf) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.olaf",
		New: func() caddy.Module { return new(Olaf) },
	}
}

// Validate implements caddy.Validator.
func (o *Olaf) Validate() error {
	if o.Type == "" {
		return fmt.Errorf("empty type")
	}
	if o.Path == "" {
		return fmt.Errorf("empty path")
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (o *Olaf) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler. Syntax:
//
//    olaf <path>
//
func (o *Olaf) UnmarshalCaddyfile(d *caddyfile.Dispenser) (err error) {
	if !d.Next() || !d.NextArg() {
		return d.ArgErr()
	}
	path := d.Val()

	if strings.HasPrefix(path, "http://") {
		o.Type = TypeHTTP
		o.Path = path
		return nil
	}

	o.Type = TypeFile

	if filepath.IsAbs(path) {
		o.Path = path
		return nil
	}

	// Make the path relative to the current Caddyfile rather than the
	// current working directory.
	absFile, err := filepath.Abs(d.File())
	if err != nil {
		return fmt.Errorf("failed to get absolute path of file: %s: %v", d.File(), err)
	}
	o.Path = filepath.Join(filepath.Dir(absFile), path)

	return nil
}

// parseCaddyfile sets up a handler for olaf from Caddyfile tokens.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	o := new(Olaf)
	if err := o.UnmarshalCaddyfile(h.Dispenser); err != nil {
		return nil, err
	}
	return o, nil
}

const (
	TypeFile = "file"
	TypeHTTP = "http"
)

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*Olaf)(nil)
	_ caddyfile.Unmarshaler       = (*Olaf)(nil)
)
