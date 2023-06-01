package gen

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/tools/badctl/cmd/cmderrors"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

//go:embed docker/*
//go:embed config/*
var embedFS embed.FS

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

var DBPorts = map[string]int{
	Cockroachdb: 26257,
	Postgres:    5432,
	MySQL:       3306,
}

var DBDialectors = map[string]configuration.DBDialector{
	Cockroachdb: configuration.PostgreSQL,
	Postgres:    configuration.PostgreSQL,
	MySQL:       configuration.MySQL,
}

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

func generateDockerFiles(cmd *cobra.Command, args []string) {
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

	copyBadaasConfig(dbProvider)
}

func copyBadaasConfig(dbProvider string) {
	configFile, err := embedFS.ReadFile(
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

	configData["database"].(map[string]any)["port"] = DBPorts[dbProvider]
	configData["database"].(map[string]any)["dialector"] = string(DBDialectors[dbProvider])

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
		configBytes, 0o0600,
	)
	if err != nil {
		cmderrors.FailErr(err)
	}
}

func copyFile(sourcePath, destPath string) {
	fileContent, err := embedFS.ReadFile(sourcePath)
	if err != nil {
		cmderrors.FailErr(err)
	}

	if err := os.WriteFile(destPath, fileContent, 0o0600); err != nil {
		cmderrors.FailErr(err)
	}
}

func copyDir(sourceDir, destDir string) {
	files, err := embedFS.ReadDir(sourceDir)
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
