package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

//go:embed docker/*
//go:embed config/*
var embedFS embed.FS

var genCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "gen",
	Short: "Generate files and configurations necessary to use BadAss",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateDockerFiles,
})

const destBadaasDir = "badaas"

const (
	DBProviderKey = "db_provider"
	Cockroachdb   = "cockroachdb"
	Postgres      = "postgres"
)

var DBProviders = []string{Cockroachdb, Postgres}

var DBPorts = map[string]int{
	Cockroachdb: 26257,
	Postgres:    5432,
}

func init() {
	rootCmd.AddSubCommand(genCmd)

	err := genCmd.LKey(
		DBProviderKey, verdeter.IsStr, "p",
		fmt.Sprintf(
			"Database provider (%s), default: %s",
			strings.Join(DBProviders, "|"),
			Cockroachdb,
		),
	)
	if err != nil {
		panic(err)
	}
	genCmd.SetDefault(DBProviderKey, Cockroachdb)
	genCmd.AddValidator(
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
		panic(err)
	}

	configData := map[string]any{}
	err = yaml.Unmarshal(configFile, &configData)
	if err != nil {
		panic(err)
	}

	configData["database"].(map[string]any)["port"] = DBPorts[dbProvider]

	configBytes, err := yaml.Marshal(&configData)
	if err != nil {
		panic(err)
	}

	destConfigDir := filepath.Join(destBadaasDir, "config")
	err = os.MkdirAll(destConfigDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(
		filepath.Join(destConfigDir, "badaas.yml"),
		configBytes, 0o0600,
	)
	if err != nil {
		panic(err)
	}
}

func copyFile(sourcePath, destPath string) {
	fileContent, err := embedFS.ReadFile(sourcePath)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(destPath, fileContent, 0o0600); err != nil {
		panic(err)
	}
}

func copyDir(sourceDir, destDir string) {
	files, err := embedFS.ReadDir(sourceDir)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(destDir)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	for _, file := range files {
		copyFile(
			filepath.Join(sourceDir, file.Name()),
			filepath.Join(destDir, file.Name()),
		)
	}
}
