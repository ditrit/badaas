package configuration

import "strings"

// maskPassword return a string of * of len(password)
func maskPassword(password string) string {
	nbChars := len(password)
	return strings.Repeat("*", nbChars)
}
