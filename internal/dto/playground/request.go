package playground

type ApplicationRegisterRequest struct {
	Slug string `json:"slug" validate:"required"`
}
