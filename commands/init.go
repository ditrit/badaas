package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/userservice"
	"go.uber.org/zap"
)

// Create a super user
func createSuperUser(
	config configuration.InitializationConfiguration,
	logger *zap.Logger,
	userService userservice.UserService,
) error {
	// Create a super admin user and exit with code 1 on error
	_, err := userService.NewUser("admin", "admin-no-reply@badaas.com", config.GetAdminPassword(), "")
	if err != nil {
		if err == repository.ErrDuplicateKey {
			logger.Sugar().Infof("The superadmin user already exists in database")
			return nil
		}
		logger.Sugar().Errorf("failed to save the super admin %w", err)
		return err
	}
	return nil
}

// Create a OIDC user
func createOIDCUser(
	logger *zap.Logger,
	userService userservice.UserService,
	userRepo repository.CRUDRepository[models.User, uint],
) error {
	// Create a super admin user and exit with code 1 on error
	_, err := userService.NewUser("John Doe", "johndoe@example.com", "lqsjkdqsjnd", "johndoe@example.com")
	if err != nil {
		if err == repository.ErrDuplicateKey {
			logger.Sugar().Infof("The oidc user already exists in database")
			return nil
		}
		logger.Sugar().Errorf("failed to save the oidc user %w", err)
		return err
	}

	logger.Sugar().Infof("The oidc user has been successfully created")
	return nil
}
