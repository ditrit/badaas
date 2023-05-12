package badorm

import (
	"github.com/ditrit/badaas/persistence/models"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CRUDService[T any, ID BadaasID] interface {
	GetEntity(id ID) (*T, error)
	GetEntities(conditions map[string]any) ([]*T, error)
	CreateEntity(attributeValues map[string]any) (*T, error)
	UpdateEntity(entityID ID, newValues map[string]any) (*T, error)
	DeleteEntity(entityID ID) error
}

// check interface compliance
var _ CRUDService[models.User, uuid.UUID] = (*crudServiceImpl[models.User, uuid.UUID])(nil)

// Implementation of the Generic CRUD Repository
type crudServiceImpl[T any, ID BadaasID] struct {
	CRUDService[T, ID]
	logger     *zap.Logger
	db         *gorm.DB
	repository CRUDRepository[T, ID]
}

func NewCRUDService[T any, ID BadaasID](
	logger *zap.Logger,
	db *gorm.DB,
	repository CRUDRepository[T, ID],
) CRUDService[T, ID] {
	return &crudServiceImpl[T, ID]{
		logger:     logger,
		db:         db,
		repository: repository,
	}
}

// Get the Entity of type with name "entityTypeName" that has the "id"
func (service *crudServiceImpl[T, ID]) GetEntity(id ID) (*T, error) {
	return service.repository.GetByID(service.db, id)
}

// Get entities of type with name "entityTypeName" that match all "conditions"
//
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
func (service *crudServiceImpl[T, ID]) GetEntities(conditions map[string]any) ([]*T, error) {
	return service.repository.GetMultiple(service.db, conditions)
}

// Creates a Entity of type "entityType" and its Values from the information provided in "attributeValues"
// not specified values are set with the default value (if exists) or a null value.
// entries in "attributeValues" that do not correspond to any attribute of "entityType" are ignored
//
// "attributeValues" are in {"attributeName": value} format
func (service *crudServiceImpl[T, ID]) CreateEntity(attributeValues map[string]any) (*T, error) {
	var entity T
	// TODO ver si dejo esto o desencodeo el json directo en la entidad
	// TODO testear lo de que se le pueden agregar relaciones a esto
	err := mapstructure.Decode(attributeValues, &entity)
	if err != nil {
		return nil, err // TODO ver que errores puede haber aca
	}

	err = service.repository.Create(service.db, &entity)
	if err != nil {
		return nil, err
	}

	// TODO eliminar esto
	// err := service.db.Model(&entity).Create(attributeValues).Error
	// if err != nil {
	// 	return nil, err
	// }
	// entity.ID = attributeValues["id"]

	return &entity, nil
}

// Updates entity with type "entityTypeName" and "id" Values to the new values provided in "newValues"
// entries in "newValues" that do not correspond to any attribute of the EntityType are ignored
//
// "newValues" are in {"attributeName": newValue} format
func (service *crudServiceImpl[T, ID]) UpdateEntity(entityID ID, newValues map[string]any) (*T, error) {
	// TODO
	return nil, nil
}

// Deletes Entity with type "entityTypeName" and id "entityID" and its values
func (service *crudServiceImpl[T, ID]) DeleteEntity(entityID ID) error {
	// TODO
	return nil
}
