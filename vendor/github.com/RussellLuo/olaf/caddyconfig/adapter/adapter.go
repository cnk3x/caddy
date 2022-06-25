package adapter

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

func init() {
	caddyconfig.RegisterAdapter("olaf", Adapter{})
}

// Adapter adapts Olaf's configuration to Caddy JSON.
type Adapter struct{}

// Adapt converts the Olaf's configuration in body to Caddy JSON.
func (Adapter) Adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	caddyfileAdapter := caddyconfig.GetAdapter("caddyfile")
	caddyfileResult, warn, err := caddyfileAdapter.Adapt(body, options)
	if err != nil {
		return nil, warn, err
	}

	result, err := patch(caddyfileResult)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func patch(caddyfileResult []byte) ([]byte, error) {
	config := make(map[string]interface{})
	if err := json.Unmarshal(caddyfileResult, &config); err != nil {
		return nil, err
	}

	expander := NewExpander(nil)
	if err := expander.Expand(config); err != nil {
		return nil, err
	}

	result, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return result, nil
}
