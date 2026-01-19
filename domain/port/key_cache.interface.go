package port

import (
	"crypto/ed25519"
	"crypto/rsa"
)

type KeyCache interface {

	// Parsed keys

	SetRSAPrivateKey(source, keyVersion string, isApp bool, key *rsa.PrivateKey)
	GetRSAPrivateKey(source, keyVersion string, isApp bool) *rsa.PrivateKey

	SetRSAPublicKey(source, keyVersion string, isApp bool, key *rsa.PublicKey)
	GetRSAPublicKey(source, keyVersion string, isApp bool) *rsa.PublicKey

	SetEd25519PrivateKey(source, keyVersion string, isApp bool, key ed25519.PrivateKey)
	GetEd25519PrivateKey(source, keyVersion string, isApp bool) ed25519.PrivateKey

	SetEd25519PublicKey(source, keyVersion string, isApp bool, key ed25519.PublicKey)
	GetEd25519PublicKey(source, keyVersion string, isApp bool) ed25519.PublicKey

	// Raw keys

	SetRawRSAPublicKey(source, keyVersion string, isApp bool, key []byte)
	GetRawRSAPublicKey(source, keyPath, keyVersion string, isApp bool) []byte

	SetRawED25519PublicKey(source, keyVersion string, isApp bool, key []byte)
	GetRawED25519PublicKey(source, keyPath, keyVersion string, isApp bool) []byte
}
