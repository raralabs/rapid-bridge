package port

import (
	"crypto/ed25519"
	"crypto/rsa"
	"time"
)

type ServerConfig interface {
	GetRapidLinksUrl() string
}

type CLIConfig interface {
	GetRegisteredBanks() []string
	GetRegisteredApplications() []string

	GetApplicationDetails(applicationSlug string) *CLIApplicationDetails

	AddBankSlug(bankSlug string)
	AddRegisteredBanks(bankSlug string)
	AddBankKeysPaths(rsaPublicKeyPath string, ed25519PublicKeyPath string)

	AddRegisteredApplications(applicationSlug string)
	AddApplicationSlug(applicationSlug string)
	AddApplicationUlid(ulid string)
	AddApplicationKeysPaths(rsaPrivateKeyPath string, rsaPublicKeyPath string, ed25519PrivateKeyPath string, ed25519PublicKeyPath string)
	AddKeysValidityPeriod(encryptionKeyValidityPeriod, signingKeyValidityPeriod int)

	SaveApplicationConfigToFile(applicationSlug string, ulid string, rsaPrivateKeyPath, rsaPublicKeyPath, ed25519PrivateKeyPath, ed25519PublicKeyPath string) error
	SaveBankConfigToFile(bankSlug string, rsaPublicKeyPath, ed25519PublicKeyPath string) error

	SaveConfigToFile() error
}

type ApplicationDetails struct {
	// for reading keys if file path specified
	RSAPrivateKeyPath     string `json:"rsa_private_key_path" mapstructure:"rsa_private_key_path"`
	RSAPublicKeyPath      string `json:"rsa_public_key_path" mapstructure:"rsa_public_key_path"`
	Ed25519PrivateKeyPath string `json:"ed25519_private_key_path" mapstructure:"ed25519_private_key_path"`
	Ed25519PublicKeyPath  string `json:"ed25519_public_key_path" mapstructure:"ed25519_public_key_path"`

	// for in-memory
	RSAPrivateKey     *rsa.PrivateKey    `json:"rsa_private_key,omitempty"`
	RSAPublicKey      *rsa.PublicKey     `json:"rsa_public_key,omitempty"`
	Ed25519PrivateKey ed25519.PrivateKey `json:"ed25519_private_key,omitempty"`
	Ed25519PublicKey  ed25519.PublicKey  `json:"ed25519_public_key,omitempty"`

	// Keys expiry / validity
	RSAKeysValidUntil time.Time `json:"rsa_keys_valid_until" mapstructure:"rsa_keys_valid_until"`

	Ed25519KeysValidUntil time.Time `json:"ed25519_keys_valid_until" mapstructure:"ed25519_keys_valid_until"`

	Slug          string `json:"slug" mapstructure:"slug"`
	KeyVersion    string `json:"key_version" mapstructure:"key_version"`
	ServerAddress string `json:"server_address,omitempty"`
}

type CLIApplicationDetails struct {
	RSAPrivateKeyPath     string `json:"rsa_private_key_path" mapstructure:"rsa_private_key_path"`
	RSAPublicKeyPath      string `json:"rsa_public_key_path" mapstructure:"rsa_public_key_path"`
	Ed25519PrivateKeyPath string `json:"ed25519_private_key_path" mapstructure:"ed25519_private_key_path"`
	Ed25519PublicKeyPath  string `json:"ed25519_public_key_path" mapstructure:"ed25519_public_key_path"`

	// Keys expiry / validity
	RSAKeysValidUntil time.Time `json:"rsa_keys_valid_until" mapstructure:"rsa_keys_valid_until"`

	Ed25519KeysValidUntil time.Time `json:"ed25519_keys_valid_until" mapstructure:"ed25519_keys_valid_until"`

	Slug       string `json:"slug" mapstructure:"slug"`
	KeyVersion string `json:"key_version" mapstructure:"key_version"`
}
