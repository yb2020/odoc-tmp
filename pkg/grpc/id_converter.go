package grpc

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/yb2020/odoc/pkg/logging"
)

// IDConverter 是一个用于在gRPC服务中处理ID转换的工具
type IDConverter struct {
	logger logging.Logger
}

// NewIDConverter 创建一个新的ID转换器
func NewIDConverter(logger logging.Logger) *IDConverter {
	return &IDConverter{
		logger: logger,
	}
}

// StringToInt64 将字符串转换为int64，处理错误并返回gRPC状态
func (c *IDConverter) StringToInt64(ctx context.Context, idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.logger.Warn("msg", "无效的ID参数", "id", idStr, "error", err.Error())
		return 0, status.Errorf(codes.InvalidArgument, "无效的ID参数: %s", idStr)
	}
	return id, nil
}

// Int64ToUint64 将int64转换为uint64
func (c *IDConverter) Int64ToUint64(id int64) uint64 {
	return uint64(id)
}

// Uint64ToInt64 将uint64转换为int64
func (c *IDConverter) Uint64ToInt64(id uint64) int64 {
	return int64(id)
}

// UnaryServerInterceptor 返回一个gRPC拦截器，用于自动处理请求中的ID转换
func (c *IDConverter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 这里可以添加请求前的ID转换逻辑
		resp, err := handler(ctx, req)
		// 这里可以添加响应后的ID转换逻辑
		return resp, err
	}
}

// ConvertResponseIDs 转换响应中的ID字段
func (c *IDConverter) ConvertResponseIDs(resp interface{}) interface{} {
	// 这里可以实现响应中ID字段的转换逻辑
	// 通过反射遍历结构体字段，将int64类型的ID字段转换为uint64
	return resp
}
