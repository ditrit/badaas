package goembedded

import "github.com/ditrit/badaas/badorm"

type ToBeEmbedded struct {
	EmbeddedInt int
}

type GoEmbedded struct {
	badorm.UIntModel

	ToBeEmbedded
}
