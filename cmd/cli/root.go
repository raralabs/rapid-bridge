package cli

import (
	"context"
	"log"
	"os"
	server "rapid-bridge/cmd/server"
	"rapid-bridge/constants"
	"rapid-bridge/internal/setup"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "rapid-bridge",
	Short: "Rapid Bridge CLI - Backend utility",
	Long:  `Rapid Bridge is a CLI tool for backend initialization and management.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		app := setup.NewCLIApplication()
		ctx := context.WithValue(cmd.Context(), constants.Application, app)
		cmd.SetContext(ctx)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Rapid Bridge",
	Long:  `Initialize Rapid Bridge.`,
}

func init() {

	initCmd.AddCommand(initAppCmd)
	initCmd.AddCommand(initBankCmd)
	initCmd.AddCommand(server.InitServerCmd)

	RootCmd.AddCommand(initCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
