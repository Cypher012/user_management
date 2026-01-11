package authhttp

type Auth struct {
	Token string `json:"token"`
}

type UserResponse struct {
	Message string `json:"message"`
	Auth    `json:"auth"`
}
