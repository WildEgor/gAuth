package auth

// OTPLoginRequestDto is a DTO for OTP login request where identity is a phone number or email
type OTPGenerateRequestDto struct {
	Identity string `json:"identity,omitempty"`
}
