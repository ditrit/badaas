package cmd

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed docker/*
//go:embed scripts/*
var embedFS embed.FS

// genCmd represents the run command
var genCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "gen",
	Short: "TODO",
	Long:  `TODO`,
	Run:   generateDockerFiles,
})

const (
	DBProviderKey = "db_provider"
	Cockroachdb   = "cockroachdb"
	Postgres      = "postgres"
)

func init() {
	rootCmd.AddSubCommand(genCmd)

	genCmd.LKey(DBProviderKey, verdeter.IsStr, "p", "Database provider")
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
		filepath.Join("scripts", "run.sh"),
		filepath.Join(destBadaasDir, "run.sh"),
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
