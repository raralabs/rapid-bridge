package handler

import (
	"net/http"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/dto/playground"
	errors "rapid-bridge/internal/error"
	"rapid-bridge/internal/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PlaygroundHandler struct {
	logger            port.Logger
	playgroundService *service.PlaygroundService
}

func NewPlaygroundHandler(logger port.Logger, playgroundService *service.PlaygroundService) *PlaygroundHandler {
	return &PlaygroundHandler{
		logger:            logger,
		playgroundService: playgroundService,
	}
}

func (h *PlaygroundHandler) HandleApplicationRegister(c echo.Context) error {
	request := playground.ApplicationRegisterRequest{}

	if err := c.Bind(&request); err != nil {
		h.logger.Error("Validation Error: Request payload does not follow proper format", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 400)
	}
	if err := c.Validate(request); err != nil {
		h.logger.Error("Validation Error: Request payload does not follow proper format", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 400)
	}

	response, err := h.playgroundService.RegisterApplication(request)
	if err != nil {
		h.logger.Error("Failed to register application", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 500)
	}

	return c.JSON(http.StatusOK, response)
}
