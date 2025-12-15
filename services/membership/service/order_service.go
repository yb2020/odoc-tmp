package service

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	msPb "github.com/yb2020/odoc-proto/gen/go/membership"
	pb "github.com/yb2020/odoc-proto/gen/go/order"
	"github.com/yb2020/odoc/internal/biz"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"

	"github.com/yb2020/odoc/services/membership/dao"
	"github.com/yb2020/odoc/services/membership/dto"
	"github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/membership/model"
	payevent "github.com/yb2020/odoc/services/pay/event"
)

// OrderService 会员订单服务实现
type OrderService struct {
	logger   logging.Logger
	tracer   opentracing.Tracer
	orderDAO *dao.OrderDAO

	msConfigService       *ConfigService
	creditService         interfaces.ICreditService
	userMembershipService interfaces.IUserMembershipService
}

func NewOrderService(logger logging.Logger, tracer opentracing.Tracer, orderDAO *dao.OrderDAO,
	membershipConfigService *ConfigService, creditService interfaces.ICreditService, userMembershipService interfaces.IUserMembershipService) *OrderService {
	return &OrderService{
		logger:                logger,
		tracer:                tracer,
		orderDAO:              orderDAO,
		msConfigService:       membershipConfigService,
		creditService:         creditService,
		userMembershipService: userMembershipService,
	}
}

func (s *OrderService) GetById(ctx context.Context, id string) (*model.Order, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.GetById")
	defer span.Finish()
	return s.orderDAO.FindExistById(ctx, id)
}

func (s *OrderService) GetOrderInfoById(ctx context.Context, id string) (*pb.OrderInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.GetOrderInfoById")
	defer span.Finish()
	order, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.Biz("order not found")
	}
	return &pb.OrderInfo{
		Id:             order.Id,
		UserId:         order.UserId,
		MembershipId:   order.MembershipId,
		OrderStatus:    pb.OrderStatus(order.OrderStatus),
		OrderType:      pb.OrderType(order.OrderType),
		SubName:        order.SubName,
		SubCredit:      uint64(order.SubCredit),
		SubAddOnCredit: uint64(order.SubAddOnCredit),
		SubStartDate:   uint64(order.SubStartDate.UnixMilli()),
		SubEndDate:     uint64(order.SubEndDate.UnixMilli()),
		Price:          uint64(order.Price),
		NumberCount:    int32(order.NumberCount),
		TotalAmount:    uint64(order.TotalAmount),
		PayAmount:      uint64(order.PayAmount),
		Currency:       order.Currency,
		PayExepiredAt:  uint64(order.PayExepiredAt.UnixMilli()),
		CreatedAt:      uint64(order.CreatedAt.UnixMilli()),
		UpdatedAt:      uint64(order.UpdatedAt.UnixMilli()),
	}, nil
}

// Subscribe 订阅会员订单
func (s *OrderService) Subscribe(ctx context.Context, userId string, orderType pb.OrderType, numberCount int32) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.Subscribe")
	defer span.Finish()

	// 1.获取用户会员
	userMembership, err := s.userMembershipService.GetByUserId(ctx, userId)
	if err != nil {
		return "0", errors.Biz("get user membership failed")
	}
	if userMembership == nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserAccountNotFound, "user account not found")
	}

	// 2.获取订阅信息
	subInfo := s.msConfigService.GetSubBaseInfo(ctx, int32(orderType))
	if subInfo == nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_SubscribeTypeNotFound, "subscription info not found")
	}

	// 3.创建订阅订单
	orderId, err := s.NewOrder(ctx, userId, userMembership.Id, subInfo, numberCount)
	if err != nil {
		return "0", err
	}

	return orderId, nil
}

// NewOrder 创建订阅订单
func (s *OrderService) NewOrder(ctx context.Context, userId string, msId string, subInfo *dto.MembershipSubBaseInfo, numberCount int32) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.NewOrder")
	defer span.Finish()

	s.logger.Info("NewOrder", "userId", userId, "msId", msId, "subInfo", subInfo, "numberCount", numberCount)

	// 1.用户会员是否过期
	expired, currentUserType, err := s.userMembershipService.CheckExpired(ctx, userId)
	if err != nil {
		return "0", errors.Biz("check user membership expired failed")
	}
	s.logger.Info("NewOrder", "expired", expired, "nowUserType", currentUserType)

	// 2.检查订阅类型
	switch subInfo.Type {
	case pb.OrderType_ORDER_TYPE_SUB_FREE:
		if currentUserType == int32(msPb.MembershipType_MEMBERSHIP_TYPE_FREE) && !expired {
			s.logger.Info("NewOrder", "user is free and not expired, can't subscribe free")
			return "0", errors.BizWithStatus(biz.Membership_Status_CanNotSubscribeFree, "user is free and not expired, can't subscribe free")
		}
		if currentUserType == int32(msPb.MembershipType_MEMBERSHIP_TYPE_PRO) && !expired {
			s.logger.Info("NewOrder", "user is pro and not expired, can't subscribe free")
			return "0", errors.BizWithStatus(biz.Membership_Status_CanNotSubscribeFree, "user is pro and not expired, can't subscribe free")
		}
	case pb.OrderType_ORDER_TYPE_SUB_PRO:
		// if currentUserType == int32(msPb.MembershipType_MEMBERSHIP_TYPE_PRO) && !expired {
		// 	s.logger.Info("NewOrder", "user is pro and not expired, can't subscribe pro")
		// 	return 0, errors.BizWithStatus(biz.Membership_Status_CanNotSubscribePro, "user is pro and not expired, can't subscribe pro")
		// }
	case pb.OrderType_ORDER_TYPE_SUB_PRO_ADD_ON_CREDIT:
		if currentUserType != int32(msPb.MembershipType_MEMBERSHIP_TYPE_PRO) {
			s.logger.Info("NewOrder", "user is not pro, can't subscribe pro add on credit")
			return "0", errors.BizWithStatus(biz.Membership_Status_CanNotSubscribeProAddOnCredit, "user is not pro, can't subscribe pro add on credit")
		}
		// TODO: 检查在用户会员有效期内充值次数是否已用完
		isUsedUp, err := s.checkMaxSubProAddOnCredit(ctx, userId)
		if err != nil {
			return "0", err
		}
		if isUsedUp {
			s.logger.Info("NewOrder", "user is used up, can't subscribe pro add on credit")
			return "0", errors.BizWithStatus(biz.Membership_Status_OverMaxAddOnCreditSubCountOfMonth, "user is used up, can't subscribe pro add on credit")
		}

	}

	// 3.创建订阅订单
	id := idgen.GenerateUUID()
	order := &model.Order{
		UserId:         userId,
		MembershipId:   msId,
		OrderStatus:    int32(pb.OrderStatus_ORDER_STATUS_PENDING),
		OrderType:      int32(subInfo.Type),
		SubName:        subInfo.Name,
		SubCredit:      subInfo.Credit * int64(numberCount),
		SubAddOnCredit: subInfo.AddOnCredit * int64(numberCount),
		SubStartDate:   time.Now(),
		SubEndDate:     time.Now().AddDate(0, subInfo.Duration, 0),
		Price:          subInfo.Price,
		NumberCount:    numberCount,
		TotalAmount:    subInfo.Price * int64(numberCount),
		PayAmount:      subInfo.Price * int64(numberCount),
		Currency:       subInfo.Currency,
		PayExepiredAt:  time.Now().Add(time.Minute * 30),
		// 折扣相关，暂无此需求待以后扩展
		IsDiscount:          false,
		DiscountPercent:     0,
		TotalDiscountAmount: 0,
		StripePayMode:       subInfo.StripePayMode,
		StripePriceId:       subInfo.StripePriceId,
	}
	order.Id = id

	return id, s.orderDAO.Save(ctx, order)
}

// CancelOrder 取消订阅订单
func (s *OrderService) CancelOrder(ctx context.Context, orderId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.CancelOrder")
	defer span.Finish()

	// TODO: 取消订阅订单
	order, err := s.GetById(ctx, orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.Biz("order not found")
	}

	if order.OrderStatus != int32(pb.OrderStatus_ORDER_STATUS_PENDING) {
		s.logger.Info("order status is not pending", "orderId", orderId, "orderStatus", order.OrderStatus)
		return nil
	}

	// 1.更新订单状态为已取消
	order.OrderStatus = int32(pb.OrderStatus_ORDER_STATUS_CANCELLED)
	order.PayTime = time.Now()
	err = s.orderDAO.Modify(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

// DoOrderPaySuccessHandler 订阅订单支付成功处理
func (s *OrderService) DoOrderPaySuccessHandler(ctx context.Context, orderId string, payOrderId string, stripeSubscriptionId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.DoOrderPaySuccessHandler")
	defer span.Finish()

	// TODO: 订阅订单支付成功处理
	order, err := s.GetById(ctx, orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.Biz("order not found")
	}

	if order.OrderStatus != int32(pb.OrderStatus_ORDER_STATUS_PENDING) {
		s.logger.Info("order status is not pending", "orderId", orderId, "orderStatus", order.OrderStatus)
		return nil
	}

	// 1.更新订单状态为已支付
	order.OrderStatus = int32(pb.OrderStatus_ORDER_STATUS_PAID)
	order.StripeSubscriptionId = stripeSubscriptionId
	order.PayOrderId = payOrderId
	order.PayTime = time.Now()
	err = s.orderDAO.Modify(ctx, order)
	if err != nil {
		return err
	}

	err = s.DoOrderProcessingHandler(ctx, orderId)
	if err != nil {
		return err
	}
	return nil
}

// DoOrderPayFailedHandler 订阅订单支付失败处理
func (s *OrderService) DoOrderPayFailedHandler(ctx context.Context, orderId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.DoOrderPayFailedHandler")
	defer span.Finish()

	// 订阅订单支付失败处理
	order, err := s.GetById(ctx, orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.Biz("order not found")
	}

	if order.OrderStatus != int32(pb.OrderStatus_ORDER_STATUS_PENDING) {
		s.logger.Info("order status is not pending", "orderId", orderId, "orderStatus", order.OrderStatus)
		return nil
	}

	// 1.更新订单状态为支付失败
	err = s.updateOrderStatus(ctx, orderId, int32(pb.OrderStatus_ORDER_STATUS_PAYMENT_FAILED))
	if err != nil {
		return err
	}
	return nil
}

// DoOrderPayExpireHandler 订阅订单支付过期处理
func (s *OrderService) DoOrderPayExpireHandler(ctx context.Context, orderId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.DoOrderPayExpireHandler")
	defer span.Finish()

	// 订阅订单支付过期处理
	order, err := s.GetById(ctx, orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.Biz("order not found")
	}

	if order.OrderStatus != int32(pb.OrderStatus_ORDER_STATUS_PENDING) {
		s.logger.Info("order status is not pending", "orderId", orderId, "orderStatus", order.OrderStatus)
		return nil
	}

	// 1.更新订单状态为已取消
	err = s.updateOrderStatus(ctx, orderId, int32(pb.OrderStatus_ORDER_STATUS_CANCELLED))
	if err != nil {
		return err
	}
	return nil
}

// DoOrderProcessingHandler 订阅订单处理中处理, 发放会员权益和Credit, 更新订单状态为已完成
func (s *OrderService) DoOrderProcessingHandler(ctx context.Context, orderId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.DoOrderProcessingHandler")
	defer span.Finish()

	// TODO: 订阅订单处理中处理
	order, err := s.GetById(ctx, orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.Biz("order not found")
	}

	if order.OrderStatus != int32(pb.OrderStatus_ORDER_STATUS_PAID) {
		s.logger.Info("order status is not paid", "orderId", orderId, "orderStatus", order.OrderStatus)
		return nil
	}

	// 1.更新订单状态为处理中
	err = s.updateOrderStatus(ctx, orderId, int32(pb.OrderStatus_ORDER_STATUS_PROCESSING))
	if err != nil {
		return err
	}

	// 2.根据订单类型处理后续逻辑
	if order.OrderType == int32(pb.OrderType_ORDER_TYPE_SUB_FREE) {
		s.logger.Info("msg", "order type is free", "orderId", orderId, "orderType", order.OrderType)
		// 更新会员账户
		err = s.userMembershipService.UpdateAccountType(ctx, order.UserId, int32(msPb.MembershipType_MEMBERSHIP_TYPE_FREE), order.StripeSubscriptionId, order.SubStartDate, order.SubEndDate)
		if err != nil {
			return err
		}

		// 更新积分账户
		payCredit := dto.CreditPayIntent{
			Type:        msPb.CreditPayType_CREDIT_PAY_TYPE_SUB_FREE,
			CreditType:  msPb.CreditType_CREDIT_TYPE_CREDIT,
			Credit:      order.SubCredit,
			AddOnCredit: order.SubAddOnCredit,
			Content:     "sub free order",
			Remark:      "",
		}
		_, err = s.creditService.InOrOutCredit(ctx, order.UserId, order.MembershipId, payCredit)
		if err != nil {
			return err
		}

	} else if order.OrderType == int32(pb.OrderType_ORDER_TYPE_SUB_PRO) {
		s.logger.Info("msg", "order type is pro", "orderId", orderId, "orderType", order.OrderType)
		err = s.userMembershipService.UpdateAccountType(ctx, order.UserId, int32(msPb.MembershipType_MEMBERSHIP_TYPE_PRO), order.StripeSubscriptionId, order.SubStartDate, order.SubEndDate)
		if err != nil {
			return err
		}

		// 更新积分账户
		payCredit := dto.CreditPayIntent{
			Type:        msPb.CreditPayType_CREDIT_PAY_TYPE_SUB_PRO,
			CreditType:  msPb.CreditType_CREDIT_TYPE_CREDIT,
			Credit:      order.SubCredit,
			AddOnCredit: order.SubAddOnCredit,
			Content:     "sub pro order",
			Remark:      "",
		}
		_, err = s.creditService.InOrOutCredit(ctx, order.UserId, order.MembershipId, payCredit)
		if err != nil {
			return err
		}

	} else if order.OrderType == int32(pb.OrderType_ORDER_TYPE_SUB_PRO_ADD_ON_CREDIT) {
		s.logger.Info("msg", "order type is sub pro add on credit", "orderId", orderId, "orderType", order.OrderType)
		// 更新积分账户
		payCredit := dto.CreditPayIntent{
			Type:        msPb.CreditPayType_CREDIT_PAY_TYPE_SUB_PRO_ADD_ON_CREDIT,
			CreditType:  msPb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT,
			Credit:      order.SubCredit,
			AddOnCredit: order.SubAddOnCredit,
			Content:     "sub pro add on credit order",
			Remark:      "",
		}
		_, err = s.creditService.InOrOutCredit(ctx, order.UserId, order.MembershipId, payCredit)
		if err != nil {
			return err
		}
	} else {
		//ignore
		s.logger.Info("msg", "ignore order type is unknown", "orderId", orderId, "orderType", order.OrderType)
	}

	// 3.权益和积分发放成功后，更新订单状态为已完成
	err = s.updateOrderStatus(ctx, orderId, int32(pb.OrderStatus_ORDER_STATUS_COMPLETED))
	if err != nil {
		return err
	}

	return nil
}

// HandlePayNotifyEventHandler 支付子系统的异步通知消息机制的支付结果回调事件
func (s *OrderService) HandlePayNotifyEventHandler(ctx context.Context, event eventbus.Event) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.HandlePayNotifyEvent")
	defer span.Finish()

	s.logger.Info("msg", "收到支付成功事件", "event", event)
	payEvent := event.Data.(payevent.PayNotifyEvent)
	switch event.Type {
	case payevent.PayNotifyEvent_PaySuccess:
		return s.DoOrderPaySuccessHandler(ctx, payEvent.OrderId, payEvent.PayRecordId, payEvent.SubscriptionId)
	case payevent.PayNotifyEvent_PayFailed:
		return s.DoOrderPayFailedHandler(ctx, payEvent.OrderId)
	case payevent.PayNotifyEvent_PayExpire:
		return s.DoOrderPayExpireHandler(ctx, payEvent.OrderId)
	case payevent.PayNotifyEvent_InvoicePaymentSucceeded:
		mOrderId, err := s.Subscribe(ctx, payEvent.UserId, pb.OrderType_ORDER_TYPE_SUB_PRO, 1)
		if err != nil {
			return err
		}
		payEvent.OrderId = mOrderId
		// 更新订单时间
		order, err := s.GetById(ctx, payEvent.OrderId)
		if err != nil {
			return err
		}
		if order == nil {
			return errors.Biz("order not found")
		}
		order.SubStartDate = payEvent.SubStartAt
		order.SubEndDate = payEvent.SubEndAt
		if err := s.orderDAO.Modify(ctx, order); err != nil {
			return err
		}

		return s.DoOrderPaySuccessHandler(ctx, payEvent.OrderId, payEvent.PayRecordId, payEvent.SubscriptionId)
	case payevent.PayNotifyEvent_CustomerSubscriptionsUpdated:
		// TODO: 订阅更新处理
		s.logger.Info("msg", "收到订阅更新事件", "event", event)
		return nil
	case payevent.PayNotifyEvent_CustomerSubscriptionsDeleted:
		// TODO: 订阅取消处理
		s.logger.Info("msg", "收到订阅取消事件", "event", event)
		return nil
	default:
		return nil
	}
}

// ========Private Methods========//
// checkMaxSubProAddOnCredit 检查用户会员的订阅订单数量是否已用完
func (s *OrderService) checkMaxSubProAddOnCredit(ctx context.Context, userId string) (bool, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.checkMaxSubProAddOnCredit")
	defer span.Finish()

	// 1.获取用户会员
	userMembership, err := s.userMembershipService.GetByUserId(ctx, userId)
	if err != nil {
		return false, err
	}
	userConfig, _ := s.userMembershipService.GetUserConfig(ctx)
	if !userConfig.Base.IsEnableSubAddOnCredit {
		return true, errors.Biz("user membership is not enable sub add on credit")
	}

	totalNum, err := s.orderDAO.GetTotalNum(ctx, userId, int32(pb.OrderStatus_ORDER_STATUS_COMPLETED), int32(pb.OrderType_ORDER_TYPE_SUB_PRO_ADD_ON_CREDIT), userMembership.StartAt, userMembership.EndAt)
	if err != nil {
		return false, err
	}
	if totalNum >= int32(userConfig.Base.MaxAddOnCreditSubCountOfMonth) {
		return true, nil
	}
	return false, nil
}

func (s *OrderService) updateOrderStatus(ctx context.Context, orderId string, orderStatus int32) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.updateOrderStatus")
	defer span.Finish()

	order, err := s.GetById(ctx, orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.Biz("order not found")
	}
	order.OrderStatus = orderStatus
	return s.orderDAO.Modify(ctx, order)
}

// updateOrderStripeSubscriptionId 更新订单的Stripe订阅ID
func (s *OrderService) updateOrderStripeSubscriptionId(ctx context.Context, orderId string, stripeSubscriptionId string) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "OrderService.updateOrderStripeSubscriptionId")
	defer span.Finish()

	order, err := s.GetById(ctx, orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.Biz("order not found")
	}
	order.StripeSubscriptionId = stripeSubscriptionId
	return s.orderDAO.Modify(ctx, order)
}
