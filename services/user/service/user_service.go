package service

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/yb2020/odoc/internal/biz"
	"github.com/yb2020/odoc/pkg/cache"
	userContext "github.com/yb2020/odoc/pkg/context"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/pkg/i18n"
	idgen "github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	baseModel "github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/pkg/paginate"
	"github.com/yb2020/odoc/pkg/utils"
	pb "github.com/yb2020/odoc/proto/gen/go/user"
	"github.com/yb2020/odoc/services/user/bean"
	"github.com/yb2020/odoc/services/user/dao"
	"github.com/yb2020/odoc/services/user/event"
	"github.com/yb2020/odoc/services/user/model"
)

// UserService 默认的用户服务实现
type UserService struct {
	cache              cache.Cache
	userDAO            *dao.UserDAO
	logger             logging.Logger
	tracer             opentracing.Tracer
	localizer          i18n.Localizer
	eventBus           *eventbus.EventBus
	transactionManager *baseDao.TransactionManager
}

// NewUserService 创建一个新的用户服务
func NewUserService(cache cache.Cache, logger logging.Logger,
	tracer opentracing.Tracer, localizer i18n.Localizer, userDAO *dao.UserDAO,
	eventBus *eventbus.EventBus, transactionManager *baseDao.TransactionManager) *UserService {
	return &UserService{
		cache:              cache,
		localizer:          localizer,
		userDAO:            userDAO,
		logger:             logger,
		tracer:             tracer,
		eventBus:           eventBus,
		transactionManager: transactionManager,
	}
}

var cacheKey = "user:%d"

// GetUserByEmail 根据电子邮件获取用户
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.GetUserByEmail")
	span.SetTag("user.email", email)
	defer span.Finish()

	s.logger.Info("msg", "根据电子邮件获取用户", "email", email)
	user, err := s.userDAO.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.BizWrap(err.Error(), err)
	}
	if user == nil {
		return nil, errors.Biz("user.error.user_email_not_found")
	}
	return user, nil
}

// FindOrCreateByGoogleID handles user lookup, creation, and binding using Google ID.
// It's designed to be called within the Google OAuth callback flow.
func (s *UserService) FindOrCreateByGoogleID(ctx context.Context, googleUserInfo *model.GoogleUserInfo) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.FindOrCreateByGoogleID")
	defer span.Finish()

	var user *model.User
	var err error

	// 步骤 1: 开启事务
	txErr := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		// 步骤 2: 用 GoogleOpenId 查��用户
		s.logger.Info("msg", "通过GoogleID查找用户", "googleId", googleUserInfo.ID)
		foundUser, findErr := s.userDAO.FindByGoogleID(txCtx, googleUserInfo.ID)
		if findErr != nil {
			s.logger.Error("msg", "通过GoogleID查找用户失败", "error", findErr)
			return findErr // 返回错误以回滚事务
		}

		// 步骤 3: 如果找到，直接赋值并返回
		if foundUser != nil {
			s.logger.Info("msg", "通过GoogleID找到现有用户", "userId", foundUser.Id)
			user = foundUser
			return nil // 事务成功
		}

		// 步骤 4: 如果没找到，再用 Email 查找用户
		s.logger.Info("msg", "通过GoogleID未找到用户，尝试通过Email查找", "email", googleUserInfo.Email)
		foundUser, findErr = s.userDAO.GetUserByEmail(txCtx, googleUserInfo.Email)
		if findErr != nil {
			s.logger.Error("msg", "通过Email查找用户失败", "error", findErr)
			return findErr
		}

		// 步骤 5: 如果用 Email 找到了
		if foundUser != nil {
			s.logger.Info("msg", "通过Email找到现有用户，准备绑定GoogleID", "userId", foundUser.Id)
			// 更新该用户的 GoogleOpenId 字段，完成绑定
			foundUser.GoogleOpenId = utils.ToValidUTF8(googleUserInfo.ID, "")
			if updateErr := s.userDAO.ModifyExcludeNull(txCtx, foundUser); updateErr != nil {
				s.logger.Error("msg", "绑定GoogleID失败", "error", updateErr)
				return updateErr
			}
			s.logger.Info("msg", "绑定GoogleID成功", "userId", foundUser.Id)
			user = foundUser
			return nil // 事务成功
		}

		// 步骤 6: 如果 email 也找不到，则创建新用户
		s.logger.Info("msg", "用户不存在，创建新用户", "email", googleUserInfo.Email)
		newUser := &model.User{
			BaseModel: baseModel.BaseModel{
				Id: idgen.GenerateUUID(),
			},
			Email:        utils.ToValidUTF8(googleUserInfo.Email, ""),
			Username:     utils.ToValidUTF8(googleUserInfo.Email, ""), // 使用 Email 作为默认用户名
			Nickname:     utils.ToValidUTF8(googleUserInfo.Name, ""),
			Avatar:       utils.ToValidUTF8(googleUserInfo.Picture, ""),
			GoogleOpenId: utils.ToValidUTF8(googleUserInfo.ID, ""),
			Status:       pb.UserStatus_STATUS_ACTIVE,                // 直接设为激活状态
			Roles:        model.UserRoleSlice{pb.UserRole_ROLE_USER}, // 赋予默认角色
			// 注意：密码字段为空，因为这是第三方登录，不需要密码
		}

		if createErr := s.userDAO.Save(txCtx, newUser); createErr != nil {
			s.logger.Error("msg", "创建新用户失败", "error", createErr)
			return createErr
		}

		s.logger.Info("msg", "创建新用户成功", "userId", newUser.Id)
		user = newUser
		return nil // 事务成功
	})

	if txErr != nil {
		// 处理事务本身或事务内部返回的错误
		if txErr.Error() == "user.error.user_not_found_for_creation" {
			// This case should now be handled by the creation logic above.
			// If we still get this, it's an unexpected state.
			return nil, errors.Biz("user.error.user_creation_failed_unexpectedly")
		}
		return nil, txErr
	}

	return user, err
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.UpdateUser")
	span.SetTag("user.id", user.Id)
	defer span.Finish()

	s.logger.Info("msg", "更新用户", "id", user.Id)

	// 检查用户是否存在
	existingUser, err := s.userDAO.FindExistById(ctx, user.Id)
	if err != nil {
		return nil, errors.BizWrap(err.Error(), err)
	}
	if existingUser == nil {
		return nil, errors.Biz("user.error.user_not_found")
	}

	if err := s.userDAO.ModifyExcludeNull(ctx, user); err != nil {
		return nil, errors.BizWrap(err.Error(), err)
	}

	// 清除缓存
	userCacheKey := fmt.Sprintf(cacheKey, user.Id)
	if err := s.cache.Delete(ctx, userCacheKey); err != nil {
		s.logger.Error("msg", "清除用户缓存失败", "error", err)
		// 缓存错误不应该影响正常流程
	} else {
		s.logger.Info("msg", "清除用户缓存成功", "id", user.Id)
	}
	return user, nil
}

// GetProfile 获取用户个人资料
func (s *UserService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.GetProfile")
	defer span.Finish()

	s.logger.Info("msg", "获取用户个人资料")

	// 从上下文中获取用户ID
	userID, ok := userContext.GetUserID(ctx)
	if !ok {
		return nil, errors.Biz("user.error.user_not_found")
	}

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.BizWrap(err.Error(), err)
	}
	if user == nil {
		return nil, errors.Biz("user.error.user_not_found")
	}

	return &pb.GetProfileResponse{
		User: user.ToProto(),
	}, nil
}

// CheckEmailExists 检查邮箱是否存在
func (s *UserService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.CheckEmailExists")
	defer span.Finish()
	s.logger.Info("msg", "检查邮箱是否存在", "email", email)
	_, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == "user.error.user_email_not_found" {
			return false, nil
		}
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.Register")
	defer span.Finish()

	s.logger.Info("msg", "用户注册", "email", req.Email)

	id := idgen.GenerateUUID()

	user := &model.User{
		BaseModel: baseModel.BaseModel{
			Id: id,
		},
		Email:    req.Email,
		Password: utils.StrengthenPassword(req.Password, fmt.Sprintf("%d", id)),
		Roles:    []pb.UserRole{pb.UserRole_ROLE_USER},
		Status:   pb.UserStatus_STATUS_INACTIVE,
	}

	// 检查邮箱是否已存在
	if _, err := s.GetUserByEmail(ctx, req.Email); err == nil {
		return nil, errors.Biz("user.error.email_exists")
	}

	if err := s.userDAO.Save(ctx, user); err != nil {
		return nil, errors.System(errors.ErrorTypeDatabase, err.Error(), err)
	}

	return &pb.RegisterResponse{
		User: user.ToProto(),
	}, nil
}

// UpdateProfile 更新用户个人资料
func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.UpdateProfile")
	defer span.Finish()

	s.logger.Info("msg", "更新用户个人资料")

	// 从请求中获取用户信息
	pbUser := req.GetUser()
	if pbUser == nil {
		return nil, errors.Biz("user.error.user_not_found")
	}

	// 从上下文中获取用户ID
	userID := ctx.Value(userContext.UserIDKey).(string)
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.BizWrap(err.Error(), err)
	}
	if user == nil {
		return nil, errors.Biz("user.error.user_not_found")
	}

	// 更新用户信息
	user.Username = pbUser.Username
	user.Email = pbUser.Email
	user.Nickname = pbUser.Nickname
	user.Avatar = pbUser.Avatar
	// 因为使用了，ModifyExcludeNull，设置为空，即可不更新密码字段
	user.Password = ""

	updatedUser, err := s.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProfileResponse{
		User: updatedUser.ToProto(),
	}, nil
}

// GetByIds 根据多个ID获取用户
func (s *UserService) GetByIds(ctx context.Context, ids []string) ([]*model.User, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.GetByIds")
	defer span.Finish()

	s.logger.Info("msg", "根据多个ID获取用户", "ids", ids)
	users := make([]*model.User, 0)
	for _, id := range ids {
		user, err := s.GetUserByID(ctx, id)
		if err != nil {
			continue
		}
		if user == nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// PaginationUsers 分页获取用户列表
func (s *UserService) PaginationUsers(ctx context.Context, req *pb.PaginationUsersRequest) (*pb.PaginationUsersResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.PaginationUsers")
	defer span.Finish()

	s.logger.Info("msg", "分页获取用户列表", "page", req.Page, "size", req.Size)

	// 使用链式调用创建并配置分页查询选项
	options := paginate.NewOptions().
		AddLike("username", req.Username).
		AddLike("email", req.Email).
		AddLike("nickname", req.Nickname).
		AddTimeRange("created_at", req.FromCreateTime, req.ToCreateTime).
		AddOrder("id", "DESC")

	// 使用数据库分页查询
	users, total, err := s.userDAO.Paginate(ctx, req.Page, req.Size, options)
	if err != nil {
		return nil, err
	}

	// 转换为 protobuf 对象
	pagedUsers := make([]*pb.User, 0, len(users))
	for _, user := range users {
		pagedUsers = append(pagedUsers, user.ToProto())
	}

	return &pb.PaginationUsersResponse{
		Users: pagedUsers,
		Total: int32(total),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// CreateUserAdmin 创建用户（管理员）
func (s *UserService) CreateUserAdmin(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.CreateUserAdmin")
	defer span.Finish()

	s.logger.Info("msg", "创建用户（管理员）")

	pbUser := req.GetUser()
	if pbUser == nil {
		return nil, errors.Biz("user.error.user_not_found")
	}

	user := &model.User{}
	user.FromProto(pbUser)

	// 检查邮箱是否已存在
	if _, err := s.GetUserByEmail(ctx, user.Email); err == nil {
		return &pb.CreateUserResponse{
			User: nil,
		}, nil
	}

	if err := s.userDAO.Save(ctx, user); err != nil {
		return &pb.CreateUserResponse{
			User: nil,
		}, nil
	}

	// 发布用户注册事件
	s.eventBus.Publish(ctx, eventbus.Event{Type: event.UserRegisterEvent, Data: user.Id}, true)

	return &pb.CreateUserResponse{
		User: user.ToProto(),
	}, nil
}

// UpdateUserAdmin 更新用户（管理员）
func (s *UserService) UpdateUserAdmin(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.UpdateUserAdmin")
	defer span.Finish()

	s.logger.Info("msg", "更新用户（管理员）")

	pbUser := req.GetUser()
	if pbUser == nil {
		return nil, errors.Biz("user.error.user_not_found")
	}

	user, err := s.GetUserByID(ctx, pbUser.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Biz("user.error.user_not_found")
	}

	// 更新用户信息
	user.FromProto(pbUser)

	updatedUser, err := s.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		User: updatedUser.ToProto(),
	}, nil
}

// DeleteUserById 删除用户
func (s *UserService) DeleteUserById(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.DeleteUserById")
	defer span.Finish()

	// 所有的用户都下线
	s.eventBus.Publish(ctx, eventbus.Event{Type: event.UserDeletedEvent, Data: id}, true)

	s.logger.Info("msg", "删除用户", "id", id)
	if err := s.userDAO.DeleteById(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	if err := s.cache.Delete(ctx, fmt.Sprintf(cacheKey, id)); err != nil {
		s.logger.Error("msg", "清除用户缓存失败", "error", err)
		// 缓存错误不应该影响正常流程
	} else {
		s.logger.Info("msg", "清除用户缓存成功", "id", id)
	}

	return nil
}

// DeleteUserByIds 批量删除用户
func (s *UserService) DeleteUserByIds(ctx context.Context, ids []string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.DeleteUserByIds")
	defer span.Finish()

	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		s.logger.Info("msg", "通过事务批量删除用户", "ids", ids)
		for _, id := range ids {
			if err := s.DeleteUserById(txCtx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetUserByID 根据 ID 获取用户（带缓存）
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.GetUserByID")
	defer span.Finish()

	cacheKey := fmt.Sprintf(cacheKey, id)

	// 尝试从缓存获取用户
	var user model.User
	found, err := s.cache.Get(ctx, cacheKey, &user)
	if err != nil {
		s.logger.Error("msg", "从缓存获取用户失败", "error", err)
		// 缓存错误不应该影响正常流程，继续从数据库获取
	} else if found {
		s.logger.Info("msg", "从缓存获取用户成功", "id", id)
		return &user, nil
	}

	// 从数据库获取用户
	user2, err := s.userDAO.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	// 将用户存入缓存
	if user2 != nil {
		if err := s.cache.Set(ctx, cacheKey, user2, 30*time.Minute); err != nil {
			s.logger.Error("msg", "缓存用户失败", "error", err)
			// 缓存错误不应该影响正常流程
		} else {
			s.logger.Info("msg", "缓存用户成功", "id", id)
		}
	}

	return user2, nil
}

// 检查用户状态是否为激活
func (s *UserService) CheckUserStatusActive(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.CheckUserStatusActive")
	defer span.Finish()

	user, err := s.GetUserByID(ctx, id)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.BizWithStatus(biz.User_StatusNotFound, "user not found")
	}
	if user.Status != pb.UserStatus_STATUS_ACTIVE {
		return false, errors.BizWithStatus(biz.User_StatusNotActive, "user not active")
	}
	return true, nil
}

func (s *UserService) GetAuthInfo(ctx context.Context, id string) (*bean.AuthorBean, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserService.getAuthInfo")
	defer span.Finish()

	// MOCK:TODO
	// 生成模拟数据
	authorInfo := &bean.AuthorBean{
		Id:               id,
		NickName:         "学术达人",
		Description:      "专注人工智能和自然语言处理领域的研究者",
		Usn:              "scholar10086",
		UsnCanModify:     false,
		Mobile:           "13812345678",
		Email:            "scholar@example.com",
		IsWxPublicBind:   true,
		Self:             true,
		ShowName:         "张教授",
		AvatarUrl:        "https://example.com/avatars/10086.jpg",
		Tags:             "AI,机器学习,自然语言处理",
		AuthorId:         "110086",
		AuthorName:       "张三",
		IsAuthentication: true,
		IsCert:           true,
		IsPaperAuthor:    true,
		Profession:       "教授",
		ResearchField:    "人工智能,自然语言处理",
		SchoolCompany:    "清华大学",
		BanInfo:          nil,
	}

	return authorInfo, nil
}
