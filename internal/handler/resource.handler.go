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
		return errors.NewRapidLinksError(response.Message, response.StatusCode)
	}

	if response.StatusCode != 200 {
		r.logger.Error("Failed to handle resource", zap.Int("status_code", response.StatusCode), zap.String("message", response.Message))
		return errors.NewRapidLinksError(response.Message, response.StatusCode)
	}

	messageVal, err := unmarshalAndUnwrap(response.Message)
	if err != nil {
		r.logger.Error("Failed to unmarshal response", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 500)
	}

	respBody := map[string]interface{}{
		"message": messageVal,
	}
	if err := c.JSON(response.StatusCode, respBody); err != nil {
		r.logger.Error("Failed to send response", zap.String("error", err.Error()))
		return errors.NewRapidLinksError(err.Error(), 500)
	}

	return nil
}

func unmarshalAndUnwrap(raw string) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if inner, exists := m["message"]; exists {
			return inner, nil
		}
	}

	return data, nil
}

func NewRapidResourceHandler(logger port.Logger, service *service.RapidResourceService) *resourceHandler {
	return &resourceHandler{
		logger:               logger,
		RapidResourceService: service,
	}
}
