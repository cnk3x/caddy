package ct

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/kinvolk/container-linux-config-transpiler/config"
	"github.com/kinvolk/container-linux-config-transpiler/config/platform"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func init() {
	caddy.RegisterModule(Ct{})
	httpcaddyfile.RegisterHandlerDirective("ct", parseCaddyfile)
}

// Ct allows to transpile YAML based configuration into a JSON ignition to be used with Flatcar or Fedora CoreOS.
type Ct struct {
	logger *zap.Logger

	// Fail on non critical errors (default: false)
	Strict bool `json:"strict,omitempty"`
	// Only transpile specific MIME types (default: all)
	MIMETypes []string `json:"mime_types,omitempty"`
	// Only for dynamic data must be one of supported types by ct (default: none)
	Platform string `json:"platform,omitempty"`
}

func (Ct) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.ct",
		New: func() caddy.Module {
			return Ct{}
		},
	}
}

func (ct Ct) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	shouldBuf := func(status int, header http.Header) bool {
		bufferThis := true

		if ct.MIMETypes != nil {
			bufferThis = false
			contentType := header.Get("Content-Type")
			for _, mt := range ct.MIMETypes {
				if strings.Contains(contentType, mt) {
					return true
				}
			}
		}

		return bufferThis
	}

	rec := caddyhttp.NewResponseRecorder(w, buf, shouldBuf)

	if err := next.ServeHTTP(rec, r); err != nil {
		return err
	}
	if !rec.Buffered() {
		return nil
	}

	cfg, astNode, report := config.Parse(rec.Buffer().Bytes())
	if len(report.Entries) > 0 {
		ct.logger.Error(report.String())
		if report.IsFatal() || ct.Strict {
			return errors.New("failed to parse config")
		}
	}

	ignCfg, report := config.Convert(cfg, ct.Platform, astNode)
	if len(report.Entries) > 0 {
		ct.logger.Error(report.String())
		if report.IsFatal() || ct.Strict {
			return errors.New("failed to convert config")
		}
	}

	data, err := json.Marshal(&ignCfg)
	if err != nil {
		return err
	}

	rec.Buffer().Reset()
	rec.Header().Del("Content-Length")
	if _, err = rec.Write(data); err != nil {
		return err
	}
	rec.Header().Set("Content-Type", "application/vnd.coreos.ignition+json")

	return rec.WriteResponse()
}

func (ct *Ct) Provision(ctx caddy.Context) error {
	ct.logger = ctx.Logger(ct)
	return nil
}

func (ct *Ct) Validate() error {
	if !platform.IsSupportedPlatform(ct.Platform) {
		return fmt.Errorf("platform %s is unsupported", ct.Platform)
	}

	return nil
}

func (ct *Ct) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "strict":
				ct.Strict = true
			case "mime":
				ct.MIMETypes = d.RemainingArgs()
				if len(ct.MIMETypes) == 0 {
					return d.ArgErr()
				}
			case "platform":
				if !d.Args(&ct.Platform) {
					return d.ArgErr()
				}
			default:
				return fmt.Errorf("unknown subdirective: %q", d.Val())
			}
		}
	}
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var ct Ct
	err := ct.UnmarshalCaddyfile(h.Dispenser)
	return ct, err
}

// Interface guards
var (
	_ caddy.Provisioner     = (*Ct)(nil)
	_ caddyfile.Unmarshaler = (*Ct)(nil)
	_ caddy.Validator       = (*Ct)(nil)
)
