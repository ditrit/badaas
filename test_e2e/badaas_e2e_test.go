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
	"github.com/elliotchance/pie/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	"github.com/ditrit/badaas/testintegration"
)

type TestContext struct {
	statusCode int
	json       any
	httpClient *http.Client
	db         *gorm.DB
}

var (
	opts = godog.Options{Output: colors.Colored(os.Stdout)}
	db   *gorm.DB
)

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(_ *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	viper.Set(configuration.DatabasePortKey, 5000)
	viper.Set(configuration.DatabaseHostKey, "localhost")
	viper.Set(configuration.DatabaseNameKey, "badaas_db")
	viper.Set(configuration.DatabaseUsernameKey, "badaas")
	viper.Set(configuration.DatabasePasswordKey, "badaas")
	viper.Set(configuration.DatabaseSslmodeKey, "disable")
	viper.Set(configuration.DatabaseRetryKey, 10)
	viper.Set(configuration.DatabaseRetryDurationKey, 5)
	viper.Set(configuration.DatabaseDialectorKey, string(configuration.PostgreSQL))

	db, err = gormdatabase.SetupDatabaseConnection(
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

	// let db cleaned
	CleanDB(db)

	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	t := &TestContext{
		db: db,
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	t.httpClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second,
		Jar:       jar,
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// clean db before each scenario
		CleanDB(db)

		adminUser := &models.User{
			Username: "admin",
			Email:    "admin-no-reply@badaas.com",
			Password: basicauth.SaltAndHashPassword("admin"),
		}
		err = db.Create(&adminUser).Error
		if err != nil {
			log.Fatalln(err)
		}

		userType := &models.EntityType{
			Name: "user",
		}
		nameAttr := &models.Attribute{
			EntityTypeID: userType.ID,
			Name:         "name",
			ValueType:    models.StringValueType,
			Required:     false,
		}
		userType.Attributes = append(userType.Attributes,
			nameAttr,
		)

		err = db.Create(&userType).Error
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
		userAttr := models.NewRelationAttribute(profileType, "userID", false, false, userType)
		profileType.Attributes = append(profileType.Attributes,
			displayNameAttr,
			yearOfBirthAttr,
			userAttr,
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
	ctx.Step(`^I request "(.+)" with method "(.+)" with json$`, t.requestWithJSON)
	ctx.Step(`^a "(.+)" object exists with attributes$`, t.objectExists)
	ctx.Step(`^I query a "(.+)" with the object id$`, t.queryWithObjectID)
	ctx.Step(`^I query all "(.+)" objects$`, t.queryAllObjects)
	ctx.Step(`^there are "(\d+)" "(.+)" objects$`, t.thereAreObjects)
	ctx.Step(`^there is a "(.+)" object with attributes$`, t.thereIsObjectWithAttributes)
	ctx.Step(`^I query all "(.+)" objects with conditions$`, t.queryObjectsWithConditions)
	ctx.Step(`^I delete a "(.+)" with the object id$`, t.deleteWithObjectID)
	ctx.Step(`^I modify a "(.+)" with attributes$`, t.modifyWithAttributes)
	ctx.Step(`^a "(.+)" object exists with property "(.+)" related to last object and properties$`, t.objectExistsWithRelation)
	ctx.Step(`^a sale object exists for product "(\d+)", code "(\d+)" and description "(.+)"$`, t.saleExists)
	ctx.Step(`^I query all sale objects with conditions$`, t.querySalesWithConditions)
	ctx.Step(`^there is a sale object with attributes$`, t.thereIsSaleWithAttributes)
}

func CleanDB(db *gorm.DB) {
	testintegration.CleanDBTables(db, append(
		pie.Reverse(testintegration.ListOfTables),
		[]any{
			models.Session{},
			models.User{},
			models.Value{},
			models.Attribute{},
			models.Entity{},
			models.EntityType{},
		}...,
	))
}
