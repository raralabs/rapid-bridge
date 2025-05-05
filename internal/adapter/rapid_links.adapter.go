package adapter

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/dto/rapid"

	"go.uber.org/zap"
)

func SendRequestToRapidLinks(logger port.Logger, rapidLinksUrl string, urlPath string, payload rapid.RapidResourceRequest, header http.Header) (rapid.RapidResourceResponse, error) {
	var response rapid.RapidResourceResponse

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Request send to rapid links: Error while marshalling json payload", zap.String("error", err.Error()))
		return response, err
	}

	req, err := http.NewRequest("POST", rapidLinksUrl+urlPath, bytes.NewBuffer(jsonPayload))
	if err != nil {
		logger.Error("Request send to rapid links: Error while creating new http request to %v", rapidLinksUrl, zap.String("error", err.Error()))
		return response, err
	}

	for name, values := range header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Request send to rapid links: Error while sending request to rapid links", zap.String("error", err.Error()))
		return response, err
	}
	defer resp.Body.Close()

	logger.Info("Successfully called to Rapid Links", zap.String("url", rapidLinksUrl+urlPath))

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Request send to rapid links: Error while reading the response body", zap.String("error", err.Error()))
		return response, err
	}

	if err := json.Unmarshal(responseBodyBytes, &response); err != nil {
		logger.Error("Request send to rapid links: Error while unmarshalling the response body", zap.String("error", err.Error()))
		return response, err
	}

	return response, nil
}
