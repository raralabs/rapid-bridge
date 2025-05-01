package application

type ResourceRequest struct {
	Message  string  `json:"message"`
	TOTPCode *string `json:"totpCode"`
	Username *string `json:"username"`

	From       string `json:"from"`
	To         string `json:"to"`
	KeyVersion string `json:"key_version"`
}
