package userservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc     string
		username string
		email    string
		password string
		err      bool
	}{
		{
			desc:     "non valid email",
			username: "whatever",
			email:    "whatever",
			password: "1234",
			err:      true,
		},
		{
			desc:     "valid email",
			username: "whatever",
			email:    "whatever@email.com",
			password: "1234",
			err:      false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			user, err := NewUser(tC.username, tC.email, tC.password)
			if tC.err {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
			assert.Equal(t, "whatever", user.Username)
			assert.NotEqual(t, "1234", user.Password)
		})
	}
}
