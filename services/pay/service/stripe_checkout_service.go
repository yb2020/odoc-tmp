package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/webhook"

	"github.com/yb2020/odoc/services/pay/dao"
	"github.com/yb2020/odoc/services/pay/event"
	"github.com/yb2020/odoc/services/pay/model"

	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
)

// StripeCheckoutConfig 保存 StripeCheckoutService 的配置。
// 这些字段应从主应用程序配置 (config.Config.Pay.Stripe) 中填充。
type StripeCheckoutConfig struct {
	PublishableKey     string // Stripe 公钥
	SecretKey          string // Stripe 私钥
	WebhookSecret      string // Checkout webhook 事件专用
	CheckoutSuccessURL string // Checkout 成功 URL
	CheckoutCancelURL  string // Checkout 取消 URL
}

// StripeCheckoutService 提供与 Stripe Checkout 交互的方法。
type StripeCheckoutService struct {
	cfg                        StripeCheckoutConfig
	paymentRecordDAO           dao.PaymentRecordDAO
	logger                     logging.Logger     // 日志记录
	eventBus                   *eventbus.EventBus // 事件总线
	paymentSubscriptionService *PaymentSubscriptionService
}

// NewStripeCheckoutService 创建 StripeCheckoutService 的一个新实例。
func NewStripeCheckoutService(cfg StripeCheckoutConfig, prDAO dao.PaymentRecordDAO, logger logging.Logger, eventBus *eventbus.EventBus, paymentSubscriptionService *PaymentSubscriptionService) *StripeCheckoutService {
	// 通常的做法是在应用程序启动时全局设置一次 Stripe API 密钥。
	stripe.Key = cfg.SecretKey
	// 但是，如果此服务可能使用不同的密钥或在特定上下文中运行，
	// 也可以在每次调用时（如下面方法中那样）或每个服务实例设置它。
	return &StripeCheckoutService{
		cfg:                        cfg,
		paymentRecordDAO:           prDAO,
		logger:                     logger,
		eventBus:                   eventBus,
		paymentSubscriptionService: paymentSubscriptionService,
	}
}

// CreateCheckoutSessionParams 定义创建 Stripe Checkout Session 的参数。
type CreateCheckoutSessionParams struct {
	PayMode  string `json:"payMode"` // 支付模式 payment: 一次性付款, subscription: 订阅
	PriceId  string `json:"priceId"` // 价格ID, subscription模式下必填
	Name     string `json:"name"`    // 产品名称
	Amount   int64  `json:"amount"`  // 价格，以最小货币单位（例如：分）
	Quantity int64  `json:"quantity"`
	Currency string `json:"currency"` // ISO 货币代码，例如："usd", "cny"
	OrderId  string `json:"orderId"`  // 您的内部订单 ID
	UserId   string `json:"userId"`   // 您的内部用户 ID
	// 如果需要，您可以添加更多字段，例如 Metadata map[string]string
}

// CreateCheckoutSessionResult 保存创建 Stripe Checkout Session 的结果。
type CreateCheckoutSessionResult struct {
	SessionId string `json:"sessionId"`
}

// PublishableKey 返回 Stripe Checkout Service 的 PublishableKey。
func (s *StripeCheckoutService) PublishableKey() string {
	return s.cfg.PublishableKey
}

// CreateCheckoutSession 创建一个新的 Stripe Checkout Session。
func (s *StripeCheckoutService) CreateCheckoutSession(ctx context.Context, params CreateCheckoutSessionParams) (*CreateCheckoutSessionResult, error) {
	// 确保为此操作设置了 Stripe API 密钥。
	// 如果 stripe.Key 在启动时已全局设置，则此行可能是多余的，但能确保安全。
	// stripe.Key = s.cfg.SecretKey

	// ==Start:旧的创建Checkout Session的方式
	// stripeParams := &stripe.CheckoutSessionParams{
	// 	Mode: stripe.String(string(stripe.CheckoutSessionModePayment)), // 用于一次性付款。对于订阅，请使用 stripe.CheckoutSessionModeSubscription。
	// 	LineItems: []*stripe.CheckoutSessionLineItemParams{
	// 		{
	// 			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
	// 				Currency: stripe.String(params.Currency),
	// 				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
	// 					Name: stripe.String(params.Name),
	// 				},
	// 				UnitAmount: stripe.Int64(params.Amount),
	// 			},
	// 			Quantity: stripe.Int64(params.Quantity),
	// 		},
	// 	},
	// 	SuccessURL:        stripe.String(s.cfg.CheckoutSuccessURL),
	// 	CancelURL:         stripe.String(s.cfg.CheckoutCancelURL),
	// 	ClientReferenceID: stripe.String(params.OrderId), // 对于将 webhook 事件与您的内部订单匹配至关重要
	// 	Metadata: map[string]string{
	// 		"user_id":  params.UserId,  // 传递自定义数据的示例
	// 		"order_id": params.OrderId, // 也可以在此处放置 order_id 以实现冗余或特定的查找需求
	// 	},
	// 	// CustomerEmail: stripe.String("customer@example.com"), // 可选：预填客户电子邮件
	// 	// PaymentMethodTypes: stripe.StringSlice([]string{"card"}), // 可选：指定允许的支付方式类型
	// }
	// ==End:旧的创建Checkout Session的方式

	// 创建基本的Checkout Session参数
	stripeParams := &stripe.CheckoutSessionParams{
		SuccessURL:        stripe.String(s.cfg.CheckoutSuccessURL + "?orderId=" + params.OrderId),
		CancelURL:         stripe.String(s.cfg.CheckoutCancelURL + "?orderId=" + params.OrderId),
		ClientReferenceID: stripe.String(params.OrderId), // 对于将 webhook 事件与您的内部订单匹配至关重要
		Metadata: map[string]string{
			"user_id":  params.UserId,  // 传递自定义数据的示例
			"order_id": params.OrderId, // 也可以在此处放置 order_id 以实现冗余或特定的查找需求
		},
		// CustomerEmail: stripe.String("customer@example.com"), // 可选：预填客户电子邮件
		// PaymentMethodTypes: stripe.StringSlice([]string{"card"}), // 可选：指定允许的支付方式类型
	}

	// 根据支付模式设置不同的参数
	switch params.PayMode {
	case "subscription":
		// 订阅模式 - 使用预先创建的Price ID
		stripeParams.Mode = stripe.String(string(stripe.CheckoutSessionModeSubscription))
		stripeParams.LineItems = []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(params.PriceId),
				Quantity: stripe.Int64(params.Quantity),
			},
		}
	default:
		// 默认为一次性付款模式
		stripeParams.Mode = stripe.String(string(stripe.CheckoutSessionModePayment))
		stripeParams.LineItems = []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(params.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(params.Name),
					},
					UnitAmount: stripe.Int64(params.Amount),
				},
				Quantity: stripe.Int64(params.Quantity),
			},
		}
	}

	// 对于订阅，Mode 和 LineItems 会有所不同，通常会引用一个 Price ID。

	sess, err := session.New(stripeParams)
	if err != nil {
		s.logger.Error("为订单ID %s 创建 Stripe Checkout 会话失败：%v", params.OrderId, err)
		return nil, errors.Biz("failed to create Stripe Checkout session: %s" + err.Error())
	}
	s.logger.Info("stripe checkout url", sess.URL)

	return &CreateCheckoutSessionResult{SessionId: sess.ID}, nil
}

// 取消订阅，在周期结束时生效
func (s *StripeCheckoutService) CancelSubscriptionAtPeriodEnd(ctx context.Context, subID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	return subscription.Update(subID, params)
}

// 立即取消订阅
func (s *StripeCheckoutService) CancelSubscriptionImmediately(ctx context.Context, subID string) (*stripe.Subscription, error) {
	// 注意：subscription.Cancel 默认就是立即取消
	// 如果想在周期末取消，必须用 subscription.Update
	params := &stripe.SubscriptionCancelParams{} // 可以留空
	return subscription.Cancel(subID, params)
}

// HandleCheckoutWebhook 处理与 Checkout 相关的传入 Stripe webhook 事件。
//
// 该函数负责处理来自 Stripe 的多种 webhook 事件，以确保支付状态的同步和订阅生命周期的正确管理。
//
// 监听的事件按支付模式区分如下：
//
// --- 通用事件 (一次性支付和订阅) ---
// - checkout.session.completed: 用户在Stripe托管页面完成支付流程。对于一次性支付，这表示支付成功；对于订阅，这表示首次支付成功。
// - checkout.session.async_payment_succeeded: 异步支付方式（如银行转账）最终确认成功。
// - checkout.session.async_payment_failed: 异步支付方式最终确认失败。
//
// --- 仅订阅模式事件 ---
// - customer.subscription.created: 成功创建了一个新的订阅（在首次支付成功后触发）。
// - customer.subscription.updated: 订阅详情发生变更（如升级、降级、变更计费周期）。
// - customer.subscription.deleted: 订阅被取消或因付款失败而终止。
// - invoice.payment_succeeded: 订阅的周期性续订付款成功。
// - invoice.payment_failed: 订阅的周期性续订付款失败。
func (s *StripeCheckoutService) HandleCheckoutWebhook(ctx context.Context, payload []byte, signatureHeader string) error {
	// 确保为此操作设置了 Stripe API 密钥（尽管 webhook.ConstructEvent 可能不直接使用它，但相关的 API 调用可能会）。
	stripe.Key = s.cfg.SecretKey

	event, err := webhook.ConstructEvent(payload, signatureHeader, s.cfg.WebhookSecret)
	if err != nil {
		s.logger.Info("Webhook 签名验证失败：%v", err)
		return errors.Biz("webhook signature verification failed: %s" + err.Error())
	}

	s.logger.Info("msg", "收到 Stripe webhook 事件", "event_type", event.Type, "event_id", event.ID)

	switch event.Type {
	case "checkout.session.completed":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			// s.logger.Errorf(ctx, "解析 checkout.session.completed 的 webhook JSON 出错 (事件ID: %s): %v", event.ID, err)
			return errors.Biz("error parsing webhook JSON for %s: %s")
		}
		// 对于同步支付，状态为 'paid'。对于异步支付，状态为 'unpaid'。
		if checkoutSession.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			s.logger.Info("msg", "Checkout session completed and paid (sync).", "session_id", checkoutSession.ID, "order_id", checkoutSession.ClientReferenceID)
			return s.processCheckoutSessionCompleted(ctx, &checkoutSession)
		} else {
			s.logger.Info("msg", "Checkout session completed, awaiting async payment.", "session_id", checkoutSession.ID, "order_id", checkoutSession.ClientReferenceID, "payment_status", checkoutSession.PaymentStatus)
			// 在这里什么都不做，等待 async_payment_succeeded 或 async_payment_failed 事件
			return nil
		}

	case "checkout.session.async_payment_succeeded":
		var checkoutSession stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
			s.logger.Error("msg", "Error parsing webhook JSON for async_payment_succeeded", "event_id", event.ID, "error", err)
			return errors.Biz("error parsing webhook JSON for %s: %s")
		}
		// s.logger.Info(ctx, "Checkout 会话异步支付成功，SessionID：%s，OrderID：%s", checkoutSession.ID, checkoutSession.ClientReferenceID)
		s.logger.Info("msg", "Checkout session async payment succeeded", "session_id", checkoutSession.ID, "order_id", checkoutSession.ClientReferenceID)
		// 如果需要，实现与 processCheckoutSessionCompleted 类似的逻辑，或更新现有记录。
		// 这对于那些不是立即确认的支付方式很重要。
		return s.processCheckoutSessionCompleted(ctx, &checkoutSession) // 或一个专用的异步处理程序

	case "checkout.session.async_payment_failed":
		var checkoutSession stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
			s.logger.Error("msg", "Error parsing webhook JSON for async_payment_failed", "event_id", event.ID, "error", err)
			return errors.Biz("error parsing webhook JSON for %s: %s")
		}
		s.logger.Warn("msg", "Checkout session async payment failed.", "session_id", checkoutSession.ID, "order_id", checkoutSession.ClientReferenceID, "payment_status", checkoutSession.PaymentStatus)
		// s.logger.Warnf(ctx, "Checkout 会话异步支付失败，SessionID：%s，OrderID：%s，PaymentStatus：%s", checkoutSession.ID, checkoutSession.ClientReferenceID, checkoutSession.PaymentStatus)
		// 更新您的内部订单/支付状态以反映失败。
		return s.processCheckoutSessionPaymentFailed(ctx, &checkoutSession)

	case "customer.subscription.created":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			s.logger.Error("msg", "Error parsing webhook JSON for customer.subscription.created", "event_id", event.ID, "error", err)
			return errors.Biz(fmt.Sprintf("error parsing webhook JSON for customer.subscription.created: %v", err))
		}
		s.logger.Info("msg", "Subscription created", "subscription_id", subscription.ID, "customer_id", subscription.Customer.ID)
		return s.processSubscriptionCreated(ctx, &subscription)
	case "customer.subscription.updated":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			s.logger.Error("msg", "Error parsing webhook JSON for customer.subscription.updated", "event_id", event.ID, "error", err)
			return errors.Biz(fmt.Sprintf("error parsing webhook JSON for customer.subscription.updated: %v", err))
		}

		s.logger.Info("msg", "Subscription updated", "subscription_id", subscription.ID, "new_status", subscription.Status, "cancel_at_period_end", subscription.CancelAtPeriodEnd)
		// ignore this event
		return s.processSubscriptionUpdated(ctx, &subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			s.logger.Error("msg", "Error parsing webhook JSON for customer.subscription.deleted", "event_id", event.ID, "error", err)
			return errors.Biz(fmt.Sprintf("error parsing webhook JSON for customer.subscription.deleted: %v", err))
		}

		s.logger.Info("msg", "Subscription deleted", "subscription_id", subscription.ID, "final_status", subscription.Status)
		return s.processSubscriptionDeleted(ctx, &subscription)

	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			s.logger.Error("msg", "Error parsing webhook JSON for invoice.payment_succeeded", "event_id", event.ID, "error", err)
			return errors.Biz(fmt.Sprintf("error parsing webhook JSON for invoice.payment_succeeded: %v", err))
		}

		return s.processInvoicePaymentSucceeded(ctx, &invoice)

	case "invoice.payment_failed":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			s.logger.Error("msg", "Error parsing webhook JSON for invoice.payment_failed", "event_id", event.ID, "error", err)
			return errors.Biz(fmt.Sprintf("error parsing webhook JSON for invoice.payment_failed: %v", err))
		}

		var subscriptionID string
		if len(invoice.Lines.Data) > 0 && invoice.Lines.Data[0].Subscription != nil {
			subscriptionID = invoice.Lines.Data[0].Subscription.ID
		}

		if subscriptionID == "" {
			s.logger.Info("msg", "Invoice payment failed event is not related to a subscription, skipping.", "invoice_id", invoice.ID)
			return nil
		}

		s.logger.Warn("msg", "Subscription renewal failed", "subscription_id", subscriptionID, "invoice_id", invoice.ID)
		return s.processInvoicePaymentFailed(ctx, &invoice)

	default:
		s.logger.Info("未处理的 Stripe webhook 事件类型：%s", event.Type)
	}

	return nil
}

// processCheckoutSessionCompleted 处理接收到 checkout.session.completed 事件时的业务逻辑。
func (s *StripeCheckoutService) processCheckoutSessionCompleted(ctx context.Context, cs *stripe.CheckoutSession) error {
	if cs.ClientReferenceID == "" {
		// s.logger.Errorf(ctx, "checkout.session.completed 事件中的 ClientReferenceID 为空，SessionID：%s", cs.ID)
		return errors.Biz("client_reference_id is empty in checkout session %s" + cs.ID)
	}

	// 检查支付状态。对于大多数直接付款，它应该是 'paid'。
	// 对于某些支付方式（例如银行转账），'paid' 状态可能通过 'checkout.session.async_payment_succeeded' 事件传来。
	if cs.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		// s.logger.Warnf(ctx, "订单ID %s 的 Checkout 会话 %s 的 PaymentStatus 是 '%s'，期望是 '%s'。可能需要异步处理。",
		// 	cs.ID, cs.ClientReferenceID, cs.PaymentStatus, stripe.CheckoutSessionPaymentStatusPaid)
		// 根据您的业务逻辑，您可能会返回错误，或者只是记录日志并等待异步事件。
		// 在此示例中，如果状态是 'paid'，我们将继续处理；否则，假定将有异步事件跟随，或者这是一个问题。
		return errors.Biz("checkout session %s payment status is '%s', not '%s'" + cs.ID)
	}

	// 幂等性：检查此支付（例如，通过 ProviderTxId 或 OrderId）是否已成功处理。
	existingRecord, err := s.paymentRecordDAO.GetByProviderTxId(ctx, cs.ID)
	if err != nil { // 根据您的 DAO 调整错误检查
		s.logger.Error(ctx, "检查 ProviderTxID %s 的现有支付记录时出错：%v", cs.ID, err)
		return errors.Biz("failed to check existing payment record: %s" + err.Error())
	}
	if existingRecord != nil && existingRecord.Status == model.PaymentStatusSucceeded {
		s.logger.Info(ctx, "ProviderTxID %s (OrderID: %s) 的支付已处理。", cs.ID, cs.ClientReferenceID)
		return nil // 已成功处理
	}

	userID := ""
	if cs.Metadata != nil {
		userID = cs.Metadata["user_id"]
	}

	var subscriptionId string
	if cs.Mode == stripe.CheckoutSessionModeSubscription && cs.Subscription != nil {
		subscriptionId = cs.Subscription.ID
	}

	var invoiceId string
	if cs.Invoice != nil {
		invoiceId = cs.Invoice.ID
	}

	var payMode string
	switch cs.Mode {
	case stripe.CheckoutSessionModeSubscription:
		payMode = model.PaymentModeSubscription
	case stripe.CheckoutSessionModePayment:
		payMode = model.PaymentModePayment
	}

	currentTime := time.Now()
	paymentRecord := &model.PaymentRecord{
		// BaseModel 字段 (Id, CreatedAt 等) 将由 GORM 或您的 BaseDAO 处理
		UserId:                 userID,
		OrderId:                cs.ClientReferenceID,
		Amount:                 cs.AmountTotal, // 金额，以最小货币单位
		Currency:               string(cs.Currency),
		Status:                 model.PaymentStatusSucceeded, // 映射到您内部的 'succeeded' 状态
		Channel:                model.PaymentChannelStripe,   // 在您的 model 包中定义为常量，例如 model.PaymentChannelStripeCheckout
		ProviderTxId:           cs.ID,                        // Stripe Checkout Session ID (cs_xxx...)
		ProviderSubscriptionId: subscriptionId,               // Stripe Subscription ID (sub_xxx...)
		InvoiceId:              invoiceId,                    // Stripe Invoice ID (in_xxx...)
		Description:            fmt.Sprintf("Stripe Checkout for Order %s", cs.ClientReferenceID),
		// Metadata: cs.Metadata, // 如果 cs.Metadata 中的某些数据有用，您可能希望存储它们
		// PaymentMethodDetails: // 如果需要，可以存储 cs.PaymentMethodDetails，它是一个 stripe.CheckoutSessionPaymentMethodDetails
		PaidAt:  currentTime, // 分配 *time.Time
		PayMode: payMode,
		// ProviderErrorCode: // 如果 cs.PaymentStatus 指示此处处理的错误，则填充
		// ProviderErrorMessage: // 如果 cs.PaymentStatus 指示此处处理的错误，则填充
	}
	paymentRecord.Id = idgen.GenerateUUID()
	paymentRecord.CreatedAt = currentTime
	paymentRecord.UpdatedAt = currentTime

	// 如果 PaymentIntent ID 可用且有用，您也可以存储它。
	// if cs.PaymentIntent != nil {
	// 	paymentRecord.Metadata["stripe_payment_intent_id"] = cs.PaymentIntent.ID
	// }

	if err := s.paymentRecordDAO.Save(ctx, paymentRecord); err != nil {
		// s.logger.Errorf(ctx, "为订单ID %s (ProviderTxID: %s) 保存支付记录失败：%v",
		// 	cs.ClientReferenceID, cs.ID, err)
		return errors.Biz("failed to save payment record: %s" + err.Error())
	}

	// 处理订阅信息
	if paymentRecord.PayMode == model.PaymentModeSubscription {
		paySub, err := s.paymentSubscriptionService.GetByProviderSubscriptionId(ctx, paymentRecord.ProviderSubscriptionId)
		if err != nil {
			return errors.Biz("failed to get payment subscription: %s" + err.Error())
		}
		if paySub == nil {
			return errors.Biz("payment subscription not found")
		}
		paySub.Status = model.PaymentSubscription_StatusActive
		paySub.UserId = userID
		paySub.ProviderSubscriptionId = subscriptionId
		paySub.Description = paymentRecord.Description
		s.paymentSubscriptionService.UpdateSubscription(ctx, paySub)
	}

	// 重要提示：在此处触发任何后续业务逻辑
	// 例如：更新订单状态、授予服务访问权限、发送确认电子邮件等。
	// 这可能涉及发布事件或调用另一个服务。

	s.eventBus.Publish(ctx, eventbus.Event{
		Type: event.PayNotifyEvent_PaySuccess,
		Data: event.PayNotifyEvent{
			OrderId:        paymentRecord.OrderId, // Use the parsed int64 value
			PayRecordId:    paymentRecord.Id,      // Use the parsed int64 value
			SubscriptionId: paymentRecord.ProviderSubscriptionId,
			InvoiceId:      paymentRecord.InvoiceId,
			PayMode:        paymentRecord.PayMode,
			UserId:         paymentRecord.UserId,
		},
	}, true)

	return nil
}

// processCheckoutSessionPaymentFailed 处理 checkout.session.async_payment_failed 事件的业务逻辑。
func (s *StripeCheckoutService) processCheckoutSessionPaymentFailed(ctx context.Context, cs *stripe.CheckoutSession) error {
	if cs.ClientReferenceID == "" {
		// 如果没有 ClientReferenceID，我们无法将其与内部订单关联，因此只能记录并返回。
		s.logger.Error("msg", "client_reference_id is empty in checkout.session.async_payment_failed event", "session_id", cs.ID)
		return nil // 返回 nil 以向 Stripe 确认 webhook 已收到，因为重试也无济于事。
	}
	s.logger.Warn("msg", "Payment failed for order", "order_id", cs.ClientReferenceID, "session_id", cs.ID)

	// 在这里，您通常会更新内部支付/订单记录为“失败”状态。
	// 例如，通过 ClientReferenceID (orderId) 查找支付记录并更新其状态。
	// 这部分取决于您的具体 DAO 实现。
	// 为简单起见，我们仅作记录。

	// 您可能还想发布一个事件来通知系统的其他部分。
	// orderIDInt, _ := strconv.ParseInt(cs.ClientReferenceID, 10, 64)
	// s.eventBus.Publish(ctx, eventbus.Event{
	// 	Type: event.PayNotifyEvent_PayFail,
	// 	Data: event.PayNotifyEvent{
	// 		OrderId:     orderIDInt,
	// 	},
	// }, true)
	s.eventBus.Publish(ctx, eventbus.Event{
		Type: event.PayNotifyEvent_PayFailed,
		Data: event.PayNotifyEvent{
			OrderId:     cs.ClientReferenceID, // Use the parsed int64 value
			PayRecordId: "0",                  // Use the parsed int64 value
		},
	}, true)

	return nil
}

// processSubscriptionCreated 处理接收到 customer.subscription.created 事件时的业务逻辑。
func (s *StripeCheckoutService) processSubscriptionCreated(ctx context.Context, subscription *stripe.Subscription) error {

	if subscription.ID == "" {
		return errors.Biz("subscription ID is empty")
	}
	extPaySub, err := s.paymentSubscriptionService.GetByProviderSubscriptionId(ctx, subscription.ID)
	if err != nil {
		s.logger.Error("msg", "Error checking existing payment subscription for subscription", "subscription_id", subscription.ID, "error", err)
		return err
	}
	if extPaySub != nil {
		s.logger.Info("msg", "Payment subscription already exists for subscription", "subscription_id", subscription.ID)
		return nil
	}

	PeriodStart := time.Unix(subscription.Items.Data[0].CurrentPeriodStart, 0)
	PeriodEnd := time.Unix(subscription.Items.Data[0].CurrentPeriodEnd, 0)

	// TODO: 处理订阅创建事件
	paySub := &model.PaymentSubscription{
		Status:                 model.PaymentSubscription_StatusActive,
		ProviderSubscriptionId: subscription.ID,
		PriceId:                subscription.Items.Data[0].Price.ID,
		StartAt:                &PeriodStart,
		EndAt:                  &PeriodEnd,
	}
	paySubId, err := s.paymentSubscriptionService.NewSubscription(ctx, paySub)
	if err != nil {
		return err
	}
	s.logger.Info("msg", "Subscription created successfully", "subscription_id", subscription.ID, "payment_subscription_id", paySubId)

	// 发送订阅创建事件
	// orderIDInt, _ := strconv.ParseInt(subscription.Metadata["order_id"], 10, 64)
	// s.eventBus.Publish(ctx, eventbus.Event{
	// 	Type: event.PayNotifyEvent_CustomerSubscriptionsCreated,
	// 	Data: event.PayNotifyEvent{
	// 		OrderId: orderIDInt,
	// 	},
	// }, true)
	return nil
}

// processSubscriptionUpdated 处理接收到 customer.subscription.updated 事件时的业务逻辑。
func (s *StripeCheckoutService) processSubscriptionUpdated(ctx context.Context, subscription *stripe.Subscription) error {

	subRecord, err := s.paymentSubscriptionService.GetByProviderSubscriptionId(ctx, subscription.ID)
	if err != nil {
		return err
	}
	if subRecord == nil {
		return nil
	}
	if subRecord.Status == model.PaymentSubscription_StatusCanceled {
		return nil
	}

	PeriodStart := time.Unix(subscription.Items.Data[0].CurrentPeriodStart, 0)
	PeriodEnd := time.Unix(subscription.Items.Data[0].CurrentPeriodEnd, 0)

	if subscription.CancelAtPeriodEnd {
		cancelAtTime := time.Unix(subscription.CancelAt, 0)
		subRecord.CancelAtPeriodEnd = true
		subRecord.CancelAt = &cancelAtTime
	} else {
		subRecord.CancelAtPeriodEnd = false
		subRecord.CancelAt = nil
	}
	subRecord.StartAt = &PeriodStart
	subRecord.EndAt = &PeriodEnd
	s.paymentSubscriptionService.UpdateSubscription(ctx, subRecord)
	// 发送订阅更新事件
	s.eventBus.Publish(ctx, eventbus.Event{
		Type: event.PayNotifyEvent_CustomerSubscriptionsUpdated,
		Data: event.PayNotifyEvent{
			SubscriptionId: subscription.ID,
		},
	}, true)
	return nil
}

// processSubscriptionDeleted 处理接收到 customer.subscription.deleted 事件时的业务逻辑。
func (s *StripeCheckoutService) processSubscriptionDeleted(ctx context.Context, subscription *stripe.Subscription) error {

	subRecord, err := s.paymentSubscriptionService.GetByProviderSubscriptionId(ctx, subscription.ID)
	if err != nil {
		return err
	}
	if subRecord == nil {
		return nil
	}

	if subRecord.Status == model.PaymentSubscription_StatusCanceled {
		s.logger.Info("msg", "Subscription already canceled", "subscription_id", subscription.ID)
		return nil
	}

	subRecord.Status = model.PaymentSubscription_StatusCanceled
	canceledAtTime := time.Unix(subscription.CanceledAt, 0)
	subRecord.CancelAt = &canceledAtTime
	subRecord.CancelReason = string(subscription.CancellationDetails.Reason)
	s.paymentSubscriptionService.UpdateSubscription(ctx, subRecord)

	// 发送订阅删除事件
	s.eventBus.Publish(ctx, eventbus.Event{
		Type: event.PayNotifyEvent_CustomerSubscriptionsDeleted,
		Data: event.PayNotifyEvent{
			SubscriptionId: subscription.ID,
		},
	}, true)
	return nil
}

// processInvoicePaymentSucceeded 处理接收到 invoice.payment_succeeded 事件时的业务逻辑。
func (s *StripeCheckoutService) processInvoicePaymentSucceeded(ctx context.Context, invoice *stripe.Invoice) error {
	var subscriptionID string
	if invoice.Parent.SubscriptionDetails.Subscription.ID != "" {
		subscriptionID = invoice.Parent.SubscriptionDetails.Subscription.ID
	}

	if subscriptionID == "" {
		s.logger.Info("msg", "Invoice payment succeeded event is not related to a subscription, skipping.", "invoice_id", invoice.ID)
		return nil
	}

	s.logger.Info("msg", "Subscription renewal successful", "subscription_id", subscriptionID, "invoice_id", invoice.ID)
	existingSubRecord, err := s.paymentSubscriptionService.GetByProviderSubscriptionId(ctx, subscriptionID)
	if err != nil {
		s.logger.Error("msg", "Error checking existing payment subscription for subscription", "subscription_id", subscriptionID, "error", err)
		return err
	}
	if existingSubRecord == nil {
		s.logger.Info("msg", "ignore invoice payment succeeded event for subscription", "subscription_id", subscriptionID)
		return nil
	}

	existingPaymentRecord, err := s.paymentRecordDAO.GetBySubscriptionIdAndInvoiceId(ctx, subscriptionID, invoice.ID)
	if err != nil {
		s.logger.Error("msg", "Error checking existing payment record for subscription", "subscription_id", subscriptionID, "invoice_id", invoice.ID, "error", err)
		return err
	}
	if existingPaymentRecord != nil && existingPaymentRecord.Status == model.PaymentStatusSucceeded {
		s.logger.Info(ctx, "Subscription %s (InvoiceID: %s) 的支付已处理。", subscriptionID, invoice.ID)
		return nil // 已成功处理
	}

	// 更新订阅订单时间
	startAt := time.Unix(invoice.Lines.Data[0].Period.Start, 0)
	endAt := time.Unix(invoice.Lines.Data[0].Period.End, 0)
	existingSubRecord.StartAt = &startAt
	existingSubRecord.EndAt = &endAt
	if err := s.paymentSubscriptionService.UpdateSubscription(ctx, existingSubRecord); err != nil {
		s.logger.Error("msg", "Error updating payment subscription for subscription", "subscription_id", subscriptionID, "error", err)
		return err
	}

	// TODO: 根据发票信息创建支付记录
	renewalPaymentRecord := &model.PaymentRecord{
		OrderId:  invoice.Metadata["order_id"],
		UserId:   existingSubRecord.UserId,
		Amount:   invoice.AmountPaid,
		Currency: string(invoice.Currency),
		Status:   model.PaymentStatusSucceeded,
		Channel:  model.PaymentChannelStripe,
		//ProviderTxId:           invoice.ID,
		ProviderSubscriptionId: subscriptionID,
		InvoiceId:              invoice.ID,
		Description:            "renewal subscription for invoice " + invoice.ID,
		PayMode:                model.PaymentModeSubscription,
		PaidAt:                 time.Unix(invoice.StatusTransitions.PaidAt, 0),
	}
	renewalPaymentRecord.Id = idgen.GenerateUUID()

	if err := s.paymentRecordDAO.Save(ctx, renewalPaymentRecord); err != nil {
		s.logger.Error("msg", "Failed to create payment record for subscription renewal", "subscription_id", subscriptionID, "invoice_id", invoice.ID, "error", err)
		return errors.Biz(fmt.Sprintf("failed to create payment record for renewal: %v", err))
	}

	s.logger.Info("msg", "Successfully created payment record for subscription renewal", "subscription_id", subscriptionID, "invoice_id", invoice.ID)

	s.eventBus.Publish(ctx, eventbus.Event{
		Type: event.PayNotifyEvent_InvoicePaymentSucceeded,
		Data: event.PayNotifyEvent{
			// OrderId:        orderIDInt,       // Use the parsed int64 value
			PayRecordId:    renewalPaymentRecord.Id, // Use the parsed int64 value
			SubscriptionId: renewalPaymentRecord.ProviderSubscriptionId,
			InvoiceId:      renewalPaymentRecord.InvoiceId,
			PayMode:        renewalPaymentRecord.PayMode,
			UserId:         renewalPaymentRecord.UserId,
			SubStartAt:     *existingSubRecord.StartAt,
			SubEndAt:       *existingSubRecord.EndAt,
		},
	}, true)

	return nil
}

// processInvoicePaymentFailed 处理接收到 invoice.payment_failed 事件时的业务逻辑。
func (s *StripeCheckoutService) processInvoicePaymentFailed(ctx context.Context, invoice *stripe.Invoice) error {

	s.logger.Info("msg", "Invoice payment failed", "invoice_id", invoice.ID)
	return nil
}
