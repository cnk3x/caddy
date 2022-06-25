package namedroutes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
)

func init() {
	caddy.RegisterModule(Middleware{})
}

// Middleware implements an HTTP handler that reuses
// a configured named route.
type Middleware struct {
	Name string `json:"name,omitempty"`

	routes caddyhttp.RouteList
	log    *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.named_route",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	m.log = ctx.Logger(m).Named(m.Name)

	{
		// check if previously visited to prevent possible cycle
		if ok := ctx.Value(m.ctxKey()); ok != nil {
			return fmt.Errorf("cycle detected for named route '%s'", m.Name)
		}
		// mark as visited
		ctx.Context = context.WithValue(ctx.Context, m.ctxKey(), struct{}{})
	}

	// fetch named route
	appIface, err := ctx.App(appModule)
	if err != nil {
		return fmt.Errorf("getting %s app: %v", appModule, err)
	}
	app := appIface.(*App)

	routes := app.get(m.Name)
	if routes == nil {
		return fmt.Errorf("no named route '%s' found", m.Name)
	}

	if err := routes.Provision(ctx); err != nil {
		return err
	}

	m.routes = routes
	return nil
}

// Validate implements caddy.Validator.
func (m *Middleware) Validate() error {
	if m.routes == nil {
		return fmt.Errorf("no named route '%s' found", m.Name)
	}
	return nil
}

func (m Middleware) ctxKey() interface{} {
	return struct{ Name string }{Name: appModule + "." + m.Name}
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	routes := m.routes.Compile(next)
	if routes != nil {
		return routes.ServeHTTP(w, r)
	}

	// this never happens
	return next.ServeHTTP(w, r)
}
