package auth

type OTPLoginRequestDto struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
