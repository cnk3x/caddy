/*
 * Copyright (c) 2020 Andreas Schneider
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package cgi

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(CGI{})
	httpcaddyfile.RegisterHandlerDirective("cgi", parseCaddyfile)
}

// CGI implements a CGI handler that executes binary files following the
// CGI protocol, passing parameters via environment variables and evaluating
// the response as the HTTP response.
type CGI struct {
	// Name of executable script or binary
	Executable string `json:"executable"`
	// Working directory (default, current Caddy working directory)
	WorkingDirectory string `json:"workingDirectory,omitempty"`
	// The script path of the uri.
	ScriptName string `json:"scriptName,omitempty"`
	// Arguments to submit to executable
	Args []string `json:"args,omitempty"`
	// Environment key value pairs (key=value) for this particular app
	Envs []string `json:"envs,omitempty"`
	// Environment keys to pass through for all apps
	PassEnvs []string `json:"passEnvs,omitempty"`
	// True to pass all environment variables to CGI executable
	PassAll bool `json:"passAllEnvs,omitempty"`
	// True to return inspection page rather than call CGI executable
	Inspect bool `json:"inspect,omitempty"`
}

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*CGI)(nil)
	_ caddyfile.Unmarshaler       = (*CGI)(nil)
)

func (c CGI) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.cgi",
		New: func() caddy.Module { return &CGI{} },
	}
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (c *CGI) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	// Consume 'em all. Matchers should be used to differentiate multiple instantiations.
	// If they are not used, we simply combine them first-to-last.
	for d.Next() {
		args := d.RemainingArgs()
		if len(args) < 1 {
			return fmt.Errorf("an executable needs to be specified")
		}
		c.Executable = args[0]
		c.Args = args[1:]

		for d.NextBlock(0) {
			switch d.Val() {
			case "dir":
				if !d.Args(&c.WorkingDirectory) {
					return d.ArgErr()
				}
			case "script_name":
				if !d.Args(&c.ScriptName) {
					return d.ArgErr()
				}
			case "env":
				c.Envs = d.RemainingArgs()
				if len(c.Envs) == 0 {
					return d.ArgErr()
				}
			case "pass_env":
				c.PassEnvs = d.RemainingArgs()
				if len(c.PassEnvs) == 0 {
					return d.ArgErr()
				}
			case "pass_all_env":
				c.PassAll = true
			case "inspect":
				c.Inspect = true
			default:
				return fmt.Errorf("unknown subdirective: %q", d.Val())
			}
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var c CGI
	err := c.UnmarshalCaddyfile(h.Dispenser)
	return c, err
}
