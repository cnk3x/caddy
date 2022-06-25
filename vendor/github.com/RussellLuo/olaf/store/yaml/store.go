package yaml

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/RussellLuo/olaf"
	"gopkg.in/yaml.v3"
)

type Store struct {
	filename string

	data *olaf.Data
}

func New(filename string) *Store {
	s := &Store{filename: filename}

	data, err := s.GetConfig(context.Background())
	if err != nil {
		log.Printf("failed to get config: %v\n", err)
		return s
	}

	s.data = data
	return s
}

func (s *Store) GetConfig(ctx context.Context) (*olaf.Data, error) {
	c, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}

	data, err := Parse(c)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Store) CreateService(ctx context.Context, svc *olaf.Service) (err error) {
	return olaf.ErrMethodNotImplemented
}

func (s *Store) ListServices(ctx context.Context) (services []*olaf.Service, err error) {
	for _, svc := range s.data.Services {
		services = append(services, svc)
	}
	return
}

func (s *Store) GetService(ctx context.Context, serviceName, routeName string) (*olaf.Service, error) {
	if routeName != "" {
		r, ok := s.data.Routes[routeName]
		if !ok {
			return nil, olaf.ErrServiceNotFound
		}
		// Overwrite serviceName since it must be empty if routeName is specified.
		serviceName = r.ServiceName
	}

	svc, ok := s.data.Services[serviceName]
	if !ok {
		return nil, olaf.ErrServiceNotFound
	}
	return svc, nil
}

func (s *Store) UpdateService(ctx context.Context, serviceName, routeName string, svc *olaf.Service) (err error) {
	return olaf.ErrMethodNotImplemented
}

func (s *Store) DeleteService(ctx context.Context, serviceName, routeName string) (err error) {
	return olaf.ErrMethodNotImplemented
}

func (s *Store) CreateRoute(ctx context.Context, serviceName string, route *olaf.Route) (err error) {
	return olaf.ErrMethodNotImplemented
}

func (s *Store) ListRoutes(ctx context.Context, serviceName string) (routes []*olaf.Route, err error) {
	for _, r := range s.data.Routes {
		if serviceName != "" {
			if r.ServiceName == serviceName {
				routes = append(routes, r)
			}
		} else {
			routes = append(routes, r)
		}
	}
	return
}

func (s *Store) GetRoute(ctx context.Context, serviceName, routeName string) (route *olaf.Route, err error) {
	route, ok := s.data.Routes[routeName]
	if !ok || (serviceName != "" && route.ServiceName != serviceName) {
		return nil, olaf.ErrRouteNotFound
	}
	return route, nil
}

func (s *Store) UpdateRoute(ctx context.Context, serviceName, routeName string, route *olaf.Route) (err error) {
	return olaf.ErrMethodNotImplemented
}

func (s *Store) DeleteRoute(ctx context.Context, serviceName, routeName string) (err error) {
	return olaf.ErrMethodNotImplemented
}

func (s *Store) CreatePlugin(ctx context.Context, serviceName, routeName string, p *olaf.Plugin) (plugin *olaf.Plugin, err error) {
	return nil, olaf.ErrMethodNotImplemented
}

func (s *Store) ListPlugins(ctx context.Context, serviceName, routeName string) (plugins []*olaf.Plugin, err error) {
	for _, p := range s.data.Plugins {
		switch {
		case serviceName != "":
			if p.ServiceName == serviceName {
				plugins = append(plugins, p)
			}
		case routeName != "":
			if p.RouteName == routeName {
				plugins = append(plugins, p)
			}
		default:
			plugins = append(plugins, p)
		}
	}
	return
}

func (s *Store) GetPlugin(ctx context.Context, serviceName, routeName, pluginName string) (plugin *olaf.Plugin, err error) {
	plugin, ok := s.data.Plugins[pluginName]
	if !ok || (serviceName != "" && plugin.ServiceName != serviceName) || (routeName != "" && plugin.RouteName != routeName) {
		return nil, olaf.ErrPluginNotFound
	}
	return plugin, nil
}

func (s *Store) UpdatePlugin(ctx context.Context, serviceName, routeName, pluginName string, plugin *olaf.Plugin) (err error) {
	return olaf.ErrMethodNotImplemented
}

func (s *Store) DeletePlugin(ctx context.Context, serviceName, routeName, pluginName string) (err error) {
	return olaf.ErrMethodNotImplemented
}

func (s *Store) ListUpstreams(ctx context.Context) (upstreams []*olaf.Upstream, err error) {
	for _, svc := range s.data.Services {
		upstreams = append(upstreams, svc.Upstream)
	}
	return
}

func (s *Store) GetUpstream(ctx context.Context, upstreamName, serviceName string) (upstream *olaf.Upstream, err error) {
	if upstreamName != "" {
		return nil, olaf.ErrMethodNotImplemented
	}

	svc, ok := s.data.Services[serviceName]
	if !ok {
		return nil, olaf.ErrUpstreamNotFound
	}
	return svc.Upstream, nil
}

func (s *Store) UpdateUpstream(ctx context.Context, upstreamName, serviceName string, upstream *olaf.Upstream) (err error) {
	return olaf.ErrMethodNotImplemented
}

// Parse recognizes and parses the YAML content.
func Parse(in []byte) (*olaf.Data, error) {
	c := new(content)
	if err := yaml.Unmarshal(in, c); err != nil {
		return nil, err
	}

	data := &olaf.Data{
		Services: make(map[string]*olaf.Service),
		Routes:   make(map[string]*olaf.Route),
		Plugins:  make(map[string]*olaf.Plugin),
	}

	for i, s := range c.Services { // global services
		if s.Name == "" {
			s.Name = fmt.Sprintf("service_%d", i)
		}

		var u *olaf.Upstream
		if s.Upstream != nil {
			var backends []*olaf.Backend
			for _, url := range s.Upstream.Backends {
				backends = append(backends, &olaf.Backend{
					Dial:        url,
					MaxRequests: s.Upstream.MaxRequests,
				})
			}
			u = &olaf.Upstream{
				Backends:   backends,
				HTTP:       &olaf.TransportHTTP{DialTimeout: s.Upstream.DialTimeout},
				HeaderUp:   s.Upstream.HeaderUp,
				HeaderDown: s.Upstream.HeaderDown,
			}
			if s.Upstream.LBPolicy != "" || s.Upstream.LBTryDuration != "" || s.Upstream.LBTryInterval != "" {
				u.LoadBalancing = &olaf.LoadBalancing{
					Policy:      s.Upstream.LBPolicy,
					TryDuration: s.Upstream.LBTryDuration,
					Interval:    s.Upstream.LBTryInterval,
				}
			}
			if s.Upstream.HealthURI != "" {
				u.ActiveHealthChecks = &olaf.ActiveHealthChecks{
					URI:        s.Upstream.HealthURI,
					Port:       s.Upstream.HealthPort,
					Interval:   s.Upstream.HealthInterval,
					Timeout:    s.Upstream.HealthTimeout,
					StatusCode: s.Upstream.HealthStatus,
				}
			}
		}

		data.Services[s.Name] = &olaf.Service{
			Name:     s.Name,
			Upstream: u,
		}

		for j, r := range s.Routes { // routes associated to a service
			if r.Route.Name == "" {
				r.Route.Name = fmt.Sprintf("%s_route_%d", s.Name, j)
			}
			r.Route.ServiceName = s.Name
			data.Routes[r.Route.Name] = r.Route

			for k, p := range r.Plugins { // plugins applied to a route
				if p.Name == "" {
					p.Name = fmt.Sprintf("%s_plugin_%d", r.Route.Name, k)
				}
				if p.OrderAfter == "" && k > 0 {
					p.OrderAfter = r.Plugins[k-1].Type
				}
				p.ServiceName = s.Name
				p.RouteName = r.Route.Name
				data.Plugins[p.Name] = p
			}
		}

		for j, p := range s.Plugins { // plugins applied to a service
			if p.Name == "" {
				p.Name = fmt.Sprintf("%s_plugin_%d", s.Name, j)
			}
			if p.OrderAfter == "" && j > 0 {
				p.OrderAfter = s.Plugins[j-1].Type
			}
			p.ServiceName = s.Name
			data.Plugins[p.Name] = p
		}
	}

	for i, p := range c.Plugins { // global plugins
		if p.Name == "" {
			p.Name = fmt.Sprintf("plugin_%d", i)
		}
		if p.OrderAfter == "" && i > 0 {
			p.OrderAfter = c.Plugins[i-1].Type
		}
		data.Plugins[p.Name] = p
	}

	return data, nil
}

type (
	upstream struct {
		Backends    []string `yaml:"backends"`
		MaxRequests int      `yaml:"max_requests"`
		DialTimeout string   `yaml:"dial_timeout"`

		LBPolicy      string `yaml:"lb_policy"`
		LBTryDuration string `yaml:"lb_try_duration"`
		LBTryInterval string `yaml:"lb_try_interval"`

		HealthURI      string `yaml:"health_uri"`
		HealthPort     int    `yaml:"health_port"`
		HealthInterval string `yaml:"health_interval"`
		HealthTimeout  string `yaml:"health_timeout"`
		HealthStatus   int    `yaml:"health_status"`

		HeaderUp   *olaf.HeaderOps `yaml:"header_up"`
		HeaderDown *olaf.HeaderOps `yaml:"header_down"`
	}

	service struct {
		Name     string    `yaml:"name"`
		Upstream *upstream `yaml:"upstream"`

		Routes  []*route       `yaml:"routes"`
		Plugins []*olaf.Plugin `yaml:"plugins"`
	}

	route struct {
		*olaf.Route `yaml:",inline"`

		Plugins []*olaf.Plugin `yaml:"plugins"`
	}

	content struct {
		Services []*service     `yaml:"services"`
		Plugins  []*olaf.Plugin `yaml:"plugins"`
	}
)
