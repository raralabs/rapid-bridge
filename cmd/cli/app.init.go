package cli

import (
	"context"
	"fmt"
	"rapid-bridge/constants"
	"rapid-bridge/internal/adapter/config"
	keymanagementfs "rapid-bridge/internal/adapter/keymanagement_fs"
	keyhandler "rapid-bridge/internal/handler"
	"rapid-bridge/internal/service"
	"rapid-bridge/internal/setup"
	"rapid-bridge/pkg/util"
	"slices"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var applicationSlug string
var encryptionKeyValidityPeriod int
var signingKeyValidityPeriod int
var applicationInitializationReason string

var initAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Initialize app configuration",
	Run: func(cmd *cobra.Command, args []string) {

		rootCtx := cmd.Context()
		app := rootCtx.Value(constants.Application).(*setup.CLIApplication)

		// set ulid in context
		ulid := util.GenerateULID().String()

		rootCtx = context.WithValue(rootCtx, constants.ApplicationUlid, ulid)

		// check if this applicationSlug is already registered in registered applications of cli config
		isApplicationRegistered := slices.Contains(app.Config.GetRegisteredApplications(), applicationSlug)

		if isApplicationRegistered {
			fmt.Printf("\nThis application: %s is already registered", applicationSlug)

			applicationConfig := config.LoadApplicationSpecificConfig(applicationSlug)
			if applicationConfig.RSAKeysValidUntil.After(time.Now()) && applicationConfig.Ed25519KeysValidUntil.After(time.Now()) {
				fmt.Println("\nYour encryption (RSA) and signing (ED25519) keys are still valid.")
				// fmt.Println("\nDo you want to rotate your keys? ")
			}

			fmt.Println("\n\nDo you want to re-initialize the application? \n[1:Yes || 2:No]")
			fmt.Print("\nEnter your choice: ")

			var choice int
			fmt.Scanln(&choice)

			if choice == 2 {
				return
			} else if choice != 1 {
				fmt.Println("Invalid choice")
				return
			}
		}

		fmt.Println("\nInitializing application...")

		fmt.Println("\nChoose an option:")
		fmt.Println("1) Generate a new key pair")
		fmt.Println("2) Use your own existing key pair")
		fmt.Print("Enter your choice: ")

		var choice int

		var rsaPublicKeyPath string
		var rsaPrivateKeyPath string
		var ed25519PublicKeyPath string
		var ed25519PrivateKeyPath string

		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("Generating new key pair...")

			encryptionKeyValidityPeriod = constants.EncryptionKeyValidityPeriod
			signingKeyValidityPeriod = constants.SigningKeyValidityPeriod

			keyLoader := keymanagementfs.NewFSKeyLoader()
			keySaver := keymanagementfs.NewFSKeySaver()
			keyConverter := keymanagementfs.NewFSKeyConverter()
			keyService := service.NewKeyService(keyLoader, keyConverter, keySaver, nil, app.Logger, app.Config)
			keyHandler := keyhandler.NewKeyHandler(keyService)
			if err := keyHandler.HandleApplicationGenerateKeyPair(applicationSlug, ulid); err != nil {
				app.Logger.Error("Error while generating key pair", zap.String("error", err.Error()))
				return
			}
			fmt.Println("Key pair generated and saved successfully")
		case 2:
			fmt.Println("Please provide the path to your application's public and private key files.")

			fmt.Print("RSA Public key path: ")
			fmt.Scanln(&rsaPublicKeyPath)

			fmt.Print("RSA Private key path: ")
			fmt.Scanln(&rsaPrivateKeyPath)

			fmt.Print("Ed25519 Public key path: ")
			fmt.Scanln(&ed25519PublicKeyPath)

			fmt.Print("Ed25519 Private key path: ")
			fmt.Scanln(&ed25519PrivateKeyPath)

			encryptionKeyValidityPeriod = constants.EncryptionKeyValidityPeriod
			signingKeyValidityPeriod = constants.SigningKeyValidityPeriod

			keyLoader := keymanagementfs.NewFSKeyLoader()
			keySaver := keymanagementfs.NewFSKeySaver()
			keyConverter := keymanagementfs.NewFSKeyConverter()
			keyService := service.NewKeyService(keyLoader, keyConverter, keySaver, nil, app.Logger, app.Config)
			keyHandler := keyhandler.NewKeyHandler(keyService)
			if err := keyHandler.HandleApplicationExistingKeyPair(applicationSlug, ulid, rsaPrivateKeyPath, rsaPublicKeyPath, ed25519PrivateKeyPath, ed25519PublicKeyPath); err != nil {
				app.Logger.Error("Error while handling existing key pair", zap.String("error", err.Error()))
				return
			}

			fmt.Println("Key pair loaded successfully")
		default:
			app.Logger.Info("Invalid choice")
		}

		if !isApplicationRegistered {
			app.Config.AddRegisteredApplications(applicationSlug)
		}

		app.Config.AddApplicationSlug(applicationSlug)
		app.Config.AddApplicationKeysPaths(constants.RapidBridgeData+"/application/"+applicationSlug+"/"+ulid+"/rsa_public_key.pem", constants.RapidBridgeData+"/application/"+applicationSlug+"/"+ulid+"/rsa_private_key.pem", constants.RapidBridgeData+"/application/"+applicationSlug+"/"+ulid+"/ed25519_public_key.pem", constants.RapidBridgeData+"/application/"+applicationSlug+"/"+ulid+"/ed25519_private_key.pem")
		app.Config.AddKeysValidityPeriod(encryptionKeyValidityPeriod, signingKeyValidityPeriod)

		fmt.Println("Do you want to use these keys for your subsequent requests? [1:Yes || 2:No]")
		fmt.Print("Enter your choice: ")

		var useNewKeys int
		fmt.Scanln(&useNewKeys)

		if useNewKeys == 2 {
			fmt.Println("Skipping...")
		} else if useNewKeys != 1 {
			fmt.Println("Invalid choice")
		} else if useNewKeys == 1 {
			app.Config.AddApplicationUlid(ulid)
		}

		if err := app.Config.SaveApplicationConfigToFile(applicationSlug, ulid, util.GetRSAPrivateKeyPath(applicationSlug, ulid, rsaPrivateKeyPath), util.GetRSAPublicKeyPath(applicationSlug, ulid, rsaPublicKeyPath), util.GetEd25519PrivateKeyPath(applicationSlug, ulid, ed25519PrivateKeyPath), util.GetEd25519PublicKeyPath(applicationSlug, ulid, ed25519PublicKeyPath)); err != nil {
			app.Logger.Error("Error while saving config", zap.String("error", err.Error()))
			return
		}

		app.Logger.Info("Application configuration initialized successfully")
	},
}

func init() {
	initAppCmd.Flags().StringVar(&applicationSlug, "slug", "", "App slug identifier (required)")
	initAppCmd.MarkFlagRequired("slug")
}
