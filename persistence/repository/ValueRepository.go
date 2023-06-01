package repository

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/elliotchance/pie/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ValueRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewValueRepository(
	logger *zap.Logger,
	db *gorm.DB,
) *ValueRepository {
	return &ValueRepository{
		logger: logger,
		db:     db,
	}
}

// Creates multiples values in the database
func (r *ValueRepository) Create(tx *gorm.DB, values []*models.Value) error {
	now := time.Now()

	pie.Each(values, func(value *models.Value) {
		value.ID = badorm.UUID(uuid.New())
	})

	query := sq.Insert("values_").
		Columns("id", "created_at", "updated_at", "is_null", "string_val", "float_val", "int_val",
			"bool_val", "relation_val", "entity_id", "attribute_id")

	for _, value := range values {
		query = query.Values(value.ID, now, now, value.IsNull, value.StringVal,
			value.FloatVal, value.IntVal, value.BoolVal,
			value.RelationVal, value.EntityID, value.Attribute.ID)
	}

	queryString, queryValues, err := query.
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	err = tx.Exec(queryString, queryValues...).Error
	if err != nil {
		return err
	}

	return nil
}
