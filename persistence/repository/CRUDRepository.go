package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BadaasID interface {
	int | uuid.UUID
}

// Generic CRUD Repository
type CRUDRepository[T models.Tabler, ID BadaasID] interface {
	// create
	Create(tx *gorm.DB, entity *T) error
	// read
	GetByID(tx *gorm.DB, id ID) (*T, error)
	Get(tx *gorm.DB, conditions map[string]any) (*T, error)
	GetOptional(tx *gorm.DB, conditions map[string]any) (*T, error)
	GetMultiple(tx *gorm.DB, conditions map[string]any) ([]*T, error)
	GetAll(tx *gorm.DB) ([]*T, error)
	Find(tx *gorm.DB, filters squirrel.Sqlizer, pagination pagination.Paginator, sort SortOption) (*pagination.Page[T], error)
	// update
	Save(tx *gorm.DB, entity *T) error
	// delete
	Delete(tx *gorm.DB, entity *T) error
}
