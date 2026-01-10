package authhttp

type Auth struct {
	Atk string `json:"token"`
	Rtk string `json:"rtk"`
}

type UserResponse struct {
	Message string `json:"message"`
	Auth
}
