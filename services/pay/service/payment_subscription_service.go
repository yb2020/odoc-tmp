package service

import (
	"context"
	"errors"

	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pay/dao"
	"github.com/yb2020/odoc/services/pay/model"
)

// PaymentSubscriptionService 支付订阅服务，协调支付订阅DAO和支付提供商
type PaymentSubscriptionService struct {
	paymentSubscriptionDAO *dao.PaymentSubscriptionDAO
	logger                 logging.Logger
}

func (p *PaymentSubscriptionService) UpdateSubscription(ctx context.Context, paySub *model.PaymentSubscription) error {
	return p.paymentSubscriptionDAO.ModifyExcludeNull(ctx, paySub)
}

// NewPaymentSubscriptionService 创建一个新的支付订阅服务
func NewPaymentSubscriptionService(
	paymentSubscriptionDAO *dao.PaymentSubscriptionDAO,
	logger logging.Logger,
) *PaymentSubscriptionService {
	return &PaymentSubscriptionService{
		paymentSubscriptionDAO: paymentSubscriptionDAO,
		logger:                 logger,
	}
}

func (p *PaymentSubscriptionService) NewSubscription(ctx context.Context, paySub *model.PaymentSubscription) (string, error) {
	if paySub == nil {
		return "0", errors.New("payment subscription is nil")
	}
	id := idgen.GenerateUUID()
	paySub.Id = id
	return id, p.paymentSubscriptionDAO.Save(ctx, paySub)
}

// GetByProviderSubscriptionId 根据支付提供商订阅ID获取支付记录
func (p *PaymentSubscriptionService) GetByProviderSubscriptionId(ctx context.Context, providerSubscriptionId string) (*model.PaymentSubscription, error) {
	return p.paymentSubscriptionDAO.GetByProviderSubscriptionId(ctx, providerSubscriptionId)
}
