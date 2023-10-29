package auth

type OTPLoginRequestDto struct {
	Phone string `json:"phone,omitempty"`
	Code  string `json:"code,omitempty"`
}
