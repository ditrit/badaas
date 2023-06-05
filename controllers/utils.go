package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/httperrors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var (
	// Sent when the request is malformed
	HTTPErrRequestMalformed httperrors.HTTPError = httperrors.NewHTTPError(
		http.StatusBadRequest,
		"Request malformed",
		"The schema of the received data is not correct",
		nil,
		false,
	)
	ErrEntityNotFound     = httperrors.NewErrorNotFound("entity", "please use a valid object id")
	ErrEntityTypeNotFound = httperrors.NewErrorNotFound("entity type", "please use a type that exists")
	ErrIDNotAnUUID        = httperrors.NewBadRequestError("id is not an uuid", "please use an uuid for the id value")
)

// Decode json present in request body
func decodeJSON(r *http.Request, to any) httperrors.HTTPError {
	err := json.NewDecoder(r.Body).Decode(to)
	if err != nil {
		return HTTPErrRequestMalformed
	}

	return nil
}

// Decode json present in request body
func decodeJSONOptional(r *http.Request) (map[string]any, httperrors.HTTPError) {
	to := map[string]any{}
	err := json.NewDecoder(r.Body).Decode(&to)
	switch {
	case err == io.EOF:
		// empty body
		return to, nil
	case err != nil:
		return nil, HTTPErrRequestMalformed
	}

	return to, nil
}

// Extract the "id" parameter from url
func getEntityIDFromRequest(r *http.Request) (badorm.UUID, httperrors.HTTPError) {
	id, present := mux.Vars(r)["id"]
	if !present {
		return badorm.NilUUID, ErrEntityNotFound
	}

	uuid, err := badorm.ParseUUID(id)
	if err != nil {
		return uuid, ErrIDNotAnUUID
	}

	return uuid, nil
}

func mapServiceError(err error) httperrors.HTTPError {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEntityNotFound
		}
		return httperrors.NewDBError(err)
	}

	return nil
}
