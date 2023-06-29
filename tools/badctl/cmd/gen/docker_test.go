package gen

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/ditrit/badaas/tools/badctl/cmd/testutils"
)

func TestGenCockroach(t *testing.T) {
	GenerateDockerFiles(nil, nil)
	checkFilesExist(t)
	checkDBPort(t, 26257)
	teardown()
}

func TestGenPostgres(t *testing.T) {
	viper.Set(DBProviderKey, Postgres)
	GenerateDockerFiles(nil, nil)
	checkFilesExist(t)
	testutils.CheckFileExists(t, "badaas/docker/db/init.sql")
	checkDBPort(t, 5432)

	teardown()
}

func checkFilesExist(t *testing.T) {
	testutils.CheckFileExists(t, ".dockerignore")
	testutils.CheckFileExists(t, "Makefile")
	testutils.CheckFileExists(t, "badaas/config/badaas.yml")
	testutils.CheckFileExists(t, "badaas/docker/api/docker-compose.yml")
	testutils.CheckFileExists(t, "badaas/docker/api/Dockerfile")
	testutils.CheckFileExists(t, "badaas/docker/db/docker-compose.yml")
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
	testutils.RemoveFile(".dockerignore")
	testutils.RemoveFile("Makefile")
	testutils.RemoveFile("badaas")
}
