package cache

import (
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"rapid-bridge/domain/port"
	"rapid-bridge/pkg/util"

	"go.uber.org/zap"
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

	loader port.KeyLoader
	logger port.Logger
}

func NewKeyCacheAdapter(keyLoader port.KeyLoader, logger port.Logger) port.KeyCache {
	return &KeyCacheAdapter{
		rsaPrivateKeys:       make(map[string]*rsa.PrivateKey),
		rsaPublicKeys:        make(map[string]*rsa.PublicKey),
		ed25519PrivateKeys:   make(map[string]ed25519.PrivateKey),
		ed25519PublicKeys:    make(map[string]ed25519.PublicKey),
		rawRSAPublicKeys:     make(map[string][]byte),
		rawED25519PublicKeys: make(map[string][]byte),

		loader: keyLoader,
		logger: logger,
	}
}

var (
	rsaPrivateKey     *rsa.PrivateKey
	bankRsaPublicKey  *rsa.PublicKey
	ed25519PrivateKey ed25519.PrivateKey
	bankEdPublicKey   ed25519.PublicKey

	rsaPublicKeyBytes []byte
	edPublicKeyBytes  []byte

	tempKey any

	err error
)

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
	if ok {
		return key
	}

	rsaPrivateKeyPath := util.GetRSAPrivateKeyPath(source, keyVersion)
	tempKey, err = c.loader.LoadPrivateKey(rsaPrivateKeyPath)
	if err != nil {
		c.logger.Error("Failed to load Application's RSA Private Key from file", zap.String("error", err.Error()))
		return nil
	}

	// Type assertion from any to *rsa.PrivateKey
	rsaPrivateKey, ok = tempKey.(*rsa.PrivateKey)
	if !ok {
		c.logger.Error("Failed to convert datatype from any to *rsa.Private")
		return nil
	}

	c.SetRSAPrivateKey(source, keyVersion, true, rsaPrivateKey)

	return rsaPrivateKey
}

// Set RSA Public Key

func (c *KeyCacheAdapter) SetRSAPublicKey(source, keyVersion string, isApp bool, key *rsa.PublicKey) {
	c.rsaPublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get RSA Public Key

func (c *KeyCacheAdapter) GetRSAPublicKey(source, keyVersion string, isApp bool) *rsa.PublicKey {
	key, ok := c.rsaPublicKeys[generateKeyID(source, keyVersion, isApp)]
	if ok {
		return key
	}

	bankRsaPublicKeyPath := util.GetBankRSAPublicKeyPath(source)
	tempKey, err = c.loader.LoadPublicKey(bankRsaPublicKeyPath)
	if err != nil {
		c.logger.Error("Failed to load Bank's RSA Public Key from file", zap.String("error", err.Error()))
		return nil
	}

	// Type assertion from any to *rsa.PublicKey
	bankRsaPublicKey, ok = tempKey.(*rsa.PublicKey)
	if !ok {
		c.logger.Error("Failed to convert datatype from any to *rsa.PublicKey")
		return nil
	}

	c.SetRSAPublicKey(source, keyVersion, false, bankRsaPublicKey)

	return bankRsaPublicKey
}

// Set ED25519 Private Key

func (c *KeyCacheAdapter) SetEd25519PrivateKey(source, keyVersion string, isApp bool, key ed25519.PrivateKey) {
	c.ed25519PrivateKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get ED25519 Private Key

func (c *KeyCacheAdapter) GetEd25519PrivateKey(source, keyVersion string, isApp bool) ed25519.PrivateKey {
	key, ok := c.ed25519PrivateKeys[generateKeyID(source, keyVersion, isApp)]
	if ok {
		return key
	}

	edPrivateKeyPath := util.GetEd25519PrivateKeyPath(source, keyVersion)
	tempKey, err = c.loader.LoadPrivateKey(edPrivateKeyPath)
	if err != nil {
		c.logger.Error("Failed to load Application's ED25519 Private Key from file", zap.String("error", err.Error()))
		return nil
	}

	// Type assertion from any to ed25519.PrivateKey
	ed25519PrivateKey, ok = tempKey.(ed25519.PrivateKey)
	if !ok {
		c.logger.Error("Failed to convert datatype from any to *rsa.Private")
		return nil
	}

	c.SetEd25519PrivateKey(source, keyVersion, true, ed25519PrivateKey)

	return ed25519PrivateKey
}

// Set ED25519 Public Key

func (c *KeyCacheAdapter) SetEd25519PublicKey(source, keyVersion string, isApp bool, key ed25519.PublicKey) {
	c.ed25519PublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get ED25519 Public Key

func (c *KeyCacheAdapter) GetEd25519PublicKey(source, keyVersion string, isApp bool) ed25519.PublicKey {
	key, ok := c.ed25519PublicKeys[generateKeyID(source, keyVersion, isApp)]
	if ok {
		return key
	}

	bankEdPublicKeyPath := util.GetBankEd25519PublicKeyPath(source)
	tempKey, err = c.loader.LoadPublicKey(bankEdPublicKeyPath)
	if err != nil {
		c.logger.Error("Failed to load Bank's ED25519 Public Key from file", zap.String("error", err.Error()))
		return nil
	}

	// Type assertion from any to ed25519.PublicKey
	bankEdPublicKey, ok = tempKey.(ed25519.PublicKey)
	if !ok {
		c.logger.Error("Failed to convert datatype from any to ed25519.PublicKey")
		return nil
	}

	c.SetEd25519PublicKey(source, keyVersion, false, bankEdPublicKey)

	return bankEdPublicKey
}

// Set and Get Functions for Raw RSA Keys

// Set Raw RSA Public Key

func (c *KeyCacheAdapter) SetRawRSAPublicKey(source, keyVersion string, isApp bool, key []byte) {
	c.rawRSAPublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get Raw RSA Public Key

func (c *KeyCacheAdapter) GetRawRSAPublicKey(source, keyPath string, keyVersion string, isApp bool) []byte {
	key, ok := c.rawRSAPublicKeys[generateKeyID(source, keyVersion, isApp)]
	if ok {
		return key
	}

	rsaPublicKeyBytes, err = util.ReadFile(keyPath)
	if err != nil {
		c.logger.Error("Failed to read RSA public keys", zap.String("error", err.Error()))
		return nil
	}

	c.SetRawRSAPublicKey(source, keyVersion, true, rsaPublicKeyBytes)

	return rsaPublicKeyBytes
}

// Set Raw EDD25519 Public Key

func (c *KeyCacheAdapter) SetRawED25519PublicKey(source, keyVersion string, isApp bool, key []byte) {
	c.rawED25519PublicKeys[generateKeyID(source, keyVersion, isApp)] = key
}

// Get Raw ED25519 Public Key

func (c *KeyCacheAdapter) GetRawED25519PublicKey(source, keyPath, keyVersion string, isApp bool) []byte {
	key, ok := c.rawED25519PublicKeys[generateKeyID(source, keyVersion, isApp)]
	if ok {
		return key
	}

	edPublicKeyBytes, err = util.ReadFile(keyPath)
	if err != nil {
		c.logger.Error("Failed to read RSA public keys", zap.String("error", err.Error()))
		return nil
	}

	c.SetRawED25519PublicKey(source, keyVersion, true, edPublicKeyBytes)

	return edPublicKeyBytes
}
