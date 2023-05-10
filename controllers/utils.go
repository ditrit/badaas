package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ditrit/badaas/httperrors"
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
