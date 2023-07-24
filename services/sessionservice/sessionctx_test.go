package sessionservice

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/badorm"
)

func TestSessionCtx(t *testing.T) {
	ctx := context.Background()
	sessionClaims := &SessionClaims{badorm.NilUUID, badorm.NewUUID()}
	ctx = SetSessionClaimsContext(ctx, sessionClaims)
	claims := GetSessionClaimsFromContext(ctx)
	assert.Equal(t, badorm.NilUUID, claims.UserID)
}

func TestSessionCtxPanic(t *testing.T) {
	ctx := context.Background()

	assert.Panics(t, func() { GetSessionClaimsFromContext(ctx) })
}
