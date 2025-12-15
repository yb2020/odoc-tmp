package interfaces

import (
	"context"
	"time"

	pb "github.com/yb2020/odoc-proto/gen/go/membership"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/services/membership/model"
)

// IUserMembershipService 用户会员服务接口
type IUserMembershipService interface {
	// GetById 根据ID获取用户会员
	GetById(ctx context.Context, id string) (*model.UserMembership, error)

	// GetByUserId 根据用户ID获取用户会员
	GetByUserId(ctx context.Context, userId string) (*model.UserMembership, error)

	// GetBaseInfo 获取用户会员基础信息
	GetBaseInfo(ctx context.Context, userId string) (*pb.UserMembershipBaseInfo, error)

	// GetInfo 获取用户会员信息
	GetInfo(ctx context.Context, userId string) (*pb.UserMembershipInfo, error)

	// GetUserPermission 获取用户会员权限
	GetUserPermission(ctx context.Context, userId string) (*pb.UserPermission, error)

	// NewAccount 创建用户会员账号
	NewAccount(ctx context.Context, userId string) (string, error)

	// DeleteAccount 删除用户会员账号
	DeleteAccount(ctx context.Context, userId string) error

	// UpdateAccountType 更新用户会员类型
	// 参数说明: userId: 用户ID, memberType: 会员类型, stripeSubscriptionId: Stripe订阅ID, startDate: 会员开始时间, endDate: 会员结束时间
	UpdateAccountType(ctx context.Context, userId string, memberType int32, stripeSubscriptionId string, startDate time.Time, endDate time.Time) error

	// GetUserConfig 获取用户会员配置信息
	GetUserConfig(ctx context.Context) (config.MembershipTypeConfig, error)

	// GetUserConfigByUserId 获取用户会员配置信息
	GetUserConfigByUserId(ctx context.Context, userId string) (config.MembershipTypeConfig, error)

	// CheckExpired 检查用户会员是否过期 返回是否过期，会员类型，错误
	CheckExpired(ctx context.Context, userId string) (bool, int32, error)

	// HasPermUseAddOnCredit 检查用户会员是否可以使用附加积分
	HasPermUseAddOnCredit(ctx context.Context, userId string) (bool, error)

	// GetExpiredList 获取过期的会员列表
	GetExpiredList(ctx context.Context, size int) ([]model.UserMembership, error)

	// GetFreeExpiredList 获取过期的Free会员列表
	GetFreeExpiredList(ctx context.Context, size int) ([]model.UserMembership, error)

	// GetProExpiredList 获取过期的Pro会员列表
	GetProExpiredList(ctx context.Context, size int) ([]model.UserMembership, error)

	// ExpiredAccountToFree 过期会员降级为Free会员
	ExpiredAccountToFree(ctx context.Context, id string) error
}
