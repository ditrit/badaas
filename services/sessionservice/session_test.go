package sessionservice

import (
	"errors"
	"testing"
	"time"

	"github.com/ditrit/badaas/httperrors"
	configurationMocks "github.com/ditrit/badaas/mocks/configuration"
	repositoryMocks "github.com/ditrit/badaas/mocks/persistence/repository"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

// make values for test
func setupTest(
	t *testing.T,
) (
	*repositoryMocks.CRUDRepository[models.Session, uuid.UUID],
	*sessionServiceImpl,
	*observer.ObservedLogs,
	*configurationMocks.SessionConfiguration,
) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	sessionRepositoryMock := repositoryMocks.NewCRUDRepository[models.Session, uuid.UUID](t)
	sessionConfiguration := configurationMocks.NewSessionConfiguration(t)
	service := &sessionServiceImpl{
		sessionRepository:    sessionRepositoryMock,
		logger:               logger,
		cache:                make(map[uuid.UUID]*models.Session),
		sessionConfiguration: sessionConfiguration,
		db:                   gormDB,
	}

	return sessionRepositoryMock, service, logs, sessionConfiguration
}

func TestNewSession(t *testing.T) {
	sessionInstance := models.NewSession(uuid.Nil, time.Second)
	assert.NotNil(t, sessionInstance)
	assert.Equal(t, uuid.Nil, sessionInstance.UserID)
}

func TestLogInUser(t *testing.T) {
	sessionRepositoryMock, service, logs, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.On("Create", gormDB, mock.Anything).Return(nil)

	sessionConfigurationMock.On("GetSessionDuration").Return(time.Minute)
	user := &models.User{
		Username: "bob",
		Email:    "bob@email.com",
	}
	_, err := service.LogUserIn(user)
	require.NoError(t, err)
	assert.Len(t, service.cache, 1)
	assert.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, "Added session", log.Message)
	require.Len(t, log.Context, 1)
}

func TestLogInUserDbError(t *testing.T) {
	sessionRepositoryMock, service, logs, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.
		On("Create", gormDB, mock.Anything).
		Return(errors.New("db err"))

	sessionConfigurationMock.On("GetSessionDuration").Return(time.Minute)

	user := &models.User{
		Username: "bob",
		Email:    "bob@email.com",
	}
	_, err := service.LogUserIn(user)
	require.Error(t, err)
	assert.Len(t, service.cache, 0)
	assert.Equal(t, 0, logs.Len())
}

func TestIsValid(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("Create", gormDB, mock.Anything).Return(nil)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuidSample,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err := service.add(session)
	require.NoError(t, err)
	assert.Len(t, service.cache, 1)
	assert.Equal(t, uuid.Nil, service.cache[uuidSample].UserID)
	isValid, claims := service.IsValid(uuidSample)
	require.True(t, isValid)
	assert.Equal(t, *claims, SessionClaims{
		UserID:      uuid.Nil,
		SessionUUID: uuidSample,
	})
}

func TestIsValid_SessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.
		On("GetOptional", gormDB, mock.Anything).
		Return(nil, nil)
	uuidSample := uuid.New()
	isValid, _ := service.IsValid(uuidSample)
	require.False(t, isValid)
}

func TestLogOutUser(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("Delete", gormDB, mock.Anything).Return(nil)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuidSample,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	err := service.LogUserOut(makeSessionClaims(session))
	require.NoError(t, err)
	assert.Len(t, service.cache, 0)
}

func TestLogOutUserDbError(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.
		On("Delete", gormDB, mock.Anything).
		Return(errors.New("db errors"))

	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuidSample,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	err := service.LogUserOut(makeSessionClaims(session))
	require.Error(t, err)
	assert.Len(t, service.cache, 1)
}

func TestLogOutUser_SessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.
		On("GetOptional", gormDB, mock.Anything).
		Return(nil, errors.New("db errors"))

	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	sessionClaims := makeSessionClaims(session)
	sessionClaims.SessionUUID = uuid.Nil
	err := service.LogUserOut(sessionClaims)
	require.Error(t, err)
	assert.Len(t, service.cache, 1)
}

func TestRollSession(t *testing.T) {
	sessionRepositoryMock, service, _, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.On("Save", gormDB, mock.Anything).Return(nil)
	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration / 4)
	uuidSample := uuid.New()
	originalExpirationTime := time.Now().Add(sessionDuration / 5)
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: originalExpirationTime,
	}
	service.cache[uuidSample] = session
	err := service.RollSession(uuidSample)
	require.NoError(t, err)
	assert.Greater(t, session.ExpiresAt, originalExpirationTime)
}

func TestRollSession_Expired(t *testing.T) {
	_, service, _, sessionConfigurationMock := setupTest(t)
	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration / 4)
	uuidSample := uuid.New()
	originalExpirationTime := time.Now().Add(-time.Hour)
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuidSample,
		},
		UserID:    uuid.Nil,
		ExpiresAt: originalExpirationTime,
	}
	service.cache[uuidSample] = session
	err := service.RollSession(uuidSample)
	require.Error(t, err)
}

func TestRollSession_falseUUID(t *testing.T) {
	repoSession, service, _, sessionConfigurationMock := setupTest(t)
	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration / 4)

	uuidSample := uuid.New()
	originalExpirationTime := time.Now().Add(-time.Hour)
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: originalExpirationTime,
	}
	service.cache[uuidSample] = session

	repoSession.
		On("GetOptional", gormDB, mock.Anything).
		Return(nil, nil)

	err := service.RollSession(uuid.New())
	require.NoError(t, err)
}

func TestRollSession_sessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.
		On("GetOptional", gormDB, map[string]any{"uuid": "00000000-0000-0000-0000-000000000000"}).
		Return(nil, nil)

	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration)

	err := service.RollSession(uuid.Nil)
	require.NoError(t, err)
}

func Test_pullFromDB(t *testing.T) {
	sessionRepositoryMock, service, logs, _ := setupTest(t)
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	sessionRepositoryMock.On("GetAll", gormDB).Return([]*models.Session{session}, nil)

	service.pullFromDB()
	assert.Len(t, service.cache, 1)
	assert.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, "Pulled sessions from DB", log.Message)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "sessionCount", Type: zapcore.Int64Type, Integer: 1},
	}, log.Context)
}

func Test_pullFromDB_repoError(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("GetAll", gormDB).Return(nil, httperrors.AnError)
	assert.PanicsWithError(t, httperrors.AnError.Error(), func() { service.pullFromDB() })
}

func Test_removeExpired(t *testing.T) {
	sessionRepositoryMock, service, logs, _ := setupTest(t)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("Delete", gormDB, session).
		Return(nil)
	service.cache[uuidSample] = session

	service.removeExpired()
	assert.Len(t, service.cache, 0)
	assert.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, "Removed expired session", log.Message)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "expiredSessionCount", Type: zapcore.Int64Type, Integer: 1},
	}, log.Context)
}

func Test_removeExpired_RepositoryError(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("Delete", gormDB, session).
		Return(httperrors.AnError)
	service.cache[uuidSample] = session

	assert.Panics(t, func() { service.removeExpired() })
}

func Test_get(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("GetOptional", gormDB, mock.Anything).
		Return(session, nil)

	sessionFound := service.get(uuidSample)
	assert.Equal(t, sessionFound, session)
}
