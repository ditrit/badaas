package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpired(t *testing.T) {
	sessionInstance := &Session{
		ExpiresAt: time.Now().Add(time.Second),
	}
	assert.False(t, sessionInstance.IsExpired())
	sessionInstance.ExpiresAt = time.Now().Add(-5 * time.Second)
	assert.True(t, sessionInstance.IsExpired())
}

func TestCanBeRolled(t *testing.T) {
	sessionDuration := 500 * time.Millisecond
	sessionInstance := &Session{
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	assert.False(t, sessionInstance.CanBeRolled(sessionDuration/4))
	time.Sleep(400 * time.Millisecond)
	assert.True(t, sessionInstance.CanBeRolled(sessionDuration))
}

func TestTableName(t *testing.T) {
	assert.Equal(t, "sessions", Session{}.TableName())
}
