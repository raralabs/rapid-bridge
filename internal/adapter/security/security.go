package securityadapter

import (
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"rapid-bridge/domain/port"
	hybridcrypto "rapid-bridge/pkg/security/crypto"
)

type HybridCryptography struct {
}

func (a *HybridCryptography) Encrypt(data []byte, applicationRSAPublicKey *rsa.PublicKey) ([]byte, []byte, []byte, error) {
	// Step 1: Generate an ephemeral AES key
	aesKey, err := hybridcrypto.GenerateAESKey()
	if err != nil {
		return nil, nil, nil, err
	}

	// Step 2: Encrypt the payload with AES-GCM
	ciphertext, nonce, err := hybridcrypto.EncryptWithAESGCM(data, aesKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Step 3: Encrypt the AES key with RSA-OAEP
	encryptedAESKey, err := hybridcrypto.EncryptWithRSA(aesKey, applicationRSAPublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	return ciphertext, encryptedAESKey, nonce, nil
}

func (r *HybridCryptography) Decrypt(rsaPrivateKey *rsa.PrivateKey, ciphertext, encryptedAESKey, nonce []byte) ([]byte, error) {

	// Decrypt the AES key using RSA-OAEP
	aesKey, err := hybridcrypto.DecryptWithRSA(encryptedAESKey, rsaPrivateKey)
	if err != nil {
		return nil, err
	}

	// Decrypt the ciphertext using AES-GCM
	plaintext, err := hybridcrypto.DecryptWithAESGCM(ciphertext, nonce, aesKey)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (a *HybridCryptography) CreateBase64Encrypted(ciphertext, encryptedAESKey, nonce []byte) (string, error) {
	base64Ciphertext := base64.StdEncoding.EncodeToString(ciphertext)
	base64EncryptedAESKey := base64.StdEncoding.EncodeToString(encryptedAESKey)
	base64Nonce := base64.StdEncoding.EncodeToString(nonce)
	return fmt.Sprintf("%s-%s-%s", base64Ciphertext, base64EncryptedAESKey, base64Nonce), nil
}

func (r *HybridCryptography) DecodeBase64Encrypted(base64EncryptedPayload string) ([]byte, []byte, []byte, error) {
	parts := hybridcrypto.SplitMessage(base64EncryptedPayload)
	if len(parts) != 3 {
		return nil, nil, nil, errors.New("invalid message format")
	}

	// Decode parts
	ciphertext, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, nil, err
	}
	encryptedAESKey, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, nil, err
	}
	return ciphertext, encryptedAESKey, nonce, nil
}

func (a *HybridCryptography) CreateDigitalSignature(ed25519PrivateKey ed25519.PrivateKey, ciphertext, aesKey, nonce []byte) (string, error) {

	messageToSign := hybridcrypto.CreateMessageToSign(ciphertext, aesKey, nonce)

	signature := hybridcrypto.SignWithEd25519(messageToSign, ed25519PrivateKey)
	base64Signature := base64.StdEncoding.EncodeToString(signature)
	return base64Signature, nil
}

func (a *HybridCryptography) VerifyDigitalSignature(base64EncryptedPayload string, signatureBase64 string, senderPublicKey ed25519.PublicKey) error {
	ciphertext, encryptedAESKey, nonce, err := a.DecodeBase64Encrypted(base64EncryptedPayload)
	if err != nil {
		return fmt.Errorf("failed to decode payload: %v", err)
	}

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	messageToSign := hybridcrypto.CreateMessageToSign(ciphertext, encryptedAESKey, nonce)
	if !ed25519.Verify(senderPublicKey, messageToSign, signature) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

func NewHybridCryptography() port.EncryptionDecryptionInterface {
	return &HybridCryptography{}
}
