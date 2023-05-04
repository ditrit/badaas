package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed docker/*
var embedFS embed.FS

var genCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "gen",
	Short: "Generate files and configurations necessary to use BadAss",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateDockerFiles,
})

const (
	DBProviderKey = "db_provider"
	Cockroachdb   = "cockroachdb"
	Postgres      = "postgres"
)

func init() {
	rootCmd.AddSubCommand(genCmd)

	genCmd.LKey(
		DBProviderKey, verdeter.IsStr, "p",
		fmt.Sprintf("Database provider (%s|%s)", Cockroachdb, Postgres),
	)
	genCmd.SetRequired(DBProviderKey)

	providerValidator := validators.AuthorizedValues(Cockroachdb, Postgres)
	genCmd.AddValidator(DBProviderKey, providerValidator)
}

func generateDockerFiles(cmd *cobra.Command, args []string) {
	sourceDockerDir := "docker"
	destBadaasDir := "badaas"
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
}

func copyFile(sourcePath, destPath string) {
	fileContent, err := embedFS.ReadFile(sourcePath)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(destPath, fileContent, 0o0666); err != nil {
		panic(err)
	}
}

func copyDir(sourceDir, destDir string) {
	files, err := embedFS.ReadDir(sourceDir)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(destDir)
	if os.IsNotExist(err) {
		os.MkdirAll(destDir, os.ModePerm)
	} else if err != nil {
		panic(err)
	}

	for _, file := range files {
		copyFile(
			filepath.Join(sourceDir, file.Name()),
			filepath.Join(destDir, file.Name()),
		)
	}
}
