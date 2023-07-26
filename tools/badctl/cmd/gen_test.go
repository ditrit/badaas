package cmd

import (
	"log"
	"os"
	"testing"
)

func TestGenCockroach(t *testing.T) {
	generateDockerFiles(nil, nil)
	checkFilesExist(t)
	teardown()
}

func checkFilesExist(t *testing.T) {
	checkFileExists(t, ".dockerignore")
	checkFileExists(t, "Makefile")
	checkFileExists(t, "badaas/config/badaas.yml")
	checkFileExists(t, "badaas/docker/api/docker-compose.yml")
	checkFileExists(t, "badaas/docker/api/Dockerfile")
	checkFileExists(t, "badaas/docker/db/docker-compose.yml")
}

func checkFileExists(t *testing.T, name string) {
	if _, err := os.Stat(name); err != nil {
		t.Error(err)
	}
}

func teardown() {
	remove(".dockerignore")
	remove("Makefile")
	remove("badaas")
}

func remove(name string) {
	err := os.RemoveAll(name)
	if err != nil {
		log.Fatal(err)
	}
}
