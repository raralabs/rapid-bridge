package rapid

// type RapidResourceResponse struct {
// 	From      string `json:"from" validate:"required"`
// 	To        string `json:"to" validate:"required"`
// 	Message   string `json:"message" validate:"required"`
// 	Signature string `json:"signature" validate:"required"`
// }

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
