package model

// UserCreateIn is the input for creating a user.
type UserCreateIn struct {
	Username string `json:"username" v:"required"`
	Password string `json:"password" v:"required"`
}
