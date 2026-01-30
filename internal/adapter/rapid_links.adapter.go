package adapter

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/adapter/config"
	"rapid-bridge/internal/dto/rapid"

	"go.uber.org/zap"
)

func SendRequestToRapidLinks(logger port.Logger, rapidLinksUrl string, urlPath string, payload rapid.RapidResourceRequest, header http.Header) (rapid.RapidResourceResponse, error) {
	var response rapid.RapidResourceResponse

	if rapidLinksUrl == "" {
		bankConfig := config.LoadBankSpecificConfig(payload.To)
		rapidLinksUrl = bankConfig.BankRapidAPIUrl
	}

	if rapidLinksUrl == "" {
		logger.Debug("Rapid URL configuration missing")
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Request send to rapid: Error while marshalling json payload", zap.String("error", err.Error()))
		return response, err
	}

	req, err := http.NewRequest("POST", rapidLinksUrl+urlPath, bytes.NewBuffer(jsonPayload))
	if err != nil {
		logger.Error("Request send to rapid: Error while creating new http request to %v", rapidLinksUrl, zap.String("error", err.Error()))
		return response, err
	}

	for name, values := range header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Request send to rapid: Error while sending request to rapid", zap.String("error", err.Error()))
		return response, err
	}
	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Request send to rapid: Error while reading the response body", zap.String("error", err.Error()))
		return response, err
	}

	if err := json.Unmarshal(responseBodyBytes, &response); err != nil {
		logger.Error("Request send to rapid: Error while unmarshalling the response body", zap.String("error", err.Error()))
		return response, err
	}

	return response, nil
}
