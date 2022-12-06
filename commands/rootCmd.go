package commands

import (
	"fmt"
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
	"github.com/ditrit/verdeter"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Run the http server for badaas
func runHTTPServer(cfg *verdeter.VerdeterCommand, args []string) error {
	fx.New(
		// Modules
		configuration.ConfigurationModule,
		router.RouterModule,
		controllers.ControllerModule,
		logger.LoggerModule,
		persistence.PersistanceModule,

		fx.Provide(userservice.NewUserService),
		fx.Provide(sessionservice.NewSessionService),
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
	return nil
}

// The command badaas
var rootCfg = verdeter.NewVerdeterCommand(
	"badaas",
	"Backend and Distribution as a Service",
	`Badaas stands for Backend and Distribution as a Service.`,
	runHTTPServer,
)

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
	PostType := &models.EntityType{
		Name: "post",
	}
	titleAttr := &models.Attribut{Name: "title", ValueType: "string", Required: true}
	bodyAttr := &models.Attribut{Name: "body", ValueType: "string", Default: false, DefaultString: "empty"}
	ownerAttr := &models.Attribut{Name: "ownerID", ValueType: "string", Required: true}

	PostType.Attributs = append(
		PostType.Attributs, titleAttr, bodyAttr, ownerAttr,
	)

	whycatslikemices := &models.Entity{EntityTypeId: PostType.ID, EntityType: PostType}
	titleVal := &models.Value{Attribut: titleAttr, StringVal: "Why cats like mices ? "}
	bodyVal, err := models.NewStringValue(bodyAttr,
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit. In consectetur, ex at hendrerit lobortis, tellus lorem blandit eros, vel ornare odio lorem eget nisi. In erat mi, pharetra ut lacinia at, facilisis vitae nunc. Fusce rhoncus id justo vitae gravida. In nisi mi, rutrum et arcu ac, gravida venenatis arcu. Nulla leo metus, molestie eu sagittis non, ultricies eu ex. Fusce a lorem eu urna porttitor molestie. Aliquam nec sapien quam. Suspendisse aliquet elementum arcu vitae interdum. Maecenas nec turpis et nulla volutpat accumsan. Pellentesque non ullamcorper leo, eu fringilla odio.
	
	Cras eu felis mauris. Pellentesque varius imperdiet suscipit. Nam tellus odio, faucibus at mattis quis, cursus at tortor. Curabitur vitae mi eu lorem feugiat pretium sed sit amet purus. Proin efficitur, magna eu malesuada fermentum, tortor tortor maximus neque, vel mattis tortor orci a ligula. Nunc nec justo ipsum. Sed fermentum, nisl eget efficitur accumsan, augue nisl sollicitudin massa, vel suscipit enim turpis nec nisi.
	
	Nam dictum risus sed leo malesuada varius. Pellentesque gravida interdum risus id vulputate. Mauris feugiat vulputate leo ut euismod. Fusce auctor at lacus eget sollicitudin. Suspendisse potenti. Aliquam dui felis, mollis quis porta a, sodales in ligula. In ac elit ornare, facilisis ex eget, tincidunt orci. Nullam eu mattis turpis, non finibus dolor.
	
	Ut aliquet laoreet risus, in ultrices purus placerat ut. Aenean eget massa et augue tristique vestibulum eu nec urna. Curabitur id scelerisque felis, ac rutrum massa. In ullamcorper ex ac turpis mattis porttitor. Nunc vitae congue ex, quis porttitor massa. Fusce ut bibendum sem. Phasellus massa dui, venenatis non erat vel, auctor fermentum dolor. Cras aliquet venenatis mauris, eu consectetur massa sodales et. Phasellus lacinia massa vel arcu suscipit congue. Vestibulum mollis tellus nisi. Phasellus at dui eget dolor sagittis vulputate. Nulla vitae est commodo, aliquam urna non, pretium neque.
	
	Maecenas sodales augue ac neque efficitur pharetra. Maecenas commodo quam magna, vel pulvinar metus condimentum eget. Phasellus malesuada ante quam. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque eu justo vel nisl fringilla scelerisque eget ut ex. Duis malesuada, risus sit amet auctor euismod, felis nisl ultrices nulla, eu porttitor arcu elit fermentum velit. Praesent ut sagittis leo. Suspendisse non condimentum nunc, eget rhoncus velit. Praesent tincidunt, arcu mattis faucibus finibus, ligula lectus sodales sapien, sed porta diam nisi vitae ex. Nullam tristique justo at laoreet varius. Ut suscipit, lacus ac ultrices ornare, nisi massa varius felis, quis condimentum dolor tellus at ante. 	
	`)
	if err != nil {
		return nil
	}

	var admin models.User
	err = db.First(&admin, "username = ?", "admin").Error
	if err != nil {
		return err
	}
	ownerVal := &models.Value{Attribut: ownerAttr, StringVal: admin.ID.String()}

	whycatslikemices.Fields = append(whycatslikemices.Fields,
		titleVal, bodyVal, ownerVal)

	db.Create(whycatslikemices)
	fmt.Println("Finished populating the database")

	return nil
}
