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

	rsaPrivateKeyPath := util.GetRSAPrivateKeyPath(request.From, request.KeyVersion, "")
	rsaPrivateKey, err := r.loader.LoadPrivateKey(rsaPrivateKeyPath)

	if err != nil {
		r.logger.Error("Failed to read private keys", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	ed25519PrivateKey, err := r.loader.LoadPrivateKey(util.GetEd25519PrivateKeyPath(request.From, request.KeyVersion, ""))

	if err != nil {
		r.logger.Error("Failed to read private keys", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	bankRsaPublicKey, err := r.loader.LoadPublicKey(util.GetBankRSAPublicKeyPath(request.To))

	if err != nil {
		r.logger.Error("Failed to read public keys", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	bankEdPublicKey, err := r.loader.LoadPublicKey(util.GetBankEd25519PublicKeyPath(request.To))

	if err != nil {
		r.logger.Error("Failed to read public keys", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// TODO: Handle key loading and saving it to config (in-memory) in server startup
	// // read public keys of bank from config file
	// bankRsaPublicKey, err := security.LoadPublicKey(r.config.BankDetails.RSAPublicKeyPath)
	// if err != nil {
	// 	r.logger.Error("Failed to read public keys", zap.String("error", err.Error()))
	// 	return application.ResourceResponse{}, err
	// }

	// // read private keys of application from config file
	// appEdPrivateKey, err := security.LoadPrivateKey(r.config.ApplicationDetails.Ed25519PrivateKeyPath)
	// if err != nil {
	// 	r.logger.Error("Failed to read private keys", zap.String("error", err.Error()))
	// 	return application.ResourceResponse{}, err
	// }

	// // read private keys of application from config file
	// appRsaPrivateKey, err := security.LoadPrivateKey(r.config.ApplicationDetails.RSAPrivateKeyPath)
	// if err != nil {
	// 	r.logger.Error("Failed to read private keys", zap.String("error", err.Error()))
	// 	return application.ResourceResponse{}, err
	// }

	// // read bank public key from config file
	// bankEd25519PublicKey, err := security.LoadPublicKey(r.config.BankDetails.Ed25519PublicKeyPath)
	// if err != nil {
	// 	r.logger.Error("Failed to read public keys", zap.String("error", err.Error()))
	// 	return application.ResourceResponse{}, err
	// }

	// // get public key of bank
	// bankEncryptionKey := bankRsaPublicKey.(*rsa.PublicKey)
	// bankSigningKey := bankEd25519PublicKey.(ed25519.PublicKey)

	// // get private key of application
	// applicationRSAPrivateKey := appRsaPrivateKey.(*rsa.PrivateKey)
	// applicationEd25519PrivateKey := appEdPrivateKey.(ed25519.PrivateKey)

	// convert request struct to bytes
	data, err := json.Marshal(request)
	if err != nil {
		r.logger.Error("Failed to marshal request", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	ciphertext, encryptedAESKey, nonce, err := r.security.Encrypt(data, bankRsaPublicKey.(*rsa.PublicKey))
	if err != nil {
		r.logger.Error("Failed to encrypt payload", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// sign payload
	signature, err := r.security.CreateDigitalSignature(ed25519PrivateKey.(ed25519.PrivateKey), ciphertext, encryptedAESKey, nonce)

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
		KeyVersion: request.KeyVersion,
	}

	// send rapid resource request to rapid links
	rapidLinksUrl := r.config.GetRapidLinksUrl()
	rapidResourceResponse, err := adapter.SendRequestToRapidLinks(r.logger, rapidLinksUrl, c.Request().URL.Path, rapidResourceRequest, c.Request().Header)
	if err != nil {
		r.logger.Error("Failed to send rapid resource request to rapid links", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	r.logger.Info("Message from rapid links", rapidResourceResponse)

	// decode message and get ciphertext, encrypted aes key and nonce
	ciphertext, encryptedAESKey, nonce, err = r.security.DecodeBase64Encrypted(rapidResourceResponse.Data.Message)
	if err != nil {
		r.logger.Error("Failed to decode message", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// decrypt payload
	decryptedPayload, err := r.security.Decrypt(rsaPrivateKey.(*rsa.PrivateKey), ciphertext, encryptedAESKey, nonce)
	if err != nil {
		r.logger.Error("Failed to decrypt payload", zap.String("error", err.Error()))
		return application.ResourceResponse{}, err
	}

	// verify signature
	err = r.security.VerifyDigitalSignature(rapidResourceResponse.Data.Message, rapidResourceResponse.Data.Signature, bankEdPublicKey.(ed25519.PublicKey))
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

func NewRapidResourceService(keyLoader port.KeyLoader, security security.Security, logger port.Logger, config port.ServerConfig) *RapidResourceService {
	return &RapidResourceService{
		loader:   keyLoader,
		security: security,
		logger:   logger,
		config:   config,
	}
}
