package hmac

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"
)

type hashAlgorithm string

func (h hashAlgorithm) valid() bool {
	switch h {
	case algSha1:
	case algSha256:
	case algMd5:
	default:
		return false
	}
	return true
}

// hash algorithms
const (
	algSha1   hashAlgorithm = "sha1"
	algSha256 hashAlgorithm = "sha256"
	algMd5    hashAlgorithm = "md5"
)

// generateSignature generates hmac signature using the hasher, secret
// and bytes.
func generateSignature(hasher func() hash.Hash, secret string, b []byte) string {
	h := hmac.New(hasher, []byte(secret))

	// no error check needed, never returns an error
	// https://github.com/golang/go/blob/go1.14.3/src/hash/hash.go#L28
	h.Write(b)

	return hex.EncodeToString(h.Sum(nil))
}
