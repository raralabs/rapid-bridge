package rapid

type RapidResourceRequest struct {
	From       string `json:"from" validate:"required"`
	To         string `json:"to" validate:"required"`
	Message    string `json:"message" validate:"required"` // Format: base64(ciphertext)-base64(encryptedAESKey)-base64(nonce)
	Signature  string `json:"signature" validate:"required"`
	KeyVersion string `json:"key_version" validate:"required"`
}
