package main

import (
	"context"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	integrationtests "github.com/ditrit/badaas/test_integration"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TestContext struct {
	statusCode int
	json       any
	httpClient *http.Client
}

var opts = godog.Options{Output: colors.Colored(os.Stdout)}
var db *gorm.DB

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(_ *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	logger, _ := zap.NewDevelopment()
	var err error

	viper.Set(configuration.DatabasePortKey, 26257)
	viper.Set(configuration.DatabaseHostKey, "localhost")
	viper.Set(configuration.DatabaseNameKey, "badaas_db")
	viper.Set(configuration.DatabaseUsernameKey, "root")
	viper.Set(configuration.DatabasePasswordKey, "postgres")
	viper.Set(configuration.DatabaseSslmodeKey, "disable")
	viper.Set(configuration.DatabaseRetryKey, 10)
	viper.Set(configuration.DatabaseRetryDurationKey, 5)
	db, err = gormdatabase.CreateDatabaseConnectionFromConfiguration(
		logger,
		configuration.NewDatabaseConfiguration(),
	)
	if err != nil {
		log.Fatalln("Unable to connect to database : ", err)
	}

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	t := &TestContext{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	t.httpClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Duration(5 * time.Second),
		Jar:       jar,
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// clean db before each scenario
		integrationtests.SetupDB(db)

		adminUser := &models.User{
			Username: "admin",
			Email:    "admin-no-reply@badaas.com",
			Password: basicauth.SaltAndHashPassword("admin"),
		}
		err = db.Create(&adminUser).Error
		if err != nil {
			log.Fatalln(err)
		}

		profileType := &models.EntityType{
			Name: "profile",
		}
		displayNameAttr := &models.Attribute{
			EntityTypeID: profileType.ID,
			Name:         "displayName",
			ValueType:    models.StringValueType,
			Required:     false,
		}
		yearOfBirthAttr := &models.Attribute{
			EntityTypeID: profileType.ID,
			Name:         "yearOfBirth",
			ValueType:    models.IntValueType,
			Required:     false,
		}
		profileType.Attributes = append(profileType.Attributes,
			displayNameAttr,
			yearOfBirthAttr,
		)

		err = db.Create(&profileType).Error
		if err != nil {
			log.Fatalln(err)
		}

		return ctx, nil
	})

	ctx.Step(`^I request "(.+)"$`, t.requestGet)
	ctx.Step(`^status code is "(\d+)"$`, t.assertStatusCode)
	ctx.Step(`^response field "(.+)" is "(.+)"$`, t.assertResponseFieldIsEquals)
	ctx.Step(`^I request "(.+)" with method "(.+)" with json$`, t.requestWithJson)
	ctx.Step(`^a "(.+)" object exists with properties$`, t.objectExists)
	ctx.Step(`^I query a "(.+)" with the object id$`, t.queryWithObjectID)
	ctx.Step(`^I query all "(.+)" objects$`, t.queryAllObjects)
	ctx.Step(`^there are "(\d+)" "(.+)" objects$`, t.thereAreObjects)
	ctx.Step(`^there is a "(.+)" object with properties$`, t.thereIsObjectWithProperties)
	ctx.Step(`^I query all "(.+)" objects with parameters$`, t.queryObjectsWithParameters)
}
