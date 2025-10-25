package errors

import (
	"errors"
	"fmt"
	"os"

	"github.com/Malwarebytes/mbvpn-linux/pkg/config"
	"github.com/Malwarebytes/mbvpn-linux/pkg/output"
	log "github.com/sirupsen/logrus"
)

// Handler is responsible for handling errors in a consistent way
type Handler struct {
	debug bool
}

// NewHandler creates a new error handler
func NewHandler() *Handler {
	return &Handler{
		debug: config.Debug(),
	}
}

// Handle handles an error by logging it and/or displaying it to the user
// Returns true if the process should exit with an error code
func (h *Handler) Handle(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's a user error (for friendly display)
	if userErr := userError(err); userErr != nil {
		output.PrintMsg(userErr.Message, output.MsgError)
		if h.debug && userErr.Err != nil {
			log.Errorf("Error details: %v", userErr.Err)
		}
		return false // Don't exit for user errors
	}

	// For other errors, log them and show a generic message to the user
	log.Errorf("Error occurred: %v", err)

	// Show a more appropriate message based on error type
	switch {
	case IsNotFound(err):
		output.PrintMsg("The requested resource was not found.", output.MsgError)
	case IsUnauthorized(err):
		output.PrintMsg("You are not authorized. Please log in and try again.", output.MsgError)
	case IsPermissionDenied(err):
		output.PrintMsg("You don't have permission to perform this action.", output.MsgError)
	case IsInvalidInput(err):
		output.PrintMsg("Invalid input provided.", output.MsgError)
	case IsServerError(err):
		output.PrintMsg("A server error occurred. Please try again later.", output.MsgError)
	case IsNetworkError(err):
		output.PrintMsg("A network error occurred. Please check your connection and try again.", output.MsgError)
	case IsConfigError(err):
		output.PrintMsg("A configuration error occurred.", output.MsgError)
	case IsVPNError(err):
		output.PrintMsg("A VPN error occurred.", output.MsgError)
	default:
		output.PrintMsg("An unexpected error occurred.", output.MsgError)
	}

	// In debug mode, print the full error to stderr
	if h.debug {
		fmt.Fprintf(os.Stderr, "Error details: %v\n", err)
	}

	return true // Exit with error code for non-user errors
}

// HandleWithMessage handles an error with a custom message
func (h *Handler) HandleWithMessage(err error, message string) bool {
	if err == nil {
		return false
	}

	// Log the error
	log.Errorf("%s: %v", message, err)

	// Display the message to the user
	output.PrintMsg(message, output.MsgError)

	// In debug mode, print the full error to stderr
	if h.debug {
		fmt.Fprintf(os.Stderr, "Error details: %v\n", err)
	}

	return true // Exit with error code
}

// userError extracts a UserError from an error chain if present
func userError(err error) *UserError {
	var ue *UserError
	if err != nil {
		if errors.As(err, &ue) {
			return ue
		}
	}
	return nil
}

// Fatal logs an error and exits the program
func (h *Handler) Fatal(err error) {
	if h.Handle(err) {
		os.Exit(1)
	}
}

// FatalWithMessage logs an error with a custom message and exits the program
func (h *Handler) FatalWithMessage(err error, message string) {
	if h.HandleWithMessage(err, message) {
		os.Exit(1)
	}
}

// Error converts an error to a user-friendly error message
func (h *Handler) Error(err error, msg string) error {
	return &UserError{
		Message: msg,
		Err:     err,
	}
}
