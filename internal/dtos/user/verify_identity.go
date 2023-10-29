package user

type VerifyIdentityRequestDto struct {
	Identity string `json:"identity,omitempty"`
	Code     string `json:"code,omitempty"`
}
