package keys

import (
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"rapid-bridge/constants"
)

func ReadAndValidateKeyFile(path string, isPrivate bool) (any, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []byte{}, fmt.Errorf("file does not exist: %s", path)
	}

	ext := filepath.Ext(path)
	if ext != ".pem" {
		return []byte{}, fmt.Errorf("unsupported file extension: %s", ext)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return []byte{}, errors.New("failed to decode PEM block")
	}

	if isPrivate {
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to parse private key: %w", err)
		}
		return privateKey, validatePrivateKey(privateKey)
	} else {
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to parse public key: %w", err)
		}
		return publicKey, validatePublicKey(publicKey)
	}
}

func validatePrivateKey(privateKey any) error {

	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		if k.N.BitLen() != constants.RSAKeyBitSize {
			return fmt.Errorf("RSA key too short: %d bits", k.N.BitLen())
		}
		return nil
	case ed25519.PrivateKey:
		if len(k) != ed25519.PrivateKeySize {
			return fmt.Errorf("invalid ed25519 key size: %d", len(k))
		}
		return nil
	default:
		return fmt.Errorf("unsupported private key type")
	}
}

func validatePublicKey(publicKey any) error {

	switch key := publicKey.(type) {
	case *rsa.PublicKey:
		if key.N.BitLen() != constants.RSAKeyBitSize {
			return fmt.Errorf("RSA public key too short: %d bits", key.N.BitLen())
		}
		return nil
	case ed25519.PublicKey:
		if len(key) != ed25519.PublicKeySize {
			return fmt.Errorf("invalid ed25519 public key size: %d", len(key))
		}
		return nil
	default:
		return fmt.Errorf("unsupported public key type: %T", key)
	}
}
