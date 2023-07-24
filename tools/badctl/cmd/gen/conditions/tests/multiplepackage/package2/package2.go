package package2

import "github.com/ditrit/badaas/badorm"

type Package2 struct {
	badorm.UUIDModel

	Package1ID badorm.UUID // Package1 HasOne Package2 (Package1 1 -> 1 Package2)
}
