package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/smallstep/certificates/authority/provisioner"
	cas "github.com/smallstep/certificates/cas/apiv1"
	"github.com/smallstep/certificates/db"
	kms "github.com/smallstep/certificates/kms/apiv1"
	"github.com/smallstep/certificates/templates"
	"go.step.sm/linkedca"
)

const (
	legacyAuthority = "step-certificate-authority"
)

var (
	// DefaultBackdate length of time to backdate certificates to avoid
	// clock skew validation issues.
	DefaultBackdate = time.Minute
	// DefaultDisableRenewal disables renewals per provisioner.
	DefaultDisableRenewal = false
	// DefaultAllowRenewalAfterExpiry allows renewals even if the certificate is
	// expired.
	DefaultAllowRenewalAfterExpiry = false
	// DefaultEnableSSHCA enable SSH CA features per provisioner or globally
	// for all provisioners.
	DefaultEnableSSHCA = false
	// GlobalProvisionerClaims default claims for the Authority. Can be overridden
	// by provisioner specific claims.
	GlobalProvisionerClaims = provisioner.Claims{
		MinTLSDur:               &provisioner.Duration{Duration: 5 * time.Minute}, // TLS certs
		MaxTLSDur:               &provisioner.Duration{Duration: 24 * time.Hour},
		DefaultTLSDur:           &provisioner.Duration{Duration: 24 * time.Hour},
		MinUserSSHDur:           &provisioner.Duration{Duration: 5 * time.Minute}, // User SSH certs
		MaxUserSSHDur:           &provisioner.Duration{Duration: 24 * time.Hour},
		DefaultUserSSHDur:       &provisioner.Duration{Duration: 16 * time.Hour},
		MinHostSSHDur:           &provisioner.Duration{Duration: 5 * time.Minute}, // Host SSH certs
		MaxHostSSHDur:           &provisioner.Duration{Duration: 30 * 24 * time.Hour},
		DefaultHostSSHDur:       &provisioner.Duration{Duration: 30 * 24 * time.Hour},
		EnableSSHCA:             &DefaultEnableSSHCA,
		DisableRenewal:          &DefaultDisableRenewal,
		AllowRenewalAfterExpiry: &DefaultAllowRenewalAfterExpiry,
	}
)

// Config represents the CA configuration and it's mapped to a JSON object.
type Config struct {
	Root             multiString          `json:"root"`
	FederatedRoots   []string             `json:"federatedRoots"`
	IntermediateCert string               `json:"crt"`
	IntermediateKey  string               `json:"key"`
	Address          string               `json:"address"`
	InsecureAddress  string               `json:"insecureAddress"`
	DNSNames         []string             `json:"dnsNames"`
	KMS              *kms.Options         `json:"kms,omitempty"`
	SSH              *SSHConfig           `json:"ssh,omitempty"`
	Logger           json.RawMessage      `json:"logger,omitempty"`
	DB               *db.Config           `json:"db,omitempty"`
	Monitoring       json.RawMessage      `json:"monitoring,omitempty"`
	AuthorityConfig  *AuthConfig          `json:"authority,omitempty"`
	TLS              *TLSOptions          `json:"tls,omitempty"`
	Password         string               `json:"password,omitempty"`
	Templates        *templates.Templates `json:"templates,omitempty"`
	CommonName       string               `json:"commonName,omitempty"`
}

// ASN1DN contains ASN1.DN attributes that are used in Subject and Issuer
// x509 Certificate blocks.
type ASN1DN struct {
	Country            string `json:"country,omitempty"`
	Organization       string `json:"organization,omitempty"`
	OrganizationalUnit string `json:"organizationalUnit,omitempty"`
	Locality           string `json:"locality,omitempty"`
	Province           string `json:"province,omitempty"`
	StreetAddress      string `json:"streetAddress,omitempty"`
	SerialNumber       string `json:"serialNumber,omitempty"`
	CommonName         string `json:"commonName,omitempty"`
}

// AuthConfig represents the configuration options for the authority. An
// underlaying registration authority can also be configured using the
// cas.Options.
type AuthConfig struct {
	*cas.Options
	AuthorityID          string                `json:"authorityId,omitempty"`
	DeploymentType       string                `json:"deploymentType,omitempty"`
	Provisioners         provisioner.List      `json:"provisioners,omitempty"`
	Admins               []*linkedca.Admin     `json:"-"`
	Template             *ASN1DN               `json:"template,omitempty"`
	Claims               *provisioner.Claims   `json:"claims,omitempty"`
	DisableIssuedAtCheck bool                  `json:"disableIssuedAtCheck,omitempty"`
	Backdate             *provisioner.Duration `json:"backdate,omitempty"`
	EnableAdmin          bool                  `json:"enableAdmin,omitempty"`
}

// init initializes the required fields in the AuthConfig if they are not
// provided.
func (c *AuthConfig) init() {
	if c.Provisioners == nil {
		c.Provisioners = provisioner.List{}
	}
	if c.Template == nil {
		c.Template = &ASN1DN{}
	}
	if c.Backdate == nil {
		c.Backdate = &provisioner.Duration{
			Duration: DefaultBackdate,
		}
	}
}

// Validate validates the authority configuration.
func (c *AuthConfig) Validate(audiences provisioner.Audiences) error {
	if c == nil {
		return errors.New("authority cannot be undefined")
	}

	// Initialize required fields.
	c.init()

	// Check that only one K8sSA is enabled
	var k8sCount int
	for _, p := range c.Provisioners {
		if p.GetType() == provisioner.TypeK8sSA {
			k8sCount++
		}
	}
	if k8sCount > 1 {
		return errors.New("cannot have more than one kubernetes service account provisioner")
	}

	if c.Backdate.Duration < 0 {
		return errors.New("authority.backdate cannot be less than 0")
	}

	return nil
}

// LoadConfiguration parses the given filename in JSON format and returns the
// configuration struct.
func LoadConfiguration(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening %s", filename)
	}
	defer f.Close()

	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, errors.Wrapf(err, "error parsing %s", filename)
	}

	c.Init()

	return &c, nil
}

// Init initializes the minimal configuration required to create an authority. This
// is mainly used on embedded authorities.
func (c *Config) Init() {
	if c.DNSNames == nil {
		c.DNSNames = []string{"localhost", "127.0.0.1", "::1"}
	}
	if c.TLS == nil {
		c.TLS = &DefaultTLSOptions
	}
	if c.AuthorityConfig == nil {
		c.AuthorityConfig = &AuthConfig{}
	}
	if c.CommonName == "" {
		c.CommonName = "Step Online CA"
	}
	c.AuthorityConfig.init()
}

// Save saves the configuration to the given filename.
func (c *Config) Save(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrapf(err, "error opening %s", filename)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return errors.Wrapf(enc.Encode(c), "error writing %s", filename)
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	switch {
	case c.Address == "":
		return errors.New("address cannot be empty")
	case len(c.DNSNames) == 0:
		return errors.New("dnsNames cannot be empty")
	case c.AuthorityConfig == nil:
		return errors.New("authority cannot be nil")
	}

	// Options holds the RA/CAS configuration.
	ra := c.AuthorityConfig.Options
	// The default RA/CAS requires root, crt and key.
	if ra.Is(cas.SoftCAS) {
		switch {
		case c.Root.HasEmpties():
			return errors.New("root cannot be empty")
		case c.IntermediateCert == "":
			return errors.New("crt cannot be empty")
		case c.IntermediateKey == "":
			return errors.New("key cannot be empty")
		}
	}

	// Validate address (a port is required)
	if _, _, err := net.SplitHostPort(c.Address); err != nil {
		return errors.Errorf("invalid address %s", c.Address)
	}

	if c.TLS == nil {
		c.TLS = &DefaultTLSOptions
	} else {
		if len(c.TLS.CipherSuites) == 0 {
			c.TLS.CipherSuites = DefaultTLSOptions.CipherSuites
		}
		if c.TLS.MaxVersion == 0 {
			c.TLS.MaxVersion = DefaultTLSOptions.MaxVersion
		}
		if c.TLS.MinVersion == 0 {
			c.TLS.MinVersion = DefaultTLSOptions.MinVersion
		}
		if c.TLS.MinVersion > c.TLS.MaxVersion {
			return errors.New("tls minVersion cannot exceed tls maxVersion")
		}
		c.TLS.Renegotiation = c.TLS.Renegotiation || DefaultTLSOptions.Renegotiation
	}

	// Validate KMS options, nil is ok.
	if err := c.KMS.Validate(); err != nil {
		return err
	}

	// Validate RA/CAS options, nil is ok.
	if err := ra.Validate(); err != nil {
		return err
	}

	// Validate ssh: nil is ok
	if err := c.SSH.Validate(); err != nil {
		return err
	}

	// Validate templates: nil is ok
	if err := c.Templates.Validate(); err != nil {
		return err
	}

	return c.AuthorityConfig.Validate(c.GetAudiences())
}

// GetAudiences returns the legacy and possible urls without the ports that will
// be used as the default provisioner audiences. The CA might have proxies in
// front so we cannot rely on the port.
func (c *Config) GetAudiences() provisioner.Audiences {
	audiences := provisioner.Audiences{
		Sign:      []string{legacyAuthority},
		Revoke:    []string{legacyAuthority},
		SSHSign:   []string{},
		SSHRevoke: []string{},
		SSHRenew:  []string{},
	}

	for _, name := range c.DNSNames {
		hostname := toHostname(name)
		audiences.Sign = append(audiences.Sign,
			fmt.Sprintf("https://%s/1.0/sign", hostname),
			fmt.Sprintf("https://%s/sign", hostname),
			fmt.Sprintf("https://%s/1.0/ssh/sign", hostname),
			fmt.Sprintf("https://%s/ssh/sign", hostname))
		audiences.Renew = append(audiences.Renew,
			fmt.Sprintf("https://%s/1.0/renew", hostname),
			fmt.Sprintf("https://%s/renew", hostname))
		audiences.Revoke = append(audiences.Revoke,
			fmt.Sprintf("https://%s/1.0/revoke", hostname),
			fmt.Sprintf("https://%s/revoke", hostname))
		audiences.SSHSign = append(audiences.SSHSign,
			fmt.Sprintf("https://%s/1.0/ssh/sign", hostname),
			fmt.Sprintf("https://%s/ssh/sign", hostname),
			fmt.Sprintf("https://%s/1.0/sign", hostname),
			fmt.Sprintf("https://%s/sign", hostname))
		audiences.SSHRevoke = append(audiences.SSHRevoke,
			fmt.Sprintf("https://%s/1.0/ssh/revoke", hostname),
			fmt.Sprintf("https://%s/ssh/revoke", hostname))
		audiences.SSHRenew = append(audiences.SSHRenew,
			fmt.Sprintf("https://%s/1.0/ssh/renew", hostname),
			fmt.Sprintf("https://%s/ssh/renew", hostname))
		audiences.SSHRekey = append(audiences.SSHRekey,
			fmt.Sprintf("https://%s/1.0/ssh/rekey", hostname),
			fmt.Sprintf("https://%s/ssh/rekey", hostname))
	}

	return audiences
}

// Audience returns the list of audiences for a given path.
func (c *Config) Audience(path string) []string {
	audiences := make([]string, len(c.DNSNames)+1)
	for i, name := range c.DNSNames {
		hostname := toHostname(name)
		audiences[i] = "https://" + hostname + path
	}
	// For backward compatibility
	audiences[len(c.DNSNames)] = path
	return audiences
}

func toHostname(name string) string {
	// ensure an IPv6 address is represented with square brackets when used as hostname
	if ip := net.ParseIP(name); ip != nil && ip.To4() == nil {
		name = "[" + name + "]"
	}
	return name
}
