package badorm_test

import (
	"testing"

	"github.com/ditrit/badaas/badorm"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseCorrectUUID(t *testing.T) {
	uuidString := uuid.New().String()
	uuid, err := badorm.ParseUUID(uuidString)
	assert.Nil(t, err)
	assert.Equal(t, uuidString, uuid.String())
}

func TestParseIncorrectUUID(t *testing.T) {
	uid, err := badorm.ParseUUID("not uuid")
	assert.Error(t, err)
	assert.Equal(t, badorm.NilUUID, uid)
}
