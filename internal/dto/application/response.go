package application

type ResourceResponse struct {
	Message      string `json:"message"`
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type OtherResponse struct {
	Data     interface{} `json:"data"`
	MetaData interface{} `json:"meta_data,omitempty"`
}
