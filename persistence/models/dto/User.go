package dto

// Data Transfert Object Package

// Login DTO
type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Output DTO
type UserDTO struct {
	BaseModelDTO
	Username string `json:"username"`
	Email    string `json:"email"`
}
