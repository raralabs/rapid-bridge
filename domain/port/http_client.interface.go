package port

import "net/url"

type HttpResponse struct {
	StatusCode int
	Data       map[string]interface{}
}

type HTTPClient interface {
	POST(baseUrl string, headers map[string]string, queryParam map[string]string, payload any, formData url.Values) (*HttpResponse, error)
	GET(baseUrl string, headers map[string]string, queryParam map[string]string) (*HttpResponse, error)
}
