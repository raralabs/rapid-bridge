package security

import (
	"crypto/ed25519"
	"crypto/rsa"
	"rapid-bridge/domain/port"
)

type Security struct {
	Cipher port.EncryptionDecryptionInterface
}

func (s *Security) Encrypt(data []byte, applicationPublicKey *rsa.PublicKey) ([]byte, []byte, []byte, error) {
	return s.Cipher.Encrypt(data, applicationPublicKey)
}

func (s *Security) Decrypt(rsaPrivateKey *rsa.PrivateKey, ciphertext, encryptedAESKey, nonce []byte) ([]byte, error) {
	return s.Cipher.Decrypt(rsaPrivateKey, ciphertext, encryptedAESKey, nonce)
}

func (s *Security) CreateDigitalSignature(ed25519PrivateKey ed25519.PrivateKey, ciphertext, aesKey, nonce []byte) (string, error) {
	return s.Cipher.CreateDigitalSignature(ed25519PrivateKey, ciphertext, aesKey, nonce)
}

func (s *Security) VerifyDigitalSignature(base64EncryptedPayload string, signatureBase64 string, senderPublicKey ed25519.PublicKey) error {
	return s.Cipher.VerifyDigitalSignature(base64EncryptedPayload, signatureBase64, senderPublicKey)
}

func (s *Security) DecodeBase64Encrypted(base64EncryptedPayload string) ([]byte, []byte, []byte, error) {
	return s.Cipher.DecodeBase64Encrypted(base64EncryptedPayload)
}

func (s *Security) CreateBase64Encrypted(ciphertext, encryptedAESKey, nonce []byte) (string, error) {
	return s.Cipher.CreateBase64Encrypted(ciphertext, encryptedAESKey, nonce)
}

func NewSecurity(cipher port.EncryptionDecryptionInterface) *Security {
	return &Security{
		Cipher: cipher,
	}
}
