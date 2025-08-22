package response

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIResponse is the standard structure for all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error response.
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginatedData is a standard structure for paginated responses.
type PaginatedData struct {
	Items    interface{} `json:"items"`
	Page     int         `json:"page"`
	Size     int         `json:"size"`
	Total    int64       `json:"total"`
}

// SuccessMessage is a standard structure for simple success messages.
type SuccessMessage struct {
	Message string `json:"message"`
}

// NewPaginatedData creates a new PaginatedData struct.
func NewPaginatedData(items interface{}, page, size int, total int64) *PaginatedData {
	return &PaginatedData{
		Items:    items,
		Page:     page,
		Size:     size,
		Total:    total,
	}
}

func successResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

func errorResponse(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// --- Success Responses ---

func Success(c *gin.Context, data interface{}) {
	successResponse(c, http.StatusOK, data)
}

func Created(c *gin.Context, data interface{}) {
	successResponse(c, http.StatusCreated, data)
}

// --- Error Responses ---

func BadRequest(c *gin.Context, code, message string) {
	errorResponse(c, http.StatusBadRequest, code, message)
}

func Unauthorized(c *gin.Context, code, message string) {
	errorResponse(c, http.StatusUnauthorized, code, message)
}

func Forbidden(c *gin.Context, code, message string) {
	errorResponse(c, http.StatusForbidden, code, message)
}

func NotFound(c *gin.Context, code, message string) {
	errorResponse(c, http.StatusNotFound, code, message)
}

func Conflict(c *gin.Context, code, message string) {
	errorResponse(c, http.StatusConflict, code, message)
}

func InternalError(c *gin.Context, code, message string) {
	errorResponse(c, http.StatusInternalServerError, code, message)
}

// ValidationError handles validation errors from the validator package.
func ValidationError(c *gin.Context, err error) {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, fmt.Sprintf("Field '%s' failed on the '%s' tag", err.Field(), err.Tag()))
	}

	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid input provided.",
			Details: errors,
		},
	})
}