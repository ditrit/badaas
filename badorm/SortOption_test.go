package badorm_test

import (
	"testing"

	"github.com/ditrit/badaas/badorm"
	"github.com/stretchr/testify/assert"
)

func TestNewSortOption(t *testing.T) {
	sortOption := badorm.NewSortOption("a", true)
	assert.Equal(t, "a", sortOption.Column())
	assert.True(t, sortOption.Desc())
}
