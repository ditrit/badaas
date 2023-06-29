package hasone

import "github.com/ditrit/badaas/badorm"

type Country struct {
	badorm.UUIDModel

	Capital City // Country HasOne City (Country 1 -> 1 City)
}

type City struct {
	badorm.UUIDModel

	Country   *Country
	CountryID badorm.UUID // Country HasOne City (Country 1 -> 1 City)
}
