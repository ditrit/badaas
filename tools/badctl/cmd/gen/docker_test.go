package gen

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func TestGenCockroach(t *testing.T) {
	generateDockerFiles(nil, nil)
	checkFilesExist(t)
	checkDBPort(t, 26257)
	teardown()
}

func TestGenPostgres(t *testing.T) {
	viper.Set(DBProviderKey, Postgres)
	generateDockerFiles(nil, nil)
	checkFilesExist(t)
	checkFileExists(t, "badaas/docker/db/init.sql")
	checkDBPort(t, 5432)

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

func checkDBPort(t *testing.T, port int) {
	yamlFile, err := ioutil.ReadFile("badaas/config/badaas.yml")
	if err != nil {
		t.Error(err)
	}

	configData := map[string]any{}

	err = yaml.Unmarshal(yamlFile, &configData)
	if err != nil {
		t.Error(err)
	}

	databaseConfigMap, ok := configData["database"].(map[string]any)
	if !ok {
		log.Fatalln("Database configuration is not a map")
	}

	assert.Equal(t, databaseConfigMap["port"], port)
}

func teardown() {
	remove(".dockerignore")
	remove("Makefile")
	remove("badaas")
}

func remove(name string) {
	if err := os.RemoveAll(name); err != nil {
		log.Fatal(err)
	}
}
