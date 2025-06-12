package service

import (
	"fmt"
	"rapid-bridge/constants"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/adapter/config"
	"rapid-bridge/internal/dto/playground"
	"rapid-bridge/internal/setup"
	"rapid-bridge/pkg/util"
	"slices"
	"strings"

	"go.uber.org/zap"
)

type ApplicationDetails struct {
	RSAPublicKey     string
	Ed25519PublicKey string
	KeyVersion       string
	Slug             string
}

type PlaygroundService struct {
	logger       port.Logger
	app          *setup.CLIApplication
	keyLoader    port.KeyLoader
	keyConverter port.KeyConverter
	keySaver     port.KeySaver
	keyService   *KeyService
}

func NewPlaygroundService(logger port.Logger, app *setup.CLIApplication, keyLoader port.KeyLoader, keyConverter port.KeyConverter, keySaver port.KeySaver, keyService *KeyService) *PlaygroundService {
	return &PlaygroundService{logger: logger, app: app, keyLoader: keyLoader, keyConverter: keyConverter, keySaver: keySaver, keyService: keyService}
}

func (s *PlaygroundService) getApplicationDetails(applicationSlug string) (ApplicationDetails, error) {

	applicationDetails := config.LoadApplicationSpecificConfig(applicationSlug)

	// rsaPrivateKeyPath := applicationDetails.RSAPrivateKeyPath
	// rsaPrivateKey, err := s.keyLoader.LoadPrivateKey(rsaPrivateKeyPath)

	// if err != nil {
	// 	s.logger.Error("Failed to read private keys", zap.String("error", err.Error()))
	// 	return ApplicationDetails{}, err
	// }

	// ed25519PrivateKey, err := s.keyLoader.LoadPrivateKey(applicationDetails.Ed25519PrivateKeyPath)

	// if err != nil {
	// 	s.logger.Error("Failed to read private keys", zap.String("error", err.Error()))
	// 	return ApplicationDetails{}, err
	// }

	rsaPublicKeyBytes, err := util.ReadFile(applicationDetails.RSAPublicKeyPath)
	if err != nil {
		s.logger.Error("Failed to read RSA public keys", zap.String("error", err.Error()))
		return ApplicationDetails{}, err
	}

	rsaPublicKey := s.sanitizePublicKey(string(rsaPublicKeyBytes))

	ed25519PublicKeyBytes, err := util.ReadFile(applicationDetails.Ed25519PublicKeyPath)
	if err != nil {
		s.logger.Error("Failed to read ED25519 public keys", zap.String("error", err.Error()))
		return ApplicationDetails{}, err
	}

	ed25519PublicKey := s.sanitizePublicKey(string(ed25519PublicKeyBytes))

	return ApplicationDetails{
		RSAPublicKey:     rsaPublicKey,
		Ed25519PublicKey: ed25519PublicKey,
		KeyVersion:       applicationDetails.KeyVersion,
		Slug:             applicationSlug,
	}, nil
}

func (s *PlaygroundService) sanitizePublicKey(pemString string) string {
	pemString = strings.ReplaceAll(pemString, "-----BEGIN PUBLIC KEY-----", "")
	pemString = strings.ReplaceAll(pemString, "-----END PUBLIC KEY-----", "")
	// Split into lines
	lines := strings.Split(pemString, "\n")

	var result []string
	result = append(result, "-----BEGIN PUBLIC KEY-----")
	for i, line := range lines {
		if strings.HasPrefix(line, "-----BEGIN PUBLIC KEY-----") || strings.HasPrefix(line, "-----END PUBLIC KEY-----") {
			//result = append(result, line)
		} else if i > 0 && i < len(lines)-1 {
			// Collect all base64 lines into one string without newline
			result = append(result, strings.ReplaceAll(strings.Join(lines[1:len(lines)-1], ""), "\n", ""))
			break
		}
	}
	result = append(result, "-----END PUBLIC KEY-----")

	// Join the result with newlines
	output := strings.Join(result, "\n")
	return output
}

func (s *PlaygroundService) RegisterApplication(request playground.ApplicationRegisterRequest) (playground.ApplicationRegisterResponse, error) {

	isApplicationRegistered := slices.Contains(s.app.Config.GetRegisteredApplications(), request.Slug)

	fmt.Println("isApplicationRegistered", isApplicationRegistered, s.app.Config.GetRegisteredApplications())

	if !isApplicationRegistered {

		ulid := util.GenerateULID().String()

		if err := s.keyService.GenerateAndSaveApplicationKeys(request.Slug, ulid); err != nil {
			s.logger.Error("Error while generating key pair", zap.String("error", err.Error()))
			return playground.ApplicationRegisterResponse{}, err
		}

		s.app.Config.AddRegisteredApplications(request.Slug)
		s.app.Config.AddApplicationSlug(request.Slug)
		s.app.Config.AddApplicationKeysPaths(constants.RapidBridgeData+"/application/"+request.Slug+"/"+ulid+"/rsa_private_key.pem", constants.RapidBridgeData+"/application/"+request.Slug+"/"+ulid+"/rsa_public_key.pem", constants.RapidBridgeData+"/application/"+request.Slug+"/"+ulid+"/ed25519_private_key.pem", constants.RapidBridgeData+"/application/"+request.Slug+"/"+ulid+"/ed25519_public_key.pem")
		s.app.Config.AddKeysValidityPeriod(constants.EncryptionKeyValidityPeriod, constants.SigningKeyValidityPeriod)
		s.app.Config.AddApplicationUlid(ulid)

		if err := s.app.Config.SaveApplicationConfigToFile(); err != nil {
			s.logger.Error("Error while saving config", zap.String("error", err.Error()))
			return playground.ApplicationRegisterResponse{}, err
		}

		s.logger.Info("Application registered successfully")
	}

	applicationDetails, err := s.getApplicationDetails(request.Slug)

	if err != nil {
		s.logger.Error("Failed to get application details", zap.String("error", err.Error()))
		return playground.ApplicationRegisterResponse{}, err
	}

	return playground.ApplicationRegisterResponse{
		KeyVersion:       applicationDetails.KeyVersion,
		Slug:             applicationDetails.Slug,
		RSAPublicKey:     applicationDetails.RSAPublicKey,
		Ed25519PublicKey: applicationDetails.Ed25519PublicKey,
		Message:          "Key pair fetched successfully",
	}, nil
}
