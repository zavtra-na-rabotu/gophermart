package dto

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Token string `json:"token"`
}

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
