package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"rapid-bridge/domain/port"
	"strings"
)

type HttpClient struct {
	logger port.Logger
	client *http.Client
}

func (hc HttpClient) POST(baseUrl string, headers map[string]string, queryParam map[string]string, payload any, formData url.Values) (*port.HttpResponse, error) {
	parsedURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	queryParams := parsedURL.Query()
	for k, v := range queryParam {
		queryParams.Add(k, v)
	}
	parsedURL.RawQuery = queryParams.Encode()

	var req *http.Request
	var reqBody io.Reader
	var contentType string

	if formData != nil {
		reqBody = strings.NewReader(formData.Encode())
	} else {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(jsonData)
		contentType = "application/json"
	}

	req, err = http.NewRequest("POST", parsedURL.String(), reqBody)
	if err != nil {
		hc.logger.Error(err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		hc.logger.Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	httpResp, err := io.ReadAll(resp.Body)

	if err != nil {
		hc.logger.Error(err.Error())
		return nil, err
	}

	var mainResponse map[string]any

	if httpResp != nil {
		err = json.Unmarshal(httpResp, &mainResponse)

		if err != nil {
			hc.logger.Error(err.Error())
			return nil, err
		}
	}

	return &port.HttpResponse{StatusCode: resp.StatusCode, Data: mainResponse}, nil
}

func (hc HttpClient) GET(baseUrl string, headers map[string]string, queryParam map[string]string) (*port.HttpResponse, error) {
	parsedURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	queryParams := parsedURL.Query()
	for k, v := range queryParam {
		queryParams.Add(k, v)
	}
	parsedURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		hc.logger.Error(err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		hc.logger.Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	httpResp, err := io.ReadAll(resp.Body)

	if err != nil {
		hc.logger.Error(err.Error())
		return nil, err
	}

	var mainResponse map[string]any

	if httpResp != nil {
		err = json.Unmarshal(httpResp, &mainResponse)

		if err != nil {
			hc.logger.Error(err.Error())
			return nil, err
		}
	}

	return &port.HttpResponse{StatusCode: resp.StatusCode, Data: mainResponse}, nil
}

func NewHttpClient(logger port.Logger) HttpClient {
	return HttpClient{
		logger: logger,
		client: &http.Client{},
	}
}
