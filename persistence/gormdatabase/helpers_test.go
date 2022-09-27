package gormdatabase_test

import (
	"errors"
	"testing"

	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/jackc/pgconn"
	"github.com/magiconair/properties/assert"
)

func TestIsDuplicateError(t *testing.T) {
	testCases := []struct {
		desc        string
		err         error
		isDuplicate bool
	}{
		{
			desc:        "classic error",
			err:         errors.New("voila"),
			isDuplicate: false,
		},
		{
			desc: "pg error not duplicate error",
			err: &pgconn.PgError{
				Code: "235252551514",
			},
			isDuplicate: false,
		},
		{
			desc: "pg error duplicate error",
			err: &pgconn.PgError{
				Code: "23505",
			},
			isDuplicate: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.isDuplicate, gormdatabase.IsDuplicateKeyError(tC.err))
		})
	}
}
