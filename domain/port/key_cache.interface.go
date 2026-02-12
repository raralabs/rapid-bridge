package port

import (
	"rapid-bridge/internal/dto/application"
)

type KeyCache interface {
	GetKeys(source, destination, keyVersion string) *application.Keys

	GetRawApplicationKeys(source, rsaPublicKeyPath string, edPublicKeyPath string, keyVersion string, isApp bool) *application.ApplicationKeys
}
