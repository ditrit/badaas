package repository

import (
	"github.com/ditrit/badaas/persistence/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func (r *EntityTypeRepository) GetByName(tx *gorm.DB, name string) (*models.EntityType, error) {
	var entityType models.EntityType
	err := tx.Preload("Attributes").First(&entityType, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &entityType, nil
}
