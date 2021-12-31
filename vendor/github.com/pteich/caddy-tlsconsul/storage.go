package storageconsul

import (
	"context"
	"net"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/caddyserver/certmagic"
	consul "github.com/hashicorp/consul/api"
	"github.com/pteich/errors"
	"go.uber.org/zap"
)

// ConsulStorage allows to store certificates and other TLS resources
// in a shared cluster environment using Consul's key/value-store.
// It uses distributed locks to ensure consistency.
type ConsulStorage struct {
	certmagic.Storage
	ConsulClient *consul.Client
	logger       *zap.SugaredLogger
	muLocks      sync.RWMutex
	locks        map[string]*consul.Lock

	Address     string `json:"address"`
	Token       string `json:"token"`
	Timeout     int    `json:"timeout"`
	Prefix      string `json:"prefix"`
	ValuePrefix string `json:"value_prefix"`
	AESKey      []byte `json:"aes_key"`
	TlsEnabled  bool   `json:"tls_enabled"`
	TlsInsecure bool   `json:"tls_insecure"`
}

// New connects to Consul and returns a ConsulStorage
func New() *ConsulStorage {
	// create ConsulStorage and pre-set values
	s := ConsulStorage{
		locks:       make(map[string]*consul.Lock),
		AESKey:      []byte(DefaultAESKey),
		ValuePrefix: DefaultValuePrefix,
		Prefix:      DefaultPrefix,
		Timeout:     DefaultTimeout,
	}

	return &s
}

func (cs *ConsulStorage) prefixKey(key string) string {
	return path.Join(cs.Prefix, key)
}

// Lock acquires a distributed lock for the given key or blocks until it gets one
func (cs *ConsulStorage) Lock(ctx context.Context, key string) error {
	cs.logger.Debugf("trying lock for %s", key)

	if _, isLocked := cs.GetLock(key); isLocked {
		return nil
	}

	// prepare the distributed lock
	cs.logger.Infof("creating Consul lock for %s", key)
	lock, err := cs.ConsulClient.LockOpts(&consul.LockOptions{
		Key:          cs.prefixKey(key),
		LockWaitTime: time.Duration(cs.Timeout) * time.Second,
		LockTryOnce:  true,
	})
	if err != nil {
		return errors.Wrapf(err, "could not create lock for %s", cs.prefixKey(key))
	}

	// acquire the lock and return a channel that is closed upon lost
	lockActive, err := lock.Lock(ctx.Done())
	if err != nil {
		return errors.Wrapf(err, "unable to lock %s", cs.prefixKey(key))
	}

	// auto-unlock and clean list of locks in case of lost
	go func() {
		<-lockActive
		cs.Unlock(key)
	}()

	// save the lock
	cs.muLocks.Lock()
	cs.locks[key] = lock
	cs.muLocks.Unlock()

	return nil
}

func (cs *ConsulStorage) GetLock(key string) (*consul.Lock, bool) {
	cs.muLocks.RLock()
	defer cs.muLocks.RUnlock()

	// if we already hold the lock, return early
	if lock, exists := cs.locks[key]; exists {
		return lock, true
	}

	return nil, false
}

// Unlock releases a specific lock
func (cs *ConsulStorage) Unlock(key string) error {
	// check if we own it and unlock
	lock, exists := cs.GetLock(key)
	if !exists {
		return errors.Errorf("lock %s not found", cs.prefixKey(key))
	}

	err := lock.Unlock()
	if err != nil {
		return errors.Wrapf(err, "unable to unlock %s", cs.prefixKey(key))
	}

	cs.muLocks.Lock()
	delete(cs.locks, key)
	cs.muLocks.Unlock()

	return nil
}

// Store saves encrypted data value for a key in Consul KV
func (cs ConsulStorage) Store(key string, value []byte) error {
	kv := &consul.KVPair{Key: cs.prefixKey(key)}

	// prepare the stored data
	consulData := &StorageData{
		Value:    value,
		Modified: time.Now(),
	}

	encryptedValue, err := cs.EncryptStorageData(consulData)
	if err != nil {
		return errors.Wrapf(err, "unable to encode data for %s", cs.prefixKey(key))
	}

	kv.Value = encryptedValue

	if _, err = cs.ConsulClient.KV().Put(kv, nil); err != nil {
		return errors.Wrapf(err, "unable to store data for %s", cs.prefixKey(key))
	}

	return nil
}

// Load retrieves the value for a key from Consul KV
func (cs ConsulStorage) Load(key string) ([]byte, error) {
	cs.logger.Debugf("loading data from Consul for %s", key)

	kv, _, err := cs.ConsulClient.KV().Get(cs.prefixKey(key), &consul.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to obtain data for %s", cs.prefixKey(key))
	} else if kv == nil {
		return nil, certmagic.ErrNotExist(errors.Errorf("key %s does not exist", cs.prefixKey(key)))
	}

	contents, err := cs.DecryptStorageData(kv.Value)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decrypt data for %s", cs.prefixKey(key))
	}

	return contents.Value, nil
}

// Delete a key from Consul KV
func (cs ConsulStorage) Delete(key string) error {
	cs.logger.Infof("deleting key %s from Consul", key)

	// first obtain existing keypair
	kv, _, err := cs.ConsulClient.KV().Get(cs.prefixKey(key), &consul.QueryOptions{RequireConsistent: true})
	if err != nil {
		return errors.Wrapf(err, "unable to obtain data for %s", cs.prefixKey(key))
	} else if kv == nil {
		return certmagic.ErrNotExist(err)
	}

	// no do a Check-And-Set operation to verify we really deleted the key
	if success, _, err := cs.ConsulClient.KV().DeleteCAS(kv, nil); err != nil {
		return errors.Wrapf(err, "unable to delete data for %s", cs.prefixKey(key))
	} else if !success {
		return errors.Errorf("failed to lock data delete for %s", cs.prefixKey(key))
	}

	return nil
}

// Exists checks if a key exists
func (cs ConsulStorage) Exists(key string) bool {
	kv, _, err := cs.ConsulClient.KV().Get(cs.prefixKey(key), &consul.QueryOptions{RequireConsistent: true})
	if kv != nil && err == nil {
		return true
	}
	return false
}

// List returns a list with all keys under a given prefix
func (cs ConsulStorage) List(prefix string, recursive bool) ([]string, error) {
	var keysFound []string

	// get a list of all keys at prefix
	keys, _, err := cs.ConsulClient.KV().Keys(cs.prefixKey(prefix), "", &consul.QueryOptions{RequireConsistent: true})
	if err != nil {
		return keysFound, err
	}

	if len(keys) == 0 {
		return keysFound, certmagic.ErrNotExist(errors.Errorf("no keys at %s", prefix))
	}

	// remove default prefix from keys
	for _, key := range keys {
		if strings.HasPrefix(key, cs.prefixKey(prefix)) {
			key = strings.TrimPrefix(key, cs.Prefix+"/")
			keysFound = append(keysFound, key)
		}
	}

	// if recursive wanted, just return all keys
	if recursive {
		return keysFound, nil
	}

	// for non-recursive split path and look for unique keys just under given prefix
	keysMap := make(map[string]bool)
	for _, key := range keysFound {
		dir := strings.Split(strings.TrimPrefix(key, prefix+"/"), "/")
		keysMap[dir[0]] = true
	}

	keysFound = make([]string, 0)
	for key := range keysMap {
		keysFound = append(keysFound, path.Join(prefix, key))
	}

	return keysFound, nil
}

// Stat returns statistic data of a key
func (cs ConsulStorage) Stat(key string) (certmagic.KeyInfo, error) {
	kv, _, err := cs.ConsulClient.KV().Get(cs.prefixKey(key), &consul.QueryOptions{RequireConsistent: true})
	if err != nil {
		return certmagic.KeyInfo{}, errors.Errorf("unable to obtain data for %s", cs.prefixKey(key))
	} else if kv == nil {
		return certmagic.KeyInfo{}, certmagic.ErrNotExist(errors.Errorf("key %s does not exist", cs.prefixKey(key)))
	}

	contents, err := cs.DecryptStorageData(kv.Value)
	if err != nil {
		return certmagic.KeyInfo{}, errors.Errorf("unable to decrypt data for %s", cs.prefixKey(key))
	}

	return certmagic.KeyInfo{
		Key:        key,
		Modified:   contents.Modified,
		Size:       int64(len(contents.Value)),
		IsTerminal: false,
	}, nil
}

func (cs *ConsulStorage) createConsulClient() error {
	// get the default config
	consulCfg := consul.DefaultConfig()
	if cs.Address != "" {
		consulCfg.Address = cs.Address
	}
	if cs.Token != "" {
		consulCfg.Token = cs.Token
	}
	if cs.TlsEnabled {
		consulCfg.Scheme = "https"
	}
	consulCfg.TLSConfig.InsecureSkipVerify = cs.TlsInsecure

	// set a dial context to prevent default keepalive
	consulCfg.Transport.DialContext = (&net.Dialer{
		Timeout:   time.Duration(cs.Timeout) * time.Second,
		KeepAlive: time.Duration(cs.Timeout) * time.Second,
	}).DialContext

	// create the Consul API client
	consulClient, err := consul.NewClient(consulCfg)
	if err != nil {
		return errors.Wrap(err, "unable to create Consul client")
	}
	if _, err := consulClient.Agent().NodeName(); err != nil {
		return errors.Wrap(err, "unable to ping Consul")
	}

	cs.ConsulClient = consulClient
	return nil
}
