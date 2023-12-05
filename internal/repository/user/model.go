package userRepository

type User struct {
	ID       string `json:"uuid"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}
