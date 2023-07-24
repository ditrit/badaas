package package1

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/multiplepackage/package2"
)

type Package1 struct {
	badorm.UUIDModel

	Package2 package2.Package2 // Package1 HasOne Package2 (Package1 1 -> 1 Package2)
}
