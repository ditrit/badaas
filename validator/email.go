package validator

import "net/mail"

// Validate email
func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
