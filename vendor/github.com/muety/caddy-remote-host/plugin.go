package caddy_remote_host

// heavily inspired by Caddy's remote_ip matcher (https://github.com/caddyserver/caddy/blob/cbb045a121464527d85cce1b56250480b0515f9a/modules/caddyhttp/matchers.go#L123)

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var hostRegex *regexp.Regexp
var cacheKey string = "hosts"

func init() {
	caddy.RegisterModule(MatchRemoteHost{})
}

// MatchRemoteHost matches based on the remote IP of the
// connection. A host name can be specified, whose A and AAAA
// DNS records will be resolved to a corresponding IP for matching.
//
// Note that IPs can sometimes be spoofed, so do not rely
// on this as a replacement for actual authentication.
type MatchRemoteHost struct {
	// Host names, whose corresponding IPs to match against
	Hosts []string `json:"hosts,omitempty"`

	// If true, prefer the first IP in the request's X-Forwarded-For
	// header, if present, rather than the immediate peer's IP, as
	// the reference IP against which to match. Note that it is easy
	// to spoof request headers. Default: false
	Forwarded bool `json:"forwarded,omitempty"`

	// By default, DNS responses are cached for 60 seconds, regardless
	// of the DNS record's TTL. Set nocache to true to disable this
	// behavior and never use caching. Default: false
	NoCache bool `json:"nocache,omitempty"`

	logger *zap.Logger
	cache  *cache.Cache
}

// CaddyModule returns the Caddy module information.
func (MatchRemoteHost) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.matchers.remote_host",
		New: func() caddy.Module { return new(MatchRemoteHost) },
	}
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *MatchRemoteHost) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextArg() {
			if d.Val() == "forwarded" {
				if len(m.Hosts) > 0 {
					return d.Err("if used, 'forwarded' must appear before 'hosts' argument")
				}
				m.Forwarded = true
				continue
			}
			if d.Val() == "nocache" {
				if len(m.Hosts) > 0 {
					return d.Err("if used, 'nocache' must appear before 'hosts' argument")
				}
				m.NoCache = true
				continue
			}
			m.Hosts = append(m.Hosts, d.Val())
		}
		if d.NextBlock(0) {
			return d.Err("malformed remote_host matcher: blocks are not supported")
		}
	}
	return nil
}

// Provision implements caddy.Provisioner.
func (m *MatchRemoteHost) Provision(ctx caddy.Context) (err error) {
	m.logger = ctx.Logger(m)
	m.cache = cache.New(1*time.Minute, 2*time.Minute)
	hostRegex, err = regexp.Compile(`^((([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))$`)
	return err
}

// Validate implements caddy.Validator.
func (m *MatchRemoteHost) Validate() error {
	for _, h := range m.Hosts {
		if matched := hostRegex.MatchString(h); !matched {
			return fmt.Errorf("'%s' is not a valid host name", h)
		}
	}
	return nil
}

// Match returns true if r matches m.
func (m *MatchRemoteHost) Match(r *http.Request) bool {
	clientIP, err := m.getClientIP(r)
	if err != nil {
		m.logger.Error("getting client IP", zap.Error(err))
		return false
	}

	allowedIPs, err := m.resolveIPs()
	if err != nil {
		m.logger.Error("resolving DNS", zap.Error(err))
		return false
	}

	for _, ip := range allowedIPs {
		if ip.Equal(clientIP) {
			return true
		}
	}

	return false
}

func (m *MatchRemoteHost) getClientIP(r *http.Request) (net.IP, error) {
	remote := r.RemoteAddr
	if m.Forwarded {
		if fwdFor := r.Header.Get("X-Forwarded-For"); fwdFor != "" {
			remote = strings.TrimSpace(strings.Split(fwdFor, ",")[0])
		}
	}
	ipStr, _, err := net.SplitHostPort(remote)
	if err != nil {
		ipStr = remote
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid client IP address: %s", ipStr)
	}
	return ip, nil
}

func (m *MatchRemoteHost) resolveIPs() ([]net.IP, error) {
	if result, ok := m.cache.Get(cacheKey); ok && !m.NoCache {
		return result.([]net.IP), nil
	}

	allIPs := make([]net.IP, 0)

	for _, h := range m.Hosts {
		ips, err := net.LookupIP(h)
		if err != nil {
			return nil, err
		}
		allIPs = append(allIPs, ips...)
	}

	m.cache.SetDefault(cacheKey, allIPs)

	return allIPs, nil
}

// Interface guards
var (
	_ caddy.Provisioner        = (*MatchRemoteHost)(nil)
	_ caddy.Validator          = (*MatchRemoteHost)(nil)
	_ caddyhttp.RequestMatcher = (*MatchRemoteHost)(nil)
	_ caddyfile.Unmarshaler    = (*MatchRemoteHost)(nil)
)
