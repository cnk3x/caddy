package caddyrl

import (
	"fmt"
	"sync"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// RateLimit describes an HTTP rate limit zone.
type RateLimit struct {
	// Request matchers, which defines the class of requests that are in the RL zone.
	MatcherSetsRaw caddyhttp.RawMatcherSets `json:"match,omitempty" caddy:"namespace=http.matchers"`

	// The key which uniquely differentiates rate limits within this zone. It could
	// be a static string (no placeholders), resulting in one and only one rate limiter
	// for the whole zone. Or, placeholders could be used to dynamically allocate
	// rate limiters. For example, a key of "foo" will create exactly one rate limiter
	// for all clients. But a key of "{http.request.remote.host}" will create one rate
	// limiter for each different client IP address.
	Key string `json:"key,omitempty"`

	// Number of events allowed within the window.
	MaxEvents int `json:"max_events,omitempty"`

	// Duration of the sliding window.
	Window caddy.Duration `json:"window,omitempty"`

	matcherSets caddyhttp.MatcherSets

	limiters *sync.Map
}

func (rl *RateLimit) provision(ctx caddy.Context, name string) error {
	if rl.Window <= 0 {
		return fmt.Errorf("window must be greater than zero")
	}
	if rl.MaxEvents < 0 {
		return fmt.Errorf("max_events must be at least zero")
	}

	if len(rl.MatcherSetsRaw) > 0 {
		matcherSets, err := ctx.LoadModule(rl, "MatcherSetsRaw")
		if err != nil {
			return err
		}
		err = rl.matcherSets.FromInterface(matcherSets)
		if err != nil {
			return err
		}
	}

	// ensure rate limiter state endures across config changes
	rl.limiters = new(sync.Map)
	if val, loaded := rateLimits.LoadOrStore(name, rl.limiters); loaded {
		rl.limiters = val.(*sync.Map)
	}

	// update existing rate limiters with new settings
	rl.limiters.Range(func(key, value interface{}) bool {
		limiter := value.(*ringBufferRateLimiter)
		limiter.SetMaxEvents(rl.MaxEvents)
		limiter.SetWindow(time.Duration(rl.Window))
		return true
	})

	return nil
}
