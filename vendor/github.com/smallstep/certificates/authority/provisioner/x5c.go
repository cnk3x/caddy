package provisioner

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/smallstep/certificates/errs"
	"go.step.sm/crypto/jose"
	"go.step.sm/crypto/sshutil"
	"go.step.sm/crypto/x509util"
)

// x5cPayload extends jwt.Claims with step attributes.
type x5cPayload struct {
	jose.Claims
	SANs   []string     `json:"sans,omitempty"`
	Step   *stepPayload `json:"step,omitempty"`
	chains [][]*x509.Certificate
}

// X5C is the default provisioner, an entity that can sign tokens necessary for
// signature requests.
type X5C struct {
	*base
	ID        string   `json:"-"`
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Roots     []byte   `json:"roots"`
	Claims    *Claims  `json:"claims,omitempty"`
	Options   *Options `json:"options,omitempty"`
	claimer   *Claimer
	audiences Audiences
	rootPool  *x509.CertPool
}

// GetID returns the provisioner unique identifier. The name and credential id
// should uniquely identify any X5C provisioner.
func (p *X5C) GetID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.GetIDForToken()
}

// GetIDForToken returns an identifier that will be used to load the provisioner
// from a token.
func (p *X5C) GetIDForToken() string {
	return "x5c/" + p.Name
}

// GetTokenID returns the identifier of the token.
func (p *X5C) GetTokenID(ott string) (string, error) {
	// Validate payload
	token, err := jose.ParseSigned(ott)
	if err != nil {
		return "", errors.Wrap(err, "error parsing token")
	}

	// Get claims w/out verification. We need to look up the provisioner
	// key in order to verify the claims and we need the issuer from the claims
	// before we can look up the provisioner.
	var claims jose.Claims
	if err = token.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return "", errors.Wrap(err, "error verifying claims")
	}
	return claims.ID, nil
}

// GetName returns the name of the provisioner.
func (p *X5C) GetName() string {
	return p.Name
}

// GetType returns the type of provisioner.
func (p *X5C) GetType() Type {
	return TypeX5C
}

// GetEncryptedKey returns the base provisioner encrypted key if it's defined.
func (p *X5C) GetEncryptedKey() (string, string, bool) {
	return "", "", false
}

// Init initializes and validates the fields of a X5C type.
func (p *X5C) Init(config Config) error {
	switch {
	case p.Type == "":
		return errors.New("provisioner type cannot be empty")
	case p.Name == "":
		return errors.New("provisioner name cannot be empty")
	case len(p.Roots) == 0:
		return errors.New("provisioner root(s) cannot be empty")
	}

	p.rootPool = x509.NewCertPool()

	var (
		block *pem.Block
		rest  = p.Roots
	)
	for rest != nil {
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return errors.Wrap(err, "error parsing x509 certificate from PEM block")
		}
		p.rootPool.AddCert(cert)
	}

	// Verify that at least one root was found.
	if len(p.rootPool.Subjects()) == 0 {
		return errors.Errorf("no x509 certificates found in roots attribute for provisioner '%s'", p.GetName())
	}

	// Update claims with global ones
	var err error
	if p.claimer, err = NewClaimer(p.Claims, config.Claims); err != nil {
		return err
	}

	p.audiences = config.Audiences.WithFragment(p.GetIDForToken())
	return nil
}

// authorizeToken performs common jwt authorization actions and returns the
// claims for case specific downstream parsing.
// e.g. a Sign request will auth/validate different fields than a Revoke request.
func (p *X5C) authorizeToken(token string, audiences []string) (*x5cPayload, error) {
	jwt, err := jose.ParseSigned(token)
	if err != nil {
		return nil, errs.Wrap(http.StatusUnauthorized, err, "x5c.authorizeToken; error parsing x5c token")
	}

	verifiedChains, err := jwt.Headers[0].Certificates(x509.VerifyOptions{
		Roots:     p.rootPool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return nil, errs.Wrap(http.StatusUnauthorized, err,
			"x5c.authorizeToken; error verifying x5c certificate chain in token")
	}
	leaf := verifiedChains[0][0]

	if leaf.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return nil, errs.Unauthorized("x5c.authorizeToken; certificate used to sign x5c token cannot be used for digital signature")
	}

	// Using the leaf certificates key to validate the claims accomplishes two
	// things:
	//   1. Asserts that the private key used to sign the token corresponds
	//      to the public certificate in the `x5c` header of the token.
	//   2. Asserts that the claims are valid - have not been tampered with.
	var claims x5cPayload
	if err = jwt.Claims(leaf.PublicKey, &claims); err != nil {
		return nil, errs.Wrap(http.StatusUnauthorized, err, "x5c.authorizeToken; error parsing x5c claims")
	}

	// According to "rfc7519 JSON Web Token" acceptable skew should be no
	// more than a few minutes.
	if err = claims.ValidateWithLeeway(jose.Expected{
		Issuer: p.Name,
		Time:   time.Now().UTC(),
	}, time.Minute); err != nil {
		return nil, errs.Wrapf(http.StatusUnauthorized, err, "x5c.authorizeToken; invalid x5c claims")
	}

	// validate audiences with the defaults
	if !matchesAudience(claims.Audience, audiences) {
		return nil, errs.Unauthorized("x5c.authorizeToken; x5c token has invalid audience "+
			"claim (aud); expected %s, but got %s", audiences, claims.Audience)
	}

	if claims.Subject == "" {
		return nil, errs.Unauthorized("x5c.authorizeToken; x5c token subject cannot be empty")
	}

	// Save the verified chains on the x5c payload object.
	claims.chains = verifiedChains
	return &claims, nil
}

// AuthorizeRevoke returns an error if the provisioner does not have rights to
// revoke the certificate with serial number in the `sub` property.
func (p *X5C) AuthorizeRevoke(ctx context.Context, token string) error {
	_, err := p.authorizeToken(token, p.audiences.Revoke)
	return errs.Wrap(http.StatusInternalServerError, err, "x5c.AuthorizeRevoke")
}

// AuthorizeSign validates the given token.
func (p *X5C) AuthorizeSign(ctx context.Context, token string) ([]SignOption, error) {
	claims, err := p.authorizeToken(token, p.audiences.Sign)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "x5c.AuthorizeSign")
	}

	// NOTE: This is for backwards compatibility with older versions of cli
	// and certificates. Older versions added the token subject as the only SAN
	// in a CSR by default.
	if len(claims.SANs) == 0 {
		claims.SANs = []string{claims.Subject}
	}

	// Certificate templates
	data := x509util.CreateTemplateData(claims.Subject, claims.SANs)
	if v, err := unsafeParseSigned(token); err == nil {
		data.SetToken(v)
	}

	templateOptions, err := TemplateOptions(p.Options, data)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "jwk.AuthorizeSign")
	}

	return []SignOption{
		templateOptions,
		// modifiers / withOptions
		newProvisionerExtensionOption(TypeX5C, p.Name, ""),
		profileLimitDuration{p.claimer.DefaultTLSCertDuration(),
			claims.chains[0][0].NotBefore, claims.chains[0][0].NotAfter},
		// validators
		commonNameValidator(claims.Subject),
		defaultSANsValidator(claims.SANs),
		defaultPublicKeyValidator{},
		newValidityValidator(p.claimer.MinTLSCertDuration(), p.claimer.MaxTLSCertDuration()),
	}, nil
}

// AuthorizeRenew returns an error if the renewal is disabled.
func (p *X5C) AuthorizeRenew(ctx context.Context, cert *x509.Certificate) error {
	if p.claimer.IsDisableRenewal() {
		return errs.Unauthorized("x5c.AuthorizeRenew; renew is disabled for x5c provisioner '%s'", p.GetName())
	}
	return nil
}

// AuthorizeSSHSign returns the list of SignOption for a SignSSH request.
func (p *X5C) AuthorizeSSHSign(ctx context.Context, token string) ([]SignOption, error) {
	if !p.claimer.IsSSHCAEnabled() {
		return nil, errs.Unauthorized("x5c.AuthorizeSSHSign; sshCA is disabled for x5c provisioner '%s'", p.GetName())
	}

	claims, err := p.authorizeToken(token, p.audiences.SSHSign)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "x5c.AuthorizeSSHSign")
	}

	if claims.Step == nil || claims.Step.SSH == nil {
		return nil, errs.Unauthorized("x5c.AuthorizeSSHSign; x5c token must be an SSH provisioning token")
	}

	opts := claims.Step.SSH
	signOptions := []SignOption{
		// validates user's SSHOptions with the ones in the token
		sshCertOptionsValidator(*opts),
		// validate users's KeyID is the token subject.
		sshCertOptionsValidator(SignSSHOptions{KeyID: claims.Subject}),
	}

	// Default template attributes.
	certType := sshutil.UserCert
	keyID := claims.Subject
	principals := []string{claims.Subject}

	// Use options in the token.
	if opts.CertType != "" {
		if certType, err = sshutil.CertTypeFromString(opts.CertType); err != nil {
			return nil, errs.Wrap(http.StatusBadRequest, err, "x5c.AuthorizeSSHSign")
		}
	}
	if opts.KeyID != "" {
		keyID = opts.KeyID
	}
	if len(opts.Principals) > 0 {
		principals = opts.Principals
	}

	// Certificate templates.
	data := sshutil.CreateTemplateData(certType, keyID, principals)
	if v, err := unsafeParseSigned(token); err == nil {
		data.SetToken(v)
	}

	templateOptions, err := TemplateSSHOptions(p.Options, data)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "x5c.AuthorizeSSHSign")
	}
	signOptions = append(signOptions, templateOptions)

	// Add modifiers from custom claims
	t := now()
	if !opts.ValidAfter.IsZero() {
		signOptions = append(signOptions, sshCertValidAfterModifier(opts.ValidAfter.RelativeTime(t).Unix()))
	}
	if !opts.ValidBefore.IsZero() {
		signOptions = append(signOptions, sshCertValidBeforeModifier(opts.ValidBefore.RelativeTime(t).Unix()))
	}

	return append(signOptions,
		// Checks the validity bounds, and set the validity if has not been set.
		&sshLimitDuration{p.claimer, claims.chains[0][0].NotAfter},
		// Validate public key.
		&sshDefaultPublicKeyValidator{},
		// Validate the validity period.
		&sshCertValidityValidator{p.claimer},
		// Require all the fields in the SSH certificate
		&sshCertDefaultValidator{},
	), nil
}
