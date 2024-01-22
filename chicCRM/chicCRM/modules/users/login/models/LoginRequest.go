package models

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequestTeamlead struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
