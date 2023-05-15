package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ditrit/badaas/httperrors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
func getEntityIDFromRequest(r *http.Request) (uuid.UUID, httperrors.HTTPError) {
	id, present := mux.Vars(r)["id"]
	if !present {
		return uuid.Nil, ErrEntityNotFound
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, ErrIDNotAnUUID
	}

	return uid, nil
}
