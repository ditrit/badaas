package badorm

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// supported types for model identifier
type BadaasID interface {
	uint | UUID
}

// Base Model for gorm
//
// Every model intended to be saved in the database must embed this badorm.UUIDModel
// reference: https://gorm.io/docs/models.html#gorm-Model
type UUIDModel struct {
	ID        UUID `gorm:"primarykey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (model *UUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	if model.ID == UUID(uuid.Nil) {
		model.ID = UUID(uuid.New())
	}
	return nil
}

type UIntModel gorm.Model

type UUID uuid.UUID

func (id UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "varchar(36)"
		// return "binary(16)"
	case "postgres":
		return "uuid"
	}
	return ""
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}

func (id UUID) URN() string {
	return uuid.UUID(id).URN()
}

func (id UUID) Variant() uuid.Variant {
	return uuid.UUID(id).Variant()
}

func (id UUID) Version() uuid.Version {
	return uuid.UUID(id).Version()
}

func (id UUID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

func (id *UUID) UnmarshalText(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalText(data)
}

func (id UUID) MarshalBinary() ([]byte, error) {
	return uuid.UUID(id).MarshalBinary()
}

func (id *UUID) UnmarshalBinary(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalBinary(data)
}

func (id *UUID) Scan(src interface{}) error {
	return (*uuid.UUID)(id).Scan(src)
}

func (id UUID) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}

func (id UUID) Time() uuid.Time {
	return uuid.UUID(id).Time()
}

func (id UUID) ClockSequence() int {
	return uuid.UUID(id).ClockSequence()
}
