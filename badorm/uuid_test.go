package badorm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/badorm"
)

func TestParseCorrectUUID(t *testing.T) {
	uuidString := badorm.NewUUID().String()
	uuid, err := badorm.ParseUUID(uuidString)
	assert.Nil(t, err)
	assert.Equal(t, uuidString, uuid.String())
}

func TestParseIncorrectUUID(t *testing.T) {
	uuid, err := badorm.ParseUUID("not uuid")
	assert.Error(t, err)
	assert.Equal(t, badorm.NilUUID, uuid)
}
