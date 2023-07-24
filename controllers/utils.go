package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ditrit/badaas/httperrors"
)

// HTTPErrRequestMalformed is sent when the request is malformed
var HTTPErrRequestMalformed httperrors.HTTPError = httperrors.NewHTTPError(
	http.StatusBadRequest,
	"Request malformed",
	"The schema of the received data is not correct",
	nil,
	false,
)

// Decode json present in request body
func decodeJSON(r *http.Request, to any) httperrors.HTTPError {
	err := json.NewDecoder(r.Body).Decode(to)
	if err != nil {
		return HTTPErrRequestMalformed
	}

	return nil
}
