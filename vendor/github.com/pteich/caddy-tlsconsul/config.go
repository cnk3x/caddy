package storageconsul

const (
	// DefaultPrefix defines the default prefix in KV store
	DefaultPrefix = "caddytls"

	// DefaultAESKey needs to be 32 bytes long
	DefaultAESKey = "consultls-1234567890-caddytls-32"

	// DefaultValuePrefix sets a prefix to KV values to check validation
	DefaultValuePrefix = "caddy-storage-consul"

	// DefaultTimeout is the default timeout for Consul connections
	DefaultTimeout = 10

	// EnvNameAESKey defines the env variable name to override AES key
	EnvNameAESKey = "CADDY_CLUSTERING_CONSUL_AESKEY"

	// EnvNamePrefix defines the env variable name to override KV key prefix
	EnvNamePrefix = "CADDY_CLUSTERING_CONSUL_PREFIX"

	// EnvValuePrefix defines the env variable name to override KV value prefix
	EnvValuePrefix = "CADDY_CLUSTERING_CONSUL_VALUEPREFIX"
)
