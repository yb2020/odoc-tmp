package exception

import "fmt"

// DifyApiException Dify API异常
type DifyApiException struct {
	StatusCode int    // HTTP状态码
	ErrorCode  string // 错误代码
	Message    string // 错误消息
}

// NewDifyApiException 创建新的Dify API异常
func NewDifyApiException(statusCode int, errorCode, message string) *DifyApiException {
	return &DifyApiException{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

// Error 实现error接口
func (e *DifyApiException) Error() string {
	return fmt.Sprintf("Dify API异常 [状态码: %d, 错误代码: %s, 消息: %s]", e.StatusCode, e.ErrorCode, e.Message)
}

// GetStatusCode 获取HTTP状态码
func (e *DifyApiException) GetStatusCode() int {
	return e.StatusCode
}

// GetErrorCode 获取错误代码
func (e *DifyApiException) GetErrorCode() string {
	return e.ErrorCode
}

// GetMessage 获取错误消息
func (e *DifyApiException) GetMessage() string {
	return e.Message
}
