package playground

type ApplicationRegisterResponse struct {
	KeyVersion string `json:"key_version"`
	Slug       string `json:"slug"`

	RSAPublicKey     string
	Ed25519PublicKey string
	Message          string `json:"message"`
}
