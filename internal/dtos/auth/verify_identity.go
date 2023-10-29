package auth

type VerifyIdentityRequestDto struct {
	Identity string `json:"identity,omitempty"`
	Code     string `json:"Code,omitempty"`
}
