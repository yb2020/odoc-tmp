package service

import (
	"context"
	"time"

	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
)

// Service is the interface that provides the basic service methods
type Service interface {
	Health(ctx context.Context) (bool, error)
	Echo(ctx context.Context, message string) (string, error)
}

// service is the implementation of the Service interface
type service struct {
	logger logging.Logger
}

// NewService creates a new service instance
func NewService(logger logging.Logger) Service {
	return &service{
		logger: logger,
	}
}

// Health returns the health status of the service
func (s *service) Health(ctx context.Context) (bool, error) {
	s.logger.Info("msg", "Health check requested")
	return true, nil
}

// Echo returns the message that was sent
func (s *service) Echo(ctx context.Context, message string) (string, error) {
	if message == "" {
		return "", errors.System(errors.ErrorTypeBadRequest, "message cannot be empty", nil)
	}

	s.logger.Info(
		"msg", "Echo requested",
		"message", message,
		"timestamp", time.Now().Format(time.RFC3339),
	)

	return message, nil
}
