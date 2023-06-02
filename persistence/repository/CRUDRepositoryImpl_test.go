package repository

import (
	"testing"

	"github.com/Masterminds/squirrel"
	mocks "github.com/ditrit/badaas/mocks/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseError(t *testing.T) {
	err := DatabaseError("test err", assert.AnError)
	require.NotNil(t, err)
	assert.True(t, err.Log())
}

type dumbModel struct{}

func (dumbModel) TableName() string {
	return "dumb_models"
}

func TestNewRepository(t *testing.T) {
	paginationConfiguration := mocks.NewPaginationConfiguration(t)
	dumbModelRepository := NewCRUDRepository[dumbModel, uint](nil, paginationConfiguration)
	assert.NotNil(t, dumbModelRepository)
}

func TestCompileSql_NoError(t *testing.T) {
	paginationConfiguration := mocks.NewPaginationConfiguration(t)
	dumbModelRepository := &CRUDRepositoryImpl[dumbModel, uint]{
		gormDatabase:            nil,
		paginationConfiguration: paginationConfiguration,
	}
	_, _, err := dumbModelRepository.compileSQL(squirrel.Eq{"name": "qsdqsd"})
	assert.Nil(t, err)
}

func TestCompileSql_Err(t *testing.T) {
	paginationConfiguration := mocks.NewPaginationConfiguration(t)
	dumbModelRepository := &CRUDRepositoryImpl[dumbModel, uint]{
		gormDatabase:            nil,
		paginationConfiguration: paginationConfiguration,
	}
	_, _, err := dumbModelRepository.compileSQL(squirrel.GtOrEq{"name": nil})

	assert.Error(t, err)
}
