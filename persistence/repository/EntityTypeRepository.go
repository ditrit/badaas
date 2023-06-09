package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/models"
)

type EntityTypeRepository struct {
	logger *zap.Logger
}

func NewEntityTypeRepository(
	logger *zap.Logger,
) *EntityTypeRepository {
	return &EntityTypeRepository{
		logger: logger,
	}
}

func (r *EntityTypeRepository) Get(tx *gorm.DB, id badorm.UUID) (*models.EntityType, error) {
	var entityType models.EntityType

	err := tx.Preload("Attributes").First(&entityType, id).Error
	if err != nil {
		return nil, err
	}

	return &entityType, nil
}

func (r *EntityTypeRepository) GetByName(tx *gorm.DB, name string) (*models.EntityType, error) {
	var entityType models.EntityType

	err := tx.Preload("Attributes").First(&entityType, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &entityType, nil
}
