package user

type ChangePasswordRequestDto struct {
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
}
