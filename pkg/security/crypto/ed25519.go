package hybridcrypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

func GenerateEd25519KeyPair() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key pair: %v", err)
	}
	return privateKey, publicKey, nil
}

func SignWithEd25519(data []byte, ed25519PrivateKey ed25519.PrivateKey) []byte {
	return ed25519.Sign(ed25519PrivateKey, data)
}
