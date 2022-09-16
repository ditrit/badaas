package commands

import (
	"errors"
	"log"

	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/registry"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"go.uber.org/zap"
)

// Create a super admin user and exit with code 1 on error
func createSuperAdminUser() {
	logg := zap.L().Sugar()
	superadmin, err := models.NewUser("superadmin", "superadmin@badaas.test", "1234")
	if err != nil {
		logg.Fatalf("failed to create superadmin %w", err)
	}
	registry.GetRegistry().UserRepo.Create(superadmin)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			logg.Debugf("The superadmin user already exists in database")
		} else {
			logg.Fatalf("failed to save the super admin %w", err)
		}
	}

}

// Run the http server for badaas
func runHTTPServer(cfg *verdeter.VerdeterCommand, args []string) error {
	err := logger.InitLogger(logger.DevelopmentLogger)
	if err != nil {
		log.Fatalf("An error happened while initializing logger (ERROR=%s)", err.Error())
	}

	zap.L().Info("The logger is initialiazed")

	registryInstance, err := registry.FactoryRegistry(registry.GormDatastore)
	if err != nil {
		zap.L().Sugar().Fatalf("An error happened while initializing datastorage layer (ERROR=%s)", err.Error())
	}
	registry.ReplaceGlobals(registryInstance)
	zap.L().Info("The datastorage layer is initialized")

	createSuperAdminUser()

	// create router
	router := router.SetupRouter()

	// create server
	srv := createServerFromConfiguration(router)

	zap.L().Sugar().Infof("Ready to serve at %s\n", srv.Addr)
	return srv.ListenAndServe()
}

var rootCfg = verdeter.NewVerdeterCommand(
	"badaas",
	"Backend and Distribution as a Service",
	`Badaas stands for Backend and Distribution as a Service.`,
	runHTTPServer,
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCfg.Execute()
}

func init() {
	rootCfg.Initialize()

	rootCfg.GKey("config_path", verdeter.IsStr, "", "Path to the config file/directory")
	rootCfg.SetDefault("config_path", ".")

	rootCfg.GKey("max_timeout", verdeter.IsInt, "", "maximum timeout (in second)")
	rootCfg.SetDefault("max_timeout", 15)

	rootCfg.GKey("host", verdeter.IsStr, "", "Address to bind (default is 0.0.0.0)")
	rootCfg.SetDefault("host", "0.0.0.0")

	rootCfg.GKey("port", verdeter.IsInt, "p", "Port to bind (default is 8000)")
	rootCfg.SetValidator("port", validators.CheckTCPHighPort)
	rootCfg.SetDefault("port", 8000)

	rootCfg.GKey("database.port", verdeter.IsInt, "", "[DB] the port of the database server")
	rootCfg.SetRequired("database.port")

	rootCfg.GKey("database.host", verdeter.IsStr, "", "[DB] the host of the database server")
	rootCfg.SetRequired("database.host")

	rootCfg.GKey("database.name", verdeter.IsStr, "", "[DB] the name of the database to use")
	rootCfg.SetRequired("database.name")

	rootCfg.GKey("database.username", verdeter.IsStr, "", "[DB] the username of the account on the database server")
	rootCfg.SetRequired("database.username")

	rootCfg.GKey("database.password", verdeter.IsStr, "", "[DB] the password of the account one the database server")
	rootCfg.SetRequired("database.password")

	rootCfg.GKey("database.sslmode", verdeter.IsStr, "", "[DB] the sslmode to use when connecting to the database server")
	rootCfg.SetRequired("database.sslmode")

}
