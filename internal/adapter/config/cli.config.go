package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"rapid-bridge/constants"
	"rapid-bridge/domain/port"
	"time"

	"github.com/spf13/viper"
)

type CLIConfig struct {
	RapidLinks         RapidLinks              `mapstructure:",squash"`
	ApplicationDetails port.ApplicationDetails `json:"application"`
	BankDetails        BankDetails             `json:"bank"`

	RegisteredApplications []string `mapstructure:"registered_applications"`
	RegisteredBanks        []string `mapstructure:"registered_banks"`
}

type FlatCLIConfig struct {
	RapidLinksURL          string   `json:"rapid_links_url"`
	RegisteredApplications []string `json:"registered_applications"`
	RegisteredBanks        []string `json:"registered_banks"`
	ApplicationKeyVersion  string   `json:"application_key_version"`

	ApplicationRSAPrivateKeyPath     string `json:"application_rsa_private_key_path"`
	ApplicationEd25519PrivateKeyPath string `json:"application_ed25519_private_key_path"`

	BankRSAPublicKeyPath     string `json:"bank_rsa_public_key_path"`
	BankEd25519PublicKeyPath string `json:"bank_ed25519_public_key_path"`
}

type FileConfigAdapter struct {
	CLIConfig
}

func (f *FileConfigAdapter) GetRapidLinksUrl() string {
	return f.CLIConfig.RapidLinks.Url
}

func (f *FileConfigAdapter) GetBankSlug(key string) string {
	return f.CLIConfig.BankDetails.Slug
}

func (f *FileConfigAdapter) GetRegisteredBanks() []string {
	return f.CLIConfig.RegisteredBanks
}

func (f *FileConfigAdapter) GetRegisteredApplications() []string {
	return f.CLIConfig.RegisteredApplications
}

func (f *FileConfigAdapter) GetBankKeysPaths() (rsaPublicKeyPath string, ed25519PublicKeyPath string) {
	return f.CLIConfig.BankDetails.RSAPublicKeyPath, f.CLIConfig.BankDetails.Ed25519PublicKeyPath
}

func (f *FileConfigAdapter) GetApplicationDetails(applicationSlug string) *port.ApplicationDetails {
	return &f.CLIConfig.ApplicationDetails
}

func (f *FileConfigAdapter) AddApplicationDetails(applicationSlug string) {
	f.CLIConfig.ApplicationDetails.Slug = applicationSlug
}

func (f *FileConfigAdapter) AddApplicationUlid(ulid string) {
	f.CLIConfig.ApplicationDetails.KeyVersion = ulid
}

func (f *FileConfigAdapter) AddApplicationKeysPaths(rsaPrivateKeyPath string, rsaPublicKeyPath string, ed25519PrivateKeyPath string, ed25519PublicKeyPath string) {
	f.CLIConfig.ApplicationDetails.RSAPrivateKeyPath = rsaPrivateKeyPath
	f.CLIConfig.ApplicationDetails.RSAPublicKeyPath = rsaPublicKeyPath
	f.CLIConfig.ApplicationDetails.Ed25519PrivateKeyPath = ed25519PrivateKeyPath
	f.CLIConfig.ApplicationDetails.Ed25519PublicKeyPath = ed25519PublicKeyPath
}

func (f *FileConfigAdapter) AddRegisteredApplications(applicationSlug string) {
	f.CLIConfig.RegisteredApplications = append(f.CLIConfig.RegisteredApplications, applicationSlug)
}

func (f *FileConfigAdapter) AddKeysValidityPeriod(encryptionKeyValidityPeriod, signingKeyValidityPeriod int) {
	f.CLIConfig.ApplicationDetails.RSAKeysValidUntil = time.Now().AddDate(0, 0, encryptionKeyValidityPeriod)
	f.CLIConfig.ApplicationDetails.Ed25519KeysValidUntil = time.Now().AddDate(0, 0, signingKeyValidityPeriod)
}

func (f *FileConfigAdapter) AddBankDetails(bankSlug string) {
	f.CLIConfig.BankDetails.Slug = bankSlug
}

func (f *FileConfigAdapter) AddRegisteredBanks(bankSlug string) {
	f.CLIConfig.RegisteredBanks = append(f.CLIConfig.RegisteredBanks, bankSlug)
}

func (f *FileConfigAdapter) AddBankKeysPaths(rsaPublicKeyPath string, ed25519PublicKeyPath string) {
	f.CLIConfig.BankDetails.RSAPublicKeyPath = rsaPublicKeyPath
	f.CLIConfig.BankDetails.Ed25519PublicKeyPath = ed25519PublicKeyPath
}

func (f *FileConfigAdapter) SaveApplicationConfigToFile(applicationSlug string, newUlid string, rsaPrivateKeyPath, rsaPublicKeyPath, ed25519PrivateKeyPath, ed25519PublicKeyPath string) error {
	folderPath := filepath.Join(constants.RapidBridgeData, "application", applicationSlug)
	filePath := filepath.Join(folderPath, applicationSlug+".json")

	// create folder
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// create file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer file.Close()

	// write to file
	f.CLIConfig.ApplicationDetails.KeyVersion = newUlid

	dataToWrite, err := json.Marshal(f.CLIConfig.ApplicationDetails)
	if err != nil {
		fmt.Println("Error marshalling data", err)
	}
	if err := os.WriteFile(filePath, dataToWrite, 0644); err != nil {
		fmt.Println("Error writing to file", err)
	}

	return f.SaveConfigToFile()
}

func (f *FileConfigAdapter) SaveBankConfigToFile(bankSlug string, rsaPublicKeyPath, ed25519PublicKeyPath string) error {

	folderPath := filepath.Join(constants.RapidBridgeData, "bank", bankSlug)
	filePath := filepath.Join(folderPath, bankSlug+".json")

	// create folder
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	// create file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer file.Close()

	// write to file
	dataToWrite, err := json.Marshal(f.CLIConfig.BankDetails)
	if err != nil {
		fmt.Println("Error marshalling data", err)
	}
	if err := os.WriteFile(filePath, dataToWrite, 0644); err != nil {
		fmt.Println("Error writing to file", err)
	}

	return f.SaveConfigToFile()
}

func (f *FileConfigAdapter) SaveConfigToFile() error {

	var flatCliConfig FlatCLIConfig

	flatCliConfig.RapidLinksURL = f.CLIConfig.RapidLinks.Url
	flatCliConfig.RegisteredApplications = f.CLIConfig.RegisteredApplications
	flatCliConfig.RegisteredBanks = f.CLIConfig.RegisteredBanks

	flatCliConfig.ApplicationEd25519PrivateKeyPath = f.CLIConfig.ApplicationDetails.Ed25519PrivateKeyPath
	flatCliConfig.ApplicationRSAPrivateKeyPath = f.CLIConfig.ApplicationDetails.RSAPrivateKeyPath

	flatCliConfig.BankEd25519PublicKeyPath = f.CLIConfig.BankDetails.Ed25519PublicKeyPath
	flatCliConfig.BankRSAPublicKeyPath = f.CLIConfig.BankDetails.RSAPublicKeyPath

	data, err := json.MarshalIndent(flatCliConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(constants.RapidBridgeData+"/core.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func LoadCLIConfig() (port.CLIConfig, error) {

	var cliConfig CLIConfig

	v := viper.New()

	v.SetConfigName("core")
	v.SetConfigType("json")
	v.AddConfigPath(constants.RapidBridgeData)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	if err := v.Unmarshal(&cliConfig); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	return &FileConfigAdapter{CLIConfig: cliConfig}, nil
}

func LoadApplicationSpecificConfig(applicationSlug string) port.ApplicationDetails {
	configPath := filepath.Join(constants.RapidBridgeData, "application", applicationSlug, applicationSlug+".json")

	configData, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	var applicationDetails port.ApplicationDetails
	if err := json.Unmarshal(configData, &applicationDetails); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	return applicationDetails
}
