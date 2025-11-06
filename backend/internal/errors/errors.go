package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
	Code    string      `json:"code,omitempty"`
}

func BadRequest(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "Bad Request",
		Message: message,
		Details: details,
		Code:    "BAD_REQUEST",
	})
}

func NotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "Not Found",
		Message: resource + " not found",
		Code:    "NOT_FOUND",
	})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "Internal Server Error",
		Message: message,
		Code:    "INTERNAL_ERROR",
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error:   "Unauthorized",
		Message: message,
		Code:    "UNAUTHORIZED",
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Error:   "Forbidden",
		Message: message,
		Code:    "FORBIDDEN",
	})
}

func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Error:   "Too Many Requests",
		Message: message,
		Code:    "RATE_LIMIT_EXCEEDED",
	})
}
