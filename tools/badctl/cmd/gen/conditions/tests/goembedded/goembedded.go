package goembedded

import "github.com/ditrit/badaas/badorm"

type ToBeEmbedded struct {
	Int int
}

type GoEmbedded struct {
	badorm.UIntModel

	Int int
	ToBeEmbedded
}
