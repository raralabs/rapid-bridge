package service

import (
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"rapid-bridge/constants"
	"rapid-bridge/domain/keys"
	"rapid-bridge/domain/port"
	hybridcrypto "rapid-bridge/pkg/security/crypto"
	"rapid-bridge/pkg/util"

	errors "rapid-bridge/internal/error"

	"go.uber.org/zap"
)

type KeyServiceInterface interface {
	GenerateAndSaveApplicationKeys(applicationSlug, ulid string) error
	UseExistingApplicationKeys(rsaPrivateKeyPath, rsaPublicKeyPath, ed25519PrivateKeyPath, ed25519PublicKeyPath string) error

	UseExistingBankKeys(rsaPublicKeyPath, ed25519PublicKeyPath string) error
	FetchAndSaveBankKeys(bankSlug string) error

	FetchBankPublicKeys() (string, string, error)
}

type KeyService struct {
	KeyLoader    port.KeyLoader
	KeyConverter port.KeyConverter
	KeySaver     port.KeySaver

	HttpClient port.HTTPClient
	Logger     port.Logger

	Config port.CLIConfig

	HybridCryptography port.EncryptionDecryptionInterface
}

func (k *KeyService) GenerateAndSaveApplicationKeys(applicationSlug, ulid string) error {
	rsaPrivateKey, rsaPublicKey, err := hybridcrypto.GenerateRSAKeyPair(constants.RSAKeyBitSize)
	if err != nil {
		k.Logger.Error("Error while generating rsa key pair", zap.String("error", err.Error()))
		return err
	}

	err = k.KeySaver.SaveRSAPrivateKeyToPEM(rsaPrivateKey, util.GetRSAPrivateKeyPath(applicationSlug, ulid))
	if err != nil {
		k.Logger.Error("Error while saving rsa private key to pem", zap.String("error", err.Error()))
		return err
	}

	err = k.KeySaver.SaveRSAPublicKeyToPEM(rsaPublicKey, util.GetRSAPublicKeyPath(applicationSlug, ulid))
	if err != nil {
		k.Logger.Error("Error while saving rsa public key to pem", zap.String("error", err.Error()))
		return err
	}

	ed25519PrivateKey, ed25519PublicKey, err := hybridcrypto.GenerateEd25519KeyPair()
	if err != nil {
		k.Logger.Error("Error while generating ed25519 key pair", zap.String("error", err.Error()))
		return err
	}

	err = k.KeySaver.SaveEd25519PrivateKeyToPEM(ed25519PrivateKey, util.GetEd25519PrivateKeyPath(applicationSlug, ulid))
	if err != nil {
		k.Logger.Error("Error while saving ed25519 private key to pem", zap.String("error", err.Error()))
		return err
	}

	err = k.KeySaver.SaveEd25519PublicKeyToPEM(ed25519PublicKey, util.GetEd25519PublicKeyPath(applicationSlug, ulid))
	if err != nil {
		k.Logger.Error("Error while saving ed25519 public key to pem", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (k *KeyService) UseExistingApplicationKeys(applicationSlug, ulid string, rsaPrivateKeyPath, rsaPublicKeyPath, ed25519PrivateKeyPath, ed25519PublicKeyPath string) error {

	if !util.FileExists(rsaPrivateKeyPath) || !util.FileExists(rsaPublicKeyPath) {
		k.Logger.Error("Rsa private or public key files do not exist")
		return fmt.Errorf("rsa private or public key files do not exist")
	}

	if !util.FileExists(ed25519PrivateKeyPath) || !util.FileExists(ed25519PublicKeyPath) {
		k.Logger.Error("Ed25519 private or public key files do not exist")
		return fmt.Errorf("ed25519 private or public key files do not exist")
	}

	rsaPrivateKey, err := keys.ReadAndValidateKeyFile(rsaPrivateKeyPath, true)
	if err != nil {
		k.Logger.Error("Error while validating rsa private key", zap.String("error", err.Error()))
		return err
	}

	rsaPublicKey, err := keys.ReadAndValidateKeyFile(rsaPublicKeyPath, false)
	if err != nil {
		k.Logger.Error("Error while validating rsa public key", zap.String("error", err.Error()))
		return err
	}

	ed25519PrivateKey, err := keys.ReadAndValidateKeyFile(ed25519PrivateKeyPath, true)
	if err != nil {
		k.Logger.Error("Error while validating ed25519 private key", zap.String("error", err.Error()))
		return err
	}

	ed25519PublicKey, err := keys.ReadAndValidateKeyFile(ed25519PublicKeyPath, false)
	if err != nil {
		k.Logger.Error("Error while validating ed25519 public key", zap.String("error", err.Error()))
		return err
	}

	if err := k.KeySaver.SaveRSAPrivateKeyToPEM(rsaPrivateKey.(*rsa.PrivateKey), util.GetRSAPrivateKeyPath(applicationSlug, ulid)); err != nil {
		fmt.Println("Error while saving rsa private key to pem file: ", err)
	}

	if err := k.KeySaver.SaveRSAPublicKeyToPEM(rsaPublicKey.(*rsa.PublicKey), util.GetRSAPublicKeyPath(applicationSlug, ulid)); err != nil {
		fmt.Println("Error while saving rsa public key to pem file: ", err)
	}

	if err := k.KeySaver.SaveEd25519PrivateKeyToPEM(ed25519PrivateKey.(ed25519.PrivateKey), util.GetEd25519PrivateKeyPath(applicationSlug, ulid)); err != nil {
		fmt.Println("Error while saving ed25519 private key to pem file: ", err)
	}

	if err := k.KeySaver.SaveEd25519PublicKeyToPEM(ed25519PublicKey.(ed25519.PublicKey), util.GetEd25519PublicKeyPath(applicationSlug, ulid)); err != nil {
		fmt.Println("Error while saving ed25519 public key to pem file: ", err)
	}

	return nil
}

func (k *KeyService) UseExistingBankKeys(bankSlug string, rsaPublicKeyPath, ed25519PublicKeyPath string) error {
	if !util.FileExists(rsaPublicKeyPath) || !util.FileExists(ed25519PublicKeyPath) {
		k.Logger.Error("Rsa or Ed25519 public key files do not exist")
		return fmt.Errorf("rsa or ed25519 public key files do not exist")
	}

	rsaPublicKey, err := keys.ReadAndValidateKeyFile(rsaPublicKeyPath, false)
	if err != nil {
		k.Logger.Error("Error while validating rsa public key", zap.String("error", err.Error()))
		return err
	}

	ed25519PublicKey, err := keys.ReadAndValidateKeyFile(ed25519PublicKeyPath, false)
	if err != nil {
		k.Logger.Error("Error while validating ed25519 public key", zap.String("error", err.Error()))
		return err
	}

	if err := k.KeySaver.SaveRSAPublicKeyToPEM(rsaPublicKey.(*rsa.PublicKey), util.GetBankRSAPublicKeyPath(bankSlug)); err != nil {
		fmt.Println("Error while saving bank rsa public key to pem file: ", err)
	}

	if err := k.KeySaver.SaveEd25519PublicKeyToPEM(ed25519PublicKey.(ed25519.PublicKey), util.GetBankEd25519PublicKeyPath(bankSlug)); err != nil {
		fmt.Println("Error while saving bank ed25519 public key to pem file: ", err)
	}

	return nil
}

func (k *KeyService) FetchAndSaveBankKeys(rapidUrl, bankSlug string) error {
	k.Logger.Info("Fetching bank's rsa and ed25519 public key from rapid")
	bankRSAPublicKey, bankED25519PublicKey, err := k.FetchBankPublicKeys(rapidUrl)
	if err != nil {
		k.Logger.Error("Error while fetching public keys of bank", zap.String("error", err.Error()))
		return err
	}

	bankRsaPublicKey, err := k.KeyConverter.ConvertBase64ToPublicKey(bankRSAPublicKey)
	if err != nil {
		k.Logger.Error("Error while converting rsa public key of bank", zap.String("error", err.Error()))
		return err
	}

	bankEdPublicKey, err := k.KeyConverter.ConvertBase64ToPublicKey(bankED25519PublicKey)
	if err != nil {
		k.Logger.Error("Error while converting ed25519 public key of bank", zap.String("error", err.Error()))
		return err
	}

	bankRsaPubKey := bankRsaPublicKey.(*rsa.PublicKey)
	bankEdPubKey := bankEdPublicKey.(ed25519.PublicKey)

	rsaPubKeyPath := util.GetBankRSAPublicKeyPath(bankSlug)
	edPubKeyPath := util.GetBankEd25519PublicKeyPath(bankSlug)

	if err := k.KeySaver.SaveRSAPublicKeyToPEM(bankRsaPubKey, rsaPubKeyPath); err != nil {
		k.Logger.Error("Error while saving rsa public key of bank", zap.String("error", err.Error()))
		return err
	}

	if err := k.KeySaver.SaveEd25519PublicKeyToPEM(bankEdPubKey, edPubKeyPath); err != nil {
		k.Logger.Error("Error while saving ed25519 public key of bank", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (k *KeyService) FetchBankPublicKeys(rapidUrl string) (string, string, error) {
	pubKeyResponse, err := k.HttpClient.GET(rapidUrl+"/public-key", map[string]string{}, map[string]string{})
	if err != nil {
		k.Logger.Error("Failed to get public keys", zap.String("error", err.Error()))
		return "", "", err
	}

	publicKeys := pubKeyResponse.Data["data"].(map[string]interface{})

	bankRsaPublicKey, ok := publicKeys["encryptingKey"].(string)
	if !ok {
		k.Logger.Error("Failed to get public keys", zap.String("error", "bank_rsa_public_key not found"))
		return "", "", errors.NewRapidLinksError("bank_rsa_public_key not found", 500)
	}

	bankEd25519PublicKey, ok := publicKeys["signingKey"].(string)
	if !ok {
		k.Logger.Error("Failed to get public keys", zap.String("error", "bank_ed25519_public_key not found"))
		return "", "", errors.NewRapidLinksError("bank_ed25519_public_key not found", 500)
	}

	k.Logger.Info("Bank public keys successfully fetched from rapid", zap.Int("status_code", pubKeyResponse.StatusCode))

	return bankRsaPublicKey, bankEd25519PublicKey, nil
}

func NewKeyService(keyLoader port.KeyLoader, keyConverter port.KeyConverter, keySaver port.KeySaver, httpClient port.HTTPClient, logger port.Logger, config port.CLIConfig) *KeyService {
	return &KeyService{
		KeyLoader:    keyLoader,
		KeyConverter: keyConverter,
		KeySaver:     keySaver,
		HttpClient:   httpClient,
		Logger:       logger,
		Config:       config,
	}
}
