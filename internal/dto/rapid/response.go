package rapid

type RapidResourceResponse struct {
	Data struct {
		From       string `json:"from"`
		To         string `json:"to"`
		Message    string `json:"message"`
		Signature  string `json:"signature"`
		KeyVersion string `json:"key_version"`
	} `json:"data"`
	Error bool `json:"error"`
}
