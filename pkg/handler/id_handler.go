package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/serializer"
)

// GinIDHandlerFunc 是一个处理ID参数的Gin处理函数类型
type GinIDHandlerFunc func(c *gin.Context, id int64)

// WithGinIDParam 包装一个Gin处理函数，自动处理ID参数的转换
func WithGinIDParam(handlerFunc GinIDHandlerFunc, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少ID参数"})
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Warn("msg", "无效的ID参数", "id", idStr, "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
			return
		}

		handlerFunc(c, id)
	}
}

// 保留原来的函数以兼容旧代码
// IDHandlerFunc 是一个处理ID参数的HTTP处理函数类型
type IDHandlerFunc func(w http.ResponseWriter, r *http.Request, id int64)

// WithIDParam 包装一个HTTP处理函数，自动处理ID参数的转换 (已弃用，请使用WithGinIDParam)
func WithIDParam(handlerFunc IDHandlerFunc, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 这个函数保留以兼容旧代码，但不再推荐使用
		logger.Warn("msg", "使用已弃用的WithIDParam函数，请迁移到WithGinIDParam")

		// 尝试从URL路径中提取ID参数
		path := r.URL.Path
		parts := strings.Split(path, "/")
		var idStr string
		for i, part := range parts {
			if i > 0 && parts[i-1] == "id" {
				idStr = part
				break
			}
		}

		if idStr == "" {
			http.Error(w, "缺少ID参数", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Warn("msg", "无效的ID参数", "id", idStr, "error", err.Error())
			http.Error(w, "无效的ID参数", http.StatusBadRequest)
			return
		}

		handlerFunc(w, r, id)
	}
}

// JSONResponse 发送JSON响应，自动将int64 ID转换为uint64
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// 使用自定义的编码器处理int64到uint64的转换
	resp := serializer.JSONResponse{Data: data}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(resp); err != nil {
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
	}
}
