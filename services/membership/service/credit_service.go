package service

import (
	"context"

	"strconv"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/membership"
	"github.com/yb2020/odoc/internal/biz"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/membership/dao"
	"github.com/yb2020/odoc/services/membership/dto"
	"github.com/yb2020/odoc/services/membership/event"
	"github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/membership/model"
)

// CreditService 会员积分账户服务实现
type CreditService struct {
	logger    logging.Logger
	tracer    opentracing.Tracer
	creditDAO *dao.CreditDAO

	creditBillService    *CreditBillService
	creditPaymentService interfaces.ICreditPaymentService
}

func NewCreditService(logger logging.Logger, tracer opentracing.Tracer, creditDAO *dao.CreditDAO, creditBillService *CreditBillService, creditPaymentService interfaces.ICreditPaymentService) *CreditService {
	return &CreditService{
		logger:               logger,
		tracer:               tracer,
		creditDAO:            creditDAO,
		creditBillService:    creditBillService,
		creditPaymentService: creditPaymentService,
	}
}

// GetById 获取用户会员积分账户
func (s *CreditService) GetById(ctx context.Context, id string) (*model.Credit, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.GetById")
	defer span.Finish()

	membershipCredit, err := s.creditDAO.FindExistById(ctx, id)
	if err != nil {
		return nil, err
	}

	return membershipCredit, nil
}

// GetByUserId 获取用户会员积分账户
func (s *CreditService) GetByUserId(ctx context.Context, userId string) (*model.Credit, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.GetByUserId")
	defer span.Finish()

	membershipCredit, err := s.creditDAO.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return membershipCredit, nil
}

// GetByUserMembershipId 获取用户会员积分账户
func (s *CreditService) GetByMembershipId(ctx context.Context, membershipId string) (*model.Credit, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.GetByMembershipId")
	defer span.Finish()

	membershipCredit, err := s.creditDAO.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return nil, err
	}

	return membershipCredit, nil
}

// NewCreditAccount 创建用户会员积分账户
func (s *CreditService) NewCreditAccount(ctx context.Context, membershipId string, userId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.NewCreditAccount")
	defer span.Finish()

	membershipCreditExit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return "0", err
	}
	if membershipCreditExit != nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserCreditAccountAlreadyExists, "credit account already exists")
	}

	membershipCredit := &model.Credit{
		MembershipId: membershipId,
		UserId:       userId,
		Credit:       0,
		AddOnCredit:  0,
	}

	return membershipCredit.Id, s.creditDAO.Save(ctx, membershipCredit)
}

// DeleteCreditAccount 删除用户会员积分账户
func (s *CreditService) DeleteCreditAccount(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.DeleteCreditAccount")
	defer span.Finish()

	return s.creditDAO.DeleteById(ctx, id)
}

// InOrOutCredit 增加或减少用户会员积分账户
func (s *CreditService) InOrOutCredit(ctx context.Context, userId string, membershipId string, payCredit dto.CreditPayIntent) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.InOrOutCredit")
	defer span.Finish()

	s.logger.Info("msg", "InOrOutCredit do pay credit actions", "payCredit", payCredit)

	if payCredit.Type == pb.CreditPayType_CREDIT_PAY_TYPE_UNKNOWN {

		return "0", errors.BizWithStatus(biz.Membership_Status_CreditBillTypeUnknown, "credit bill type unknown")
	}
	if payCredit.Type == pb.CreditPayType_CREDIT_PAY_TYPE_EXPIRED {
		_, err := s.resetCreditZero(ctx, userId, membershipId, int32(pb.CreditPayType_CREDIT_PAY_TYPE_EXPIRED), payCredit.Content, payCredit.Remark) //积分过期
		if err != nil {
			return "0", err
		}
		_, err = s.resetAddOnCreditZero(ctx, userId, membershipId, int32(pb.CreditPayType_CREDIT_PAY_TYPE_EXPIRED), payCredit.Content, payCredit.Remark) //附加积分过期
		if err != nil {
			return "0", err
		}
		return "0", nil
	}
	if payCredit.Type == pb.CreditPayType_CREDIT_PAY_TYPE_SUB_FREE {
		s.resetCreditZero(ctx, userId, membershipId, int32(pb.CreditPayType_CREDIT_PAY_TYPE_EXPIRED), payCredit.Content, payCredit.Remark) //积分过期
		return s.inCredit(ctx, userId, membershipId, int32(payCredit.Type), payCredit.Credit, payCredit.Content, payCredit.Remark)         //订阅Free积分
	}
	if payCredit.Type == pb.CreditPayType_CREDIT_PAY_TYPE_SUB_PRO {
		s.resetCreditZero(ctx, userId, membershipId, int32(pb.CreditPayType_CREDIT_PAY_TYPE_EXPIRED), payCredit.Content, payCredit.Remark) //积分过期
		return s.inCredit(ctx, userId, membershipId, int32(payCredit.Type), payCredit.Credit, payCredit.Content, payCredit.Remark)         //订阅Pro积分
	}
	if payCredit.Type == pb.CreditPayType_CREDIT_PAY_TYPE_SUB_PRO_ADD_ON_CREDIT {
		return s.inAddOnCredit(ctx, userId, membershipId, int32(payCredit.Type), payCredit.AddOnCredit, payCredit.Content, payCredit.Remark) //订阅Pro附加积分
	}
	if payCredit.Type == pb.CreditPayType_CREDIT_PAY_TYPE_SERVICE_COST {
		if payCredit.CreditType == pb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT {
			return s.outAddOnCredit(ctx, userId, membershipId, int32(payCredit.Type), payCredit.AddOnCredit, payCredit.Content, payCredit.Remark) //服务消费附加积分
		}
		return s.outCredit(ctx, userId, membershipId, int32(payCredit.Type), payCredit.Credit, payCredit.Content, payCredit.Remark) //服务消费积分
	}
	if payCredit.Type == pb.CreditPayType_CREDIT_PAY_TYPE_SERVICE_RETRIEVE {
		if payCredit.CreditType == pb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT {
			return s.inAddOnCredit(ctx, userId, membershipId, int32(payCredit.Type), payCredit.AddOnCredit, payCredit.Content, payCredit.Remark) //服务回退附加积分
		}
		return s.inCredit(ctx, userId, membershipId, int32(payCredit.Type), payCredit.Credit, payCredit.Content, payCredit.Remark) //服务回退积分
	}

	s.logger.Error("msg", "InOrOutCredit do pay credit actions", "payCredit", payCredit)
	return "0", errors.BizWithStatus(biz.Membership_Status_CreditBillTypeUnknown, "credit bill type unknown")
}

// resetCreditZero 重置用户会员积分账户为0
func (s *CreditService) resetCreditZero(ctx context.Context, userId string, membershipId string, billType int32, content string, remark string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.resetCreditZero")
	defer span.Finish()

	membershipCredit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return "0", err
	}
	if membershipCredit == nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}

	//计算积分变化
	beforeCredit := membershipCredit.Credit
	afterCredit := int64(0)
	changeCredit := beforeCredit - afterCredit

	membershipCredit.Credit = afterCredit
	_ = s.creditDAO.Modify(ctx, membershipCredit)

	//记录积分流水
	billId, err := s.creditBillService.NewBill(ctx, &model.CreditBill{
		MembershipId:      membershipId,
		UserId:            userId,
		CreditId:          membershipCredit.Id,
		Type:              billType,
		CreditType:        int32(pb.CreditType_CREDIT_TYPE_CREDIT),
		InOutType:         int32(pb.CreditInOutType_CREDIT_IN_OUT_TYPE_EXPENSE),
		Credit:            changeCredit,
		BeforeCredit:      beforeCredit,
		AfterCredit:       afterCredit,
		AddOnCredit:       0,
		BeforeAddOnCredit: membershipCredit.AddOnCredit,
		AfterAddOnCredit:  membershipCredit.AddOnCredit,
		Content:           content,
		Remark:            remark,
	})

	return billId, err
}

// inCredit 增加用户会员积分账户
func (s *CreditService) inCredit(ctx context.Context, userId string, membershipId string, billType int32, credit int64, content string, remark string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.inCredit")
	defer span.Finish()

	membershipCredit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return "0", err
	}
	if membershipCredit == nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}

	//计算积分变化
	beforeCredit := membershipCredit.Credit
	afterCredit := beforeCredit + credit
	changeCredit := credit

	membershipCredit.Credit = afterCredit
	_ = s.creditDAO.Modify(ctx, membershipCredit)

	//记录积分流水
	billId, err := s.creditBillService.NewBill(ctx, &model.CreditBill{
		MembershipId:      membershipId,
		UserId:            userId,
		CreditId:          membershipCredit.Id,
		Type:              billType,
		CreditType:        int32(pb.CreditType_CREDIT_TYPE_CREDIT),
		InOutType:         int32(pb.CreditInOutType_CREDIT_IN_OUT_TYPE_INCOME),
		Credit:            changeCredit,
		BeforeCredit:      beforeCredit,
		AfterCredit:       afterCredit,
		AddOnCredit:       0,
		BeforeAddOnCredit: membershipCredit.AddOnCredit,
		AfterAddOnCredit:  membershipCredit.AddOnCredit,
		Content:           content,
		Remark:            remark,
	})

	return billId, err
}

// outCredit 减少用户会员积分账户
func (s *CreditService) outCredit(ctx context.Context, userId string, membershipId string, billType int32, credit int64, content string, remark string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.outCredit")
	defer span.Finish()

	membershipCredit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return "0", err
	}
	if membershipCredit == nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}

	if membershipCredit.Credit < credit {
		return "0", errors.BizWithStatus(biz.Membership_Status_CreditNotEnough, "credit not enough")
	}

	//计算积分变化
	beforeCredit := membershipCredit.Credit
	afterCredit := beforeCredit - credit
	changeCredit := credit

	membershipCredit.Credit = afterCredit
	_ = s.creditDAO.Modify(ctx, membershipCredit)

	//记录积分流水
	billId, err := s.creditBillService.NewBill(ctx, &model.CreditBill{
		MembershipId:      membershipId,
		UserId:            userId,
		CreditId:          membershipCredit.Id,
		Type:              billType,
		CreditType:        int32(pb.CreditType_CREDIT_TYPE_CREDIT),
		InOutType:         int32(pb.CreditInOutType_CREDIT_IN_OUT_TYPE_EXPENSE),
		Credit:            changeCredit,
		BeforeCredit:      beforeCredit,
		AfterCredit:       afterCredit,
		AddOnCredit:       0,
		BeforeAddOnCredit: membershipCredit.AddOnCredit,
		AfterAddOnCredit:  membershipCredit.AddOnCredit,
		Content:           content,
		Remark:            remark,
	})

	return billId, err
}

// resetAddOnCreditZero 重置用户会员附加积分账户为0
func (s *CreditService) resetAddOnCreditZero(ctx context.Context, userId string, membershipId string, billType int32, content string, remark string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.resetAddOnCreditZero")
	defer span.Finish()

	membershipCredit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return "0", err
	}
	if membershipCredit == nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}

	//计算积分变化
	beforeAddOnCredit := membershipCredit.AddOnCredit
	afterAddOnCredit := int64(0)
	changeAddOnCredit := beforeAddOnCredit - afterAddOnCredit

	membershipCredit.AddOnCredit = afterAddOnCredit
	_ = s.creditDAO.Modify(ctx, membershipCredit)

	//记录积分流水
	billId, err := s.creditBillService.NewBill(ctx, &model.CreditBill{
		MembershipId:      membershipId,
		UserId:            userId,
		CreditId:          membershipCredit.Id,
		Type:              billType,
		CreditType:        int32(pb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT),
		InOutType:         int32(pb.CreditInOutType_CREDIT_IN_OUT_TYPE_EXPENSE),
		Credit:            membershipCredit.Credit,
		BeforeCredit:      membershipCredit.Credit,
		AfterCredit:       membershipCredit.Credit,
		AddOnCredit:       changeAddOnCredit,
		BeforeAddOnCredit: beforeAddOnCredit,
		AfterAddOnCredit:  afterAddOnCredit,
		Content:           content,
		Remark:            remark,
	})

	return billId, err
}

// inAddOnCredit 增加用户会员积分账户
func (s *CreditService) inAddOnCredit(ctx context.Context, userId string, membershipId string, billType int32, credit int64, content string, remark string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.inAddOnCredit")
	defer span.Finish()

	membershipCredit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return "0", err
	}
	if membershipCredit == nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}

	//计算积分变化
	beforeCredit := membershipCredit.AddOnCredit
	afterCredit := beforeCredit + credit
	changeCredit := credit

	membershipCredit.AddOnCredit = afterCredit
	_ = s.creditDAO.Modify(ctx, membershipCredit)

	//记录积分流水
	billId, err := s.creditBillService.NewBill(ctx, &model.CreditBill{
		MembershipId:      membershipId,
		UserId:            userId,
		CreditId:          membershipCredit.Id,
		Type:              billType,
		CreditType:        int32(pb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT),
		InOutType:         int32(pb.CreditInOutType_CREDIT_IN_OUT_TYPE_INCOME),
		Credit:            0,
		BeforeCredit:      membershipCredit.Credit,
		AfterCredit:       membershipCredit.Credit,
		AddOnCredit:       changeCredit,
		BeforeAddOnCredit: beforeCredit,
		AfterAddOnCredit:  afterCredit,
		Content:           content,
		Remark:            remark,
	})

	return billId, err
}

// outAddOnCredit 减少用户会员积分账户
func (s *CreditService) outAddOnCredit(ctx context.Context, userId string, membershipId string, billType int32, credit int64, content string, remark string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.outAddOnCredit")
	defer span.Finish()

	membershipCredit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return "0", err
	}
	if membershipCredit == nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}

	if membershipCredit.AddOnCredit < credit {
		return "0", errors.BizWithStatus(biz.Membership_Status_CreditAddOnNotEnough, "credit add on not enough")
	}

	//计算积分变化
	beforeCredit := membershipCredit.AddOnCredit
	afterCredit := beforeCredit - credit
	changeCredit := credit

	membershipCredit.AddOnCredit = afterCredit
	_ = s.creditDAO.Modify(ctx, membershipCredit)

	//记录积分流水
	billId, err := s.creditBillService.NewBill(ctx, &model.CreditBill{
		MembershipId:      membershipId,
		UserId:            userId,
		CreditId:          membershipCredit.Id,
		Type:              billType,
		CreditType:        int32(pb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT),
		InOutType:         int32(pb.CreditInOutType_CREDIT_IN_OUT_TYPE_EXPENSE),
		Credit:            0,
		BeforeCredit:      membershipCredit.Credit,
		AfterCredit:       membershipCredit.Credit,
		AddOnCredit:       changeCredit,
		BeforeAddOnCredit: beforeCredit,
		AfterAddOnCredit:  afterCredit,
		Content:           content,
		Remark:            remark,
	})

	return billId, err
}

// CheckCreditEnough 检查用户会员积分账户是否足够
func (s *CreditService) CheckCreditEnough(ctx context.Context, userId string, membershipId string, credit int64) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.CheckCreditEnough")
	defer span.Finish()

	membershipCredit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return false, err
	}
	if membershipCredit == nil {
		return false, errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}
	if membershipCredit.Credit < credit {
		return false, errors.BizWithStatus(biz.Membership_Status_CreditNotEnough, "credit not enough")
	}
	return true, nil
}

// CheckAddOnCreditEnough 检查用户会员附加积分账户是否足够
func (s *CreditService) CheckAddOnCreditEnough(ctx context.Context, userId string, membershipId string, credit int64) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.CheckAddOnCreditEnough")
	defer span.Finish()

	membershipCredit, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return false, err
	}
	if membershipCredit == nil {
		return false, errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}
	if membershipCredit.AddOnCredit < credit {
		return false, errors.BizWithStatus(biz.Membership_Status_CreditAddOnNotEnough, "credit add on not enough")
	}
	return true, nil
}

// Pay 支付积分
func (s *CreditService) Pay(ctx context.Context, userId string, membershipId string, paymentRecordId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.Pay")
	defer span.Finish()

	record, err := s.creditPaymentService.GetById(ctx, paymentRecordId)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordNotFound, "payment record not found")
	}
	if record.Status != int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_PENDING) {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordStatusNotPending, "payment record status not pending")
	}

	var credit, addOnCredit int64
	if pb.CreditType(record.CreditType) == pb.CreditType_CREDIT_TYPE_CREDIT {
		credit = record.Credit
	} else if pb.CreditType(record.CreditType) == pb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT {
		addOnCredit = record.Credit
	}

	payCredit := dto.CreditPayIntent{
		Type:        pb.CreditPayType_CREDIT_PAY_TYPE_SERVICE_COST,
		CreditType:  pb.CreditType(record.CreditType),
		ServiceType: pb.CreditServiceType(record.ServiceType),
		Credit:      credit,
		AddOnCredit: addOnCredit,
		Content:     record.Content,
		Remark:      record.Remark,
	}

	// //4.会员账户支付支付订单
	billId, err := s.InOrOutCredit(ctx, userId, membershipId, payCredit)
	if err != nil {
		// 支付失败
		s.logger.Error("msg", "Pay failed", "userId", userId, "paymentRecordId", paymentRecordId)

		// Determine the error code for the event using a closure
		eventErrorCodeStr := func() string {
			if bizErr, ok := err.(*errors.BizError); ok {
				return strconv.Itoa(int(bizErr.Status))
			} else {
				return ""
			}
		}()

		updateErr := s.creditPaymentService.HandlePayEventHandler(ctx, &event.CreditPayNotifyEvent{
			RecordId:     paymentRecordId,
			EventType:    event.CreditPayNotifyEvent_PrePayFailed,
			UserId:       userId,
			ErrorCode:    eventErrorCodeStr,
			ErrorMessage: err.Error(),
		})
		if updateErr != nil {
			s.logger.Error("msg", "Update payment order status failed", "userId", userId, "paymentRecordId", paymentRecordId)
			return updateErr
		}
		return err
	}

	// 支付成功
	err = s.creditPaymentService.HandlePayEventHandler(ctx, &event.CreditPayNotifyEvent{
		EventType: event.CreditPayNotifyEvent_PrePaySuccess,
		RecordId:  paymentRecordId,
		UserId:    userId,
		BillId:    billId,
	})
	if err != nil {
		s.logger.Error("msg", "Update payment order status failed", "userId", userId, "paymentRecordId", paymentRecordId)
		return err
	}

	s.logger.Info("msg", "Pay success", "userId", userId, "paymentRecordId", paymentRecordId)

	return nil
}

// Confirm 确认支付积分
func (s *CreditService) Confirm(ctx context.Context, paymentRecordId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.Confirm")
	defer span.Finish()

	if paymentRecordId == "" {
		return nil
	}
	record, err := s.creditPaymentService.GetById(ctx, paymentRecordId)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordNotFound, "payment record not found")
	}
	if record.Status == int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_SUCCESS) ||
		record.Status == int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_FAILED) ||
		record.Status == int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_CANCELLED) {

		s.logger.Info("msg", "ignore and return", "paymentRecordId", paymentRecordId)
		return nil
	}
	if record.Status != int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_AWAITING_CONFIRMATION) {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordStatusNotAwaitingConfirmation, "payment record status not pending")
	}

	// 确认支付订单通知
	err = s.creditPaymentService.HandlePayEventHandler(ctx, &event.CreditPayNotifyEvent{
		EventType: event.CreditPayNotifyEvent_PayConfirm,
		RecordId:  paymentRecordId,
		UserId:    record.UserId,
		BillId:    record.RelCreditBid,
	})
	if err != nil {
		return err
	}
	return nil
}

// Retrieve 回滚积分
func (s *CreditService) Retrieve(ctx context.Context, paymentRecordId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.Retrieve")
	defer span.Finish()

	if paymentRecordId == "" {
		return nil
	}
	record, err := s.creditPaymentService.GetById(ctx, paymentRecordId)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordNotFound, "payment record not found")
	}
	if record.Status == int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_SUCCESS) ||
		record.Status == int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_FAILED) ||
		record.Status == int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_CANCELLED) {
		s.logger.Info("msg", "ignore and return", "paymentRecordId", paymentRecordId)
		return nil
	}
	if record.Status != int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_AWAITING_CONFIRMATION) {
		return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordStatusNotAwaitingConfirmation, "payment record status not pending")
	}

	// creditAccount, err := s.GetByUserId(ctx, record.UserId)
	creditAccount, err := s.GetByMembershipId(ctx, record.MembershipId)
	if err != nil {
		return err
	}
	if creditAccount == nil {
		return errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}

	// 1.回退支付积分
	var credit, addOnCredit int64
	if pb.CreditType(record.CreditType) == pb.CreditType_CREDIT_TYPE_CREDIT {
		credit = record.Credit
	} else if pb.CreditType(record.CreditType) == pb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT {
		addOnCredit = record.Credit
	}
	creditPayIntent := dto.CreditPayIntent{
		Type:        pb.CreditPayType_CREDIT_PAY_TYPE_SERVICE_RETRIEVE,
		CreditType:  pb.CreditType(record.CreditType),
		ServiceType: pb.CreditServiceType(record.ServiceType),
		Credit:      credit,
		AddOnCredit: addOnCredit,
		Content:     record.Content,
		Remark:      record.Remark,
	}
	billId, err := s.InOrOutCredit(ctx, record.UserId, creditAccount.MembershipId, creditPayIntent)
	if err != nil {
		return err
	}

	// 2.回滚支付订单通知
	err = s.creditPaymentService.HandlePayEventHandler(ctx, &event.CreditPayNotifyEvent{
		EventType: event.CreditPayNotifyEvent_PayRetrieve,
		RecordId:  paymentRecordId,
		UserId:    record.UserId,
		BillId:    billId,
	})
	if err != nil {
		return err
	}
	return nil
}

// CreditAccountExpiredAllCredit 会员账户积分和附加积分过期清零
func (s *CreditService) CreditAccountExpiredAllCredit(ctx context.Context, membershipId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipCreditService.CreditAccountExpiredAllCredit")
	defer span.Finish()

	creditAccount, err := s.GetByMembershipId(ctx, membershipId)
	if err != nil {
		return err
	}
	if creditAccount == nil {
		return errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "credit account not found")
	}

	payCredit := dto.CreditPayIntent{
		Type:       pb.CreditPayType_CREDIT_PAY_TYPE_EXPIRED,
		CreditType: pb.CreditType_CREDIT_TYPE_CREDIT,
		// ServiceType: pb.CreditServiceType_CREDIT_SERVICE_TYPE_MEMBERSHIP_EXPIRED,
		Credit:      creditAccount.Credit,
		AddOnCredit: creditAccount.AddOnCredit,
		Content:     "会员账户积分和附加积分过期清零",
		Remark:      "会员账户积分和附加积分过期清零",
	}

	// 3.会员账户积分和附加积分过期清零
	_, err = s.InOrOutCredit(ctx, creditAccount.UserId, creditAccount.MembershipId, payCredit)
	if err != nil {
		return err
	}
	return nil
}
