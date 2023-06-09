package badorm

import (
	"context"
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type UUID uuid.UUID

var NilUUID = UUID(uuid.Nil)

func (id UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "binary(16)"
	case "postgres":
		return "uuid"
	case "sqlite":
		return "varchar(36)"
	case "sqlserver":
		return "uniqueidentifier"
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

func (id UUID) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(id) == 0 {
		return gorm.Expr("NULL")
	}

	switch db.Dialector.Name() {
	case "mysql", "sqlserver":
		binary, _ := id.MarshalBinary()
		return gorm.Expr("?", binary)
	default:
		return gorm.Expr("?", id.String())
	}
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

func NewUUID() UUID {
	return UUID(uuid.New())
}

func ParseUUID(s string) (UUID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return UUID(uuid.Nil), err
	}

	return UUID(uid), nil
}