package cli

import (
	"context"
	"fmt"
	"rapid-bridge/constants"
	keymanagementfs "rapid-bridge/internal/adapter/keymanagement_fs"
	keyhandler "rapid-bridge/internal/handler"
	"rapid-bridge/internal/service"
	"rapid-bridge/internal/setup"
	"rapid-bridge/pkg/util"
	"slices"

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

		app.Config.AddApplicationUlid(ulid)

		if err := app.Config.SaveApplicationConfigToFile(); err != nil {
			app.Logger.Error("Error while saving config", zap.String("error", err.Error()))
			return
		}

		app.Logger.Info("Application configuration saved successfully")

	},
}

func init() {
	initAppCmd.Flags().StringVar(&applicationSlug, "slug", "", "App slug identifier (required)")
	initAppCmd.MarkFlagRequired("slug")
}
