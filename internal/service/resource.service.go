package service

import (
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/json"
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
}

func (r *RapidResourceService) HandleResource(c echo.Context, request application.ResourceRequest) (application.ResourceResponse, error) {
	ctx := util.GetReqCtxFromEchoCtx(c)

	from := ctx.Value(constants.From).(string)
	to := ctx.Value(constants.To).(string)
	keyVersion := ctx.Value(constants.KeyVersion).(string)

	// Load RSA private key
	rsaPrivateKeyPath := util.GetRSAPrivateKeyPath(from, keyVersion)
	rsaPrivateKey, err := r.loader.LoadPrivateKey(rsaPrivateKeyPath)
	if err != nil {
		r.logger.Error("Failed to read RSA private key", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Failed to load RSA private key: " + err.Error(),
			Message:      "",
		}, err
	}

	// Load Ed25519 private key
	ed25519PrivateKey, err := r.loader.LoadPrivateKey(util.GetEd25519PrivateKeyPath(from, keyVersion))
	if err != nil {
		r.logger.Error("Failed to read Ed25519 private key", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Failed to load Ed25519 private key: " + err.Error(),
			Message:      "",
		}, err
	}

	// Load bank public keys
	bankRsaPublicKey, err := r.loader.LoadPublicKey(util.GetBankRSAPublicKeyPath(to))
	if err != nil {
		r.logger.Error("Failed to read bank RSA public key", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Failed to load bank RSA public key: " + err.Error(),
			Message:      "",
		}, err
	}

	bankEdPublicKey, err := r.loader.LoadPublicKey(util.GetBankEd25519PublicKeyPath(to))
	if err != nil {
		r.logger.Error("Failed to read bank Ed25519 public key", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Failed to load bank Ed25519 public key: " + err.Error(),
			Message:      "",
		}, err
	}

	// Marshal request into JSON
	data, err := json.Marshal(request)
	if err != nil {
		r.logger.Error("Failed to marshal request", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   400,
			ErrorMessage: "Invalid request payload: " + err.Error(),
			Message:      "",
		}, err
	}

	// Encrypt payload
	ciphertext, encryptedAESKey, nonce, err := r.security.Encrypt(data, bankRsaPublicKey.(*rsa.PublicKey))
	if err != nil {
		r.logger.Error("Failed to encrypt payload", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Encryption failed: " + err.Error(),
			Message:      "",
		}, err
	}

	// Sign payload
	signature, err := r.security.CreateDigitalSignature(ed25519PrivateKey.(ed25519.PrivateKey), ciphertext, encryptedAESKey, nonce)
	if err != nil {
		r.logger.Error("Failed to create digital signature", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Digital signature failed: " + err.Error(),
			Message:      "",
		}, err
	}

	// Create base64 encrypted payload
	base64EncryptedPayload, err := r.security.CreateBase64Encrypted(ciphertext, encryptedAESKey, nonce)
	if err != nil {
		r.logger.Error("Failed to create base64 encrypted payload", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Base64 encryption failed: " + err.Error(),
			Message:      "",
		}, err
	}

	// Build Rapid resource request
	rapidResourceRequest := rapid.RapidResourceRequest{
		From:       from,
		To:         to,
		Message:    base64EncryptedPayload,
		Signature:  signature,
		KeyVersion: keyVersion,
	}

	// Send request to Rapid Links
	rapidLinksUrl := r.config.GetRapidLinksUrl()
	rapidResourceResponse, err := adapter.SendRequestToRapidLinks(r.logger, rapidLinksUrl, c.Request().URL.Path, rapidResourceRequest, c.Request().Header)
	if err != nil {
		r.logger.Error("Failed to send request to Rapid Links", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   502,
			ErrorMessage: "Failed to communicate with Rapid Links: " + err.Error(),
			Message:      "",
		}, err
	}

	from = rapidResourceResponse.GetFrom()
	to = rapidResourceResponse.GetTo()
	message := rapidResourceResponse.GetMessage()
	signature = rapidResourceResponse.GetSignature()

	r.logger.Info("Message from Rapid Links", zap.String("from", from), zap.String("to", to))

	// Decode base64 encrypted message
	ciphertext, encryptedAESKey, nonce, err = r.security.DecodeBase64Encrypted(message)
	if err != nil {
		r.logger.Error("Failed to decode message", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Failed to decode message: " + err.Error(),
			Message:      "",
		}, err
	}

	// Decrypt payload
	decryptedPayload, err := r.security.Decrypt(rsaPrivateKey.(*rsa.PrivateKey), ciphertext, encryptedAESKey, nonce)
	if err != nil {
		r.logger.Error("Failed to decrypt payload", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   500,
			ErrorMessage: "Failed to decrypt message: " + err.Error(),
			Message:      "",
		}, err
	}

	// Verify digital signature
	err = r.security.VerifyDigitalSignature(message, signature, bankEdPublicKey.(ed25519.PublicKey))
	if err != nil {
		r.logger.Error("Failed to verify digital signature", zap.String("error", err.Error()))
		return application.ResourceResponse{
			StatusCode:   400,
			ErrorMessage: "Invalid digital signature: " + err.Error(),
			Message:      "",
		}, err
	}

	// Check status code from Rapid Links
	if rapidResourceResponse.StatusCode != 200 {
		r.logger.Error("Rapid Links returned error", zap.Int("status_code", rapidResourceResponse.StatusCode))
		return application.ResourceResponse{
			StatusCode:   rapidResourceResponse.StatusCode,
			ErrorMessage: string(decryptedPayload),
			Message:      "",
		}, nil
	}

	// Success response
	return application.ResourceResponse{
		Message:      string(decryptedPayload),
		StatusCode:   200,
		ErrorMessage: "",
	}, nil
}

func NewRapidResourceService(keyLoader port.KeyLoader, security security.Security, logger port.Logger, config port.ServerConfig) *RapidResourceService {
	return &RapidResourceService{
		loader:   keyLoader,
		security: security,
		logger:   logger,
		config:   config,
	}
}
