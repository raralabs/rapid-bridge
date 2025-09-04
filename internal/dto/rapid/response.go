package rapid

type RapidResourceResponse struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Message    string `json:"message"`
	Signature  string `json:"signature"`
	KeyVersion string `json:"key_version"`
	Error      bool   `json:"error"`
}
