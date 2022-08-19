package openid_connect

// This script defines the structures to be used by the backend endpoints

type AuthenticatedJson struct {
	Value string `json:"authenticated"`
}

// To decode the OIDC code sent by the frontend
type Code struct {
	Value string `json:"code"`
}
