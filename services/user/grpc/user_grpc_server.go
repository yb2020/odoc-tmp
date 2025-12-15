package grpc

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/yb2020/odoc-proto/gen/go/user"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/user/service"
)

// UserGRPCServer 用户服务的gRPC服务器
type UserGRPCServer struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
	logger      logging.Logger
	tracer      opentracing.Tracer
}

// NewUserGRPCServer 创建一个新的用户gRPC服务器
func NewUserGRPCServer(logger logging.Logger, tracer opentracing.Tracer, userService *service.UserService) *UserGRPCServer {
	return &UserGRPCServer{
		userService: userService,
		logger:      logger,
		tracer:      tracer,
	}
}

// RegisterServer 注册gRPC服务
func (s *UserGRPCServer) RegisterServer(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, s)
}

// GetProfile 获取当前用户信息
func (s *UserGRPCServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	// 创建一个新的span
	span := s.tracer.StartSpan("grpc-user-get-profile-endpoint")
	defer span.Finish()

	// 将span添加到上下文
	ctx = opentracing.ContextWithSpan(ctx, span)

	// 调用服务层获取用户信息
	response, err := s.userService.GetProfile(ctx, req)
	if err != nil {
		s.logger.Error("msg", "gRPC获取用户信息失败", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "获取用户信息失败: %v", err)
	}

	return response, nil
}

// CheckEmailExists 检查邮箱是否已被注册
func (s *UserGRPCServer) CheckEmailExists(ctx context.Context, req *pb.GetExistsByEmailRequest) (*pb.GetExistsByEmailResponse, error) {
	// 创建一个新的span
	span := s.tracer.StartSpan("grpc-user-check-email-exists-endpoint")
	defer span.Finish()

	// 将span添加到上下文
	ctx = opentracing.ContextWithSpan(ctx, span)

	// 调用服务层检查邮箱是否存在
	response, err := s.userService.CheckEmailExists(ctx, req.Email)
	if err != nil {
		s.logger.Error("msg", "gRPC检查邮箱是否存在失败", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "检查邮箱是否存在失败: %v", err)
	}

	return &pb.GetExistsByEmailResponse{
		Exists: response,
	}, nil
}

// Register 用户注册
func (s *UserGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 创建一个新的span
	span := s.tracer.StartSpan("grpc-user-register-endpoint")
	defer span.Finish()

	// 将span添加到上下文
	ctx = opentracing.ContextWithSpan(ctx, span)

	// 调用服务层注册用户
	response, err := s.userService.Register(ctx, req)
	if err != nil {
		s.logger.Error("msg", "gRPC注册用户失败", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "注册用户失败: %v", err)
	}

	return response, nil
}

// UpdateProfile 更新用户信息
func (s *UserGRPCServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	// 创建一个新的span
	span := s.tracer.StartSpan("grpc-user-update-profile-endpoint")
	defer span.Finish()

	// 将span添加到上下文
	ctx = opentracing.ContextWithSpan(ctx, span)

	// 调用服务层更新用户信息
	response, err := s.userService.UpdateProfile(ctx, req)
	if err != nil {
		s.logger.Error("msg", "gRPC更新用户信息失败", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "更新用户信息失败: %v", err)
	}

	return response, nil
}
