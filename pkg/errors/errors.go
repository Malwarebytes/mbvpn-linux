package errors

import (
	"errors"
	"fmt"
)

// Standard errors that can be used for comparison
var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrUnauthorized is returned when the user is not authenticated
	ErrUnauthorized = errors.New("unauthorized access")

	// ErrPermissionDenied is returned when the user doesn't have permission
	ErrPermissionDenied = errors.New("permission denied")

	// ErrInvalidInput is returned when the input is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrServerError is returned when a server error occurs
	ErrServerError = errors.New("server error")

	// ErrNetworkError is returned when a network error occurs
	ErrNetworkError = errors.New("network error")

	// ErrConfigError is returned when there's an issue with configuration
	ErrConfigError = errors.New("configuration error")

	// ErrVPNError is returned when there's an issue with VPN operations
	ErrVPNError = errors.New("vpn operation error")
)

// Error types for specific domains

// ConfigError represents errors related to configuration
type ConfigError struct {
	Op  string // Operation that was being performed
	Err error  // Underlying error
}

func (e *ConfigError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("config error during %s", e.Op)
	}
	return fmt.Sprintf("config error during %s: %v", e.Op, e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new ConfigError
func NewConfigError(op string, err error) *ConfigError {
	return &ConfigError{Op: op, Err: err}
}

// AuthError represents authentication and authorization errors
type AuthError struct {
	Op  string // Operation that was being performed
	Err error  // Underlying error
}

func (e *AuthError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("auth error during %s", e.Op)
	}
	return fmt.Sprintf("auth error during %s: %v", e.Op, e.Err)
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// NewAuthError creates a new AuthError
func NewAuthError(op string, err error) *AuthError {
	return &AuthError{Op: op, Err: err}
}

// NetworkError represents network and API communication errors
type NetworkError struct {
	Op  string // Operation that was being performed
	Err error  // Underlying error
}

func (e *NetworkError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("network error during %s", e.Op)
	}
	return fmt.Sprintf("network error during %s: %v", e.Op, e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// NewNetworkError creates a new NetworkError
func NewNetworkError(op string, err error) *NetworkError {
	return &NetworkError{Op: op, Err: err}
}

// VPNError represents errors related to VPN operations
type VPNError struct {
	Op  string // Operation that was being performed
	Err error  // Underlying error
}

func (e *VPNError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("vpn error during %s", e.Op)
	}
	return fmt.Sprintf("vpn error during %s: %v", e.Op, e.Err)
}

func (e *VPNError) Unwrap() error {
	return e.Err
}

// NewVPNError creates a new VPNError
func NewVPNError(op string, err error) *VPNError {
	return &VPNError{Op: op, Err: err}
}

// UserError represents errors that should be displayed to the user
type UserError struct {
	Message string // User-friendly error message
	Err     error  // Underlying error (for logging)
}

func (e *UserError) Error() string {
	return e.Message
}

func (e *UserError) Unwrap() error {
	return e.Err
}

// NewUserError creates a new UserError
func NewUserError(message string, err error) *UserError {
	return &UserError{Message: message, Err: err}
}

// IsUserError checks if an error is a UserError
func IsUserError(err error) bool {
	var ue *UserError
	return errors.As(err, &ue)
}

// GetUserErrorMessage gets the user-friendly message from an error
func GetUserErrorMessage(err error) string {
	var ue *UserError
	if errors.As(err, &ue) {
		return ue.Message
	}
	return "An unexpected error occurred"
}

// IsNotFound checks if an error is a NotFound error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorized checks if an error is an Unauthorized error
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsPermissionDenied checks if an error is a PermissionDenied error
func IsPermissionDenied(err error) bool {
	return errors.Is(err, ErrPermissionDenied)
}

// IsInvalidInput checks if an error is an InvalidInput error
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsServerError checks if an error is a ServerError error
func IsServerError(err error) bool {
	return errors.Is(err, ErrServerError)
}

// IsNetworkError checks if an error is a NetworkError error
func IsNetworkError(err error) bool {
	return errors.Is(err, ErrNetworkError)
}

// IsConfigError checks if an error is a ConfigError error
func IsConfigError(err error) bool {
	return errors.Is(err, ErrConfigError)
}

// IsVPNError checks if an error is a VPNError error
func IsVPNError(err error) bool {
	return errors.Is(err, ErrVPNError)
}