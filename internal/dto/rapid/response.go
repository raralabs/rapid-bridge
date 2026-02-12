package rapid

type RapidResourceResponse struct {
	From          string `json:"from"`
	To            string `json:"to"`
	Message       string `json:"message"`
	Signature     string `json:"signature"`
	KeyVersion    string `json:"key_version"`
	Error         bool   `json:"error"`
	StatusCode    int    `json:"status_code"`
	EncryptedFlag bool   `json:"encrypted_flag"`

	// Nested data (for second format)
	Data *struct {
		From       string `json:"from"`
		To         string `json:"to"`
		Message    string `json:"message"`
		Signature  string `json:"signature"`
		KeyVersion string `json:"key_version"`
	} `json:"data,omitempty"`
}

func (r *RapidResourceResponse) GetFrom() string {
	if r.Data != nil {
		return r.Data.From
	}
	return r.From
}

func (r *RapidResourceResponse) GetTo() string {
	if r.Data != nil {
		return r.Data.To
	}
	return r.To
}

func (r *RapidResourceResponse) GetMessage() string {
	if r.Data != nil {
		return r.Data.Message
	}
	return r.Message
}

func (r *RapidResourceResponse) GetSignature() string {
	if r.Data != nil {
		return r.Data.Signature
	}
	return r.Signature
}

func (r *RapidResourceResponse) GetKeyVersion() string {
	if r.Data != nil {
		return r.Data.KeyVersion
	}
	return r.KeyVersion
}
