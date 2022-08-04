package models

// Represents a user session
type User struct {
	Code   string
	Email  string
	Tokens Tokens
}

// OpenIDConnect tokens
type Tokens struct {
	Id_token      string `json:"id_token"`
	Refresh_token string `json:"refresh_token"`
	Access_token  string `json:"access_token"`
}
