package server

import (
	"net/http"
	"rapid-bridge/internal/route"
	"rapid-bridge/internal/setup"
	"rapid-bridge/pkg/util"

	rmiddleware "rapid-bridge/pkg/middleware"

	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"

	"github.com/labstack/echo/v4"
)

var InitServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Initialize backend server configuration",
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}

func StartServer() {

	app := setup.NewApplication()
	defer app.Logger.Sync()

	e := echo.New()
	e.Validator = util.NewCustomValidator()

	e.Use(middleware.Secure())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(rmiddleware.CreateEchoLogger(app.Logger))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
	}))

	route.SetupRoutes(e, app)

	// // call rapid url to sync public keys of bank and application
	// rapidResourceService := service.NewRapidResourceService(app)
	// rapidResourceService.SyncApplicationPublicKeys()
	// rapidResourceService.SyncBankPublicKeys()

	app.Logger.Info("Server started successfully", app.Config.GetRapidLinksUrl())
	e.Start(":8080")
}
