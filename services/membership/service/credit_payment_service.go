package service

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/internal/biz"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	pb "github.com/yb2020/odoc/proto/gen/go/membership"
	"github.com/yb2020/odoc/services/membership/dao"
	"github.com/yb2020/odoc/services/membership/dto"
	"github.com/yb2020/odoc/services/membership/event"
	"github.com/yb2020/odoc/services/membership/model"
)

// CreditPaymentService Cedit支付服务实现
type CreditPaymentService struct {
	logger                 logging.Logger
	tracer                 opentracing.Tracer
	creditPaymentRecordDAO *dao.CreditPaymentRecordDAO
}

func NewCreditPaymentService(logger logging.Logger, tracer opentracing.Tracer, creditPaymentRecordDAO *dao.CreditPaymentRecordDAO) *CreditPaymentService {
	return &CreditPaymentService{
		logger:                 logger,
		tracer:                 tracer,
		creditPaymentRecordDAO: creditPaymentRecordDAO,
	}
}

// GetById 根据ID获取支付订单
func (s *CreditPaymentService) GetById(ctx context.Context, id string) (*model.CreditPaymentRecord, error) {
	return s.creditPaymentRecordDAO.FindExistById(ctx, id)
}

// NewPaymentOrder 创建支付订单
func (s *CreditPaymentService) NewPaymentOrder(ctx context.Context, membershipId string, userId string, payCredit *dto.NewCreditPayOrder) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreditPaymentService.NewPaymentOrder")
	defer span.Finish()

	record := &model.CreditPaymentRecord{
		MembershipId: membershipId,
		UserId:       userId,
		ServiceType:  int32(payCredit.ServiceType),
		CreditType:   int32(payCredit.CreditType),
		Credit:       payCredit.Credit,
		Status:       int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_PENDING),
		Content:      payCredit.Content,
		Remark:       payCredit.Remark,
	}
	record.Id = idgen.GenerateUUID()

	err := s.creditPaymentRecordDAO.Save(ctx, record)
	if err != nil {
		return "0", err
	}

	return record.Id, nil
}

// HandlePayEventHandler 处理支付回调事件
func (s *CreditPaymentService) HandlePayEventHandler(ctx context.Context, notifyEvent *event.CreditPayNotifyEvent) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreditPaymentService.HandlePayEventHandler")
	defer span.Finish()

	switch notifyEvent.EventType {
	case event.CreditPayNotifyEvent_PrePaySuccess:

		return s.handlePrePaySuccess(ctx, notifyEvent)
	case event.CreditPayNotifyEvent_PrePayFailed:
		return s.handlePrePayFailed(ctx, notifyEvent)
	case event.CreditPayNotifyEvent_PayConfirm:
		return s.handlePayConfirm(ctx, notifyEvent)
	case event.CreditPayNotifyEvent_PayRetrieve:
		return s.handlePayRetrieve(ctx, notifyEvent)
	default:
		return errors.Biz("unsupported event type")
	}
}

// GetConfirmExpiredList 获取确认积分支付过期的列表
func (s *CreditPaymentService) GetConfirmExpiredList(ctx context.Context, size int) ([]model.CreditPaymentRecord, error) {
	return s.creditPaymentRecordDAO.GetConfirmExpiredList(ctx, size)
}

// =========private method=========//
// 预扣成功事件处理
func (s *CreditPaymentService) handlePrePaySuccess(ctx context.Context, notifyEvent *event.CreditPayNotifyEvent) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreditPaymentService.handlePrePaySuccess")
	defer span.Finish()

	record, err := s.GetById(ctx, notifyEvent.RecordId)
	if err != nil {
		return err
	}
	record.Status = int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_AWAITING_CONFIRMATION)
	now := time.Now()
	record.TransactionAt = &now
	record.PayAt = &now
	record.RelCreditBid = notifyEvent.BillId
	confirmExpiredAt := now.Add(time.Hour)
	record.ConfirmExpiredAt = &confirmExpiredAt
	return s.creditPaymentRecordDAO.Modify(ctx, record)
}

// 预扣失败事件处理
func (s *CreditPaymentService) handlePrePayFailed(ctx context.Context, notifyEvent *event.CreditPayNotifyEvent) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreditPaymentService.handlePrePayFailed")
	defer span.Finish()

	record, err := s.GetById(ctx, notifyEvent.RecordId)
	if err != nil {
		return err
	}
	record.Status = int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_FAILED)
	now := time.Now()
	record.TransactionAt = &now
	record.PayAt = &now
	record.ErrorCode = notifyEvent.ErrorCode
	record.ErrorMessage = notifyEvent.ErrorMessage
	return s.creditPaymentRecordDAO.Modify(ctx, record)
}

// ConfirmPaymentOrder 确认支付订单
func (s *CreditPaymentService) handlePayConfirm(ctx context.Context, notifyEvent *event.CreditPayNotifyEvent) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreditPaymentService.ConfirmPaymentOrder")
	defer span.Finish()

	if notifyEvent.RecordId == "0" {
		// 为空，不需要确认返回成功
		return nil
	}

	record, err := s.GetById(ctx, notifyEvent.RecordId)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordNotFound, "payment record not exists")
	}
	if record.Status != int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_AWAITING_CONFIRMATION) {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordStatusNotAwaitingConfirmation, "payment record status not pending")
	}
	record.Status = int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_SUCCESS)
	now := time.Now()
	record.ConfirmAt = &now
	return s.creditPaymentRecordDAO.Modify(ctx, record)
}

// RetrievePaymentOrder 回滚支付订单
func (s *CreditPaymentService) handlePayRetrieve(ctx context.Context, notifyEvent *event.CreditPayNotifyEvent) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreditPaymentService.RetrievePaymentOrder")
	defer span.Finish()

	if notifyEvent.RecordId == "" {
		// 为空，不需要确认返回成功
		return nil
	}

	record, err := s.GetById(ctx, notifyEvent.RecordId)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordNotFound, "payment record not exists")
	}
	if record.Status != int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_AWAITING_CONFIRMATION) {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordStatusNotAwaitingConfirmation, "payment record status not pending")
	}
	record.Status = int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_CANCELLED)
	now := time.Now()
	record.ConfirmAt = &now
	record.RetrieveCreditBid = notifyEvent.BillId
	return s.creditPaymentRecordDAO.Modify(ctx, record)
}
