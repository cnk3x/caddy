package authority

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/binary"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/smallstep/certificates/authority/config"
	"github.com/smallstep/certificates/authority/provisioner"
	"github.com/smallstep/certificates/db"
	"github.com/smallstep/certificates/errs"
	"github.com/smallstep/certificates/templates"
	"go.step.sm/crypto/randutil"
	"go.step.sm/crypto/sshutil"
	"golang.org/x/crypto/ssh"
)

const (
	// SSHAddUserPrincipal is the principal that will run the add user command.
	// Defaults to "provisioner" but it can be changed in the configuration.
	SSHAddUserPrincipal = "provisioner"

	// SSHAddUserCommand is the default command to run to add a new user.
	// Defaults to "sudo useradd -m <principal>; nc -q0 localhost 22" but it can be changed in the
	// configuration. The string "<principal>" will be replace by the new
	// principal to add.
	SSHAddUserCommand = "sudo useradd -m <principal>; nc -q0 localhost 22"
)

// GetSSHRoots returns the SSH User and Host public keys.
func (a *Authority) GetSSHRoots(context.Context) (*config.SSHKeys, error) {
	return &config.SSHKeys{
		HostKeys: a.sshCAHostCerts,
		UserKeys: a.sshCAUserCerts,
	}, nil
}

// GetSSHFederation returns the public keys for federated SSH signers.
func (a *Authority) GetSSHFederation(context.Context) (*config.SSHKeys, error) {
	return &config.SSHKeys{
		HostKeys: a.sshCAHostFederatedCerts,
		UserKeys: a.sshCAUserFederatedCerts,
	}, nil
}

// GetSSHConfig returns rendered templates for clients (user) or servers (host).
func (a *Authority) GetSSHConfig(ctx context.Context, typ string, data map[string]string) ([]templates.Output, error) {
	if a.sshCAUserCertSignKey == nil && a.sshCAHostCertSignKey == nil {
		return nil, errs.NotFound("getSSHConfig: ssh is not configured")
	}

	if a.templates == nil {
		return nil, errs.NotFound("getSSHConfig: ssh templates are not configured")
	}

	var ts []templates.Template
	switch typ {
	case provisioner.SSHUserCert:
		if a.templates != nil && a.templates.SSH != nil {
			ts = a.templates.SSH.User
		}
	case provisioner.SSHHostCert:
		if a.templates != nil && a.templates.SSH != nil {
			ts = a.templates.SSH.Host
		}
	default:
		return nil, errs.BadRequest("getSSHConfig: type %s is not valid", typ)
	}

	// Merge user and default data
	var mergedData map[string]interface{}

	if len(data) == 0 {
		mergedData = a.templates.Data
	} else {
		mergedData = make(map[string]interface{}, len(a.templates.Data)+1)
		mergedData["User"] = data
		for k, v := range a.templates.Data {
			mergedData[k] = v
		}
	}

	// Render templates
	output := []templates.Output{}
	for _, t := range ts {
		if err := t.Load(); err != nil {
			return nil, err
		}

		// Check for required variables.
		if err := t.ValidateRequiredData(data); err != nil {
			return nil, errs.BadRequestErr(err, errs.WithMessage("%v, please use `--set <key=value>` flag", err))
		}

		o, err := t.Output(mergedData)
		if err != nil {
			return nil, err
		}
		output = append(output, o)
	}
	return output, nil
}

// GetSSHBastion returns the bastion configuration, for the given pair user,
// hostname.
func (a *Authority) GetSSHBastion(ctx context.Context, user, hostname string) (*config.Bastion, error) {
	if a.sshBastionFunc != nil {
		bs, err := a.sshBastionFunc(ctx, user, hostname)
		return bs, errs.Wrap(http.StatusInternalServerError, err, "authority.GetSSHBastion")
	}
	if a.config.SSH != nil {
		if a.config.SSH.Bastion != nil && a.config.SSH.Bastion.Hostname != "" {
			// Do not return a bastion for a bastion host.
			//
			// This condition might fail if a different name or IP is used.
			// Trying to resolve hostnames to IPs and compare them won't be a
			// complete solution because it depends on the network
			// configuration, of the CA and clients and can also return false
			// positives. Although not perfect, this simple solution will work
			// in most cases.
			if !strings.EqualFold(hostname, a.config.SSH.Bastion.Hostname) {
				return a.config.SSH.Bastion, nil
			}
		}
		return nil, nil
	}
	return nil, errs.NotFound("authority.GetSSHBastion; ssh is not configured")
}

// SignSSH creates a signed SSH certificate with the given public key and options.
func (a *Authority) SignSSH(ctx context.Context, key ssh.PublicKey, opts provisioner.SignSSHOptions, signOpts ...provisioner.SignOption) (*ssh.Certificate, error) {
	var (
		certOptions []sshutil.Option
		mods        []provisioner.SSHCertModifier
		validators  []provisioner.SSHCertValidator
	)

	// Validate given options.
	if err := opts.Validate(); err != nil {
		return nil, errs.Wrap(http.StatusBadRequest, err, "authority.SignSSH")
	}

	// Set backdate with the configured value
	opts.Backdate = a.config.AuthorityConfig.Backdate.Duration

	for _, op := range signOpts {
		switch o := op.(type) {
		// add options to NewCertificate
		case provisioner.SSHCertificateOptions:
			certOptions = append(certOptions, o.Options(opts)...)

		// modify the ssh.Certificate
		case provisioner.SSHCertModifier:
			mods = append(mods, o)

		// validate the ssh.Certificate
		case provisioner.SSHCertValidator:
			validators = append(validators, o)

		// validate the given SSHOptions
		case provisioner.SSHCertOptionsValidator:
			if err := o.Valid(opts); err != nil {
				return nil, errs.Wrap(http.StatusForbidden, err, "authority.SignSSH")
			}

		default:
			return nil, errs.InternalServer("authority.SignSSH: invalid extra option type %T", o)
		}
	}

	// Simulated certificate request with request options.
	cr := sshutil.CertificateRequest{
		Type:       opts.CertType,
		KeyID:      opts.KeyID,
		Principals: opts.Principals,
		Key:        key,
	}

	// Create certificate from template.
	certificate, err := sshutil.NewCertificate(cr, certOptions...)
	if err != nil {
		if _, ok := err.(*sshutil.TemplateError); ok {
			return nil, errs.NewErr(http.StatusBadRequest, err,
				errs.WithMessage(err.Error()),
				errs.WithKeyVal("signOptions", signOpts),
			)
		}
		return nil, errs.Wrap(http.StatusInternalServerError, err, "authority.SignSSH")
	}

	// Get actual *ssh.Certificate and continue with provisioner modifiers.
	certTpl := certificate.GetCertificate()

	// Use SignSSHOptions to modify the certificate validity. It will be later
	// checked or set if not defined.
	if err := opts.ModifyValidity(certTpl); err != nil {
		return nil, errs.Wrap(http.StatusBadRequest, err, "authority.SignSSH")
	}

	// Use provisioner modifiers.
	for _, m := range mods {
		if err := m.Modify(certTpl, opts); err != nil {
			return nil, errs.Wrap(http.StatusForbidden, err, "authority.SignSSH")
		}
	}

	// Get signer from authority keys
	var signer ssh.Signer
	switch certTpl.CertType {
	case ssh.UserCert:
		if a.sshCAUserCertSignKey == nil {
			return nil, errs.NotImplemented("authority.SignSSH: user certificate signing is not enabled")
		}
		signer = a.sshCAUserCertSignKey
	case ssh.HostCert:
		if a.sshCAHostCertSignKey == nil {
			return nil, errs.NotImplemented("authority.SignSSH: host certificate signing is not enabled")
		}
		signer = a.sshCAHostCertSignKey
	default:
		return nil, errs.InternalServer("authority.SignSSH: unexpected ssh certificate type: %d", certTpl.CertType)
	}

	// Sign certificate.
	cert, err := sshutil.CreateCertificate(certTpl, signer)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "authority.SignSSH: error signing certificate")
	}

	// User provisioners validators.
	for _, v := range validators {
		if err := v.Valid(cert, opts); err != nil {
			return nil, errs.Wrap(http.StatusForbidden, err, "authority.SignSSH")
		}
	}

	if err = a.storeSSHCertificate(cert); err != nil && err != db.ErrNotImplemented {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "authority.SignSSH: error storing certificate in db")
	}

	return cert, nil
}

// RenewSSH creates a signed SSH certificate using the old SSH certificate as a template.
func (a *Authority) RenewSSH(ctx context.Context, oldCert *ssh.Certificate) (*ssh.Certificate, error) {
	if oldCert.ValidAfter == 0 || oldCert.ValidBefore == 0 {
		return nil, errs.BadRequest("renewSSH: cannot renew certificate without validity period")
	}

	if err := a.authorizeSSHCertificate(ctx, oldCert); err != nil {
		return nil, err
	}

	backdate := a.config.AuthorityConfig.Backdate.Duration
	duration := time.Duration(oldCert.ValidBefore-oldCert.ValidAfter) * time.Second
	now := time.Now()
	va := now.Add(-1 * backdate)
	vb := now.Add(duration - backdate)

	// Build base certificate with the old key.
	// Nonce and serial will be automatically generated on signing.
	certTpl := &ssh.Certificate{
		Key:             oldCert.Key,
		CertType:        oldCert.CertType,
		KeyId:           oldCert.KeyId,
		ValidPrincipals: oldCert.ValidPrincipals,
		Permissions:     oldCert.Permissions,
		Reserved:        oldCert.Reserved,
		ValidAfter:      uint64(va.Unix()),
		ValidBefore:     uint64(vb.Unix()),
	}

	// Get signer from authority keys
	var signer ssh.Signer
	switch certTpl.CertType {
	case ssh.UserCert:
		if a.sshCAUserCertSignKey == nil {
			return nil, errs.NotImplemented("renewSSH: user certificate signing is not enabled")
		}
		signer = a.sshCAUserCertSignKey
	case ssh.HostCert:
		if a.sshCAHostCertSignKey == nil {
			return nil, errs.NotImplemented("renewSSH: host certificate signing is not enabled")
		}
		signer = a.sshCAHostCertSignKey
	default:
		return nil, errs.InternalServer("renewSSH: unexpected ssh certificate type: %d", certTpl.CertType)
	}

	// Sign certificate.
	cert, err := sshutil.CreateCertificate(certTpl, signer)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "signSSH: error signing certificate")
	}

	if err = a.storeSSHCertificate(cert); err != nil && err != db.ErrNotImplemented {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "renewSSH: error storing certificate in db")
	}

	return cert, nil
}

// RekeySSH creates a signed SSH certificate using the old SSH certificate as a template.
func (a *Authority) RekeySSH(ctx context.Context, oldCert *ssh.Certificate, pub ssh.PublicKey, signOpts ...provisioner.SignOption) (*ssh.Certificate, error) {
	var validators []provisioner.SSHCertValidator

	for _, op := range signOpts {
		switch o := op.(type) {
		// validate the ssh.Certificate
		case provisioner.SSHCertValidator:
			validators = append(validators, o)
		default:
			return nil, errs.InternalServer("rekeySSH; invalid extra option type %T", o)
		}
	}

	if oldCert.ValidAfter == 0 || oldCert.ValidBefore == 0 {
		return nil, errs.BadRequest("rekeySSH; cannot rekey certificate without validity period")
	}

	if err := a.authorizeSSHCertificate(ctx, oldCert); err != nil {
		return nil, err
	}

	backdate := a.config.AuthorityConfig.Backdate.Duration
	duration := time.Duration(oldCert.ValidBefore-oldCert.ValidAfter) * time.Second
	now := time.Now()
	va := now.Add(-1 * backdate)
	vb := now.Add(duration - backdate)

	// Build base certificate with the new key.
	// Nonce and serial will be automatically generated on signing.
	cert := &ssh.Certificate{
		Key:             pub,
		CertType:        oldCert.CertType,
		KeyId:           oldCert.KeyId,
		ValidPrincipals: oldCert.ValidPrincipals,
		Permissions:     oldCert.Permissions,
		Reserved:        oldCert.Reserved,
		ValidAfter:      uint64(va.Unix()),
		ValidBefore:     uint64(vb.Unix()),
	}

	// Get signer from authority keys
	var signer ssh.Signer
	switch cert.CertType {
	case ssh.UserCert:
		if a.sshCAUserCertSignKey == nil {
			return nil, errs.NotImplemented("rekeySSH; user certificate signing is not enabled")
		}
		signer = a.sshCAUserCertSignKey
	case ssh.HostCert:
		if a.sshCAHostCertSignKey == nil {
			return nil, errs.NotImplemented("rekeySSH; host certificate signing is not enabled")
		}
		signer = a.sshCAHostCertSignKey
	default:
		return nil, errs.BadRequest("rekeySSH; unexpected ssh certificate type: %d", cert.CertType)
	}

	var err error
	// Sign certificate.
	cert, err = sshutil.CreateCertificate(cert, signer)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "signSSH: error signing certificate")
	}

	// Apply validators from provisioner.
	for _, v := range validators {
		if err := v.Valid(cert, provisioner.SignSSHOptions{Backdate: backdate}); err != nil {
			return nil, errs.Wrap(http.StatusForbidden, err, "rekeySSH")
		}
	}

	if err = a.storeSSHCertificate(cert); err != nil && err != db.ErrNotImplemented {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "rekeySSH; error storing certificate in db")
	}

	return cert, nil
}

func (a *Authority) storeSSHCertificate(cert *ssh.Certificate) error {
	type sshCertificateStorer interface {
		StoreSSHCertificate(crt *ssh.Certificate) error
	}
	if s, ok := a.adminDB.(sshCertificateStorer); ok {
		return s.StoreSSHCertificate(cert)
	}
	return a.db.StoreSSHCertificate(cert)
}

// IsValidForAddUser checks if a user provisioner certificate can be issued to
// the given certificate.
func IsValidForAddUser(cert *ssh.Certificate) error {
	if cert.CertType != ssh.UserCert {
		return errors.New("certificate is not a user certificate")
	}

	switch len(cert.ValidPrincipals) {
	case 0:
		return errors.New("certificate does not have any principals")
	case 1:
		return nil
	case 2:
		// OIDC provisioners adds a second principal with the email address.
		// @ cannot be the first character.
		if strings.Index(cert.ValidPrincipals[1], "@") > 0 {
			return nil
		}
		return errors.New("certificate does not have only one principal")
	default:
		return errors.New("certificate does not have only one principal")
	}
}

// SignSSHAddUser signs a certificate that provisions a new user in a server.
func (a *Authority) SignSSHAddUser(ctx context.Context, key ssh.PublicKey, subject *ssh.Certificate) (*ssh.Certificate, error) {
	if a.sshCAUserCertSignKey == nil {
		return nil, errs.NotImplemented("signSSHAddUser: user certificate signing is not enabled")
	}
	if err := IsValidForAddUser(subject); err != nil {
		return nil, errs.Wrap(http.StatusForbidden, err, "signSSHAddUser")
	}

	nonce, err := randutil.ASCII(32)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "signSSHAddUser")
	}

	var serial uint64
	if err := binary.Read(rand.Reader, binary.BigEndian, &serial); err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "signSSHAddUser: error reading random number")
	}

	signer := a.sshCAUserCertSignKey
	principal := subject.ValidPrincipals[0]
	addUserPrincipal := a.getAddUserPrincipal()

	cert := &ssh.Certificate{
		Nonce:           []byte(nonce),
		Key:             key,
		Serial:          serial,
		CertType:        ssh.UserCert,
		KeyId:           principal + "-" + addUserPrincipal,
		ValidPrincipals: []string{addUserPrincipal},
		ValidAfter:      subject.ValidAfter,
		ValidBefore:     subject.ValidBefore,
		Permissions: ssh.Permissions{
			CriticalOptions: map[string]string{
				"force-command": a.getAddUserCommand(principal),
			},
		},
		SignatureKey: signer.PublicKey(),
	}

	// Get bytes for signing trailing the signature length.
	data := cert.Marshal()
	data = data[:len(data)-4]

	// Sign the certificate
	sig, err := signer.Sign(rand.Reader, data)
	if err != nil {
		return nil, err
	}
	cert.Signature = sig

	if err = a.storeSSHCertificate(cert); err != nil && err != db.ErrNotImplemented {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "signSSHAddUser: error storing certificate in db")
	}

	return cert, nil
}

// CheckSSHHost checks the given principal has been registered before.
func (a *Authority) CheckSSHHost(ctx context.Context, principal, token string) (bool, error) {
	if a.sshCheckHostFunc != nil {
		exists, err := a.sshCheckHostFunc(ctx, principal, token, a.GetRootCertificates())
		if err != nil {
			return false, errs.Wrap(http.StatusInternalServerError, err,
				"checkSSHHost: error from injected checkSSHHost func")
		}
		return exists, nil
	}
	exists, err := a.db.IsSSHHost(principal)
	if err != nil {
		if err == db.ErrNotImplemented {
			return false, errs.Wrap(http.StatusNotImplemented, err,
				"checkSSHHost: isSSHHost is not implemented")
		}
		return false, errs.Wrap(http.StatusInternalServerError, err,
			"checkSSHHost: error checking if hosts exists")
	}

	return exists, nil
}

// GetSSHHosts returns a list of valid host principals.
func (a *Authority) GetSSHHosts(ctx context.Context, cert *x509.Certificate) ([]config.Host, error) {
	if a.sshGetHostsFunc != nil {
		hosts, err := a.sshGetHostsFunc(ctx, cert)
		return hosts, errs.Wrap(http.StatusInternalServerError, err, "getSSHHosts")
	}
	hostnames, err := a.db.GetSSHHostPrincipals()
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, err, "getSSHHosts")
	}

	hosts := make([]config.Host, len(hostnames))
	for i, hn := range hostnames {
		hosts[i] = config.Host{Hostname: hn}
	}
	return hosts, nil
}

func (a *Authority) getAddUserPrincipal() (cmd string) {
	if a.config.SSH.AddUserPrincipal == "" {
		return SSHAddUserPrincipal
	}
	return a.config.SSH.AddUserPrincipal
}

func (a *Authority) getAddUserCommand(principal string) string {
	var cmd string
	if a.config.SSH.AddUserCommand == "" {
		cmd = SSHAddUserCommand
	} else {
		cmd = a.config.SSH.AddUserCommand
	}
	return strings.ReplaceAll(cmd, "<principal>", principal)
}
