package openid_connect

// This script defines the structures to be used by the backend endpoints

type AuthenticatedJson struct {
	Value string `json:"authenticated"`
}

// To decode the OIDC code sent by the frontend
type Code struct {
	Value string `json:"code"`
}

// OpenIDConnect tokens
type Tokens struct {
	Id_token      string `json:"id_token"`
	Refresh_token string `json:"refresh_token"`
	Access_token  string `json:"access_token"`
}

// Represents a user session
type User struct {
	Code   string
	Email  string
	Tokens Tokens
}
