package dto

// Data Transfert Object Package

// Describe the login payload
type DTOLoginJSONStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
