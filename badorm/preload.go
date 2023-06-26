package badorm

import "errors"

var ErrRelationNotLoaded = errors.New("Relation not loaded")

func VerifyStructLoaded[T Model](toVerify *T) (*T, error) {
	if toVerify == nil || !(*toVerify).IsLoaded() {
		return nil, ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyPointerLoaded[TModel Model, TID ModelID](id *TID, toVerify *TModel) (*TModel, error) {
	// if id == nil the relation is null
	// if (*id).IsNil(), id is loaded from a null
	// if toVerify != nil, the relation is loaded and not null
	if id != nil && !(*id).IsNil() && toVerify == nil {
		return nil, ErrRelationNotLoaded
	}

	return toVerify, nil
}
