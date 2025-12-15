package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/response"
)

func ErrorHandler(errorReporter *errors.ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err, errorReporter)
			c.Abort()
		}
	}
}

func handleError(c *gin.Context, err error, errorReporter *errors.ErrorHandler) {
	// 获取本地化器
	localizer := i18n.GetLocalizer()

	// 尝试转换为业务错误
	var bizErr *errors.BizError
	if errors.As(err, &bizErr) {
		// 本地化消息
		msg := localizer.Localize(bizErr.MsgID, c)
		response.BizErrorNoData(c, bizErr.Status, msg)
		return
	}

	// 尝试转换为系统错误
	var systemErr *errors.SystemError
	if errors.As(err, &systemErr) {
		// 本地化消息
		msg := localizer.Localize(systemErr.MessageId, c)
		// 记录错误
		errorReporter.Handle(c, err)
		response.BizErrorNoData(c, response.Status_Fail, msg)
		return
	}

	// 默认错误处理
	response.BizErrorNoData(c, response.Status_Fail, err.Error())
}
