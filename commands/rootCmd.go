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
	HumanType := &models.EntityType{
		Name: "human",
	}
	nameAttr := &models.Attribut{Name: "name", ValueType: "string", Required: true}
	HumanType.Attributs = append(
		HumanType.Attributs, nameAttr,
	)
	bobName, err := models.NewStringValue(nameAttr, "bob")
	if err != nil {
		panic(err)
	}
	bob := &models.Entity{EntityType: HumanType}
	bob.Fields = append(bob.Fields, bobName)
	db.Create(bob)

	jeanName, err := models.NewStringValue(nameAttr, "jean")
	if err != nil {
		panic(err)
	}
	jean := &models.Entity{EntityType: HumanType}
	jean.Fields = append(jean.Fields, jeanName)

	db.Create(jean)

	// Defining a bird
	colorAttr := &models.Attribut{Name: "color", ValueType: "string", Required: true}
	specieAttr := &models.Attribut{Name: "specie", ValueType: "string", Required: true}
	heightAttr := &models.Attribut{Name: "height", ValueType: "int", Default: true, DefaultInt: 12, Required: false}
	weightAttr := &models.Attribut{Name: "weight", ValueType: "float", Default: true, DefaultFloat: 12.500, Required: false}
	ownerAttr := &models.Attribut{Name: "owner", ValueType: "relation", Required: false, TargetEntityTypeId: HumanType.ID}

	BirdType := &models.EntityType{
		Name: "bird",
	}
	BirdType.Attributs = append(
		BirdType.Attributs, colorAttr, specieAttr, heightAttr, weightAttr, ownerAttr,
	)

	val1, err := models.NewStringValue(colorAttr, "blue")
	if err != nil {
		panic(err)
	}
	val2, err := models.NewStringValue(specieAttr, "mésange")
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

	mesange := &models.Entity{EntityType: BirdType}
	mesange.Fields = append(mesange.Fields, val1, val2, val3, val4, val5)

	db.Create(mesange)
	fmt.Println("Finished populating the database")

	return nil
}

func PopulateDatabase2(db *gorm.DB) error {
	// GETTING THE SUSER ADMIN FOR REFERENCE
	USERID := "wowasupercooluserID"

	// CREATION OF THE PROFILE TYPE AND ASSOCIATED ATTRIBUTES
	ProfileType := &models.EntityType{
		Name: "profile",
	}
	displayNameAttr := &models.Attribut{
		EntityTypeId: ProfileType.ID,
		Name:         "displayName",
		ValueType:    "string",
		Required:     true,
	}
	urlPicAttr := &models.Attribut{
		EntityTypeId:  ProfileType.ID,
		Name:          "urlPic",
		ValueType:     "string",
		Required:      false,
		Default:       true,
		DefaultString: "https://www.startpage.com/av/proxy-image?piurl=https%3A%2F%2Fimg.favpng.com%2F17%2F19%2F1%2Fbusiness-google-account-organization-service-png-favpng-sUuKmS4aDNRzxDKx8kJciXdFp.jpg&sp=1672915826Tc106d9b5cab08d9d380ce6fdc9564b199a49e494a069e1923c21aa202ba3ed73", //nolint:lll
	}
	userIdAttr := &models.Attribut{
		EntityTypeId: ProfileType.ID,
		Name:         "userId",
		ValueType:    "string",
		Required:     true,
	}
	ProfileType.Attributs = append(ProfileType.Attributs,
		displayNameAttr,
		urlPicAttr,
		userIdAttr,
	)

	// INSTANCIATION OF A Profile
	adminProfile := &models.Entity{
		EntityTypeId: ProfileType.ID,
		EntityType:   ProfileType,
	}
	displayNameVal := &models.Value{Attribut: urlPicAttr, StringVal: "The Super Admin"}
	userPicVal := &models.Value{Attribut: urlPicAttr, IsNull: true}
	userIdVal := &models.Value{Attribut: userIdAttr, StringVal: USERID}
	adminProfile.Fields = append(adminProfile.Fields,
		displayNameVal,
		userPicVal,
		userIdVal,
	)

	// CREATION OF THE POST TYPE AND ASSOCIATED ATTRIBUTES
	PostType := &models.EntityType{
		Name: "post",
	}
	titleAttr := &models.Attribut{
		EntityTypeId: PostType.ID,
		Name:         "title",
		ValueType:    "string",
		Required:     true,
	}
	bodyAttr := &models.Attribut{
		Name:          "body",
		ValueType:     "string",
		Default:       false,
		DefaultString: "empty",
	}
	ownerAttr := &models.Attribut{
		Name:      "ownerID",
		ValueType: "string",
		Required:  true,
	}

	PostType.Attributs = append(
		PostType.Attributs, titleAttr, bodyAttr, ownerAttr,
	)
	// INSTANCIATION OF A POST
	whycatslikemices := &models.Entity{
		EntityTypeId: PostType.ID,
		EntityType:   PostType,
	}
	titleVal := &models.Value{
		Attribut:  titleAttr,
		StringVal: "Why cats like mices ?",
	}
	bodyVal, err := models.NewStringValue(bodyAttr,
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit. In consectetur, ex at hendrerit lobortis, tellus lorem blandit eros, vel ornare odio lorem eget nisi. In erat mi, pharetra ut lacinia at, facilisis vitae nunc. Fusce rhoncus id justo vitae gravida. In nisi mi, rutrum et arcu ac, gravida venenatis arcu. Nulla leo metus, molestie eu sagittis non, ultricies eu ex. Fusce a lorem eu urna porttitor molestie. Aliquam nec sapien quam. Suspendisse aliquet elementum arcu vitae interdum. Maecenas nec turpis et nulla volutpat accumsan. Pellentesque non ullamcorper leo, eu fringilla odio.
	
	Cras eu felis mauris. Pellentesque varius imperdiet suscipit. Nam tellus odio, faucibus at mattis quis, cursus at tortor. Curabitur vitae mi eu lorem feugiat pretium sed sit amet purus. Proin efficitur, magna eu malesuada fermentum, tortor tortor maximus neque, vel mattis tortor orci a ligula. Nunc nec justo ipsum. Sed fermentum, nisl eget efficitur accumsan, augue nisl sollicitudin massa, vel suscipit enim turpis nec nisi.
	
	Nam dictum risus sed leo malesuada varius. Pellentesque gravida interdum risus id vulputate. Mauris feugiat vulputate leo ut euismod. Fusce auctor at lacus eget sollicitudin. Suspendisse potenti. Aliquam dui felis, mollis quis porta a, sodales in ligula. In ac elit ornare, facilisis ex eget, tincidunt orci. Nullam eu mattis turpis, non finibus dolor.
	`) //nolint:lll
	if err != nil {
		return err
	}
	ownerVal := &models.Value{
		Attribut:  ownerAttr,
		StringVal: USERID,
	}

	whycatslikemices.Fields = append(whycatslikemices.Fields,
		titleVal, bodyVal, ownerVal,
	)

	err = db.Create(whycatslikemices).Error
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
