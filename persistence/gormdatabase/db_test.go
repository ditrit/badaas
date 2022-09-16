package gormdatabase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDsnFromconf(t *testing.T) {
	assert.NotEmpty(t, createDsnFromConf(), "no dsn should be empty")
}

func TestCreateDsn(t *testing.T) {
	assert.NotEmpty(t,
		createDsn(
			"192.168.2.5",
			"username",
			"password",
			"disable",
			"badaas_db",
			1225,
		),
		"no dsn should be empty",
	)
}

func TestInitializeDBFromDsn(t *testing.T) {
	db, err := InitializeDBFromDsn(createDsn(
		"192.168.2.5",
		"username",
		"password",
		"disable",
		"badaas_db",
		1225,
	))
	if err == nil {
		t.Errorf("should return an error on invalid dsn")
	}
	if db != nil {
		t.Errorf("should return a null value")
	}
}
