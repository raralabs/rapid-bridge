package cli

import (
	"fmt"
	"rapid-bridge/constants"
	httpclient "rapid-bridge/internal/adapter/http_client"
	keymanagementfs "rapid-bridge/internal/adapter/keymanagement_fs"
	"rapid-bridge/internal/handler"
	"rapid-bridge/internal/service"

	"rapid-bridge/internal/setup"
	"slices"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var bankSlug string

var initBankCmd = &cobra.Command{
	Use:   "bank",
	Short: "Initialize bank configuration",
	Run: func(cmd *cobra.Command, args []string) {

		app := cmd.Context().Value(constants.Application).(*setup.CLIApplication)

		// check if this bankSlug is already registered in registered banks of cli config
		isBankRegistered := slices.Contains(app.Config.GetRegisteredBanks(), bankSlug)

		if isBankRegistered {
			fmt.Printf("\nThis bank: %s is already registered", bankSlug)
			fmt.Println("\n\nDo you want to re-initialize the bank? \n[1:Yes || 2:No]")
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

		fmt.Println("\nInitializing Bank...")

		fmt.Println("Choose an option:")
		fmt.Println("1) Fetch Bank Public Keys")
		fmt.Println("2) You already have the Bank Public Keys")
		fmt.Print("Enter your choice: ")

		var choice int

		var rsaPublicKeyPath string
		var ed25519PublicKeyPath string

		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("Fetching Bank Public Keys...")

			http_client := httpclient.NewHttpClient(app.Logger)
			keyService := service.NewKeyService(keymanagementfs.NewFSKeyLoader(), keymanagementfs.NewFSKeyConverter(), keymanagementfs.NewFSKeySaver(), http_client, app.Logger, app.Config)
			keyHandler := handler.NewKeyHandler(keyService)
			if err := keyHandler.HandleBankFetchKeys(bankSlug); err != nil {
				app.Logger.Error("Error while fetching bank public keys", zap.String("error", err.Error()))
				return
			}

			fmt.Println("Bank Public Keys fetched successfully")
		case 2:
			fmt.Println("Please provide the path to your bank's public and private key files.")

			fmt.Print("RSA Public key path: ")
			fmt.Scanln(&rsaPublicKeyPath)

			fmt.Print("Ed25519 Public key path: ")
			fmt.Scanln(&ed25519PublicKeyPath)

			http_client := httpclient.NewHttpClient(app.Logger)
			keyService := service.NewKeyService(keymanagementfs.NewFSKeyLoader(), keymanagementfs.NewFSKeyConverter(), keymanagementfs.NewFSKeySaver(), http_client, app.Logger, app.Config)
			keyHandler := handler.NewKeyHandler(keyService)
			if err := keyHandler.HandleBankExistingKeys(bankSlug, rsaPublicKeyPath, ed25519PublicKeyPath); err != nil {
				app.Logger.Error("Error while handling existing rsa and ed25519 keys of bank", zap.String("error", err.Error()))
				return
			}

			fmt.Println("Bank RSA and ED25519 keys loaded successfully")
		default:
			app.Logger.Info("Invalid choice")
		}

		// update config and save to core.json

		app.Config.AddRegisteredBanks(bankSlug)
		app.Config.AddBankDetails(bankSlug)

		app.Config.AddBankKeysPaths(constants.RapidBridgeData+"/bank/"+bankSlug+"/rsa_public_key.pem", constants.RapidBridgeData+"/bank/"+bankSlug+"/ed25519_public_key.pem")

		// TODO: Create a util function to create a file path without manually appending names to a string

		if err := app.Config.SaveBankConfigToFile(bankSlug, rsaPublicKeyPath, ed25519PublicKeyPath); err != nil {
			app.Logger.Error("Error while saving config to file", zap.String("error", err.Error()))
			return
		}

		app.Logger.Info("Bank configuration initialized successfully")
	},
}

func init() {
	initBankCmd.Flags().StringVar(&bankSlug, "slug", "", "Bank slug identifier (required)")
	initBankCmd.MarkFlagRequired("slug")
}
