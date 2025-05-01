package hybridcrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

func EncryptWithRSA(data []byte, applicationRSAPublicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, applicationRSAPublicKey, data, nil)
}

func DecryptWithRSA(encryptedAESKey []byte, applicationRSAPrivateKey *rsa.PrivateKey) ([]byte, error) {
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, applicationRSAPrivateKey, encryptedAESKey, nil)
	if err != nil {
		return nil, err
	}

	return aesKey, nil
}
