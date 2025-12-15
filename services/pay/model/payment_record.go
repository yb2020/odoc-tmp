package model

import (
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// PaymentRecord represents a payment transaction in the system.
// It stores details about the payment, its status, and associated metadata.
type PaymentRecord struct {
	model.BaseModel                  // Embedded BaseModel
	UserId                 string    `json:"user_id" gorm:"column:user_id;type:varchar(255);index;not null"`                                       // 用户ID
	OrderId                string    `json:"order_id" gorm:"column:order_id;type:varchar(255);index;not null"`                                     // 关联的业务订单ID
	Amount                 int64     `json:"amount" gorm:"column:amount;not null"`                                                                 // 支付金额 (最小货币单位，分)
	Currency               string    `json:"currency" gorm:"column:currency;type:varchar(10);not null"`                                            // ISO 4217 货币代码 (例如 "CNY", "USD")
	Status                 string    `json:"status" gorm:"column:status;type:varchar(50);index;not null"`                                          // 支付状态 (使用下方定义的常量)
	ProviderTxId           string    `json:"provider_tx_id" gorm:"column:provider_tx_id;type:varchar(255);index"`                                  // 支付渠道返回的交易ID
	Channel                string    `json:"channel" gorm:"column:channel;type:varchar(50);index;not null"`                                        // 支付服务提供商 (例如 "STRIPE", "WECHAT_PAY", "ALIPAY")
	Description            string    `json:"description,omitempty" gorm:"column:description;type:text"`                                            // 支付描述
	Metadata               string    `json:"metadata,omitempty" gorm:"column:metadata;type:text"`                                                  // 存储自定义元数据
	PaymentMethodType      string    `json:"payment_method_type,omitempty" gorm:"column:payment_method_type;type:varchar(50)"`                     // 支付方式类型
	ProviderErrorCode      string    `json:"provider_error_code,omitempty" gorm:"column:provider_error_code;type:varchar(100)"`                    // 支付渠道返回的错误码
	ProviderErrorMessage   string    `json:"provider_error_message,omitempty" gorm:"column:provider_error_message;type:text"`                      // 支付渠道返回的错误信息
	ProviderSubscriptionId string    `json:"provider_subscription_id" gorm:"column:provider_subscription_id;type:varchar(255);index;comment:订阅ID"` // 订阅ID （订阅模式独有字段）
	InvoiceId              string    `json:"invoice_id" gorm:"column:invoice_id;type:varchar(255);index;comment:发票ID"`                             // 发票ID
	PayMode                string    `json:"pay_mode" gorm:"column:pay_mode;type:varchar(50);index;comment:支付模式"`                                  // 支付模式 payment: 一次性支付模式, subscription: 订阅支付模式
	PaidAt                 time.Time `json:"paid_at,omitempty" gorm:"column:paid_at;type:timestamptz"`                                             // 支付成功时间
}

// TableName returns the database table name for the PaymentRecord model.
func (PaymentRecord) TableName() string {
	return "t_pay_payment_record"
}

// Constants for PaymentRecord Status
// 这些常量定义了支付记录可能处于的各种状态。
const (
	PaymentStatusPending        = "PENDING"         // 待处理：支付已发起，但尚未最终完成（可能包括等待用户操作、渠道处理或银行确认等）
	PaymentStatusRequiresAction = "REQUIRES_ACTION" // 需要操作：支付需要用户进行额外操作（如3D安全验证）
	PaymentStatusSucceeded      = "SUCCEEDED"       // 成功：支付已成功完成
	PaymentStatusFailed         = "FAILED"          // 失败：支付未能成功完成
	PaymentStatusCanceled       = "CANCELED"        // 已取消：支付被用户或系统取消
	PaymentStatusRefunded       = "REFUNDED"        // 已退款：支付已全额退款
)

// Constants for PaymentRecord Channel
// 这些常量定义了支付记录可能通过的支付服务提供商。
const (
	PaymentChannelStripe = "STRIPE" // Stripe 支付
	// PaymentChannelWechatPay      = "WECHAT_PAY"      // 微信支付
	// PaymentChannelAlipay         = "ALIPAY"          // 支付宝
	// 可以根据需要添加更多支付渠道，例如:
	// PaymentChannelPayPal     = "PAYPAL"      // PayPal 支付
	// PaymentChannelApplePay   = "APPLE_PAY"   // Apple Pay
)

// Constants for PaymentRecord PayMode
// 这些常量定义了支付记录可能的支付模式。
const (
	PaymentModePayment      = "PAYMENT"      // 一次性支付模式
	PaymentModeSubscription = "SUBSCRIPTION" // 订阅支付模式
)
