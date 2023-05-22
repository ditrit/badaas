package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/ditrit/verdeter"
	"github.com/spf13/cobra"
)

var runCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "run",
	Short: "Run your BadAss application",
	Long:  `run is the command that will allow you to run your application once you have generated the necessary files with gen.`,
	Run:   runApp,
})

func runApp(cmd *cobra.Command, args []string) {
	dockerCommand := exec.Command(
		"docker", "compose",
		"-f", "badaas/docker/db/docker-compose.yml",
		"-f", "badaas/docker/api/docker-compose.yml",
		"up", "--build", "-d",
	)
	dockerCommand.Stdout = os.Stdout
	dockerCommand.Stderr = os.Stderr

	err := dockerCommand.Run()
	if err != nil {
		panic(err)
	}

	log.Println("Your api is available at http://localhost:8000")
}
