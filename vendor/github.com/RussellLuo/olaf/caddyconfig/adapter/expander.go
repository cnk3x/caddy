package adapter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/RussellLuo/olaf"
	"github.com/RussellLuo/olaf/caddyconfig/builder"
	"github.com/RussellLuo/olaf/caddymodule"
	"github.com/RussellLuo/olaf/store/yaml"
	"github.com/mitchellh/mapstructure"
)

type Apps struct {
	HTTP struct {
		Servers map[string]struct {
			Listen []string                 `json:"listen"`
			Routes []map[string]interface{} `json:"routes"`
		} `json:"servers"`
	} `json:"http"`
}

type Loader interface {
	Load(mod *caddymodule.Olaf) (*olaf.Data, error)
}

type Expander struct {
	loader Loader
}

func NewExpander(loader Loader) *Expander {
	if loader == nil {
		loader = defaultLoader{}
	}
	return &Expander{loader: loader}
}

// Expand expands all `olaf` handlers in config. This is done by
// replacing each `olaf` handler with a `subroute` handler, whose routes
// are built by parsing the Olaf's configuration.
func (e *Expander) Expand(config map[string]interface{}) error {
	apps := new(Apps)
	if err := mapstructure.Decode(config["apps"], apps); err != nil {
		return err
	}

	for _, server := range apps.HTTP.Servers {
		if err := e.expand(server.Routes); err != nil {
			return err
		}
	}

	return nil
}

func (e *Expander) expand(routes []map[string]interface{}) error {
NextRoute:
	for _, r := range routes {
		handle := r["handle"].([]interface{})
		for _, h := range handle {
			h := h.(map[string]interface{})

			switch h["handler"] {
			case "olaf":
				mod, err := decodeOlafModule(h)
				if err != nil {
					return err
				}

				data, err := e.loader.Load(mod)
				if err != nil {
					return err
				}

				// Replace the `olaf` handler with a `subroute` handler.
				delete(h, "type")
				delete(h, "path")
				delete(h, "timeout")
				h["handler"] = "subroute"
				h["routes"] = builder.Build(data)

				// We assume that there is only one `olaf` handler in the list.
				continue NextRoute

			case "subroute":
				var subRoutes []map[string]interface{}
				if err := mapstructure.Decode(h["routes"], &subRoutes); err != nil {
					return err
				}
				if err := e.expand(subRoutes); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func decodeOlafModule(h map[string]interface{}) (*caddymodule.Olaf, error) {
	o := new(caddymodule.Olaf)
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   o,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(h); err != nil {
		return nil, err
	}

	return o, nil
}

type defaultLoader struct{}

func (l defaultLoader) Load(mod *caddymodule.Olaf) (*olaf.Data, error) {
	var data *olaf.Data

	switch mod.Type {
	case caddymodule.TypeFile:
		content, err := ioutil.ReadFile(mod.Path)
		if err != nil {
			return nil, err
		}

		data, err = yaml.Parse(content)
		if err != nil {
			return nil, err
		}

	case caddymodule.TypeHTTP:
		client := &http.Client{Timeout: time.Duration(mod.Timeout)}
		resp, err := client.Get(mod.Path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			msg, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("code: %d, err: %s", resp.StatusCode, msg)
		}

		data = new(olaf.Data)
		if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
			return nil, err
		}
	}

	return data, nil
}
