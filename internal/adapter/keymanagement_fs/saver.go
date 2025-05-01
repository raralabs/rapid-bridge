package keymanagementfs

import (
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

type FSKeySaver struct{}

func NewFSKeySaver() *FSKeySaver {
	return &FSKeySaver{}
}

func MarshalPrivateKey(privateKey any) (*pem.Block, error) {

	switch privateKey.(type) {
	case *rsa.PrivateKey, ed25519.PrivateKey:
		privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal private key: %w", err)
		}

		privateKeyPEM := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKeyBytes,
		}
		return privateKeyPEM, nil

	default:
		return nil, fmt.Errorf("unknown private key type")
	}

}

func MarshalPublicKey(publicKey any) (*pem.Block, error) {

	switch publicKey.(type) {
	case *rsa.PublicKey, ed25519.PublicKey:
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal public key: %w", err)
		}

		publicKeyPEM := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}

		return publicKeyPEM, nil

	default:
		return nil, fmt.Errorf("unknown public key type")
	}
}

func (s *FSKeySaver) SaveToFile(filePath string, pemBlock *pem.Block) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, pemBlock); err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}

	return nil
}

func (s *FSKeySaver) SaveRSAPrivateKeyToPEM(privateKey *rsa.PrivateKey, filePath string) error {
	privateKeyPEM, err := MarshalPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	err = s.SaveToFile(filePath, privateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to save private key to file: %w", err)
	}

	return nil
}

func (s *FSKeySaver) SaveRSAPublicKeyToPEM(publicKey *rsa.PublicKey, filePath string) error {
	publicKeyPEM, err := MarshalPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	err = s.SaveToFile(filePath, publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to save public key to file: %w", err)
	}

	return nil
}

func (s *FSKeySaver) SaveEd25519PrivateKeyToPEM(privateKey ed25519.PrivateKey, filePath string) error {
	privateKeyPEM, err := MarshalPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	err = s.SaveToFile(filePath, privateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to save private key to file: %w", err)
	}

	return nil
}

func (s *FSKeySaver) SaveEd25519PublicKeyToPEM(publicKey ed25519.PublicKey, filePath string) error {

	publicKeyPEM, err := MarshalPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	err = s.SaveToFile(filePath, publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to save public key to file: %w", err)
	}

	return nil
}
