package badorm

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// supported types for model identifier
type BadaasID interface {
	uint | uuid.UUID
}

// Base Model for gorm
//
// Every model intended to be saved in the database must embed this badorm.UUIDModel
// reference: https://gorm.io/docs/models.html#gorm-Model
type UUIDModel struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UIntModel gorm.Model
