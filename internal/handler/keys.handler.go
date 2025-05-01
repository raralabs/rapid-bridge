package handler

import (
	"fmt"
	"rapid-bridge/internal/service"
)

type KeyHandler struct {
	Service *service.KeyService
}

func (k *KeyHandler) HandleApplicationGenerateKeyPair(applicationSlug string, ulid string) error {
	fmt.Println("Generating key pair...", ulid)
	return k.Service.GenerateAndSaveApplicationKeys(applicationSlug, ulid)
}

func (k *KeyHandler) HandleApplicationExistingKeyPair(applicationSlug, ulid string, rsaPrivateKeyPath, rsaPublicKeyPath, ed25519PrivateKeyPath, ed25519PublicKeyPath string) error {
	return k.Service.UseExistingApplicationKeys(applicationSlug, ulid, rsaPrivateKeyPath, rsaPublicKeyPath, ed25519PrivateKeyPath, ed25519PublicKeyPath)
}

func (k *KeyHandler) HandleBankExistingKeys(bankSlug, rsaPublicKeyPath, ed25519PublicKeyPath string) error {
	return k.Service.UseExistingBankKeys(bankSlug, rsaPublicKeyPath, ed25519PublicKeyPath)
}

func (k *KeyHandler) HandleBankFetchKeys(bankSlug string) error {
	return k.Service.FetchAndSaveBankKeys(bankSlug)
}

func NewKeyHandler(service *service.KeyService) *KeyHandler {
	return &KeyHandler{
		Service: service,
	}
}
