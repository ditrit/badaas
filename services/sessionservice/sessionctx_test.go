package sessionservice

import (
	"context"
	"testing"

	"github.com/ditrit/badaas/badorm"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionCtx(t *testing.T) {
	ctx := context.Background()
	sessionClaims := &SessionClaims{badorm.NilUUID, badorm.UUID(uuid.New())}
	ctx = SetSessionClaimsContext(ctx, sessionClaims)
	claims := GetSessionClaimsFromContext(ctx)
	assert.Equal(t, badorm.NilUUID, claims.UserID)
}

func TestSessionCtxPanic(t *testing.T) {
	ctx := context.Background()
	assert.Panics(t, func() { GetSessionClaimsFromContext(ctx) })
}
