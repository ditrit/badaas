package repository

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/elliotchance/pie/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EntityRepository struct {
	logger          *zap.Logger
	db              *gorm.DB
	valueRepository *ValueRepository
}

func NewEntityRepository(
	logger *zap.Logger,
	db *gorm.DB,
	valueRepository *ValueRepository,
) *EntityRepository {
	return &EntityRepository{
		logger:          logger,
		db:              db,
		valueRepository: valueRepository,
	}
}

// Creates an entity and its values in the database
// must be used in place of gorm's db.Save(entity) because of the bug
// when using gorm with cockroachDB. For more info refer to:
// https://github.com/FrancoLiberali/cockroachdb_gorm_bug
func (r *EntityRepository) Create(entity *models.Entity) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		query, values, err := sq.Insert("entities").
			Columns("created_at", "updated_at", "entity_type_id").
			Values(now, now, entity.EntityType.ID).
			Suffix("RETURNING \"id\"").
			PlaceholderFormat(sq.Dollar).ToSql()

		if err != nil {
			return err
		}

		var result string
		err = tx.Raw(query, values...).Scan(&result).Error
		if err != nil {
			return err
		}

		uuid, err := uuid.Parse(result)
		if err != nil {
			return err
		}

		pie.Each(entity.Fields, func(value *models.Value) {
			value.EntityID = uuid
		})

		if len(entity.Fields) > 0 {
			err = r.valueRepository.Create(tx, entity.Fields)
			if err != nil {
				return err
			}
		}

		entity.ID = uuid
		return nil
	})
}
