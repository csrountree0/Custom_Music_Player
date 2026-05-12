// Package apierror defines the standard error shape returned by this API.
// Every error response looks like:
//
//	{ "error": { "code": "NOT_FOUND", "message": "...", "status": 404 } }
package apierror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *APIError) Error() string {
	return e.Message
}


func (e *APIError) Respond(c *gin.Context) {
	c.AbortWithStatusJSON(e.Status, gin.H{"error": e})
}


var (
	ErrNotFound = &APIError{
		Code:    "NOT_FOUND",
		Message: "the requested resource was not found",
		Status:  http.StatusNotFound,
	}
	ErrUnauthorized = &APIError{
		Code:    "UNAUTHORIZED",
		Message: "authentication is required",
		Status:  http.StatusUnauthorized,
	}
	ErrForbidden = &APIError{
		Code:    "FORBIDDEN",
		Message: "you do not have permission to perform this action",
		Status:  http.StatusForbidden,
	}
	ErrBadRequest = &APIError{
		Code:    "BAD_REQUEST",
		Message: "invalid request body or parameters",
		Status:  http.StatusBadRequest,
	}
	ErrInternal = &APIError{
		Code:    "INTERNAL_ERROR",
		Message: "an unexpected error occurred",
		Status:  http.StatusInternalServerError,
	}
	ErrConflict = &APIError{
		Code:    "CONFLICT",
		Message: "the resource already exists",
		Status:  http.StatusConflict,
	}
)

func New(status int, code, message string) *APIError {
	return &APIError{Code: code, Message: message, Status: status}
}
