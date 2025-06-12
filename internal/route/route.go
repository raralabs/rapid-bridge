package route

import (
	"rapid-bridge/domain/security"
	keymanagementfs "rapid-bridge/internal/adapter/keymanagement_fs"
	securityadapter "rapid-bridge/internal/adapter/security"
	"rapid-bridge/internal/handler"
	"rapid-bridge/internal/service"
	"rapid-bridge/internal/setup"
	"rapid-bridge/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *setup.Application) {

	swaggerRoutes(e)

	api := e.Group("/api/v1")

	rapidResource := api.Group("/resource", middleware.APIContractMiddleware())
	resourceForwardingRoutes(rapidResource, app)

	// Route to register new application in bridge
	// This is just for playground and not for production
	cliApp := setup.NewCLIApplication()
	keyLoader := keymanagementfs.NewFSKeyLoader()
	keyConverter := keymanagementfs.NewFSKeyConverter()
	keySaver := keymanagementfs.NewFSKeySaver()
	keyService := service.NewKeyService(keyLoader, keyConverter, keySaver, nil, cliApp.Logger, cliApp.Config)
	playgroundService := service.NewPlaygroundService(cliApp.Logger, cliApp, keyLoader, keyConverter, keySaver, keyService)
	playgroundHandler := handler.NewPlaygroundHandler(cliApp.Logger, playgroundService)
	api.POST("/application/register", playgroundHandler.HandleApplicationRegister)

}

func resourceForwardingRoutes(resourceRoutes *echo.Group, app *setup.Application) {

	newCipher := securityadapter.NewHybridCryptography()
	newSecurity := security.NewSecurity(newCipher)

	keyLoader := keymanagementfs.NewFSKeyLoader()

	service := service.NewRapidResourceService(keyLoader, *newSecurity, app.Logger, app.Config)
	handler := handler.NewRapidResourceHandler(app.Logger, service)

	resourceRoutes.POST("/balance", handler.HandleResource)
	resourceRoutes.POST("/statement", handler.HandleResource)
	resourceRoutes.POST("/payment/initiate", handler.HandleResource)
	resourceRoutes.POST("/payment/approve", handler.HandleResource)
	resourceRoutes.POST("/account/open", handler.HandleResource)
}
