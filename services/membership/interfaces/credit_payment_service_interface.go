package interfaces

import (
	"context"

	"github.com/yb2020/odoc/services/membership/dto"
	"github.com/yb2020/odoc/services/membership/event"
	"github.com/yb2020/odoc/services/membership/model"
)

// CreditPaymentServiceInterface Cedit支付服务接口
type ICreditPaymentService interface {

	// GetById 根据ID获取支付订单
	GetById(ctx context.Context, id string) (*model.CreditPaymentRecord, error)

	// NewPaymentOrder 创建支付订单 返回订单ID
	NewPaymentOrder(ctx context.Context, membershipId string, userId string, payCredit *dto.NewCreditPayOrder) (string, error)

	// HandlePayEventHandler 处理支付回调事件
	HandlePayEventHandler(ctx context.Context, notifyEvent *event.CreditPayNotifyEvent) error

	// GetConfirmExpiredList 获取确认积分支付过期的列表
	GetConfirmExpiredList(ctx context.Context, size int) ([]model.CreditPaymentRecord, error)
}
