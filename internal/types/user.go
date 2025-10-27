package types

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserSignup struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
