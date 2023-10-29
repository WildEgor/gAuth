package auth

type OTPGenerateRequestDto struct {
	Identity string `json:"identity,omitempty"`
}
