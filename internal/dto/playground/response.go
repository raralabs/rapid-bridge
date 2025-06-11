package playground

import (
	"crypto/ed25519"
	"crypto/rsa"
)

type ApplicationRegisterResponse struct {
	KeyVersion string `json:"key_version"`
	Slug       string `json:"slug"`

	RSAPublicKey     *rsa.PublicKey
	Ed25519PublicKey ed25519.PublicKey
	Message          string `json:"message"`
}
