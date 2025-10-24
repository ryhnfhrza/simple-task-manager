// model/web/user_request.go
package web

type UserRegisterRequest struct {
	Username string `json:"username" validate:"required,alphanum,min=5,max=30"`
	Password string `json:"password" validate:"required,passComplex"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
