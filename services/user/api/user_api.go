package api

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	pb "github.com/yb2020/odoc-proto/gen/go/user"
	"github.com/yb2020/odoc/internal/biz"
	"github.com/yb2020/odoc/pkg/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/services/user/service"
)

// UserAPI 用户API处理器
type UserAPI struct {
	userService *service.UserService
	logger      logging.Logger
	tracer      opentracing.Tracer
	localizer   i18n.Localizer
}

// NewUserAPI 创建用户API处理器
func NewUserAPI(logger logging.Logger, tracer opentracing.Tracer,
	localizer i18n.Localizer,
	userService *service.UserService,
) *UserAPI {
	return &UserAPI{
		userService: userService,
		logger:      logger,
		tracer:      tracer,
		localizer:   localizer,
	}
}

// GetProfile 获取用户个人资料
func (api *UserAPI) GetProfile(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.GetProfile")
	defer span.Finish()

	req := &pb.GetProfileRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析获取用户个人资料请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	resp, err := api.userService.GetProfile(ctx, req)
	if err != nil {
		api.logger.Warn("msg", "获取用户个人资料失败", "error", err.Error())
		c.Error(err)
		return
	}

	response.Success(c, "success", resp)
}

// CheckEmailExists 检查邮箱是否存在
func (api *UserAPI) CheckEmailExists(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.CheckEmailExists")
	defer span.Finish()

	req := &pb.GetExistsByEmailRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析检查邮箱是否存在请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 直接使用请求的上下文，它已经在中间件中被更新
	resp, err := api.userService.CheckEmailExists(ctx, req.Email)
	if err != nil {
		api.logger.Warn("msg", "检查邮箱是否存在失败", "error", err.Error())
		c.Error(err)
		return
	}

	response.Success(c, "success", &pb.GetExistsByEmailResponse{
		Exists: resp,
	})
}

// Register 用户注册
func (api *UserAPI) Register(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.Register")
	defer span.Finish()

	req := &pb.RegisterRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析用户注册请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	resp, err := api.userService.Register(ctx, req)
	if err != nil {
		api.logger.Warn("msg", "用户注册失败", "error", err.Error())
		c.Error(err)
		return
	}

	response.Success(c, "success", resp)
}

// 自定义验证
//
//	if err := validateUsername(req.Email); err != nil {
//		api.logger.Warn("msg", "用户注册失败", "error", err.Error())
//		c.Error(err)
//		return
//	}
//
// 改用户密码时，需要使用
func validateUsername(username string) error {
	// 检查是否以下划线开头
	if strings.HasPrefix(username, "_") {
		return errors.New("username cannot start with an underscore")
	}

	// 检查是否以下划线结尾
	if strings.HasSuffix(username, "_") {
		return errors.New("username cannot end with an underscore")
	}

	// 检查是否全是数字
	if regexp.MustCompile(`^[0-9]+$`).MatchString(username) {
		return errors.New("username cannot be all numbers")
	}

	// 检查是否是重复字符
	if regexp.MustCompile(`^(.)(\1)+$`).MatchString(username) {
		return errors.New("username is exists")
	}

	// 检查是否是连续递增数字
	if regexp.MustCompile(`^0*1*2*3*4*5*6*7*8*9*$`).MatchString(username) {
		return errors.New("username is exists")
	}

	// 检查是否是连续递减数字
	if regexp.MustCompile(`^9*8*7*6*5*4*3*2*1*0*$`).MatchString(username) {
		return errors.New("username is exists")
	}

	return nil
}

// UpdateProfile 更新用户个人资料
func (api *UserAPI) UpdateProfile(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.UpdateProfile")
	defer span.Finish()

	req := &pb.UpdateProfileRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "解析更新个人资料请求失败", "error", err.Error())
		c.Error(err)
		return
	}

	resp, err := api.userService.UpdateProfile(ctx, req)
	if err != nil {
		api.logger.Warn("msg", "更新个人资料失败", "error", err.Error())
		c.Error(err)
		return
	}

	response.Success(c, "success", resp)
}

// GetById 根据ID获取用户
func (api *UserAPI) GetById(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.GetById")
	defer span.Finish()

	req := &pb.GetByIdRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "无效的用户ID", "error", err.Error())
		c.Error(err)
		return
	}
	user, err := api.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		c.Error(err)
		return
	}
	if user == nil {
		response.BizErrorNoData(c, biz.User_StatusInvalidUserData, "invalid_request")
		return
	}

	response.Success(c, "success", &pb.GetByIdResponse{User: user.ToProto()})
}

// GetByIds 根据多个ID获取用户
func (api *UserAPI) GetByIds(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.GetByIds")
	defer span.Finish()

	req := &pb.GetByIdsRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "无效的用户ID", "error", err.Error())
		c.Error(err)
		return
	}

	resp, err := api.userService.GetByIds(ctx, req.Ids)
	if err != nil {
		c.Error(err)
		return
	}

	// 将用户模型转换为 protobuf 消息
	pbUsers := make([]*pb.User, 0, len(resp))
	for _, user := range resp {
		pbUsers = append(pbUsers, user.ToProto())
	}

	response.Success(c, "success", &pb.GetByIdsResponse{
		Users: pbUsers,
	})
}

// PaginationUsers 分页获取用户列表
func (api *UserAPI) PaginationUsers(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.PaginationUsers")
	defer span.Finish()

	req := &pb.PaginationUsersRequest{}

	if err := transport.BindProto(c, req); err != nil {
		api.logger.Warn("msg", "无效的分页参数", "error", err.Error())
		c.Error(err)
		return
	}
	resp, err := api.userService.PaginationUsers(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, "success", resp)
}

// CreateUser 创建用户
func (api *UserAPI) CreateUser(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.CreateUser")
	defer span.Finish()

	var req pb.CreateUserRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Warn("msg", "无效的用户数据", "error", err.Error())
		c.Error(err)
		return
	}

	resp, err := api.userService.CreateUserAdmin(ctx, &req)
	if err != nil {
		api.logger.Warn("msg", "创建用户失败", "error", err.Error())
		c.Error(err)
		return
	}
	response.Success(c, "success", resp.User)
}

// UpdateUser 更新用户
func (api *UserAPI) UpdateUser(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.UpdateUser")
	defer span.Finish()

	var req pb.UpdateUserRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Warn("msg", "无效的用户数据", "error", err.Error())
		c.Error(err)
		return
	}

	resp, err := api.userService.UpdateUserAdmin(ctx, &req)
	if err != nil {
		api.logger.Warn("msg", "更新用户失败", "error", err.Error())
		c.Error(err)
		return
	}

	response.Success(c, "success", resp.User)
}

// DeleteUserById 删除用户
func (api *UserAPI) DeleteUserById(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.DeleteUserById")
	defer span.Finish()

	var req pb.DeleteUserByIdRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Warn("msg", "无效的用户ID", "error", err.Error())
		c.Error(err)
		return
	}

	err := api.userService.DeleteUserById(ctx, req.Id)
	if err != nil {
		api.logger.Warn("msg", "删除用户失败", "error", err.Error())
		c.Error(err)
		return
	}

	response.SuccessNoData(c, "success")
}

// DeleteUserByIds 批量删除用户
func (api *UserAPI) DeleteUserByIds(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "UserAPI.DeleteUserByIds")
	defer span.Finish()

	var req pb.DeleteUserByIdsRequest
	if err := transport.BindProto(c, &req); err != nil {
		api.logger.Warn("msg", "无效的用户ID", "error", err.Error())
		c.Error(err)
		return
	}

	err := api.userService.DeleteUserByIds(ctx, req.Ids)
	if err != nil {
		api.logger.Warn("msg", "批量删除用户失败", "error", err.Error())
		c.Error(err)
		return
	}

	response.SuccessNoData(c, "success")
}
