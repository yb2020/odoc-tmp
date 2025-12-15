package interfaces

import (
	"context"

	"github.com/yb2020/odoc/services/membership/dto"
	"github.com/yb2020/odoc/services/membership/model"
)

type ICreditService interface {

	// GetById 获取用户会员积分账户
	GetById(ctx context.Context, id string) (*model.Credit, error)

	// GetByUserId 获取用户会员积分账户
	GetByUserId(ctx context.Context, userId string) (*model.Credit, error)

	// GetByMembershipId 获取用户会员积分账户
	GetByMembershipId(ctx context.Context, membershipId string) (*model.Credit, error)

	// NewCreditAccount 创建用户会员积分账户
	NewCreditAccount(ctx context.Context, membershipId string, userId string) (string, error)

	// DeleteCreditAccount 删除用户会员积分账户
	DeleteCreditAccount(ctx context.Context, id string) error

	// InOrOutCredit 增加或减少用户会员积分账户
	InOrOutCredit(ctx context.Context, userId string, membershipId string, payCredit dto.CreditPayIntent) (string, error)

	// CheckCreditEnough 检查用户会员积分账户是否足够
	CheckCreditEnough(ctx context.Context, userId string, membershipId string, credit int64) (bool, error)

	// CheckAddOnCreditEnough 检查用户会员附加积分账户是否足够
	CheckAddOnCreditEnough(ctx context.Context, userId string, membershipId string, credit int64) (bool, error)

	// Pay 支付积分
	Pay(ctx context.Context, userId string, membershipId string, paymentRecordId string) error

	// Confirm 确认支付积分
	Confirm(ctx context.Context, paymentRecordId string) error

	// Retrieve 回滚支付积分
	Retrieve(ctx context.Context, paymentRecordId string) error

	// CreditAccountExpiredAllCredit 会员账户积分和附加积分全部过期清零
	CreditAccountExpiredAllCredit(ctx context.Context, membershipId string) error
}
