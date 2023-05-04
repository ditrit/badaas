package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	executableDir := filepath.Dir(executablePath)

	dbProvider := viper.GetString(DBProviderKey)
	copyFolder(
		filepath.Join(executableDir, "docker", "api"),
		filepath.Join("badaas", "docker", "api"),
	)
	copyFolder(
		filepath.Join(executableDir, "docker", dbProvider),
		filepath.Join("badaas", "docker", "db"),
	)

	copyFile(
		filepath.Join(executableDir, "scripts", "run.sh"),
		"badaas",
	)
	copyFile(
		filepath.Join(executableDir, "docker", ".dockerignore"),
		".",
	)
}

func copyFile(sourcePath, destPath string) {
	err := exec.Command("cp", "-f", sourcePath, destPath).Run()
	if err != nil {
		panic(err)
	}
}

func copyFolder(sourcePath, destPath string) {
	err := exec.Command("mkdir", "-p", destPath).Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("cp", "-rf", sourcePath+"/.", destPath).Run()
	if err != nil {
		panic(err)
	}
}
