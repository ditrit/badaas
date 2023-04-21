package repository

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/persistence/models"
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

func (r *ValueRepository) Create(tx *gorm.DB, values []*models.Value) error {
	now := time.Now()

	query := sq.Insert("values").
		Columns("created_at", "updated_at", "is_null", "string_val", "float_val", "int_val",
			"bool_val", "relation_val", "entity_id", "attribute_id")

	for _, value := range values {
		query = query.Values(now, now, value.IsNull, value.StringVal,
			value.FloatVal, value.IntVal, value.BoolVal,
			value.RelationVal, value.EntityID, value.Attribute.ID)
	}

	queryString, queryValues, err := query.
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	var results []string
	err = tx.Raw(queryString, queryValues...).Scan(&results).Error
	if err != nil {
		return err
	}

	for i, result := range results {
		uuid, err := uuid.Parse(result)
		if err != nil {
			return err
		}

		values[i].ID = uuid
	}

	return nil
}
