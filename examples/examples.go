package examples

import (
	"fmt"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	BirdsExample string = "birds"
	PostsExample string = "posts"
)

func StartExample(logger *zap.Logger, db *gorm.DB) error {
	example := viper.GetString(configuration.ServerExampleKey)
	switch example {
	case "":
		return nil
	case BirdsExample:
		return startBirdsExample(logger, db)
	case PostsExample:
		return startPostsExample(logger, db)
	default:
		return fmt.Errorf("unknown example %s", example)
	}
}

func startBirdsExample(logger *zap.Logger, db *gorm.DB) error {
	logger.Sugar().Info("Setting up Birds example")
	//defining
	humanType := &models.EntityType{
		Name: "human",
	}
	nameAttr := &models.Attribute{Name: "name", ValueType: models.StringValueType, Required: true}
	humanType.Attributes = append(
		humanType.Attributes, nameAttr,
	)
	bobName, err := models.NewStringValue(nameAttr, "bob")
	if err != nil {
		return err
	}
	bob := &models.Entity{EntityType: humanType}
	bob.Fields = append(bob.Fields, bobName)
	db.Create(bob)

	jeanName, err := models.NewStringValue(nameAttr, "jean")
	if err != nil {
		return err
	}
	jean := &models.Entity{EntityType: humanType}
	jean.Fields = append(jean.Fields, jeanName)

	db.Create(jean)

	// Defining a bird
	colorAttr := &models.Attribute{Name: "color", ValueType: models.StringValueType, Required: true}
	specieAttr := &models.Attribute{Name: "specie", ValueType: models.StringValueType, Required: true}
	heightAttr := &models.Attribute{Name: "height", ValueType: models.IntValueType, Default: true, DefaultInt: 12, Required: false}
	weightAttr := &models.Attribute{Name: "weight", ValueType: models.FloatValueType, Default: true, DefaultFloat: 12.500, Required: false}
	ownerAttr := &models.Attribute{Name: "owner", ValueType: models.RelationValueType, Required: false, RelationTargetEntityTypeID: humanType.ID}

	BirdType := &models.EntityType{
		Name: "bird",
	}
	BirdType.Attributes = append(
		BirdType.Attributes, colorAttr, specieAttr, heightAttr, weightAttr, ownerAttr,
	)

	val1, err := models.NewStringValue(colorAttr, "blue")
	if err != nil {
		return err
	}
	val2, err := models.NewStringValue(specieAttr, "chickadee")
	if err != nil {
		return err
	}
	val3, err := models.NewIntValue(heightAttr, 8)
	if err != nil {
		return err
	}
	val4, err := models.NewNullValue(weightAttr)
	if err != nil {
		return err
	}

	val5, err := models.NewRelationValue(ownerAttr, bob)
	if err != nil {
		return err
	}

	chickadee := &models.Entity{EntityType: BirdType}
	chickadee.Fields = append(chickadee.Fields, val1, val2, val3, val4, val5)

	db.Create(chickadee)
	logger.Sugar().Info("Finished populating the database")

	return nil
}

func startPostsExample(logger *zap.Logger, db *gorm.DB) error {
	logger.Sugar().Info("Setting up Posts example")
	// GETTING THE USER ADMIN FOR REFERENCE
	userID := "wowASuperCoolUserID"

	// CREATION OF THE PROFILE TYPE AND ASSOCIATED ATTRIBUTES
	profileType := &models.EntityType{
		Name: "profile",
	}
	displayNameAttr := &models.Attribute{
		EntityTypeID: profileType.ID,
		Name:         "displayName",
		ValueType:    models.StringValueType,
		Required:     true,
	}
	urlPicAttr := &models.Attribute{
		EntityTypeID:  profileType.ID,
		Name:          "urlPic",
		ValueType:     models.StringValueType,
		Required:      false,
		Default:       true,
		DefaultString: "https://www.startpage.com/av/proxy-image?piurl=https%3A%2F%2Fimg.favpng.com%2F17%2F19%2F1%2Fbusiness-google-account-organization-service-png-favpng-sUuKmS4aDNRzxDKx8kJciXdFp.jpg&sp=1672915826Tc106d9b5cab08d9d380ce6fdc9564b199a49e494a069e1923c21aa202ba3ed73", //nolint:lll
	}
	userIDAttr := &models.Attribute{
		EntityTypeID: profileType.ID,
		Name:         "userId",
		ValueType:    models.StringValueType,
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
	displayNameVal, _ := models.NewStringValue(displayNameAttr, "The Super Admin")
	userPicVal, _ := models.NewNullValue(urlPicAttr)
	userIDVal, _ := models.NewStringValue(userIDAttr, userID)
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
		ValueType:    models.StringValueType,
		Required:     true,
	}
	bodyAttr := &models.Attribute{
		Name:          "body",
		ValueType:     models.StringValueType,
		Default:       false,
		DefaultString: "empty",
	}
	ownerAttr := &models.Attribute{
		Name:      "ownerID",
		ValueType: models.StringValueType,
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
	titleVal, _ := models.NewStringValue(titleAttr, "Why cats like mice ?")
	bodyVal, _ := models.NewStringValue(bodyAttr,
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit.

		In consectetur, ex at hendrerit lobortis, tellus lorem blandit eros, vel ornare odio lorem eget nisi.

		In erat mi, pharetra ut lacinia at, facilisis vitae nunc.
	`)
	ownerVal, _ := models.NewStringValue(ownerAttr, userID)

	whyCatsLikeMice.Fields = append(whyCatsLikeMice.Fields,
		titleVal, bodyVal, ownerVal,
	)

	err := db.Create(whyCatsLikeMice).Error
	if err != nil {
		return err
	}

	err = db.Create(adminProfile).Error
	if err != nil {
		return err
	}

	logger.Sugar().Info("Finished populating the database")

	return nil
}
