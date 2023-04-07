package commands

import (
	"fmt"
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/resources"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services/eavservice"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
	"github.com/ditrit/verdeter"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Run the http server for badaas
func runHTTPServer(cmd *cobra.Command, args []string) {
	fx.New(
		// Modules
		configuration.ConfigurationModule,
		router.RouterModule,
		controllers.ControllerModule,
		logger.LoggerModule,
		persistence.PersistanceModule,

		fx.Provide(userservice.NewUserService),
		fx.Provide(sessionservice.NewSessionService),
		fx.Provide(eavservice.NewEAVService),
		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		fx.Provide(NewHTTPServer),

		// Finally: we invoke the newly created server
		fx.Invoke(func(*http.Server) { /* we need this function to be empty*/ }),
		fx.Invoke(createSuperUser),
		fx.Invoke(PopulateDatabase2),
	).Run()
}

// The command badaas
var rootCfg = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:     "badaas",
	Short:   "Backend and Distribution as a Service",
	Long:    "Badaas stands for Backend and Distribution as a Service.",
	Version: resources.Version,
	Run:     runHTTPServer,
})

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCfg.Execute()
}

func init() {
	rootCfg.GKey("config_path", verdeter.IsStr, "", "Path to the config file/directory")
	rootCfg.SetDefault("config_path", ".")

	initServerCommands(rootCfg)
	initLoggerCommands(rootCfg)
	initDatabaseCommands(rootCfg)
	initInitialisationCommands(rootCfg)
	initSessionCommands(rootCfg)
}

func PopulateDatabase(db *gorm.DB) error {
	//defining
	humanType := &models.EntityType{
		Name: "human",
	}
	nameAttr := &models.Attribute{Name: "name", ValueType: "string", Required: true}
	humanType.Attributes = append(
		humanType.Attributes, nameAttr,
	)
	bobName, err := models.NewStringValue(nameAttr, "bob")
	if err != nil {
		panic(err)
	}
	bob := &models.Entity{EntityType: humanType}
	bob.Fields = append(bob.Fields, bobName)
	db.Create(bob)

	jeanName, err := models.NewStringValue(nameAttr, "jean")
	if err != nil {
		panic(err)
	}
	jean := &models.Entity{EntityType: humanType}
	jean.Fields = append(jean.Fields, jeanName)

	db.Create(jean)

	// Defining a bird
	colorAttr := &models.Attribute{Name: "color", ValueType: "string", Required: true}
	specieAttr := &models.Attribute{Name: "specie", ValueType: "string", Required: true}
	heightAttr := &models.Attribute{Name: "height", ValueType: "int", Default: true, DefaultInt: 12, Required: false}
	weightAttr := &models.Attribute{Name: "weight", ValueType: "float", Default: true, DefaultFloat: 12.500, Required: false}
	ownerAttr := &models.Attribute{Name: "owner", ValueType: "relation", Required: false, TargetEntityTypeID: humanType.ID}

	BirdType := &models.EntityType{
		Name: "bird",
	}
	BirdType.Attributes = append(
		BirdType.Attributes, colorAttr, specieAttr, heightAttr, weightAttr, ownerAttr,
	)

	val1, err := models.NewStringValue(colorAttr, "blue")
	if err != nil {
		panic(err)
	}
	val2, err := models.NewStringValue(specieAttr, "chickadee")
	if err != nil {
		panic(err)
	}
	val3, err := models.NewIntValue(heightAttr, 8)
	if err != nil {
		panic(err)
	}
	val4, err := models.NewNullValue(weightAttr)
	if err != nil {
		panic(err)
	}

	val5, err := models.NewRelationValue(ownerAttr, bob)
	if err != nil {
		panic(err)
	}

	chickadee := &models.Entity{EntityType: BirdType}
	chickadee.Fields = append(chickadee.Fields, val1, val2, val3, val4, val5)

	db.Create(chickadee)
	fmt.Println("Finished populating the database")

	return nil
}

func PopulateDatabase2(db *gorm.DB) error {
	// GETTING THE USER ADMIN FOR REFERENCE
	userID := "wowASuperCoolUserID"

	// CREATION OF THE PROFILE TYPE AND ASSOCIATED ATTRIBUTES
	profileType := &models.EntityType{
		Name: "profile",
	}
	displayNameAttr := &models.Attribute{
		EntityTypeID: profileType.ID,
		Name:         "displayName",
		ValueType:    "string",
		Required:     true,
	}
	urlPicAttr := &models.Attribute{
		EntityTypeID:  profileType.ID,
		Name:          "urlPic",
		ValueType:     "string",
		Required:      false,
		Default:       true,
		DefaultString: "https://www.startpage.com/av/proxy-image?piurl=https%3A%2F%2Fimg.favpng.com%2F17%2F19%2F1%2Fbusiness-google-account-organization-service-png-favpng-sUuKmS4aDNRzxDKx8kJciXdFp.jpg&sp=1672915826Tc106d9b5cab08d9d380ce6fdc9564b199a49e494a069e1923c21aa202ba3ed73", //nolint:lll
	}
	userIDAttr := &models.Attribute{
		EntityTypeID: profileType.ID,
		Name:         "userId",
		ValueType:    "string",
		Required:     true,
	}
	profileType.Attributes = append(profileType.Attributes,
		displayNameAttr,
		urlPicAttr,
		userIDAttr,
	)

	// INSTANTIATION OF A Profile
	adminProfile := &models.Entity{
		EntityTypeID: profileType.ID,
		EntityType:   profileType,
	}
	displayNameVal := &models.Value{Attribute: urlPicAttr, StringVal: "The Super Admin"}
	userPicVal := &models.Value{Attribute: urlPicAttr, IsNull: true}
	userIDVal := &models.Value{Attribute: userIDAttr, StringVal: userID}
	adminProfile.Fields = append(adminProfile.Fields,
		displayNameVal,
		userPicVal,
		userIDVal,
	)

	// CREATION OF THE POST TYPE AND ASSOCIATED ATTRIBUTES
	postType := &models.EntityType{
		Name: "post",
	}
	titleAttr := &models.Attribute{
		EntityTypeID: postType.ID,
		Name:         "title",
		ValueType:    "string",
		Required:     true,
	}
	bodyAttr := &models.Attribute{
		Name:          "body",
		ValueType:     "string",
		Default:       false,
		DefaultString: "empty",
	}
	ownerAttr := &models.Attribute{
		Name:      "ownerID",
		ValueType: "string",
		Required:  true,
	}

	postType.Attributes = append(
		postType.Attributes, titleAttr, bodyAttr, ownerAttr,
	)
	// INSTANTIATION OF A POST
	whyCatsLikeMice := &models.Entity{
		EntityTypeID: postType.ID,
		EntityType:   postType,
	}
	titleVal := &models.Value{
		Attribute: titleAttr,
		StringVal: "Why cats like mice ?",
	}
	bodyVal, err := models.NewStringValue(bodyAttr,
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit.

		In consectetur, ex at hendrerit lobortis, tellus lorem blandit eros, vel ornare odio lorem eget nisi.
		
		In erat mi, pharetra ut lacinia at, facilisis vitae nunc.
	`)
	if err != nil {
		return err
	}
	ownerVal := &models.Value{
		Attribute: ownerAttr,
		StringVal: userID,
	}

	whyCatsLikeMice.Fields = append(whyCatsLikeMice.Fields,
		titleVal, bodyVal, ownerVal,
	)

	err = db.Create(whyCatsLikeMice).Error
	if err != nil {
		return err
	}

	err = db.Create(adminProfile).Error
	if err != nil {
		return err
	}

	fmt.Println("Finished populating the database")

	return nil
}
