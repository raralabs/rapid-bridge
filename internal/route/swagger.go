package route

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Rapid Bridge API
// @version 1.0
// @description API for Rapid Bridge

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide your Bearer token in the format 'Bearer {token}'

// @securityDefinitions.basic BasicAuth
// @description Enter your username and password to authenticate

// @BasePath /api/v1
func swaggerRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.PersistAuthorization(true)))
}
