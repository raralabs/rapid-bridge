package application

import (
	"crypto/ed25519"
	"crypto/rsa"
)

type ResourceResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type OtherResponse struct {
	Data     interface{} `json:"data"`
	MetaData interface{} `json:"meta_data,omitempty"`
}

type ApplicationKeys struct {
	// Private Key Objects
	RSAPrivateKey *rsa.PrivateKey
	EDPrivateKey  ed25519.PrivateKey

	// Public Key Bytes
	RSAPublicKey []byte
	EDPublicKey  []byte
}

type BankKeys struct {
	//Bank
	RSAPublicKey *rsa.PublicKey
	EDPublicKey  ed25519.PublicKey
}

type Keys struct {
	ApplicationKeys *ApplicationKeys
	BankKeys        *BankKeys
}
