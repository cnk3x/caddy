package quantitylimiter

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func init() {
	caddy.RegisterModule(QuantityLimiter{})
	httpcaddyfile.RegisterHandlerDirective("quantity_limiter", parseCaddyfile)
}

// QuantityLimiter limits the number of successful requests for a token and allows the counter to be reset.
type QuantityLimiter struct {
	logger *zap.Logger

	// Contains tokens and their number of allowed requests left.
	counter map[string]uint64

	// Parameter used to set a token.
	paramSet string
	// Parameter used to get a limited resource.
	paramGet string

	// Prefix to be used for GET parameters for set and get tokens.
	ParameterNamePrefix string `json:"parameterNamePrefix,omitempty"`
	// Number of successful requests that can be made using a token.
	Quantity uint64 `json:"quantity,omitempty"`
}

func (QuantityLimiter) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.quantity_limiter",
		New: func() caddy.Module {
			return QuantityLimiter{}
		},
	}
}

func (ql QuantityLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	if r.URL.Query().Has(ql.paramSet) {
		ql.counter[r.URL.Query().Get(ql.paramSet)] = ql.Quantity
		w.WriteHeader(http.StatusAccepted)
		return nil
	}

	if r.URL.Query().Has(ql.paramGet) {
		if ql.counter[r.URL.Query().Get(ql.paramGet)] == 0 {
			delete(ql.counter, ql.paramGet)
			w.WriteHeader(http.StatusNotFound)
			return nil
		}
		ql.counter[r.URL.Query().Get(ql.paramGet)] -= 1
		r.Header.Del(ql.paramGet)
	}

	return next.ServeHTTP(w, r)
}

func (ql *QuantityLimiter) Provision(ctx caddy.Context) error {
	ql.logger = ctx.Logger(ql)

	ql.counter = make(map[string]uint64)

	if ql.ParameterNamePrefix == "" {
		ql.ParameterNamePrefix = "ql_"
	}

	ql.paramSet = ql.ParameterNamePrefix + "set"
	ql.paramGet = ql.ParameterNamePrefix + "get"

	if ql.Quantity == 0 {
		ql.Quantity = 1
	}

	return nil
}

func (ql *QuantityLimiter) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "parameterNamePrefix":
				if !d.Args(&ql.ParameterNamePrefix) {
					return d.ArgErr()
				}
			case "quantity":
				var quantity string
				if !d.Args(&quantity) {
					return d.ArgErr()
				}
				var err error
				ql.Quantity, err = strconv.ParseUint(quantity, 10, 32)
				if err != nil {
					return d.Err(err.Error())
				}
			default:
				return fmt.Errorf("unknown subdirective: %q", d.Val())
			}
		}
	}
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var ql QuantityLimiter
	err := ql.UnmarshalCaddyfile(h.Dispenser)
	return ql, err
}

// Interface guards
var (
	_ caddy.Provisioner     = (*QuantityLimiter)(nil)
	_ caddyfile.Unmarshaler = (*QuantityLimiter)(nil)
)
