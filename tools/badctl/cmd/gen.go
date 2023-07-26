package cmd

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ditrit/verdeter"
)

//go:embed docker/*
//go:embed config/*
var genEmbedFS embed.FS

var genCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "gen",
	Short: "Generate files and configurations necessary to use BadAss",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateDockerFiles,
})

const destBadaasDir = "badaas"

func init() {
	rootCmd.AddSubCommand(genCmd)
}

func generateDockerFiles(cmd *cobra.Command, args []string) {
	sourceDockerDir := "docker"

	copyDir(
		filepath.Join(sourceDockerDir, "db"),
		filepath.Join(destBadaasDir, "docker", "db"),
	)

	copyDir(
		filepath.Join(sourceDockerDir, "api"),
		filepath.Join(destBadaasDir, "docker", "api"),
	)

	copyFile(
		filepath.Join(sourceDockerDir, ".dockerignore"),
		".dockerignore",
	)

	copyFile(
		filepath.Join(sourceDockerDir, "Makefile"),
		"Makefile",
	)

	copyDir(
		"config",
		filepath.Join(destBadaasDir, "config"),
	)
}

func copyFile(sourcePath, destPath string) {
	fileContent, err := genEmbedFS.ReadFile(sourcePath)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(destPath, fileContent, 0o0600); err != nil {
		panic(err)
	}
}

func copyDir(sourceDir, destDir string) {
	_, err := os.Stat(destDir)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	files, err := genEmbedFS.ReadDir(sourceDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		copyFile(
			filepath.Join(sourceDir, file.Name()),
			filepath.Join(destDir, file.Name()),
		)
	}
}
