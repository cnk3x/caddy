package linkedca

import "context"

type contextKeyType int

const (
	_ contextKeyType = iota
	adminContextKey
	provisionerContextKey
	externalAccountKeyContextKey
)

// NewContextWithAdmin returns a copy of ctx which carries an Admin.
func NewContextWithAdmin(ctx context.Context, admin *Admin) context.Context {
	return context.WithValue(ctx, adminContextKey, admin)
}

// AdminFromContext returns the Admin ctx carries.
//
// AdminFromContext panics in case ctx carries no Admin.
func AdminFromContext(ctx context.Context) *Admin {
	return ctx.Value(adminContextKey).(*Admin)
}

// NewContextWithProvisioner returns a copy of ctx which carries a Provisioner.
func NewContextWithProvisioner(ctx context.Context, provisioner *Provisioner) context.Context {
	return context.WithValue(ctx, provisionerContextKey, provisioner)
}

// ProvisionerFromContext returns the Provisioner ctx carries.
//
// ProvisionerFromContext panics in case ctx carries no Provisioner.
func ProvisionerFromContext(ctx context.Context) *Provisioner {
	return ctx.Value(provisionerContextKey).(*Provisioner)
}

// NewContextWithExternalAccountKey returns a copy of ctx which carries an EABKey.
func NewContextWithExternalAccountKey(ctx context.Context, k *EABKey) context.Context {
	return context.WithValue(ctx, externalAccountKeyContextKey, k)
}

// ExternalAccountKeyFromContext returns the EABKey ctx carries.
//
// ExternalAccountKeyFromContext panics in case ctx carries no EABKey.
func ExternalAccountKeyFromContext(ctx context.Context) *EABKey {
	return ctx.Value(externalAccountKeyContextKey).(*EABKey)
}
