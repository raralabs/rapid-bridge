package port

import (
	"crypto/ed25519"
	"crypto/rsa"
)

type EncryptionDecryptionInterface interface {
	Encrypt(data []byte, applicationRSAPublicKey *rsa.PublicKey) ([]byte, []byte, []byte, error)
	Decrypt(rsaPrivateKey *rsa.PrivateKey, ciphertext, encryptedAESKey, nonce []byte) ([]byte, error)
	CreateDigitalSignature(ed25519PrivateKey ed25519.PrivateKey, ciphertext, aesKey, nonce []byte) (string, error)
	VerifyDigitalSignature(base64EncryptedPayload string, signatureBase64 string, senderPublicKey ed25519.PublicKey) error
	DecodeBase64Encrypted(base64EncryptedPayload string) ([]byte, []byte, []byte, error)
	CreateBase64Encrypted(ciphertext, encryptedAESKey, nonce []byte) (string, error)
}
