package cache

import (
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"rapid-bridge/domain/port"
)

type KeyCacheAdapter struct {
	// Parsed keys
	rsaPrivateKeys     map[string]*rsa.PrivateKey
	rsaPublicKeys      map[string]*rsa.PublicKey
	ed25519PrivateKeys map[string]ed25519.PrivateKey
	ed25519PublicKeys  map[string]ed25519.PublicKey

	// Raw []byte keys
	rawPrivateKeys map[string][]byte
	rawPublicKeys  map[string][]byte
}

func NewKeyCacheAdapter() port.KeyCache {
	return &KeyCacheAdapter{
		rsaPrivateKeys:     make(map[string]*rsa.PrivateKey),
		rsaPublicKeys:      make(map[string]*rsa.PublicKey),
		ed25519PrivateKeys: make(map[string]ed25519.PrivateKey),
		ed25519PublicKeys:  make(map[string]ed25519.PublicKey),
		rawPrivateKeys:     make(map[string][]byte),
		rawPublicKeys:      make(map[string][]byte),
	}
}

func generateKeyID(source, keyVersion string, isApp bool) string {
	if isApp {
		return fmt.Sprintf("%s:%s", source, keyVersion)
	}
	return source
}

// Set and Get Functions for Parsed Keys

// Set RSA Private Key

func (c *KeyCacheAdapter) SetRSAPrivateKey(source, keyVersion string, isApp bool, key *rsa.PrivateKey) {
	c.rsaPrivateKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get RSA Private Key

func (c *KeyCacheAdapter) GetRSAPrivateKey(source, keyVersion string, isApp bool) (*rsa.PrivateKey, error) {
	key, ok := c.rsaPrivateKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil, fmt.Errorf("RSA private key not found for %s", generateKeyID(source, keyVersion, isApp))
	}
	return key, nil
}

// Set RSA Public Key

func (c *KeyCacheAdapter) SetRSAPublicKey(source, keyVersion string, isApp bool, key *rsa.PublicKey) {
	c.rsaPublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get RSA Public Key

func (c *KeyCacheAdapter) GetRSAPublicKey(source, keyVersion string, isApp bool) (*rsa.PublicKey, error) {
	key, ok := c.rsaPublicKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil, fmt.Errorf("RSA public key not found for %s", generateKeyID(source, keyVersion, isApp))
	}
	return key, nil
}

// Set ED25519 Private Key

func (c *KeyCacheAdapter) SetEd25519PrivateKey(source, keyVersion string, isApp bool, key ed25519.PrivateKey) {
	c.ed25519PrivateKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get ED25519 Private Key

func (c *KeyCacheAdapter) GetEd25519PrivateKey(source, keyVersion string, isApp bool) (ed25519.PrivateKey, error) {
	key, ok := c.ed25519PrivateKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil, fmt.Errorf("Ed25519 private key not found for %s", generateKeyID(source, keyVersion, isApp))
	}
	return key, nil
}

// Set ED25519 Public Key

func (c *KeyCacheAdapter) SetEd25519PublicKey(source, keyVersion string, isApp bool, key ed25519.PublicKey) {
	c.ed25519PublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get ED25519 Public Key

func (c *KeyCacheAdapter) GetEd25519PublicKey(source, keyVersion string, isApp bool) (ed25519.PublicKey, error) {
	key, ok := c.ed25519PublicKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil, fmt.Errorf("Ed25519 public key not found for %s", generateKeyID(source, keyVersion, isApp))
	}
	return key, nil
}

// Set and Get Functions for Raw Keys

// Set Raw Private Key

func (c *KeyCacheAdapter) SetRawPrivateKey(source, keyVersion string, isApp bool, key []byte) {
	c.rawPrivateKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get Raw Public Key

func (c *KeyCacheAdapter) GetRawPrivateKey(source, keyVersion string, isApp bool) ([]byte, error) {
	key, ok := c.rawPrivateKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil, fmt.Errorf("Raw private key not found for %s", generateKeyID(source, keyVersion, isApp))
	}
	return key, nil
}

// Set Raw Public Key

func (c *KeyCacheAdapter) SetRawPublicKey(source, keyVersion string, isApp bool, key []byte) {
	c.rawPublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get Raw Public Key

func (c *KeyCacheAdapter) GetRawPublicKey(source, keyVersion string, isApp bool) ([]byte, error) {
	key, ok := c.rawPublicKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil, fmt.Errorf("Raw public key not found for %s", generateKeyID(source, keyVersion, isApp))
	}
	return key, nil
}
