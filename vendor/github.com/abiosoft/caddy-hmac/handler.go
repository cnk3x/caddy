package hmac

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m HMAC) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	if r.Body == nil {
		// nothing to do
		return next.ServeHTTP(w, r)
	}
	body, err := copyRequestBody(r)
	if err != nil {
		return err
	}

	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

	secret := repl.ReplaceAll(m.Secret, "")
	signature := generateSignature(m.hasher, secret, body)
	if err != nil {
		return err
	}

	repl.Set(m.replacerKey(), signature)
	return next.ServeHTTP(w, r)
}

// copyRequestBody copies the request body while making it reusable.
// It returns the copied []byte.
func copyRequestBody(r *http.Request) ([]byte, error) {
	bodyCopy := bytes.Buffer{}
	tee := io.TeeReader(r.Body, &bodyCopy)
	body, err := ioutil.ReadAll(tee)
	if err != nil {
		return nil, err
	}

	// replace the body
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	// return the copy
	return bodyCopy.Bytes(), nil
}
