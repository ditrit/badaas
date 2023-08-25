package controllers

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ditrit/badaas/httperrors"
	mockRepository "github.com/ditrit/badaas/mocks/persistence/repository"
	mockSAMLService "github.com/ditrit/badaas/mocks/services/auth/protocols/samlservice"
	mockSessionService "github.com/ditrit/badaas/mocks/services/sessionservice"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewSAMLController(t *testing.T) {
	samlService := mockSAMLService.NewSAMLService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	samlController := NewSAMLController(zap.L(), samlService, userRepository, sessionService)
	assert.NotNil(t, samlController)
}
// Il faut tester les méthodes SAML on peut trouver des documments XML sur internet 
func TestSAMLController_SpToIdp(t *testing.T) {
	samlService := mockSAMLService.NewSAMLService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	samlController := NewSAMLController(zap.L(), samlService, userRepository, sessionService)


}

func TestSAMLController_IdpToSp(t *testing.T) {
	samlService := mockSAMLService.NewSAMLService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	samlController := NewSAMLController(zap.L(), samlService, userRepository, sessionService)


}
func TestSAMLController_BuildSPMetadata(t *testing.T) {
	samlService := mockSAMLService.NewSAMLService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	samlController := NewSAMLController(zap.L(), samlService, userRepository, sessionService)


}
func TestSAMLController_HandleLogoutFromIDP(t *testing.T) {
	samlService := mockSAMLService.NewSAMLService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	samlController := NewSAMLController(zap.L(), samlService, userRepository, sessionService)


}
func TestSAMLController_HandleLogoutFromSP(t *testing.T) {
	samlService := mockSAMLService.NewSAMLService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	samlController := NewSAMLController(zap.L(), samlService, userRepository, sessionService)


}
func TestSAMLController_GenerateMetadataWithSLO(t *testing.T) {
	samlService := mockSAMLService.NewSAMLService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	samlController := NewSAMLController(zap.L(), samlService, userRepository, sessionService)


}
