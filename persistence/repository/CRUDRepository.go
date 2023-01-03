package repository

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
)

var (
	// ErrDuplicateKey is returned a conflict happen. (Most of the time because a field is unique.)
	ErrDuplicateKey = httperrors.NewHTTPError(
		http.StatusConflict,
		"conflict",
		"conflict on write",
		nil, false,
	)
)

// Generic CRUD Repository
type CRUDRepository[T models.Tabler, ID any] interface {
	Create(*T) httperrors.HTTPError
	Delete(*T) httperrors.HTTPError
	Save(*T) httperrors.HTTPError
	GetByID(ID) (*T, httperrors.HTTPError)
	GetAll(SortOption) ([]*T, httperrors.HTTPError)
	Count(squirrel.Sqlizer) (uint, httperrors.HTTPError)
	Find(squirrel.Sqlizer, pagination.Paginator, SortOption) (*pagination.Page[T], httperrors.HTTPError)
	Transaction(fn func(CRUDRepository[T, ID]) (any, error)) (any, httperrors.HTTPError)
}
