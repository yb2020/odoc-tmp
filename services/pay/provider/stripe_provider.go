package provider

import (
	"context"
	"encoding/json"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/refund"
	"github.com/stripe/stripe-go/v82/webhook"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pay/model"
)

// StripeProvider 实现PaymentProvider接口，封装Stripe SDK和逻辑
type StripeProvider struct {
	apiKey        string
	webhookSecret string
	logger        logging.Logger
}

// NewStripeProvider 创建一个新的Stripe支付提供者
func NewStripeProvider(apiKey, webhookSecret string, logger logging.Logger) *StripeProvider {
	// 初始化Stripe客户端
	stripe.Key = apiKey

	return &StripeProvider{
		apiKey:        apiKey,
		webhookSecret: webhookSecret,
		logger:        logger,
	}
}

// GetName 返回支付渠道名称
func (p *StripeProvider) GetName() string {
	return "STRIPE"
}

// CreateCharge 创建Stripe支付意向
func (p *StripeProvider) CreateCharge(ctx context.Context, params *ChargeParams) (*ChargeResult, error) {
	// 构建Stripe支付意向参数
	piParams := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(params.Amount),
		Currency:      stripe.String(params.Currency),
		PaymentMethod: stripe.String(params.PaymentMethodId),
		Description:   stripe.String(params.Description),
		Params: stripe.Params{
			Context: ctx,
		},
		ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodManual)),
		Confirm:            stripe.Bool(true), // 您设置了立即确认
		// ReturnURL:          stripe.String(params.ReturnURL), //这个URL就是用户完成银行验证后返回的地方
	}

	// 添加元数据
	piParams.AddMetadata("order_id", params.OrderId)
	piParams.AddMetadata("user_id", params.UserId)
	for key, value := range params.Metadata {
		piParams.AddMetadata(key, value)
	}

	// 创建支付意向
	pi, err := paymentintent.New(piParams)
	if err != nil {
		p.logger.Error("msg", "创建Stripe支付意向失败", "error", err.Error())
		return &ChargeResult{
			Status:       model.PaymentStatusFailed,
			ErrorCode:    "stripe_create_payment_intent_error",
			ErrorMessage: err.Error(),
		}, errors.Biz("创建Stripe支付意向失败")
	}

	// 构建返回结果
	result := &ChargeResult{
		ProviderTxId:      pi.ID,
		Status:            mapStripeStatusToInternal(pi.Status),
		ClientSecret:      pi.ClientSecret,
		RequiresAction:    pi.Status == stripe.PaymentIntentStatusRequiresAction,
		PaymentMethodType: string(pi.PaymentMethodTypes[0]),
		Metadata:          make(map[string]string),
	}

	// 如果需要额外操作，添加下一步操作URL
	if pi.NextAction != nil && pi.NextAction.RedirectToURL != nil {
		result.RedirectURL = pi.NextAction.RedirectToURL.URL
	}

	// 复制元数据
	for key, value := range pi.Metadata {
		result.Metadata[key] = value
	}

	return result, nil
}

// HandleWebhook 解析并验证Stripe Webhook数据
func (p *StripeProvider) HandleWebhook(ctx context.Context, requestData []byte, signature string) (*WebhookEvent, error) {
	// 验证webhook签名
	event, err := webhook.ConstructEvent(requestData, signature, p.webhookSecret)
	if err != nil {
		p.logger.Error("msg", "Stripe webhook签名验证失败", "error", err.Error())
		return nil, err
	}

	// 解析事件数据
	webhookEvent := &WebhookEvent{
		EventType: string(event.Type),
		RawData:   make(map[string]interface{}),
		Metadata:  make(map[string]string),
	}

	// 将原始数据转换为map
	if err := json.Unmarshal(event.Data.Raw, &webhookEvent.RawData); err != nil {
		p.logger.Error("msg", "解析Stripe事件数据失败", "error", err.Error())
		return nil, err
	}

	// 根据事件类型处理不同的数据结构
	switch event.Type {
	case "payment_intent.succeeded", "payment_intent.payment_failed", "payment_intent.canceled":
		var pi stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &pi)
		if err != nil {
			p.logger.Error("msg", "解析Stripe PaymentIntent数据失败", "error", err.Error())
			return nil, err
		}

		webhookEvent.ProviderTxId = pi.ID
		webhookEvent.Status = mapStripeStatusToInternal(pi.Status)
		webhookEvent.Amount = pi.Amount
		webhookEvent.Currency = string(pi.Currency)

		// 复制元数据
		for key, value := range pi.Metadata {
			webhookEvent.Metadata[key] = value
		}

	case "charge.refunded":
		var charge stripe.Charge
		err := json.Unmarshal(event.Data.Raw, &charge)
		if err != nil {
			p.logger.Error("msg", "解析Stripe Charge数据失败", "error", err.Error())
			return nil, err
		}

		webhookEvent.ProviderTxId = charge.PaymentIntent.ID
		webhookEvent.Status = model.PaymentStatusRefunded
		webhookEvent.Amount = charge.Amount
		webhookEvent.Currency = string(charge.Currency)

		// 复制元数据
		for key, value := range charge.Metadata {
			webhookEvent.Metadata[key] = value
		}

	default:
		// 对于其他事件类型，只记录事件类型，不做特殊处理
		p.logger.Info("msg", "收到未处理的Stripe事件类型", "eventType", event.Type)
	}

	return webhookEvent, nil
}

// CreateRefund 发起Stripe退款
func (p *StripeProvider) CreateRefund(ctx context.Context, chargeId string, amount int64, reason string) (*RefundResult, error) {
	// 构建退款参数
	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(chargeId),
		Amount:        stripe.Int64(amount),
		Reason:        stripe.String(mapRefundReasonToStripe(reason)),
		Params: stripe.Params{
			Context: ctx,
		},
	}

	// 创建退款
	r, err := refund.New(refundParams)
	if err != nil {
		p.logger.Error("msg", "创建Stripe退款失败", "error", err.Error())
		return &RefundResult{
			Status:       "failed",
			ErrorCode:    "stripe_create_refund_error",
			ErrorMessage: err.Error(),
		}, err
	}

	// 构建返回结果
	result := &RefundResult{
		RefundId: r.ID,
		Status:   string(r.Status),
	}

	return result, nil
}

// GetChargeStatus 查询Stripe支付状态
func (p *StripeProvider) GetChargeStatus(ctx context.Context, chargeId string) (string, error) {
	// 查询支付意向
	pi, err := paymentintent.Get(chargeId, &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Context: ctx,
		},
	})
	if err != nil {
		p.logger.Error("msg", "查询Stripe支付状态失败", "chargeId", chargeId, "error", err.Error())
		return "", err
	}

	return mapStripeStatusToInternal(pi.Status), nil
}

// mapStripeStatusToInternal 将Stripe支付状态映射为内部状态
func mapStripeStatusToInternal(status stripe.PaymentIntentStatus) string {
	switch status {
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		return model.PaymentStatusPending
	case stripe.PaymentIntentStatusRequiresConfirmation:
		return model.PaymentStatusPending
	case stripe.PaymentIntentStatusRequiresAction:
		return model.PaymentStatusRequiresAction
	case stripe.PaymentIntentStatusProcessing:
		return model.PaymentStatusPending
	case stripe.PaymentIntentStatusSucceeded:
		return model.PaymentStatusSucceeded
	case stripe.PaymentIntentStatusCanceled:
		return model.PaymentStatusCanceled
	default:
		return model.PaymentStatusFailed
	}
}

// mapRefundReasonToStripe 将内部退款原因映射为Stripe退款原因
func mapRefundReasonToStripe(reason string) string {
	switch reason {
	case "requested_by_customer":
		return string(stripe.RefundReasonRequestedByCustomer)
	case "duplicate":
		return string(stripe.RefundReasonDuplicate)
	case "fraudulent":
		return string(stripe.RefundReasonFraudulent)
	default:
		return string(stripe.RefundReasonRequestedByCustomer)
	}
}
