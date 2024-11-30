package models

type User struct {
	Guid     string `json:"guid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
