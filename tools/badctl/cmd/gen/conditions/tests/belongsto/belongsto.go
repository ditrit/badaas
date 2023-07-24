package belongsto

import "github.com/ditrit/badaas/badorm"

type Owner struct {
	badorm.UUIDModel
}
type Owned struct {
	badorm.UUIDModel

	// Owned belongsTo Owner (Owned 0..* -> 1 Owner)
	Owner   Owner
	OwnerID badorm.UUID
}
