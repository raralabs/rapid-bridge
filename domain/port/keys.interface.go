package port

import (
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/pem"
)

type KeyLoader interface {
	LoadPrivateKey(privateKeyPath string) (any, error)
	LoadPublicKey(publicKeyPath string) (any, error)
}

type KeyConverter interface {
	ConvertPublicKeyToBase64(publicKey any) (string, error)
	ConvertBase64ToPublicKey(encodedPublicKey string) (any, error)
}

type KeySaver interface {
	SaveToFile(filePath string, pemBlock *pem.Block) error
	SaveRSAPrivateKeyToPEM(privateKey *rsa.PrivateKey, filePath string) error
	SaveRSAPublicKeyToPEM(publicKey *rsa.PublicKey, filePath string) error
	SaveEd25519PrivateKeyToPEM(privateKey ed25519.PrivateKey, filePath string) error
	SaveEd25519PublicKeyToPEM(publicKey ed25519.PublicKey, filePath string) error
}
