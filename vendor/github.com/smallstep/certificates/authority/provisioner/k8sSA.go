package provisioner

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"

	"github.com/pkg/errors"
	"github.com/smallstep/certificates/errs"
	"go.step.sm/crypto/jose"
	"go.step.sm/crypto/pemutil"
	"go.step.sm/crypto/sshutil"
	"go.step.sm/crypto/x509util"
)

// NOTE: There can be at most one kubernetes service account provisioner configured
// per instance of step-ca. This is due to a lack of distinguishing information
// contained in kubernetes service account tokens.

const (
	// K8sSAName is the default name used for kubernetes service account provisioners.
	K8sSAName = "k8sSA-default"
	// K8sSAID is the default ID for kubernetes service account provisioners.
	K8sSAID     = "k8ssa/" + K8sSAName
	k8sSAIssuer = "kubernetes/serviceaccount"
)

// jwtPayload extends jwt.Claims with step attributes.
type k8sSAPayload struct {
	jose.Claims
	Namespace          string `json:"kubernetes.io/serviceaccount/namespace,omitempty"`
	SecretName         string `json:"kubernetes.io/serviceaccount/secret.name,omitempty"`
	ServiceAccountName string `json:"kubernetes.io/serviceaccount/service-account.name,omitempty"`
	ServiceAccountUID  string `json:"kubernetes.io/serviceaccount/service-account.uid,omitempty"`
}

// K8sSA represents a Kubernetes ServiceAccount provisioner; an
// entity trusted to make signature requests.
type K8sSA struct {
	*base
	ID        string   `json:"-"`
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	PubKeys   []byte   `json:"publicKeys,omitempty"`
	Claims    *Claims  `json:"claims,omitempty"`
	Options   *Options `json:"options,omitempty"`
	claimer   *Claimer
	audiences Audiences
	//kauthn    kauthn.AuthenticationV1Interface
	pubKeys []interface{}
}

// GetID returns the provisioner unique identifier. The name and credential id
// should uniquely identify any K8sSA provisioner.
func (p *K8sSA) GetID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.GetIDForToken()
}

// GetIDForToken returns an identifier that will be used to load the provisioner
// from a token.
func (p *K8sSA) GetIDForToken() string {
	return K8sSAID
}

// GetTokenID returns an unimplemented error and does not use the input ott.
func (p *K8sSA) GetTokenID(ott string) (string, error) {
	return "", errors.New("not implemented")
}

// GetName returns the name of the provisioner.
func (p *K8sSA) GetName() string {
	return p.Name
}

// GetType returns the type of provisioner.
func (p *K8sSA) GetType() Type {
	return TypeK8sSA
}

// GetEncryptedKey returns false, because the kubernetes provisioner does not
// have access to the private key.
func (p *K8sSA) GetEncryptedKey() (string, string, bool) {
	return "", "", false
}

// Init initializes and validates the fields of a K8sSA type.
func (p *K8sSA) Init(config Config) (err error) {
	switch {
	case p.Type == "":
		return errors.New("provisioner type cannot be empty")
	case p.Name == "":
		return errors.New("provisioner name cannot be empty")
	}

	if p.PubKeys != nil {
		var (
			block *pem.Block
			rest  = p.PubKeys
		)
		for rest != nil {
			block, rest = pem.Decode(rest)
			if block == nil {
				break
			}
			key, err := pemutil.ParseKey(pem.EncodeToMemory(block))
			if err != nil {
				return errors.Wrapf(err, "error parsing public key in provisioner '%s'", p.GetName())
			}
			switch q := key.(type) {
			case *rsa.PublicKey, *ecdsa.PublicKey, ed25519.PublicKey:
			default:
				return errors.Errorf("Unexpected public key type %T in provisioner '%s'", q, p.GetName())
			}
			p.pubKeys = append(p.pubKeys, key)
		}
	} else {
		// TODO: Use the TokenReview API if no pub keys provided. This will need to
		// be configured with additional attributes in the K8sSA struct for
		// connecting to the kubernetes API server.
		return errors.New("K8s Service Account provisioner cannot be initialized without pub keys")
	}
	/*
		// NOTE: Not sure if we should be doing this initialization here ...
		// If you have a k8sSA provisioner defined in your config, but you're not
		// in a kubernetes pod then your CA will fail to startup. Maybe we just postpone
		// creating the authn until token validation time?
		if err := checkAccess(k8s.AuthorizationV1()); err != nil {
			return errors.Wrapf(err, "error verifying access to kubernetes authz service for provisioner %s", p.GetID())
		}

		p.kauthn = k8s.AuthenticationV1()
	*/

	// Update claims with global ones
	if p.claimer, err = NewClaimer(p.Claims, config.Claims); err != nil {
		return err
	}

	p.audiences = config.Audiences
	return err
}

// authorizeToken performs common jwt authorization actions and returns the
// claims for case specific downstream parsing.
// e.g. a Sign request will auth/validate different fields than a Revoke request.
func (p *K8sSA) authorizeToken(token string, audiences []string) (*k8sSAPayload, error) {
	jwt, err := jose.ParseSigned(token)
	if err != nil {
		return nil, errs.Wrap(http.StatusUnauthorized, err,
			"k8ssa.authorizeToken; error parsing k8sSA token")
	}

	var (
		valid  bool
		claims k8sSAPayload
	)
	if p.pubKeys == nil {
		return nil, errs.Unauthorized("k8ssa.authorizeToken; k8sSA TokenReview API integration not implemented")
		/* NOTE: We plan to support the TokenReview API in a future release.
		         Below is some code that should be useful when we prioritize
				 this integration.

			tr := kauthnApi.TokenReview{Spec: kauthnApi.TokenReviewSpec{Token: string(token)}}
			rvw, err := p.kauthn.TokenReviews().Create(&tr)
			if err != nil {
				return nil, errors.Wrap(err, "error using kubernetes TokenReview API")
			}
			if rvw.Status.Error != "" {
				return nil, errors.Errorf("error from kubernetes TokenReviewAPI: %s", rvw.Status.Error)
			}
			if !rvw.Status.Authenticated {
				return nil, errors.New("error from kubernetes TokenReviewAPI: token could not be authenticated")
			}
			if err = jwt.UnsafeClaimsWithoutVerification(&claims); err != nil {
				return nil, errors.Wrap(err, "error parsing claims")
			}
		*/
	}
	for _, pk := range p.pubKeys {
		if err = jwt.Claims(pk, &claims); err == nil {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errs.Unauthorized("k8ssa.authorizeToken; error validating k8sSA token and extracting claims")
	}

	// According to "rfc7519 JSON Web Token" acceptable skew should be no
	// more than a few minutes.
	if err = claims.Validate(jose.Expected{
		Issuer: k8sSAIssuer,
	}); err != nil {
		return nil, errs.Wrap(http.StatusUnauthorized, err, "k8ssa.authorizeToken; invalid k8sSA token claims")
	}

	if claims.Subject == "" {
		return nil, errs.Unauthorized("k8ssa.authorizeToken; k8sSA token subject cannot be empty")
	}

	return &claims, nil
}

// AuthorizeRevoke returns an error if the provisioner does not have rights to
// revoke the certificate with serial number in the `sub` property.
func (p *K8sSA) AuthorizeRevoke(ctx context.Context, token string) error {
	_, err := p.authorizeToken(token, p.audiences.Revoke)
	return errs.Wrap(http.StatusInternalServerError, err, "k8ssa.AuthorizeRevoke")
}

// AuthorizeSign validates the given token.
func (p *K8sSA) AuthorizeSign(ctx context.Context, token string) ([]SignOption, error) {
	claims, err := p.authorizeToken(token, p.audiences.Sign)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "k8ssa.AuthorizeSign")
	}

	// Add some values to use in custom templates.
	data := x509util.NewTemplateData()
	data.SetCommonName(claims.ServiceAccountName)
	if v, err := unsafeParseSigned(token); err == nil {
		data.SetToken(v)
	}

	// Certificate templates: on K8sSA the default template is the certificate
	// request.
	templateOptions, err := CustomTemplateOptions(p.Options, data, x509util.DefaultAdminLeafTemplate)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "k8ssa.AuthorizeSign")
	}

	return []SignOption{
		templateOptions,
		// modifiers / withOptions
		newProvisionerExtensionOption(TypeK8sSA, p.Name, ""),
		profileDefaultDuration(p.claimer.DefaultTLSCertDuration()),
		// validators
		defaultPublicKeyValidator{},
		newValidityValidator(p.claimer.MinTLSCertDuration(), p.claimer.MaxTLSCertDuration()),
	}, nil
}

// AuthorizeRenew returns an error if the renewal is disabled.
func (p *K8sSA) AuthorizeRenew(ctx context.Context, cert *x509.Certificate) error {
	if p.claimer.IsDisableRenewal() {
		return errs.Unauthorized("k8ssa.AuthorizeRenew; renew is disabled for k8sSA provisioner '%s'", p.GetName())
	}
	return nil
}

// AuthorizeSSHSign validates an request for an SSH certificate.
func (p *K8sSA) AuthorizeSSHSign(ctx context.Context, token string) ([]SignOption, error) {
	if !p.claimer.IsSSHCAEnabled() {
		return nil, errs.Unauthorized("k8ssa.AuthorizeSSHSign; sshCA is disabled for k8sSA provisioner '%s'", p.GetName())
	}
	claims, err := p.authorizeToken(token, p.audiences.SSHSign)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "k8ssa.AuthorizeSSHSign")
	}

	// Certificate templates.
	// Set some default variables to be used in the templates.
	data := sshutil.CreateTemplateData(sshutil.HostCert, claims.ServiceAccountName, []string{claims.ServiceAccountName})
	if v, err := unsafeParseSigned(token); err == nil {
		data.SetToken(v)
	}

	templateOptions, err := CustomSSHTemplateOptions(p.Options, data, sshutil.CertificateRequestTemplate)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "k8ssa.AuthorizeSSHSign")
	}
	signOptions := []SignOption{templateOptions}

	return append(signOptions,
		// Require type, key-id and principals in the SignSSHOptions.
		&sshCertOptionsRequireValidator{CertType: true, KeyID: true, Principals: true},
		// Set the validity bounds if not set.
		&sshDefaultDuration{p.claimer},
		// Validate public key
		&sshDefaultPublicKeyValidator{},
		// Validate the validity period.
		&sshCertValidityValidator{p.claimer},
		// Require and validate all the default fields in the SSH certificate.
		&sshCertDefaultValidator{},
	), nil
}

/*
func checkAccess(authz kauthz.AuthorizationV1Interface) error {
	r := &kauthzApi.SelfSubjectAccessReview{
		Spec: kauthzApi.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &kauthzApi.ResourceAttributes{
				Group:    "authentication.k8s.io",
				Version:  "v1",
				Resource: "tokenreviews",
				Verb:     "create",
			},
		},
	}
	rvw, err := authz.SelfSubjectAccessReviews().Create(r)
	if err != nil {
		return err
	}
	if !rvw.Status.Allowed {
		return fmt.Errorf("Unable to create kubernetes token reviews: %s", rvw.Status.Reason)
	}

	return nil
}
*/
