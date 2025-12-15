package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/membership/dao"
	"github.com/yb2020/odoc/services/membership/model"
)

// CreditBillService 会员积分账单服务实现
type CreditBillService struct {
	logger        logging.Logger
	tracer        opentracing.Tracer
	creditBillDAO *dao.CreditBillDAO
}

func NewCreditBillService(logger logging.Logger, tracer opentracing.Tracer, creditBillDAO *dao.CreditBillDAO) *CreditBillService {
	return &CreditBillService{
		logger:        logger,
		tracer:        tracer,
		creditBillDAO: creditBillDAO,
	}
}

// NewBill 创建积分账单
func (s *CreditBillService) NewBill(ctx context.Context, bill *model.CreditBill) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditBillService.NewBill")
	defer span.Finish()

	if bill == nil {
		return "0", errors.Biz("bill obj is nil")
	}
	bill.Id = idgen.GenerateUUID()

	return bill.Id, s.creditBillDAO.Save(ctx, bill)
}
