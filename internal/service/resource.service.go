package service

import (
	"encoding/json"
	"fmt"
	"rapid-bridge/constants"
	"rapid-bridge/domain/port"
	"rapid-bridge/domain/security"
	"rapid-bridge/internal/adapter"
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
	keyCache port.KeyCache
}

func (r *RapidResourceService) HandleResource(c echo.Context, request application.ResourceRequest) (application.ResourceResponse, error) {

	ctx := util.GetReqCtxFromEchoCtx(c)

	from := ctx.Value(constants.From).(string)
	to := ctx.Value(constants.To).(string)
	keyVersion := ctx.Value(constants.KeyVersion).(string)

	// Load keys from cache instead of file
	// If loading from cache fails then loading keys from the file structure

	// RSA Private Key (Application)
	rsaPrivateKey := r.keyCache.GetRSAPrivateKey(from, keyVersion, true)
	if rsaPrivateKey == nil {
		r.logger.Error("Failed to fetch Application's RSA Private Key.")
		return application.ResourceResponse{}, fmt.Errorf("Failed to fetch Application's RSA Private Key from cache.")
	}

	// ED25519 Private Key (Application)
	ed25519PrivateKey := r.keyCache.GetEd25519PrivateKey(from, keyVersion, true)
	if ed25519PrivateKey == nil {
		r.logger.Error("Failed to fetch Application's ED25519 Private Key from cache.")
		return application.ResourceResponse{}, fmt.Errorf("Failed to fetch Application's ED25519 Private Key from cache.")
	}

	// RSA Public Key (Bank)
	bankRsaPublicKey := r.keyCache.GetRSAPublicKey(to, keyVersion, false)
	if bankRsaPublicKey == nil {
		r.logger.Error("Failed to fetch Bank's RSA Public Key from cache.")
		return application.ResourceResponse{}, fmt.Errorf("Failed to fetch Bank's RSA Public Key from cache.")
	}

	// ED25519 Public Key (Bank)
	bankEdPublicKey := r.keyCache.GetEd25519PublicKey(to, keyVersion, false)
	if bankEdPublicKey == nil {
		r.logger.Error("Failed to fetch Bank's EDD25519 Public Key from cache.")
		return application.ResourceResponse{}, fmt.Errorf("Failed to fetch Bank's EDD25519 Public Key from cache.")
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
		r.logger.Error("Failed to send rapid resource request to rapid links", zap.String("error`", err.Error()))
		return application.ResourceResponse{}, err
	}

	from = rapidResourceResponse.GetFrom()
	to = rapidResourceResponse.GetTo()
	message := rapidResourceResponse.GetMessage()
	signature = rapidResourceResponse.GetSignature()

	r.logger.Info("Successfully called to Rapid", zap.String("URL", rapidLinksUrl+c.Request().URL.Path), zap.String("Source Slug", rapidResourceRequest.From), zap.String("Destination Slug", rapidResourceRequest.To), zap.String("Key Version", rapidResourceRequest.KeyVersion))

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

func NewRapidResourceService(keyLoader port.KeyLoader, security security.Security, logger port.Logger, config port.ServerConfig, keyCache port.KeyCache) *RapidResourceService {
	return &RapidResourceService{
		loader:   keyLoader,
		security: security,
		logger:   logger,
		config:   config,
		keyCache: keyCache,
	}
}
