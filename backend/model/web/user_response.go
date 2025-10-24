package web

type UserResponse struct {
	Username string `json:"username"`
}

type UserLoginResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
