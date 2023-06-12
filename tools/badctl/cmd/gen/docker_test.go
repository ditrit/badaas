package gen

import (
	"io/ioutil"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/ditrit/badaas/tools/badctl/cmd/utils"
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
	utils.CheckFileExists(t, "badaas/docker/db/init.sql")
	checkDBPort(t, 5432)

	teardown()
}

func checkFilesExist(t *testing.T) {
	utils.CheckFileExists(t, ".dockerignore")
	utils.CheckFileExists(t, "Makefile")
	utils.CheckFileExists(t, "badaas/config/badaas.yml")
	utils.CheckFileExists(t, "badaas/docker/api/docker-compose.yml")
	utils.CheckFileExists(t, "badaas/docker/api/Dockerfile")
	utils.CheckFileExists(t, "badaas/docker/db/docker-compose.yml")
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

	assert.Equal(t, configData["database"].(map[string]any)["port"], port)
}

func teardown() {
	utils.RemoveFile(".dockerignore")
	utils.RemoveFile("Makefile")
	utils.RemoveFile("badaas")
}
