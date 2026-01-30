package cache

import (
	"crypto/ed25519"
	"crypto/rsa"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/dto/application"
	"rapid-bridge/pkg/util"

	"go.uber.org/zap"
)

type KeyCacheAdapter struct {
	applicationCache map[string]*application.ApplicationKeys
	bankCache        map[string]*application.BankKeys

	loader port.KeyLoader
	logger port.Logger
}

func NewKeyCacheAdapter(keyLoader port.KeyLoader, logger port.Logger) port.KeyCache {
	return &KeyCacheAdapter{
		applicationCache: make(map[string]*application.ApplicationKeys),
		bankCache:        make(map[string]*application.BankKeys),

		loader: keyLoader,
		logger: logger,
	}
}

// Single Get Keys Function

func (c *KeyCacheAdapter) GetKeys(source, destination, keyVersion string) *application.Keys {
	var applicationKey *application.ApplicationKeys
	var ok bool
	// Get or load ApplicationKeys
	applicationKey, ok = c.applicationCache[source]
	if applicationKey == nil {
		c.applicationCache[source] = &application.ApplicationKeys{}
	}

	if c.applicationCache[source].RSAPrivateKey == nil || c.applicationCache[source].EDPrivateKey == nil {
		// RSA Private Key (Application)
		rsaPrivateKeyPath := util.GetRSAPrivateKeyPath(source, keyVersion)
		tempKey, err := c.loader.LoadPrivateKey(rsaPrivateKeyPath)
		if err != nil {
			c.logger.Error("Failed to load Application's RSA Private Key from file", zap.String("error", err.Error()))
			return nil
		}
		// Type assertion from any to *rsa.PrivateKey
		rsaPrivateKey, ok := tempKey.(*rsa.PrivateKey)
		if !ok {
			c.logger.Error("Failed to convert datatype from any to *rsa.Private")
			return nil
		}

		// ED25519 Private Key (Application)
		edPrivateKeyPath := util.GetEd25519PrivateKeyPath(source, keyVersion)
		tempKey, err = c.loader.LoadPrivateKey(edPrivateKeyPath)
		if err != nil {
			c.logger.Error("Failed to load Application's ED25519 Private Key from file", zap.String("error", err.Error()))
			return nil
		}
		// Type assertion from any to ed25519.PrivateKey
		ed25519PrivateKey, ok := tempKey.(ed25519.PrivateKey)
		if !ok {
			c.logger.Error("Failed to convert datatype from any to *rsa.Private")
			return nil
		}

		// Setting RSA and ED Private Keys in Application Cache
		c.applicationCache[source].RSAPrivateKey = rsaPrivateKey
		c.applicationCache[source].EDPrivateKey = ed25519PrivateKey

		applicationKey = &application.ApplicationKeys{
			RSAPrivateKey: rsaPrivateKey,
			EDPrivateKey:  ed25519PrivateKey,
		}
	}

	// Get or load BankKeys
	bankKey, ok := c.bankCache[destination]
	if bankKey == nil {
		c.bankCache[destination] = &application.BankKeys{}
	}

	if !ok {
		// RSA Public Key (Bank)
		bankRsaPublicKeyPath := util.GetBankRSAPublicKeyPath(destination)
		tempKey, err := c.loader.LoadPublicKey(bankRsaPublicKeyPath)
		if err != nil {
			c.logger.Error("Failed to load Bank's RSA Public Key from file", zap.String("error", err.Error()))
			return nil
		}
		// Type assertion from any to *rsa.PublicKey
		bankRsaPublicKey, ok := tempKey.(*rsa.PublicKey)
		if !ok {
			c.logger.Error("Failed to convert datatype from any to *rsa.PublicKey")
			return nil
		}

		// ED25519 Public Key (Bank)
		bankEdPublicKeyPath := util.GetBankEd25519PublicKeyPath(destination)
		tempKey, err = c.loader.LoadPublicKey(bankEdPublicKeyPath)
		if err != nil {
			c.logger.Error("Failed to load Bank's ED25519 Public Key from file", zap.String("error", err.Error()))
			return nil
		}
		// Type assertion from any to ed25519.PublicKey
		bankEdPublicKey, ok := tempKey.(ed25519.PublicKey)
		if !ok {
			c.logger.Error("Failed to convert datatype from any to ed25519.PublicKey")
			return nil
		}

		// Setting RSA and ED Public Keys in Bank Cache
		c.bankCache[destination].RSAPublicKey = bankRsaPublicKey
		c.bankCache[destination].EDPublicKey = bankEdPublicKey

		bankKey = &application.BankKeys{
			RSAPublicKey: bankRsaPublicKey,
			EDPublicKey:  bankEdPublicKey,
		}
	}

	// Return Application Private Keys and Bank Public Keys
	return &application.Keys{
		ApplicationKeys: applicationKey,
		BankKeys:        bankKey,
	}
}

// Single Get Raw Keys Function

func (c *KeyCacheAdapter) GetRawApplicationKeys(source, rsaPublicKeyPath string, edPublicKeyPath string, keyVersion string, isApp bool) *application.ApplicationKeys {
	var appKeys *application.ApplicationKeys

	// Get or load ApplicationKeys
	appKeys, _ = c.applicationCache[source]
	if appKeys == nil {
		c.applicationCache[source] = &application.ApplicationKeys{}
	}

	if appKeys.RSAPublicKey == nil || appKeys.EDPublicKey == nil {
		rsaPublicKeyBytes, err := util.ReadFile(rsaPublicKeyPath)
		if err != nil {
			c.logger.Error("Failed to read RSA public keys", zap.String("error", err.Error()))
			return nil
		}

		edPublicKeyBytes, err := util.ReadFile(edPublicKeyPath)
		if err != nil {
			c.logger.Error("Failed to read ED25519 public keys", zap.String("error", err.Error()))
			return nil
		}

		// Setting RSA and ED Public Keys in Application Cache
		c.applicationCache[source].RSAPublicKey = rsaPublicKeyBytes
		c.applicationCache[source].EDPublicKey = edPublicKeyBytes

		// Return Application Private Keys
		appKeys = &application.ApplicationKeys{
			RSAPublicKey: rsaPublicKeyBytes,
			EDPublicKey:  edPublicKeyBytes,
		}

	}

	return appKeys
}
