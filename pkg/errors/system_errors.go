package errors

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
)

// ErrorType represents the type of an error
type ErrorType string

const (
	// ErrorTypeInternal represents an internal server error
	ErrorTypeInternal ErrorType = "INTERNAL"
	// ErrorTypeNotFound represents a not found error
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
	// ErrorTypeBadRequest represents a bad request error
	ErrorTypeBadRequest ErrorType = "BAD_REQUEST"
	// ErrorTypeUnauthorized represents an unauthorized error
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	// ErrorTypeForbidden represents a forbidden error
	ErrorTypeForbidden ErrorType = "FORBIDDEN"
	// ErrorTypeConflict represents a conflict error
	ErrorTypeConflict ErrorType = "CONFLICT"
	// ErrorTypeBiz represents a business error
	ErrorTypeBiz ErrorType = "BUSINESS"
	// ErrorTypeInvalidArgument represents an invalid argument error
	ErrorTypeInvalidArgument ErrorType = "INVALID_ARGUMENT"
	// ErrorTypeRateLimit represents a rate limit error
	ErrorTypeRateLimit ErrorType = "RATE_LIMIT"
	// ErrorTypeDatabase represents a database error
	ErrorTypeDatabase ErrorType = "DATABASE"
	// ErrorTypeMqHandle represents a message queue error
	ErrorTypeMqHandle ErrorType = "MQ"
	// ErrorTypeExternalInterface represents an external interface error
	ErrorTypeExternalInterface ErrorType = "EXTERNAL_INTERFACE"
)

// SystemError represents an application error
type SystemError struct {
	Type      ErrorType
	MessageId string
	Err       error
	Stack     string
}

// Error returns the error message
func (e *SystemError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.MessageId, e.Err.Error())
	}
	return e.MessageId
}

// Unwrap returns the wrapped error
func (e *SystemError) Unwrap() error {
	return e.Err
}

// System creates a new AppError
func System(errType ErrorType, messageId string, err error) *SystemError {
	stack := captureStack(2)
	return &SystemError{
		Type:      errType,
		MessageId: messageId,
		Err:       err,
		Stack:     stack,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf wraps an error with additional formatted context
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// As finds the first error in err's chain that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// HttpCode returns the HTTP status code for an error type
func HttpCode(errType ErrorType) int {
	switch errType {
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// captureStack captures the current stack trace
func captureStack(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			fmt.Fprintf(&builder, "%s:%d %s\n", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}
	return builder.String()
}

// ErrorNotifier is the interface for error notification
type ErrorNotifier interface {
	Notify(ctx context.Context, err error) error
}

// ErrorHandler handles errors and sends notifications
type ErrorHandler struct {
	logger    logging.Logger
	notifiers []ErrorNotifier
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger logging.Logger, notifiers ...ErrorNotifier) *ErrorHandler {
	return &ErrorHandler{
		logger:    logger,
		notifiers: notifiers,
	}
}

// Handle handles an error and sends notifications
func (h *ErrorHandler) Handle(ctx context.Context, err error) {
	var appErr *SystemError
	if errors.As(err, &appErr) {
		h.logger.Error(
			"msg", "Application error",
			"type", string(appErr.Type),
			"error", appErr.Error(),
			"stack", appErr.Stack,
		)
	} else {
		h.logger.Error("msg", "Unexpected error", "error", err.Error())
	}

	// Send error notifications
	for _, notifier := range h.notifiers {
		if notifyErr := notifier.Notify(ctx, err); notifyErr != nil {
			h.logger.Error("msg", "Failed to send error notification", "error", notifyErr.Error())
		}
	}
}
