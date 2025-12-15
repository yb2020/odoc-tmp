package provider

import (
	"context"

	"github.com/yb2020/odoc/services/pay/model"
)

// PaymentProvider 定义支付渠道必须实现的通用操作接口
// 这个接口隔离了具体支付渠道的实现细节，使PaymentService不直接依赖特定支付SDK
type PaymentProvider interface {
	// GetName 返回支付渠道名称
	GetName() string

	// CreateCharge 创建支付/扣款
	// params 包含创建支付所需的所有参数
	// 返回支付渠道的交易ID和客户端需要的信息（如Stripe的client_secret）
	CreateCharge(ctx context.Context, params *ChargeParams) (*ChargeResult, error)

	// HandleWebhook 解析并验证Webhook数据
	// requestData 是原始的webhook请求数据
	// signature 是webhook请求头中的签名（如果有）
	// 返回解析后的事件信息
	HandleWebhook(ctx context.Context, requestData []byte, signature string) (*WebhookEvent, error)

	// CreateRefund 发起退款
	// chargeId 是原始支付的ID
	// amount 是退款金额，以最小货币单位表示（例如：分）
	// reason 是退款原因
	// 返回退款的ID和状态
	CreateRefund(ctx context.Context, chargeId string, amount int64, reason string) (*RefundResult, error)

	// GetChargeStatus 查询支付状态
	// chargeId 是支付渠道的交易ID
	// 返回支付的当前状态
	GetChargeStatus(ctx context.Context, chargeId string) (string, error)
}

// ChargeParams 创建支付所需的参数
type ChargeParams struct {
	Amount          int64             // 金额，以最小货币单位表示（例如：分）
	Currency        string            // ISO 4217 货币代码 (例如 "CNY", "USD")
	Description     string            // 支付描述
	OrderId         string            // 关联的业务订单ID
	UserId          string            // 用户ID
	PaymentMethodId string            // 支付方式ID（如Stripe的PaymentMethodID）
	Metadata        map[string]string // 元数据，可用于存储自定义信息
	ReturnURL       string            // 支付完成后的返回URL（对于需要重定向的支付方式）
}

// ChargeResult 创建支付的结果
type ChargeResult struct {
	ProviderTxId      string            // 支付渠道的交易ID
	Status            string            // 支付状态
	ClientSecret      string            // 客户端密钥（如Stripe的client_secret）
	RedirectURL       string            // 重定向URL（对于需要重定向的支付方式）
	RequiresAction    bool              // 是否需要额外的客户操作
	PaymentMethodType string            // 支付方式类型
	ErrorCode         string            // 错误码（如果有）
	ErrorMessage      string            // 错误信息（如果有）
	Metadata          map[string]string // 支付渠道返回的元数据
}

// WebhookEvent Webhook事件信息
type WebhookEvent struct {
	EventType    string                 // 事件类型（如payment_intent.succeeded）
	ProviderTxId string                 // 支付渠道的交易ID
	Status       string                 // 支付状态
	Amount       int64                  // 金额
	Currency     string                 // 货币代码
	Metadata     map[string]string      // 元数据
	RawData      map[string]interface{} // 原始事件数据
}

// RefundResult 退款结果
type RefundResult struct {
	RefundId     string // 退款ID
	Status       string // 退款状态
	ErrorCode    string // 错误码（如果有）
	ErrorMessage string // 错误信息（如果有）
}

// 支付状态常量，与model.PaymentRecord中的状态常量保持一致
const (
	PaymentStatusPending        = model.PaymentStatusPending
	PaymentStatusRequiresAction = model.PaymentStatusRequiresAction
	PaymentStatusSucceeded      = model.PaymentStatusSucceeded
	PaymentStatusFailed         = model.PaymentStatusFailed
	PaymentStatusCanceled       = model.PaymentStatusCanceled
	PaymentStatusRefunded       = model.PaymentStatusRefunded
)
