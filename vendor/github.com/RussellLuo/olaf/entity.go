package olaf

import (
	"errors"
)

var (
	ErrServiceExists   = errors.New("service already exists")
	ErrServiceNotFound = errors.New("service not found")

	ErrRouteExists   = errors.New("route already exists")
	ErrRouteNotFound = errors.New("route not found")

	ErrPluginExists   = errors.New("plugin already exists")
	ErrPluginNotFound = errors.New("plugin not found")

	ErrUpstreamNotFound = errors.New("upstream not found")

	ErrMethodNotImplemented = errors.New("method not implemented")
)

const (
	PluginTypeCanary = "canary"
)

type Service struct {
	Name     string    `json:"name" yaml:"name"`
	Upstream *Upstream `json:"upstream" yaml:"upstream"`
}

type Upstream struct {
	Backends []*Backend `json:"backends" yaml:"backends"`

	HTTP *TransportHTTP `json:"http" yaml:"http"`

	LoadBalancing      *LoadBalancing      `json:"lb" yaml:"lb"`
	ActiveHealthChecks *ActiveHealthChecks `json:"active_hc" yaml:"active_hc"`

	HeaderUp   *HeaderOps `json:"header_up" yaml:"header_up"`
	HeaderDown *HeaderOps `json:"header_down" yaml:"header_down"`
}

type Backend struct {
	Dial        string `json:"dial" yaml:"dial"`
	MaxRequests int    `json:"max_requests" yaml:"max_requests"`
}

type TransportHTTP struct {
	DialTimeout string `json:"dial_timeout" yaml:"dial_timeout"`
}

type LoadBalancing struct {
	Policy      string `json:"policy" yaml:"policy"`
	TryDuration string `json:"try_duration" yaml:"try_duration"`
	Interval    string `json:"interval" yaml:"interval"`
}

type ActiveHealthChecks struct {
	URI        string `json:"uri" yaml:"uri"`
	Port       int    `json:"port" yaml:"port"`
	Interval   string `json:"interval" yaml:"interval"`
	Timeout    string `json:"timeout" yaml:"timeout"`
	StatusCode int    `json:"status_code" yaml:"status_code"`
}

// Header manipulations.
type HeaderOps struct {
	// Add new header fields or overwrite existing ones.
	Set map[string][]string `json:"set" yaml:"set"`
	// Add new header fields.
	Add map[string][]string `json:"add" yaml:"add"`
	// Remove header fields.
	Delete []string `json:"delete" yaml:"delete"`
}

// Matching rules for a route.
type Matcher struct {
	Protocol string              `json:"protocol" yaml:"protocol"`
	Methods  []string            `json:"methods" yaml:"methods"`
	Hosts    []string            `json:"hosts" yaml:"hosts"`
	Paths    []string            `json:"paths" yaml:"paths"`
	Headers  map[string][]string `json:"headers" yaml:"headers"`
}

// URI manipulations for a route.
type URI struct {
	StripPrefix string `json:"strip_prefix" yaml:"strip_prefix" mapstructure:"strip_prefix"`
	StripSuffix string `json:"strip_suffix" yaml:"strip_suffix" mapstructure:"strip_suffix"`
	TargetPath  string `json:"target_path" yaml:"target_path" mapstructure:"target_path"`
	// TODO: Deprecate AddPrefix
	AddPrefix string `json:"add_prefix" yaml:"add_prefix" mapstructure:"add_prefix"`
}

type StaticResponse struct {
	StatusCode int                 `json:"status_code" yaml:"status_code"`
	Headers    map[string][]string `json:"headers" yaml:"headers"`
	Body       string              `json:"body" yaml:"body"`
	Close      bool                `json:"close" yaml:"close"`
}

type Route struct {
	ServiceName string `json:"service_name" yaml:"service_name"`

	// Route name must be unique.
	Name     string `json:"name" yaml:"name"`
	Matcher  `yaml:",inline"`
	URI      `yaml:",inline"`
	Response *StaticResponse `json:"response" yaml:"response"`

	// Routes will be matched from highest priority to lowest.
	Priority float64 `json:"priority" yaml:"priority"`
}

type Plugin struct {
	Disabled bool `json:"disabled" yaml:"disabled"`

	Name       string                 `json:"name" yaml:"name"`
	Type       string                 `json:"type" yaml:"type"`
	OrderAfter string                 `json:"order_after" yaml:"order_after"`
	Config     map[string]interface{} `json:"config" yaml:"config"`

	RouteName   string `json:"route_name" yaml:"route_name"`
	ServiceName string `json:"service_name" yaml:"service_name"`
}

type PluginCanaryConfig struct {
	UpstreamServiceName string `json:"upstream" yaml:"upstream" mapstructure:"upstream"`

	KeyName   string `json:"key" yaml:"key" mapstructure:"key"`
	KeyType   string `json:"type" yaml:"type" mapstructure:"type"`
	Whitelist string `json:"whitelist" yaml:"whitelist" mapstructure:"whitelist"`

	// The advanced matcher.
	// See https://caddyserver.com/docs/json/apps/http/servers/routes/match/
	Matcher map[string]interface{} `json:"matcher" yaml:"matcher" mapstructure:"matcher"`

	URI `yaml:",inline" mapstructure:",squash"`
}

type Data struct {
	Version  string              `json:"version" yaml:"version"`
	Services map[string]*Service `json:"services" yaml:"services"`
	Routes   map[string]*Route   `json:"routes" yaml:"routes"`
	Plugins  map[string]*Plugin  `json:"plugins" yaml:"plugins"`
}
