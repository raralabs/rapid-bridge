package server

import (
	"net/http"
	customErr "rapid-bridge/internal/error"
	"rapid-bridge/internal/route"
	"rapid-bridge/internal/setup"
	"rapid-bridge/pkg/config"
	"rapid-bridge/pkg/util"

	"go.uber.org/zap"

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

	config, err := config.LoadConfig()
	if err != nil {
		app.Logger.Fatal("Failed to load config", zap.Error(err))
	}

	e := echo.New()
	e.Validator = util.NewCustomValidator()

	// Custom error handler for RapidLinksError
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if rapidErr, ok := err.(customErr.RapidLinksError); ok {
			c.JSON(rapidErr.GetStatusCode(), map[string]interface{}{
				"message": rapidErr.Message,
				"code":    rapidErr.GetStatusCode(),
			})
			return
		}
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Use(middleware.Secure())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(rmiddleware.CreateEchoLogger(app.Logger))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
	}))

	route.SetupRoutes(e, app)

	app.Logger.Info("Server started successfully")
	e.Start(config.ServerPort)
}
