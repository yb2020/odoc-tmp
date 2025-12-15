package response

import "net/http"

// HTTP 状态码常量
const (
	Status_Success = 1 // 业务成功
	Status_Fail    = 0 // 业务失败

	// 系统级别状态码 (9000-9999)
	Status_SystemError         = 9000 // 系统错误
	Status_DatabaseError       = 9001 // 数据库错误
	Status_CacheError          = 9002 // 缓存错误
	Status_InvalidRequest      = 9003 // 无效的请求
	Status_Unauthorized        = 9004 // 未授权
	Status_Forbidden           = 9005 // 禁止访问
	Status_NotFound            = 9006 // 资源不存在
	Status_MethodNotAllowed    = 9007 // 方法不允许
	Status_Conflict            = 9008 // 资源冲突
	Status_InternalServerError = 9009 // 内部服务器错误
	Status_ServiceUnavailable  = 9010 // 服务不可用
	Status_RequestTimeout      = 9011 // 请求超时
	Status_TooManyRequests     = 9012 // 请求过多

	// 认证状态码
	Status_AuthFailed       = 60001 // 认证失败
	Status_AuthInvalidToken = 50010 // 认证无效令牌
	Status_AuthExpiredToken = 60003 // 认证过期令牌
	Status_AuthInvalidRole  = 60004 // 认证无效角色
	Status_ServiceTokenNot  = 60005 // 服务令牌未找到

	// HTTP 状态码
	Code_Success       = http.StatusOK                  // HTTP成功
	Code_InvalidParams = http.StatusBadRequest          // HTTP无效参数
	Code_Unauthorized  = http.StatusUnauthorized        // HTTP未授权
	Code_Forbidden     = http.StatusForbidden           // HTTP禁止访问
	Code_NotFound      = http.StatusNotFound            // HTTP资源不存在
	Code_InternalError = http.StatusInternalServerError // HTTP内部服务器错误
)
