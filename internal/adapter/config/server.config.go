package config

import (
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"rapid-bridge/constants"
	"rapid-bridge/domain/port"

	"github.com/spf13/viper"
)

type RapidLinks struct {
	Url string `mapstructure:"rapid_links_url"`
}

type ApplicationDetails struct {
	// for reading keys if file path specified
	RSAPrivateKeyPath     string `json:"rsa_private_key_path"`
	RSAPublicKeyPath      string `json:"rsa_public_key_path"`
	Ed25519PrivateKeyPath string `json:"ed25519_private_key_path"`
	Ed25519PublicKeyPath  string `json:"ed25519_public_key_path"`

	// for in-memory
	RSAPrivateKey     *rsa.PrivateKey
	RSAPublicKey      *rsa.PublicKey
	Ed25519PrivateKey ed25519.PrivateKey
	Ed25519PublicKey  ed25519.PublicKey

	KeyVersion string `json:"key_version"`

	Slug          string `json:"slug"`
	ServerAddress string
}

type BankDetails struct {
	// for reading keys if file path specified
	RSAPublicKeyPath     string `json:"rsa_public_key_path"`
	Ed25519PublicKeyPath string `json:"ed25519_public_key_path"`

	// for in-memory
	RSAPublicKey     *rsa.PublicKey    `json:"rsa_public_key,omitempty"`
	Ed25519PublicKey ed25519.PublicKey `json:"ed25519_public_key,omitempty"`

	Slug string `json:"slug"`
}

type ServerConfig struct {
	RapidLinks         RapidLinks
	ApplicationDetails ApplicationDetails
	BankDetails        BankDetails
}

type ServerConfigAdapter struct {
	ServerConfig
}

func (s *ServerConfigAdapter) GetRapidLinksUrl() string {
	return s.ServerConfig.RapidLinks.Url
}

func LoadServerConfig() (port.ServerConfig, error) {

	v := viper.New()

	v.SetConfigName("core")
	v.SetConfigType("json")
	v.AddConfigPath(constants.RapidBridgeData)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	cfg := ServerConfig{
		RapidLinks: RapidLinks{
			Url: v.GetString("rapid_links_url"),
		},
	}
	serverConfig := &ServerConfigAdapter{ServerConfig: cfg}

	return serverConfig, nil
}
