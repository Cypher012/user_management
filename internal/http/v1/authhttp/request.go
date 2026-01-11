package authhttp

type UserAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
}
