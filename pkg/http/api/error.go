package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"kiln-exercice/pkg/api"
)

type Error struct {
	Status  int    `json:"-"`
	Err     error  `json:"-"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	if e.Err == nil {
		return e.Message
	}

	if e.Message == "" {
		return e.Err.Error()
	}

	return e.Err.Error() + ": " + e.Message
}

// writeError writes an error response to the client.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	var httpErr *Error

	// Check the type of the error and convert it to an HTTP error.
	// No need to use errors.As here, as err will not be wrapped.
	switch e := err.(type) {
	case *Error:
		httpErr = e
	case *api.Error:
		httpErr = &Error{
			Status:  httpStatusFromCode(e.Code),
			Err:     e.Err,
			Message: e.Message,
		}
	default:
		httpErr = InternalServerError(err)
	}

	if errors.Is(r.Context().Err(), context.Canceled) && r.Method == http.MethodGet {
		log.Ctx(r.Context()).Info().Msg(httpErr.Error())
	} else {
		if httpErr.Status >= http.StatusInternalServerError {
			log.Ctx(r.Context()).Info().Err(httpErr.Err).Msg(httpErr.Message)
		} else if httpErr.Status >= http.StatusBadRequest {
			log.Ctx(r.Context()).Info().Msg(httpErr.Error())
		}
	}
	writeErr := JSONResponse(w, httpErr.Status, httpErr)
	if writeErr != nil {
		log.Ctx(r.Context()).Error().Err(writeErr).Msg("Failed to write error response")
	}
}

// httpStatusFromCode returns the HTTP status code for the given API code.
func httpStatusFromCode(code api.Code) int {
	switch code {
	case api.OK:
		return http.StatusOK
	case api.Unknown:
		return http.StatusInternalServerError
	case api.InvalidArgument:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func BadRequestError(message string, errs ...error) *Error {
	return &Error{
		Status:  http.StatusBadRequest,
		Err:     errors.Join(errs...),
		Message: message,
	}
}

func InternalServerError(err error) *Error {
	return &Error{
		Status:  http.StatusInternalServerError,
		Err:     err,
		Message: "Internal Server Error",
	}
}
