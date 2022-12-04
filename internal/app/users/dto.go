package users

type UserResponse struct {
	Id              uint   `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	EmailVerifiedAt string `json:"email_verified_at"`
	GroupId         uint   `json:"group_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type PatchUser struct {
	Email string `json:"email" validate:"email"`
}
