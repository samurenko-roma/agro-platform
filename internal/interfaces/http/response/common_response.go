package response

import (
	"encoding/json"
	"net/http"
)

// CommandResponse стандартный ответ для CQRS команд
type CommandResponse[Data any, E comparable] struct {
	Success bool `json:"success"`
	Data    Data `json:"data,omitempty"`
	Error   E    `json:"error,omitempty"`
}

// ErrorInfo информация об ошибке
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Predefined error codes
const (
	CodeBadRequest   = "BAD_REQUEST"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeValidation   = "VALIDATION_ERROR"
	CodeInternal     = "INTERNAL_ERROR"
	CodeBusinessRule = "BUSINESS_RULE_VIOLATION"
)

// Success создает успешный ответ
func Success[T any](data T) *CommandResponse[T, interface{}] {
	return &CommandResponse[T, interface{}]{
		Success: true,
		Data:    data,
	}
}

// Error создает ответ с ошибкой
func Error(code, message string, details ...string) *CommandResponse[any, any] {
	errInfo := &ErrorInfo{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 && details[0] != "" {
		errInfo.Details = details[0]
	}
	return &CommandResponse[any, any]{
		Success: false,
		Error:   errInfo,
	}
}

// FromError создает ответ из стандартной ошибки
func FromError(err error) *CommandResponse[any, any] {
	// Попытка определить тип ошибки по строке
	errMsg := err.Error()

	switch {
	case contains(errMsg, "not found", "does not exist"):
		return Error(CodeNotFound, errMsg)
	case contains(errMsg, "already exists", "duplicate"):
		return Error(CodeConflict, errMsg)
	case contains(errMsg, "validation", "invalid", "required"):
		return Error(CodeValidation, errMsg)
	case contains(errMsg, "forbidden", "permission", "access denied"):
		return Error(CodeForbidden, errMsg)
	case contains(errMsg, "unauthorized", "authentication"):
		return Error(CodeUnauthorized, errMsg)
	default:
		return Error(CodeInternal, errMsg)
	}
}

// WriteJSON записывает CommandResponse в http.ResponseWriter
func (r *CommandResponse[any, any]) WriteJSON(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(r)
}

// WriteSuccess записывает успешный ответ
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	Success(data).WriteJSON(w, http.StatusOK)
}

// WriteCreated записывает ответ о создании ресурса
func WriteCreated(w http.ResponseWriter, data interface{}) {
	Success(data).WriteJSON(w, http.StatusCreated)
}

// WriteError записывает ответ с ошибкой
func WriteError(w http.ResponseWriter, statusCode int, code, message string, details ...string) {
	Error(code, message, details...).WriteJSON(w, statusCode)
}

// WriteValidationError записывает ошибку валидации
func WriteValidationError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, CodeValidation, message)
}

// WriteNotFound записывает ошибку "не найдено"
func WriteNotFound(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, CodeNotFound, message)
}

// WriteConflict записывает ошибку конфликта
func WriteConflict(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusConflict, CodeConflict, message)
}

// WriteInternalError записывает внутреннюю ошибку сервера
func WriteInternalError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, CodeInternal, message)
}

// WriteUnauthorized записывает ошибку авторизации
func WriteUnauthorized(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusUnauthorized, CodeUnauthorized, message)
}

// WriteForbidden записывает ошибку доступа
func WriteForbidden(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusForbidden, CodeForbidden, message)
}

// helper functions
func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) && (s == substr || len(substr) > 0) {
			// simple contains check
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}
