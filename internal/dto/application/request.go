package application

type ResourceRequest struct {
	Message  string  `json:"message"`
	TOTPCode *string `json:"totpCode"`
	Username *string `json:"username"`
}
