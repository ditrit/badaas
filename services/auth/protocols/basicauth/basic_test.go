package basicauth_test

import (
	"testing"

	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
)

func TestSaltAndHashPassword(t *testing.T) {
	password := "voila"
	hash := basicauth.SaltAndHashPassword(password)
	if string(hash) == password {
		t.Error("the password is not hashed")
	}
}
