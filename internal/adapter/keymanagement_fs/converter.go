package keymanagementfs

import (
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

type FSKeyConverter struct{}

func NewFSKeyConverter() *FSKeyConverter {
	return &FSKeyConverter{}
}

func (l *FSKeyConverter) ConvertPublicKeyToBase64(publicKey any) (string, error) {
	switch publicKeyType := publicKey.(type) {
	case *rsa.PublicKey, ed25519.PublicKey:
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKeyType)
		if err != nil {
			return "", fmt.Errorf("failed to marshal public key: %v", err)
		}
		publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)
		return publicKeyBase64, nil
	default:
		return "", fmt.Errorf("unsupported public key type: %T", publicKey)
	}
}

func (l *FSKeyConverter) ConvertBase64ToPublicKey(encodedPublicKey string) (any, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(encodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 public key: %w", err)
	}

	pubKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	switch key := pubKey.(type) {
	case *rsa.PublicKey:
		return key, nil
	case ed25519.PublicKey:
		return key, nil
	default:
		return nil, fmt.Errorf("unsupported public key type")
	}
}
