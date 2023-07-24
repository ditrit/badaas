package gormembedded

import "github.com/ditrit/badaas/badorm"

type ToBeGormEmbedded struct {
	Int int
}

type GormEmbedded struct {
	badorm.UIntModel

	Int                  int
	GormEmbedded         ToBeGormEmbedded `gorm:"embedded;embeddedPrefix:gorm_embedded_"`
	GormEmbeddedNoPrefix ToBeGormEmbedded `gorm:"embedded"`
}
