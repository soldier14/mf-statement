package domain

import (
	"fmt"
	"strings"
)

type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeNotFound   ErrorType = "not_found"
	ErrorTypeParse      ErrorType = "parse"
	ErrorTypeIO         ErrorType = "io"
	ErrorTypeInternal   ErrorType = "internal"
)

type DomainError struct {
	Type    ErrorType
	Message string
	Details map[string]interface{}
	Cause   error
}

func (e DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e DomainError) Unwrap() error {
	return e.Cause
}

func NewValidationError(message string, details map[string]interface{}) DomainError {
	return DomainError{
		Type:    ErrorTypeValidation,
		Message: message,
		Details: details,
	}
}

func NewParseError(message string, cause error) DomainError {
	return DomainError{
		Type:    ErrorTypeParse,
		Message: message,
		Cause:   cause,
	}
}

func NewIOError(message string, cause error) DomainError {
	return DomainError{
		Type:    ErrorTypeIO,
		Message: message,
		Cause:   cause,
	}
}

func NewNotFoundError(resource string) DomainError {
	return DomainError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func NewInternalError(message string, cause error) DomainError {
	return DomainError{
		Type:    ErrorTypeInternal,
		Message: message,
		Cause:   cause,
	}
}

func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	if domainErr, ok := err.(DomainError); ok {
		return domainErr.Type == ErrorTypeValidation
	}
	return false
}

func IsParseError(err error) bool {
	if err == nil {
		return false
	}
	if domainErr, ok := err.(DomainError); ok {
		return domainErr.Type == ErrorTypeParse
	}
	return false
}

func IsIOError(err error) bool {
	if err == nil {
		return false
	}
	if domainErr, ok := err.(DomainError); ok {
		return domainErr.Type == ErrorTypeIO
	}
	return false
}

type ErrorSummary struct {
	Errors []DomainError
}

func (es ErrorSummary) Error() string {
	if len(es.Errors) == 0 {
		return "no errors"
	}

	var messages []string
	for _, err := range es.Errors {
		messages = append(messages, err.Error())
	}

	return fmt.Sprintf("multiple errors: %s", strings.Join(messages, "; "))
}

func (es *ErrorSummary) AddError(err DomainError) {
	es.Errors = append(es.Errors, err)
}

func (es ErrorSummary) HasErrors() bool {
	return len(es.Errors) > 0
}
