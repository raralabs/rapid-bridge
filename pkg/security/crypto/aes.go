package hybridcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key for AES-256
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func EncryptWithAESGCM(data []byte, aesKey []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	return ciphertext, nonce, nil
}

func DecryptWithAESGCM(ciphertext, nonce, aesKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
