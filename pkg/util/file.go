package util

import (
	"math/rand"
	"os"
	"path/filepath"
	"rapid-bridge/constants"
	"time"

	"github.com/oklog/ulid/v2"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetRSAPrivateKeyPath(applicationSlug, newUlid, rsaPrivateKeyPath string) string {
	if rsaPrivateKeyPath == "" {
		rsaPrivateKeyPath = filepath.Join(constants.RapidBridgeData, "application", applicationSlug, newUlid, constants.RSAPrivateKeyFile)
	}
	return rsaPrivateKeyPath
}

func GetRSAPublicKeyPath(applicationSlug, newUlid, rsaPublicKeyPath string) string {
	if rsaPublicKeyPath == "" {
		rsaPublicKeyPath = filepath.Join(constants.RapidBridgeData, "application", applicationSlug, newUlid, constants.RSAPublicKeyFile)
	}
	return rsaPublicKeyPath
}

func GetEd25519PrivateKeyPath(applicationSlug, newUlid, ed25519PrivateKeyPath string) string {
	if ed25519PrivateKeyPath == "" {
		ed25519PrivateKeyPath = filepath.Join(constants.RapidBridgeData, "application", applicationSlug, newUlid, constants.Ed25519PrivateKeyFile)
	}
	return ed25519PrivateKeyPath
}

func GetEd25519PublicKeyPath(applicationSlug, newUlid, ed25519PublicKeyPath string) string {
	if ed25519PublicKeyPath == "" {
		ed25519PublicKeyPath = filepath.Join(constants.RapidBridgeData, "application", applicationSlug, newUlid, constants.Ed25519PublicKeyFile)
	}
	return ed25519PublicKeyPath
}

func GetBankRSAPublicKeyPath(bankSlug string) string {
	return filepath.Join(constants.RapidBridgeData, "bank", bankSlug, constants.RSAPublicKeyFile)

}

func GetBankEd25519PublicKeyPath(bankSlug string) string {
	return filepath.Join(constants.RapidBridgeData, "bank", bankSlug, constants.Ed25519PublicKeyFile)
}

func GenerateULID() ulid.ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

// CompareULIDs returns the oldest ULID (minimum) based on their embedded timestamp
func CompareULIDs(ulids ...ulid.ULID) ulid.ULID {
	if len(ulids) == 0 {
		return ulid.ULID{}
	}

	min := ulids[0]
	for _, u := range ulids[1:] {
		if u.Compare(min) < 0 {
			min = u
		}
	}
	return min
}
