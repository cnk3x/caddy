// Usage:
//
//     gopkg [path] [vcs-type] [uri]
//     gopkg [path] [uri]

package gopkg

import (
	"html/template"
	"net/http"
	"regexp"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("gopkg", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

type Config struct {
	Path string
	Vcs  string
	Uri  string
}

type GopkgHandler struct {
	Next    httpserver.Handler
	Configs []Config
}

var tmpl = template.Must(template.New("").Parse(`<html>
<head>
<meta name="go-import" content="{{.Host}}{{.Path}} {{.Vcs}} {{.Uri}}">
</head>
<body>
go get {{.Host}}{{.Path}}
<hr />
{{ .Ua }}
</body>
</html>
`))

func (g GopkgHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	for _, cfg := range g.Configs {
		rExp, err := regexp.Compile(cfg.Path)
		if err != nil || !rExp.MatchString(r.URL.Path) {
			continue
		}

		uri := rExp.ReplaceAllString(r.URL.Path, cfg.Uri)

		// Check if the request path contains go-get=1
		if r.FormValue("go-get") != "1" {
			http.Redirect(w, r, uri, http.StatusTemporaryRedirect)
			return 0, nil
		}

		host := r.Host

		err = tmpl.Execute(w, struct {
			Host string
			Path string
			Vcs  string
			Uri  string
			Ua   string
		}{
			Host: host,
			Path: r.URL.Path,
			Vcs:  cfg.Vcs,
			Uri:  uri,
			Ua:   r.Header.Get("User-Agent"),
		})
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}

	return g.Next.ServeHTTP(w, r)
}

func setup(c *caddy.Controller) error {
	configs, err := parse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return GopkgHandler{
			Configs: configs,
			Next:    next,
		}
	})
	return nil
}

func parse(c *caddy.Controller) ([]Config, error) {
	var configs []Config

	for c.Next() {

		args := c.RemainingArgs()

		if len(args) != 2 && len(args) != 3 {
			return configs, c.ArgErr()
		}

		cfg := Config{
			Vcs:  "git",
			Path: args[0],
		}

		if len(args) == 2 {
			cfg.Uri = args[1]
		} else {
			cfg.Vcs = args[1]
			cfg.Uri = args[2]
		}

		configs = append(configs, cfg)
	}

	return configs, nil
}
