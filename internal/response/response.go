package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response structure
// @Description Standard API response structure
// @Success 200 {object} Response
type Response struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Message string      `json:"message,omitempty" example:"Operation completed successfully"`
}

// ErrorInfo contains detailed error information
// @Description Error information structure
type ErrorInfo struct {
	Code    string                 `json:"code" example:"VALIDATION_ERROR"`
	Message string                 `json:"message" example:"Invalid request data"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ErrorResponse represents an error response structure for Swagger documentation
// @Description Standard error response structure
type ErrorResponse struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorInfo `json:"error"`
	Message string    `json:"message,omitempty"`
}

// Error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeInvalidToken   = "INVALID_TOKEN"
	ErrCodeExpiredToken   = "EXPIRED_TOKEN"
	ErrCodeInvalidAuth    = "INVALID_AUTH"
	ErrCodeDuplicateEmail = "DUPLICATE_EMAIL"
)

// Success sends a successful response with data
func Success(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// SuccessNoContent sends a successful response without data
func SuccessNoContent(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code string, message string, details map[string]interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest sends a 400 Bad Request error response
func BadRequest(c *gin.Context, message string, details map[string]interface{}) {
	Error(c, http.StatusBadRequest, ErrCodeBadRequest, message, details)
}

// Unauthorized sends a 401 Unauthorized error response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, ErrCodeUnauthorized, message, nil)
}

// Forbidden sends a 403 Forbidden error response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, ErrCodeForbidden, message, nil)
}

// NotFound sends a 404 Not Found error response
func NotFound(c *gin.Context, resource string) {
	Error(c, http.StatusNotFound, ErrCodeNotFound, resource+" not found", nil)
}

// Conflict sends a 409 Conflict error response
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, ErrCodeConflict, message, nil)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrCodeInternal, message, nil)
}

// ValidationError sends a 400 Bad Request error response for validation errors
func ValidationError(c *gin.Context, message string, details map[string]interface{}) {
	if details == nil {
		details = make(map[string]interface{})
	}
	Error(c, http.StatusBadRequest, ErrCodeValidation, message, details)
}

// InvalidToken sends a 401 Unauthorized error response for invalid token
func InvalidToken(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, ErrCodeInvalidToken, message, nil)
}

// ExpiredToken sends a 401 Unauthorized error response for expired token
func ExpiredToken(c *gin.Context) {
	Error(c, http.StatusUnauthorized, ErrCodeExpiredToken, "Token has expired", nil)
}

// InvalidAuth sends a 401 Unauthorized error response for invalid authentication
func InvalidAuth(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, ErrCodeInvalidAuth, message, nil)
}

// DuplicateEmail sends a 409 Conflict error response for duplicate email
func DuplicateEmail(c *gin.Context) {
	Error(c, http.StatusConflict, ErrCodeDuplicateEmail, "Email already registered", nil)
}
