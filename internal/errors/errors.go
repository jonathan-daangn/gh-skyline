// Package errors provides custom error types and utilities for the Skyline application.
package errors

import (
	"fmt"
)

// ErrorType represents categories of errors that can occur in the application
type ErrorType string

// Predefined error types for consistent error categorization
const (
	ValidationError ErrorType = "VALIDATION" // Input validation errors
	IOError         ErrorType = "IO"         // File/network I/O errors
	NetworkError    ErrorType = "NETWORK"    // Network communication errors
	GraphQLError    ErrorType = "GRAPHQL"    // GitHub GraphQL API errors
	STLError        ErrorType = "STL"        // STL file generation errors
)

// SkylineError provides structured error information including type and context
type SkylineError struct {
	Type    ErrorType // Category of the error
	Message string    // Human-readable error description
	Err     error     // Original error if wrapping another error
}

// Error implements the error interface for SkylineError
func (e *SkylineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// New creates a new SkylineError with the specified type, message, and wrapped error
func New(errType ErrorType, message string, err error) *SkylineError {
	return &SkylineError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// Wrap enhances an existing error with additional context while preserving its type
// If the original error is nil, returns nil
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	// If it's already a SkylineError, preserve the error type
	if skylineErr, ok := err.(*SkylineError); ok {
		return &SkylineError{
			Type:    skylineErr.Type,
			Message: message + ": " + skylineErr.Message,
			Err:     skylineErr.Err,
		}
	}

	// For other errors, treat as a generic error
	return &SkylineError{
		Type:    STLError, // Default to STLError for wrapped errors
		Message: message,
		Err:     err,
	}
}

// Is implements error matching for SkylineError
func (e *SkylineError) Is(target error) bool {
	t, ok := target.(*SkylineError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

// Unwrap implements error unwrapping for SkylineError
func (e *SkylineError) Unwrap() error {
	return e.Err
}
