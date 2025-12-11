package cache

import (
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
)

type KeyCache struct {
	// Parsed keys
	rsaPrivateKeys     map[string]*rsa.PrivateKey
	rsaPublicKeys      map[string]*rsa.PublicKey
	ed25519PrivateKeys map[string]ed25519.PrivateKey
	ed25519PublicKeys  map[string]ed25519.PublicKey

	// Raw []byte keys (from filesystem)
	rawPrivateKeys map[string][]byte
	rawPublicKeys  map[string][]byte
}

func NewKeyCache() *KeyCache {
	return &KeyCache{
		rsaPrivateKeys:     make(map[string]*rsa.PrivateKey),
		rsaPublicKeys:      make(map[string]*rsa.PublicKey),
		ed25519PrivateKeys: make(map[string]ed25519.PrivateKey),
		ed25519PublicKeys:  make(map[string]ed25519.PublicKey),

		rawPrivateKeys: make(map[string][]byte),
		rawPublicKeys:  make(map[string][]byte),
	}
}

func generateKeyID(source, keyVersion string, isApp bool) string {
	if isApp {
		return fmt.Sprintf("%s:%s", source, keyVersion)
	}
	return source
}

// Cache Value Setter Functions

// Parsed Keys

func (c *KeyCache) SetRSAPrivateKey(source, keyVersion string, isApp bool, key *rsa.PrivateKey) {
	c.rsaPrivateKeys[generateKeyID(source, keyVersion, isApp)] = key
}

func (c *KeyCache) SetRSAPublicKey(source, keyVersion string, isApp bool, key *rsa.PublicKey) {
	c.rsaPublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

func (c *KeyCache) SetEd25519PrivateKey(source, keyVersion string, isApp bool, key ed25519.PrivateKey) {
	c.ed25519PrivateKeys[generateKeyID(source, keyVersion, isApp)] = key
}

func (c *KeyCache) SetEd25519PublicKey(source, keyVersion string, isApp bool, key ed25519.PublicKey) {
	c.ed25519PublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Raw Keys []byte

func (c *KeyCache) SetRawPrivateKey(source, keyVersion string, isApp bool, key []byte) {
	c.rawPrivateKeys[generateKeyID(source, keyVersion, isApp)] = key
}

func (c *KeyCache) SetRawPublicKey(source, keyVersion string, isApp bool, key []byte) {
	c.rawPublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Cache Value Get Functions

// Parsed Keys

func (c *KeyCache) GetRSAPrivateKey(source, keyVersion string, isApp bool) (*rsa.PrivateKey, error) {
	keyID := generateKeyID(source, keyVersion, isApp)
	key, exists := c.rsaPrivateKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("RSA private key not found for %s", keyID)
	}
	return key, nil
}

func (c *KeyCache) GetRSAPublicKey(source, keyVersion string, isApp bool) (*rsa.PublicKey, error) {
	keyID := generateKeyID(source, keyVersion, isApp)
	key, exists := c.rsaPublicKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("RSA public key not found for %s", keyID)
	}
	return key, nil
}

func (c *KeyCache) GetEd25519PrivateKey(source, keyVersion string, isApp bool) (ed25519.PrivateKey, error) {
	keyID := generateKeyID(source, keyVersion, isApp)
	key, exists := c.ed25519PrivateKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("ed25519 private key not found for %s", keyID)
	}
	return key, nil
}

func (c *KeyCache) GetEd25519PublicKey(source, keyVersion string, isApp bool) (ed25519.PublicKey, error) {
	keyID := generateKeyID(source, keyVersion, isApp)
	key, exists := c.ed25519PublicKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("ed25519 public key not found for %s", keyID)
	}
	return key, nil
}

// Raw Keys []byte

func (c *KeyCache) GetRawPrivateKey(source, keyVersion string, isApp bool) ([]byte, error) {
	keyID := generateKeyID(source, keyVersion, isApp)
	key, exists := c.rawPrivateKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("raw private key not found for %s", keyID)
	}
	return key, nil
}

func (c *KeyCache) GetRawPublicKey(source, keyVersion string, isApp bool) ([]byte, error) {
	keyID := generateKeyID(source, keyVersion, isApp)
	key, exists := c.rawPublicKeys[keyID]
	if !exists {
		return nil, fmt.Errorf("raw public key not found for %s", keyID)
	}
	return key, nil
}
