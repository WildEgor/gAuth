package auth

type LoginRequestDto struct {
	Login    string `json:"login" validate:"required,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}
