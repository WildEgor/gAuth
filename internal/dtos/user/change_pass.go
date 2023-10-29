package user

type ChangePasswordRequestDto struct {
	OldPassword string `json:"old,omitempty"`
	NewPassword string `json:"new,omitempty"`
}
