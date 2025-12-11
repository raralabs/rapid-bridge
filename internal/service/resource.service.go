package service

import (
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"rapid-bridge/constants"
	"rapid-bridge/domain/port"
	"rapid-bridge/domain/security"
	"rapid-bridge/internal/adapter"
	"rapid-bridge/internal/adapter/cache"
	"rapid-bridge/internal/dto/application"
	"rapid-bridge/internal/dto/rapid"
	"rapid-bridge/pkg/util"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RapidResourceService struct {
	loader   port.KeyLoader
	security security.Security
	logger   port.Logger
	config   port.ServerConfig
	keyCache *cache.KeyCache
}

func (r *RapidResourceService) HandleResource(c echo.Context, request application.ResourceRequest) (application.ResourceResponse, error) {

	ctx := util.GetReqCtxFromEchoCtx(c)

	from := ctx.Value(constants.From).(string)
	to := ctx.Value(constants.To).(string)
	keyVersion := ctx.Value(constants.KeyVersion).(string)

	// Load keys from cache instead of file
	// If loading from cache fails then loading keys from the file structure

	// RSA Private Key (Application)
	rsaPrivateKey, err := r.keyCache.GetRSAPrivateKey(from, keyVersion, true)
	if err != nil {
		rsaPrivateKeyPath := util.GetRSAPrivateKeyPath(from, keyVersion)
		keyAny, err := r.loader.LoadPrivateKey(rsaPrivateKeyPath)
		if err != nil {
			r.logger.Error("Failed to load Application's RSA Private Key from file", zap.String("error", err.Error()))
			return application.ResourceResponse{}, err
		}

		// Type assertion from any to *rsa.PrivateKey
		key, ok := keyAny.(*rsa.PrivateKey)
		if !ok {
			r.logger.Error("Failed to convert datatype from any to *rsa.Private")
			return application.ResourceResponse{}, fmt.Errorf("loaded key is not an RSA private key")
		}

		r.keyCache.SetRSAPrivateKey(from, keyVersion, true, key)
		rsaPrivateKey = key
	}

	// ED25519 Private Key (Application)
	ed25519PrivateKey, err := r.keyCache.GetEd25519PrivateKey(from, keyVersion, true)
	if err != nil {
		edPrivateKeyPath := util.GetEd25519PrivateKeyPath(from, keyVersion)
		keyAny, err := r.loader.LoadPrivateKey(edPrivateKeyPath)
		if err != nil {
			r.logger.Error("Failed to load Application's ED25519 Private Key from file", zap.String("error", err.Error()))
			return application.ResourceResponse{}, err
		}

		// Type assertion from any to ed25519.PrivateKey
		key, ok := keyAny.(ed25519.PrivateKey)
		if !ok {
			r.logger.Error("Failed to convert datatype from any to *rsa.Private")
			return application.ResourceResponse{}, fmt.Errorf("loaded key is not an ED25519 private key")
		}

		r.keyCache.SetEd25519PrivateKey(from, keyVersion, true, key)
		ed25519PrivateKey = key
	}

	// RSA Public Key (Bank)
	bankRsaPublicKey, err := r.keyCache.GetRSAPublicKey(to, keyVersion, false)
	if err != nil {
		bankRsaPublicKeyPath := util.GetRSAPublicKeyPath(to, keyVersion)
		keyAny, err := r.loader.LoadPublicKey(bankRsaPublicKeyPath)
		if err != nil {
			r.logger.Error("Failed to load Bank's RSA Public Key from file", zap.String("error", err.Error()))
			return application.ResourceResponse{}, err
		}

		// Type assertion from any to *rsa.PublicKey
		key, ok := keyAny.(*rsa.PublicKey)
		if !ok {
			r.logger.Error("Failed to convert datatype from any to *rsa.PublicKey")
			return application.ResourceResponse{}, fmt.Errorf("loaded key is not an RSA public key")
		}

		r.keyCache.SetRSAPublicKey(to, keyVersion, false, key)
		bankRsaPublicKey = key
	}

	// ED25519 Public Key (Bank)
	bankEdPublicKey, err := r.keyCache.GetEd25519PublicKey(to, keyVersion, false)
	if err != nil {
		bankEdPublicKeyPath := util.GetBankEd25519PublicKeyPath(to)
		keyAny, err := r.loader.LoadPublicKey(bankEdPublicKeyPath)
		if err != nil {
			r.logger.Error("Failed to load Bank's ED25519 Public Key from file", zap.String("error", err.Error()))
			return application.ResourceResponse{}, err
		}

		// Type assertion from any to ed25519.PublicKey
		key, ok := keyAny.(ed25519.PublicKey)
		if !ok {
			r.logger.Error("Failed to convert datatype from any to ed25519.PublicKey")
			return application.ResourceResponse{}, fmt.Errorf("loaded key is not an ED25519 public key")
		}

		r.keyCache.SetEd25519PublicKey(to, keyVersion, false, key)
		bankEdPublicKey = key
	}

	// convert request struct to bytes
	data, err := json.Marshal(request)
	if err != nil {
		r.logger.Error("Failed to marshal request", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	ciphertext, encryptedAESKey, nonce, err := r.security.Encrypt(data, bankRsaPublicKey)
	if err != nil {
		r.logger.Error("Failed to encrypt payload", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// sign payload
	signature, err := r.security.CreateDigitalSignature(ed25519PrivateKey, ciphertext, encryptedAESKey, nonce)

	// create base64 encrypted payload
	base64EncryptedPayload, err := r.security.CreateBase64Encrypted(ciphertext, encryptedAESKey, nonce)
	if err != nil {
		r.logger.Error("Failed to create base64 encrypted payload", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// create rapid resource request
	rapidResourceRequest := rapid.RapidResourceRequest{
		From:       ctx.Value(constants.From).(string),
		To:         ctx.Value(constants.To).(string),
		Message:    base64EncryptedPayload,
		Signature:  signature,
		KeyVersion: keyVersion,
	}

	// send rapid resource request to rapid links
	rapidLinksUrl := r.config.GetRapidLinksUrl()
	rapidResourceResponse, err := adapter.SendRequestToRapidLinks(r.logger, rapidLinksUrl, c.Request().URL.Path, rapidResourceRequest, c.Request().Header)
	if err != nil {
		r.logger.Error("Failed to send rapid resource request to rapid links", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	from = rapidResourceResponse.GetFrom()
	to = rapidResourceResponse.GetTo()
	message := rapidResourceResponse.GetMessage()
	signature = rapidResourceResponse.GetSignature()

	r.logger.Info("Message from rapid links", zap.String("from", from), zap.String("to", to))

	// decode message and get ciphertext, encrypted aes key and nonce
	ciphertext, encryptedAESKey, nonce, err = r.security.DecodeBase64Encrypted(message)
	if err != nil {
		r.logger.Error("Failed to decode message", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// decrypt payload
	decryptedPayload, err := r.security.Decrypt(rsaPrivateKey, ciphertext, encryptedAESKey, nonce)
	if err != nil {
		r.logger.Error("Failed to decrypt payload", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// verify signature
	err = r.security.VerifyDigitalSignature(message, signature, bankEdPublicKey)
	if err != nil {
		r.logger.Error("Failed to verify digital signature", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// create rapid resource response
	applicationResponse := application.ResourceResponse{
		Message: string(decryptedPayload),
	}

	return applicationResponse, nil
}

func NewRapidResourceService(keyLoader port.KeyLoader, security security.Security, logger port.Logger, config port.ServerConfig, keyCache *cache.KeyCache) *RapidResourceService {
	return &RapidResourceService{
		loader:   keyLoader,
		security: security,
		logger:   logger,
		config:   config,
		keyCache: keyCache,
	}
}
