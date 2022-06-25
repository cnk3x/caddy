package namedroutes

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// Interface guards
var (
	_ caddy.App = (*App)(nil)
)

const appModule = "named_routes"

func init() {
	caddy.RegisterModule(App{})
}

// App handles a list route of named routes.
// With named routes, it is easier to compose
// cleaner configuration files with less nesting.
// The routes can be reused in http route handler as
// `named_route`
type App map[string]caddyhttp.RouteList

func (a *App) get(name string) caddyhttp.RouteList {
	if a == nil {
		return nil
	}

	if routes, ok := (*a)[name]; ok {
		return routes
	}
	return nil
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  appModule,
		New: func() caddy.Module { return new(App) },
	}
}

// Start implements caddy.App.
func (a *App) Start() error {
	// nothing to do
	return nil
}

// Stop implements caddy.App.
func (a *App) Stop() error {
	// the app simply consists of middleware handlers.
	// we can rely on the http server to shutdown properly.
	return nil
}
