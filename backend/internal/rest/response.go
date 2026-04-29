package rest

import (
	"net/http"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if appErr, ok := domain.AsAppError(err); ok {
		switch appErr.Code {
		case "internal_server_error":
			return http.StatusInternalServerError
		case "not_found":
			return http.StatusNotFound
		case "conflict", "username_conflict":
			return http.StatusConflict
		case "unauthorized":
			return http.StatusUnauthorized
		case "forbidden":
			return http.StatusForbidden
		case "bad_request", "validation_failed":
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrForbidden:
		return http.StatusForbidden
	case domain.ErrBadParamInput, domain.ErrUserNameConflict:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func toResponseError(err error) ResponseError {
	if appErr, ok := domain.AsAppError(err); ok {
		return ResponseError{
			Code:    appErr.Code,
			Message: appErr.Message,
			Fields:  appErr.Fields,
		}
	}

	return ResponseError{
		Code:    domain.ErrInternalServerError.Code,
		Message: domain.ErrInternalServerError.Message,
	}
}
