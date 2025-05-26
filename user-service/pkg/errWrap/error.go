package errWrap

import (
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

type ErrorType string

const (
	ErrUniqueViolation      ErrorType = "unique_violation"
	ErrForeignKey           ErrorType = "foreign_key"
	ErrNotNullViolation     ErrorType = "not_null_violation"
	ErrCheckViolation       ErrorType = "check_violation"
	ErrUnauthorized         ErrorType = "unauthorized"
	ErrPaymentRequired      ErrorType = "payment_required"
	ErrValidation           ErrorType = "validation"
	ErrTooManyRequests      ErrorType = "too_many_requests"
	ErrOrderAlreadyExists   ErrorType = "order_already_exists"
	ErrOrderAlreadyUploaded ErrorType = "order_already_uploaded"
	ErrNoResponseData       ErrorType = "no_response_data"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func NewAppError(t ErrorType, msg string, err error) *AppError {
	return &AppError{
		Type:    t,
		Message: msg,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func WrapError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return NewAppError(ErrUniqueViolation, "unique violation", err)
		case "23503":
			return NewAppError(ErrForeignKey, "foreign key violation", err)
		case "23502":
			return NewAppError(ErrNotNullViolation, "not null violation", err)
		case "23514":
			return NewAppError(ErrCheckViolation, "check violation", err)
		}
	}
	return err
}

type ErrorResponse struct {
	Type    ErrorType
	Message string
}

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *AppError
	if errors.As(err, &appErr) {
		response := ErrorResponse{
			Type:    appErr.Type,
			Message: appErr.Message,
		}

		switch appErr.Type {
		case ErrUniqueViolation:
			w.WriteHeader(http.StatusConflict)
		case ErrForeignKey, ErrNotNullViolation, ErrCheckViolation:
			w.WriteHeader(http.StatusBadRequest)
		case ErrUnauthorized:
			w.WriteHeader(http.StatusUnauthorized)
		case ErrPaymentRequired:
			w.WriteHeader(http.StatusPaymentRequired)
		case ErrValidation:
			w.WriteHeader(http.StatusUnprocessableEntity)
		case ErrTooManyRequests:
			w.WriteHeader(http.StatusTooManyRequests)
		case ErrOrderAlreadyExists:
			w.WriteHeader(http.StatusConflict)
		case ErrOrderAlreadyUploaded:
			w.WriteHeader(http.StatusAccepted)
			return
		case ErrNoResponseData:
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponse{
		Type:    "internal_error",
		Message: err.Error(),
	})
}
