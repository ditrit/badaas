package httperrors

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Describe an HTTP error
type HTTPError struct {
	Status      int
	Err         string
	Message     string
	GolangError error
	toLog       bool
}

// HTTPError constructor
func NewHTTPError(status int, err string, message string, golangError error, toLog bool) *HTTPError {
	return &HTTPError{
		Status:      status,
		Err:         err,
		Message:     message,
		GolangError: golangError,
		toLog:       toLog,
	}
}

// Convert an HTTPError to a json string
func (httpError *HTTPError) ToJSON() string {
	return fmt.Sprintf(`{"error": %q, "msg":%q, "status": %q}`, httpError.Err, httpError.Message, http.StatusText(httpError.Status))
}

// Write the HTTPError to the [http.ResponseWriter] passed as argument.
func (httpError *HTTPError) Write(httpResponse http.ResponseWriter) {
	if httpError.toLog {
		logHTTPError(httpError)
	}
	http.Error(httpResponse, httpError.ToJSON(), httpError.Status)
}

func logHTTPError(httpError *HTTPError) {
	zap.L().Info(
		"http error",
		zap.String("error", httpError.Err),
		zap.String("msg", httpError.Message),
		zap.Int("status", httpError.Status),
	)
}

// A contructor for an HttpError "Not Found"
func NewErrorNotFound(ressourceName string, msg string) *HTTPError {
	return NewHTTPError(
		http.StatusNotFound,
		fmt.Sprintf("%s not found", ressourceName),
		msg,
		nil,
		false,
	)
}

// A contructor for an HttpError "Internal Server Error"
func NewInternalServerError(errorName string, msg string, err error) *HTTPError {
	return NewHTTPError(
		http.StatusNotFound,
		errorName,
		msg,
		err,
		true,
	)
}

// A contructor for an HttpError "Unauthorized Error"
func NewUnauthorizedError(errorName string, msg string) *HTTPError {
	return NewHTTPError(
		http.StatusNotFound,
		errorName,
		msg,
		nil,
		true,
	)
}
