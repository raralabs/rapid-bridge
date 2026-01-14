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
	rawRSAPublicKeys     map[string][]byte
	rawED25519PublicKeys map[string][]byte
}

func NewKeyCacheAdapter() port.KeyCache {
	return &KeyCacheAdapter{
		rsaPrivateKeys:       make(map[string]*rsa.PrivateKey),
		rsaPublicKeys:        make(map[string]*rsa.PublicKey),
		ed25519PrivateKeys:   make(map[string]ed25519.PrivateKey),
		ed25519PublicKeys:    make(map[string]ed25519.PublicKey),
		rawRSAPublicKeys:     make(map[string][]byte),
		rawED25519PublicKeys: make(map[string][]byte),
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

func (c *KeyCacheAdapter) GetRSAPrivateKey(source, keyVersion string, isApp bool) *rsa.PrivateKey {
	key, ok := c.rsaPrivateKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil
	}
	return key
}

// Set RSA Public Key

func (c *KeyCacheAdapter) SetRSAPublicKey(source, keyVersion string, isApp bool, key *rsa.PublicKey) {
	c.rsaPublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get RSA Public Key

func (c *KeyCacheAdapter) GetRSAPublicKey(source, keyVersion string, isApp bool) *rsa.PublicKey {
	key, ok := c.rsaPublicKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil
	}
	return key
}

// Set ED25519 Private Key

func (c *KeyCacheAdapter) SetEd25519PrivateKey(source, keyVersion string, isApp bool, key ed25519.PrivateKey) {
	c.ed25519PrivateKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get ED25519 Private Key

func (c *KeyCacheAdapter) GetEd25519PrivateKey(source, keyVersion string, isApp bool) ed25519.PrivateKey {
	key, ok := c.ed25519PrivateKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil
	}
	return key
}

// Set ED25519 Public Key

func (c *KeyCacheAdapter) SetEd25519PublicKey(source, keyVersion string, isApp bool, key ed25519.PublicKey) {
	c.ed25519PublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get ED25519 Public Key

func (c *KeyCacheAdapter) GetEd25519PublicKey(source, keyVersion string, isApp bool) ed25519.PublicKey {
	key, ok := c.ed25519PublicKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil
	}
	return key
}

// Set and Get Functions for Raw RSA Keys

// Set Raw RSA Public Key

func (c *KeyCacheAdapter) SetRawRSAPublicKey(source, keyVersion string, isApp bool, key []byte) {
	c.rawRSAPublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get Raw RSA Public Key

func (c *KeyCacheAdapter) GetRawRSAPublicKey(source, keyVersion string, isApp bool) []byte {
	key, ok := c.rawRSAPublicKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil
	}
	return key
}

// Set Raw EDD25519 Public Key

func (c *KeyCacheAdapter) SetRawED25519PublicKey(source, keyVersion string, isApp bool, key []byte) {
	c.rawED25519PublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get Raw ED25519 Public Key

func (c *KeyCacheAdapter) GetRawED25519PublicKey(source, keyVersion string, isApp bool) []byte {
	key, ok := c.rawED25519PublicKeys[generateKeyID(source, keyVersion, isApp)]
	if !ok {
		return nil
	}
	return key
}
