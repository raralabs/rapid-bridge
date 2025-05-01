package handler

import (
	"encoding/json"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/dto/application"
	errors "rapid-bridge/internal/error"
	service "rapid-bridge/internal/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type resourceHandler struct {
	logger               port.Logger
	RapidResourceService *service.RapidResourceService
}

func (r *resourceHandler) HandleResource(c echo.Context) error {
	request := application.ResourceRequest{}

	if err := c.Bind(&request); err != nil {
		r.logger.Error("Validation Error: Request payload does not follow proper format", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 400)
	}
	if err := c.Validate(request); err != nil {
		r.logger.Error("Validation Error: Request payload does not follow proper format", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 400)
	}

	response, err := r.RapidResourceService.HandleResource(c, request)
	if err != nil {
		r.logger.Error("Failed to handle resource", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 500)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response.Message), &data)
	if err != nil {
		r.logger.Error("Failed to unmarshal request", zap.String("error", err.Error()))
		return err
	}

	if err := c.JSON(200, data); err != nil {
		r.logger.Error("Failed to send response", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 500)
	}
	return nil
}

func NewRapidResourceHandler(logger port.Logger, service *service.RapidResourceService) *resourceHandler {
	return &resourceHandler{
		logger:               logger,
		RapidResourceService: service,
	}
}
