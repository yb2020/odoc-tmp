// pkg/squid_proxy/errors.go
package squid_proxy

import "fmt"

// ErrorCode 错误代码类型
type ErrorCode int

// 错误代码常量
const (
	ErrUnknown ErrorCode = iota
	ErrInvalidConfig
	ErrInvalidURL
	ErrInvalidProxy
	ErrProxyConnection
	ErrRequestFailed
	ErrHTTPError
	ErrFileTooLarge
	ErrReadFailed
	ErrWriteFailed
	ErrInvalidResponse
)

// SquidError Squid代理错误
type SquidError struct {
	Code    ErrorCode // 错误代码
	Message string    // 错误消息
	Cause   error     // 原始错误
}

// Error 实现error接口
func (e *SquidError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code.String(), e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code.String(), e.Message)
}

// Unwrap 支持错误链
func (e *SquidError) Unwrap() error {
	return e.Cause
}

// String 错误代码字符串表示
func (code ErrorCode) String() string {
	switch code {
	case ErrUnknown:
		return "UNKNOWN"
	case ErrInvalidConfig:
		return "INVALID_CONFIG"
	case ErrInvalidURL:
		return "INVALID_URL"
	case ErrInvalidProxy:
		return "INVALID_PROXY"
	case ErrProxyConnection:
		return "PROXY_CONNECTION"
	case ErrRequestFailed:
		return "REQUEST_FAILED"
	case ErrHTTPError:
		return "HTTP_ERROR"
	case ErrFileTooLarge:
		return "FILE_TOO_LARGE"
	case ErrReadFailed:
		return "READ_FAILED"
	case ErrWriteFailed:
		return "WRITE_FAILED"
	case ErrInvalidResponse:
		return "INVALID_RESPONSE"
	default:
		return "UNKNOWN"
	}
}

// NewSquidError 创建新的Squid错误
func NewSquidError(code ErrorCode, message string) *SquidError {
	return &SquidError{
		Code:    code,
		Message: message,
	}
}

// NewSquidErrorWithCause 创建带原因的Squid错误
func NewSquidErrorWithCause(code ErrorCode, message string, cause error) *SquidError {
	return &SquidError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsSquidError 检查是否为Squid错误
func IsSquidError(err error) bool {
	_, ok := err.(*SquidError)
	return ok
}

// GetSquidError 获取Squid错误
func GetSquidError(err error) *SquidError {
	if squidErr, ok := err.(*SquidError); ok {
		return squidErr
	}
	return nil
}

// IsRetryableError 判断错误是否可重试
func IsRetryableError(err error) bool {
	squidErr := GetSquidError(err)
	if squidErr == nil {
		return true // 未知错误，尝试重试
	}

	switch squidErr.Code {
	case ErrProxyConnection, ErrRequestFailed, ErrHTTPError:
		return true
	case ErrInvalidURL, ErrInvalidConfig, ErrInvalidProxy, ErrFileTooLarge:
		return false
	default:
		return true
	}
}

// IsNetworkError 判断是否为网络相关错误
func IsNetworkError(err error) bool {
	squidErr := GetSquidError(err)
	if squidErr == nil {
		return false
	}

	switch squidErr.Code {
	case ErrProxyConnection, ErrRequestFailed:
		return true
	default:
		return false
	}
}

// IsConfigError 判断是否为配置错误
func IsConfigError(err error) bool {
	squidErr := GetSquidError(err)
	if squidErr == nil {
		return false
	}

	switch squidErr.Code {
	case ErrInvalidConfig, ErrInvalidURL, ErrInvalidProxy:
		return true
	default:
		return false
	}
}
