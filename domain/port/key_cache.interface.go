package port

import (
	"crypto/ed25519"
	"crypto/rsa"
)

type KeyCache interface {

	// Parsed keys

	SetRSAPrivateKey(source, keyVersion string, isApp bool, key *rsa.PrivateKey)
	GetRSAPrivateKey(source, keyVersion string, isApp bool) (*rsa.PrivateKey, error)

	SetRSAPublicKey(source, keyVersion string, isApp bool, key *rsa.PublicKey)
	GetRSAPublicKey(source, keyVersion string, isApp bool) (*rsa.PublicKey, error)

	SetEd25519PrivateKey(source, keyVersion string, isApp bool, key ed25519.PrivateKey)
	GetEd25519PrivateKey(source, keyVersion string, isApp bool) (ed25519.PrivateKey, error)

	SetEd25519PublicKey(source, keyVersion string, isApp bool, key ed25519.PublicKey)
	GetEd25519PublicKey(source, keyVersion string, isApp bool) (ed25519.PublicKey, error)

	// Raw keys

	SetRawPrivateKey(source, keyVersion string, isApp bool, key []byte)
	GetRawPrivateKey(source, keyVersion string, isApp bool) ([]byte, error)

	SetRawPublicKey(source, keyVersion string, isApp bool, key []byte)
	GetRawPublicKey(source, keyVersion string, isApp bool) ([]byte, error)
}
