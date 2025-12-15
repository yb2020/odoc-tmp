package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	pb "github.com/yb2020/odoc-proto/gen/go/common"
)

// NewAPIResponse 创建一个新的 API 响应
func Success(c *gin.Context,
	message string, data proto.Message) {
	// 设置响应头
	handleHeaderAndHttpStatus(c, Code_Success)
	handleResponse(c, Code_Success, Status_Success, message, data)
}

// NewAPIResponse 创建一个新的 API 响应
func SuccessNoData(c *gin.Context,
	message string) {
	// 设置响应头
	handleHeaderAndHttpStatus(c, Code_Success)
	handleResponse(c, Code_Success, Status_Success, message, nil)
}

// Error 创建一个业务错误响应，通用的业务状态码（兼容原来的业务）
func Error(c *gin.Context,
	message string, data proto.Message) {
	// 设置响应头
	handleHeaderAndHttpStatus(c, Code_Success)
	handleResponse(c, Code_Success, Status_Fail, message, data)
}

// ErrorNoData 创建一个业务错误响应，通用的业务状态码（兼容原来的业务）
func ErrorNoData(c *gin.Context,
	message string) {
	// 设置响应头
	handleHeaderAndHttpStatus(c, Code_Success)
	handleResponse(c, Code_Success, Status_Fail, message, nil)
}

// BizError 创建一个业务错误响应, 自定义业务状态码（不兼容原来的业务）
func BizError(c *gin.Context,
	status int32, message string, data proto.Message) {
	// 设置响应头
	handleHeaderAndHttpStatus(c, Code_Success)
	handleResponse(c, Code_Success, status, message, data)
}

// BizErrorNoData 创建一个业务错误响应, 自定义业务状态码（不兼容原来的业务）
func BizErrorNoData(c *gin.Context,
	status int32, message string) {
	// 设置响应头
	handleHeaderAndHttpStatus(c, Code_Success)
	handleResponse(c, Code_Success, status, message, nil)
}

// SystemErrorNoData 创建一个系统错误响应，返回的状态码，http状态码都由自已定义
func SystemErrorNoData(c *gin.Context,
	httpCode int, status int32, message string) {
	// 设置响应头
	handleHeaderAndHttpStatus(c, httpCode)
	handleResponse(c, httpCode, status, message, nil)
}

// SystemError 创建一个系统错误响应，返回的状态码，http状态码都由自已定义
func SystemError(c *gin.Context,
	httpCode int, status int32, message string, data proto.Message) {
	// 设置响应头
	handleHeaderAndHttpStatus(c, httpCode)
	handleResponse(c, httpCode, status, message, data)
}

// SSESuccess 创建一个SSE成功响应
func SSESuccess(c *gin.Context, message string, data proto.Message) {
	// 创建响应
	resp := &pb.APIResponse{
		Status:  Status_Success,
		Message: message,
	}

	// 如果有数据，将其包装到 Any 中
	if data != nil {
		anyData, err := anypb.New(data)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Data = anyData
	}

	// 序列化为JSON
	opts := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
	}
	jsonData, err := opts.Marshal(resp)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 解析为map以处理@type字段
	var responseMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &responseMap); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 移除@type字段
	removeAnyTypeField(responseMap)

	// 重新序列化
	finalJSON, err := json.Marshal(responseMap)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Fprintf(c.Writer, "data: %s\n\n", string(finalJSON))
}

// SSEError 创建一个SSE错误响应
func SSEError(c *gin.Context, message string) {
	// 设置SSE头信息
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	// 创建错误响应
	resp := &pb.APIResponse{
		Status:  Status_Fail,
		Message: message,
	}

	// 序列化为JSON
	opts := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
	}
	jsonData, err := opts.Marshal(resp)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Fprintf(c.Writer, "data: %s\n\n", string(jsonData))
	c.Writer.Flush() // 确保数据立即发送到客户端
}

func handleHeaderAndHttpStatus(c *gin.Context, httpStatusCode int) {
	c.Header("Content-Type", "application/json")
	c.Status(httpStatusCode)
}

func handleResponse(c *gin.Context, httpCode int, status int32, message string, data proto.Message) {
	// 创建响应
	resp := &pb.APIResponse{
		Status:  status,
		Message: message,
	}

	// 如果有数据，将其包装到 Any 中
	if data != nil {
		anyData, err := anypb.New(data)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Data = anyData
	}

	sendResponse(c, httpCode, resp)
}

// sendResponse 将Proto响应序列化为JSON并发送给客户端
func sendResponse(c *gin.Context, httpCode int, resp *pb.APIResponse) {
	// 配置 protojson 的序列化选项
	opts := protojson.MarshalOptions{
		UseProtoNames:   true, // 使用 proto 字段名而不是驼峰命名
		EmitUnpopulated: true, // 包含未设置的字段
		UseEnumNumbers:  true, // 使用枚举数字而不是名称
	}

	// 序列化 Proto 消息
	respBytes, err := opts.Marshal(resp)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 解析为 map 以便于处理
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBytes, &responseMap); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 移除@type字段
	removeAnyTypeField(responseMap)

	// 写入响应
	c.JSON(httpCode, responseMap)
}

// removeAnyTypeField 处理响应中的 Any 类型字段，移除 @type
func removeAnyTypeField(responseMap map[string]interface{}) {
	if dataAny, ok := responseMap["data"]; ok && dataAny != nil {
		if dataMap, ok := dataAny.(map[string]interface{}); ok {
			if _, ok := dataMap["@type"].(string); ok {
				// 移除 @type 字段，它只是 Any 类型的元数据
				delete(dataMap, "@type")

				// 如果有其他字段，则保留它们
				if len(dataMap) > 0 {
					responseMap["data"] = dataMap
				} else {
					// 如果没有其他字段，则删除 data 字段
					delete(responseMap, "data")
				}
			}
		}
	} else {
		// 如果data为nil或不存在，删除data字段
		delete(responseMap, "data")
	}
}
