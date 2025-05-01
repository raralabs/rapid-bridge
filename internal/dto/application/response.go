package application

type ResourceResponse struct {
	Message string `json:"message"`
}

type OtherResponse struct {
	Data     interface{} `json:"data"`
	MetaData interface{} `json:"meta_data,omitempty"`
}
