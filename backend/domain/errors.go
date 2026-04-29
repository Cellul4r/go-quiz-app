package domain

import "errors"

type AppError struct {
	Code    string
	Message string
	Fields  map[string]string
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok || e == nil || t == nil {
		return false
	}

	return e.Code == t.Code
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		return nil, false
	}

	return appErr, true
}

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = &AppError{Code: "internal_server_error", Message: "internal server error"}
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = &AppError{Code: "not_found", Message: "requested resource not found"}
	// ErrConflict will throw if the current action already exists
	ErrConflict = &AppError{Code: "conflict", Message: "resource already exists"}
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput    = &AppError{Code: "bad_request", Message: "invalid request input"}
	ErrUnauthorized     = &AppError{Code: "unauthorized", Message: "you are not authorized to perform this action"}
	ErrForbidden        = &AppError{Code: "forbidden", Message: "you don't have permission to access this resource"}
	ErrUserNameConflict = &AppError{Code: "username_conflict", Message: "username already exists"}
)
