package apperror

import "net/http"

type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func BadRequest(message string) *AppError {
	return &AppError{Code: "bad_request", Message: message, Status: http.StatusBadRequest}
}

func NotFound(message string) *AppError {
	return &AppError{Code: "not_found", Message: message, Status: http.StatusNotFound}
}

func Conflict(message string) *AppError {
	return &AppError{Code: "conflict", Message: message, Status: http.StatusConflict}
}

func Internal(message string, err error) *AppError {
	return &AppError{Code: "internal_error", Message: message, Status: http.StatusInternalServerError, Err: err}
}
