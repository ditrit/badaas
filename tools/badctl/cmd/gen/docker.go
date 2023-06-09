package gen

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/tools/badctl/cmd/cmderrors"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

//go:embed docker/*
//go:embed config/*
var genEmbedFS embed.FS

var genDockerCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "docker",
	Short: "Generate files and configurations necessary to use BadAss over Docker",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateDockerFiles,
})

const destBadaasDir = "badaas"

const (
	DBProviderKey = "db_provider"
	Cockroachdb   = "cockroachdb"
	Postgres      = "postgres"
	MySQL         = "mysql"
)

var DBProviders = []string{Cockroachdb, Postgres, MySQL}

const (
	CockroachdbDefaultPort = 26257
	PostgresDefaultPort    = 5432
	MySQLDefaultPort       = 3306
)

var DBPorts = map[string]int{
	Cockroachdb: CockroachdbDefaultPort,
	Postgres:    PostgresDefaultPort,
	MySQL:       MySQLDefaultPort,
}

var DBDialectors = map[string]configuration.DBDialector{
	Cockroachdb: configuration.PostgreSQL,
	Postgres:    configuration.PostgreSQL,
	MySQL:       configuration.MySQL,
}

const FilePermissions = 0o0600

func init() {
	err := genDockerCmd.LKey(
		DBProviderKey, verdeter.IsStr, "p",
		fmt.Sprintf(
			"Database provider (%s), default: %s",
			strings.Join(DBProviders, "|"),
			Cockroachdb,
		),
	)
	if err != nil {
		cmderrors.FailErr(err)
	}

	genDockerCmd.SetDefault(DBProviderKey, Cockroachdb)
	genDockerCmd.AddValidator(
		DBProviderKey,
		validators.AuthorizedValues(DBProviders...),
	)
}

func generateDockerFiles(_ *cobra.Command, _ []string) {
	sourceDockerDir := "docker"
	destDockerDir := filepath.Join(destBadaasDir, "docker")

	copyDir(
		filepath.Join(sourceDockerDir, "api"),
		filepath.Join(destDockerDir, "api"),
	)

	dbProvider := viper.GetString(DBProviderKey)
	copyDir(
		filepath.Join(sourceDockerDir, dbProvider),
		filepath.Join(destDockerDir, "db"),
	)

	copyFile(
		filepath.Join(sourceDockerDir, ".dockerignore"),
		".dockerignore",
	)

	copyFile(
		filepath.Join(sourceDockerDir, "Makefile"),
		"Makefile",
	)

	copyBadaasConfig(dbProvider)
}

func copyBadaasConfig(dbProvider string) {
	configFile, err := genEmbedFS.ReadFile(
		filepath.Join("config", "badaas.yml"),
	)
	if err != nil {
		cmderrors.FailErr(err)
	}

	configData := map[string]any{}

	err = yaml.Unmarshal(configFile, &configData)
	if err != nil {
		cmderrors.FailErr(err)
	}

	databaseConfigMap, ok := configData["database"].(map[string]any)
	if !ok {
		cmderrors.FailErr(errors.New("database configuration is not a map"))
	}

	databaseConfigMap["port"] = DBPorts[dbProvider]
	databaseConfigMap["dialector"] = string(DBDialectors[dbProvider])
	configData["database"] = databaseConfigMap

	configBytes, err := yaml.Marshal(&configData)
	if err != nil {
		cmderrors.FailErr(err)
	}

	destConfigDir := filepath.Join(destBadaasDir, "config")

	err = os.MkdirAll(destConfigDir, os.ModePerm)
	if err != nil {
		cmderrors.FailErr(err)
	}

	err = os.WriteFile(
		filepath.Join(destConfigDir, "badaas.yml"),
		configBytes, FilePermissions,
	)
	if err != nil {
		cmderrors.FailErr(err)
	}
}

func copyFile(sourcePath, destPath string) {
	fileContent, err := genEmbedFS.ReadFile(sourcePath)
	if err != nil {
		cmderrors.FailErr(err)
	}

	if err := os.WriteFile(destPath, fileContent, FilePermissions); err != nil {
		cmderrors.FailErr(err)
	}
}

func copyDir(sourceDir, destDir string) {
	files, err := genEmbedFS.ReadDir(sourceDir)
	if err != nil {
		cmderrors.FailErr(err)
	}

	_, err = os.Stat(destDir)
	if err != nil {
		if !os.IsNotExist(err) {
			cmderrors.FailErr(err)
		}

		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			cmderrors.FailErr(err)
		}
	}

	for _, file := range files {
		copyFile(
			filepath.Join(sourceDir, file.Name()),
			filepath.Join(destDir, file.Name()),
		)
	}
}
