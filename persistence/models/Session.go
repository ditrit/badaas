package models

import (
	"time"

	"github.com/google/uuid"
)

// Represent a user session
type Session struct {
	BaseModel
	UserID    uuid.UUID `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

// Create a new session
func NewSession(userID uuid.UUID, sessionDuration time.Duration) *Session {
	return &Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
}

// Return true is expired
func (session *Session) IsExpired() bool {
	return time.Now().After(session.ExpiresAt)
}

// Return true if the session is expired in less than an hour
func (session *Session) CanBeRolled(rollInterval time.Duration) bool {
	return time.Now().After(session.ExpiresAt.Add(-rollInterval))
}

// Return the pluralized table name
//
// Satisfy the Tabler interface
func (Session) TableName() string {
	return "sessions"
}
