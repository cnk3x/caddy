package storageconsul

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"

	"github.com/pteich/errors"
)

func (cs *ConsulStorage) encrypt(bytes []byte) ([]byte, error) {
	// No key? No encrypt
	if len(cs.AESKey) == 0 {
		return bytes, nil
	}

	c, err := aes.NewCipher(cs.AESKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create AES cipher")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create GCM cipher")
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate nonce")
	}

	return gcm.Seal(nonce, nonce, bytes, nil), nil
}

func (cs *ConsulStorage) EncryptStorageData(data *StorageData) ([]byte, error) {
	// JSON marshal, then encrypt if key is there
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal")
	}

	// Prefix with simple prefix and then encrypt
	bytes = append([]byte(cs.ValuePrefix), bytes...)
	return cs.encrypt(bytes)
}

func (cs *ConsulStorage) decrypt(bytes []byte) ([]byte, error) {
	// No key? No decrypt
	if len(cs.AESKey) == 0 {
		return bytes, nil
	}
	if len(bytes) < aes.BlockSize {
		return nil, errors.New("invalid contents")
	}

	block, err := aes.NewCipher(cs.AESKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create AES cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create GCM cipher")
	}

	out, err := gcm.Open(nil, bytes[:gcm.NonceSize()], bytes[gcm.NonceSize():], nil)
	if err != nil {
		return nil, errors.Wrap(err, "decryption failure")
	}

	return out, nil
}

func (cs *ConsulStorage) DecryptStorageData(bytes []byte) (*StorageData, error) {
	// We have to decrypt if there is an AES key and then JSON unmarshal
	bytes, err := cs.decrypt(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decrypt data")
	}

	// Simple sanity check of the beginning of the byte array just to check
	if len(bytes) < len(cs.ValuePrefix) || string(bytes[:len(cs.ValuePrefix)]) != cs.ValuePrefix {
		return nil, errors.New("invalid data format")
	}

	// Now just json unmarshal
	data := &StorageData{}
	if err := json.Unmarshal(bytes[len(cs.ValuePrefix):], data); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal result")
	}
	return data, nil
}
