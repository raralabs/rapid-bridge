package keymanagementfs

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

type FSKeyLoader struct{}

func NewFSKeyLoader() *FSKeyLoader {
	return &FSKeyLoader{}
}

func (l *FSKeyLoader) LoadPrivateKey(privateKeyPath string) (any, error) {

	_, err := filepath.Abs(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	privateKeyBytes, err := os.ReadFile(privateKeyPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privateKey, nil
}

func (l *FSKeyLoader) LoadPublicKey(publicKeyPath string) (any, error) {

	_, err := filepath.Abs(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %v", err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyBytes)
	if publicKeyBlock == nil {
		return nil, fmt.Errorf("failed to decode public key PEM")
	}
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return publicKey, nil
}
