package model

import (
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// PaymentSubscription represents a payment transaction in the system.
// It stores details about the payment, its status, and associated metadata.
type PaymentSubscription struct {
	model.BaseModel                   // Embedded BaseModel
	UserId                 string     `json:"user_id" gorm:"column:user_id;type:varchar(255);index;not null"`                                       // 用户ID
	Status                 string     `json:"status" gorm:"column:status;type:varchar(50);index;not null"`                                          // 支付状态 (使用下方定义的常量)
	Description            string     `json:"description,omitempty" gorm:"column:description;type:text"`                                            // 支付描述
	ProviderSubscriptionId string     `json:"provider_subscription_id" gorm:"column:provider_subscription_id;type:varchar(255);index;comment:订阅ID"` // 订阅ID
	PriceId                string     `json:"price_id" gorm:"column:price_id;type:varchar(255);index;comment:价格ID"`                                 // 价格ID
	StartAt                *time.Time `json:"start_at" gorm:"column:start_at;type:timestamptz"`                                                     // 订阅开始时间
	EndAt                  *time.Time `json:"end_at" gorm:"column:end_at;type:timestamptz"`                                                         // 订阅结束时间
	CancelAtPeriodEnd      bool       `json:"cancel_at_period_end" gorm:"column:cancel_at_period_end;default:false;comment:是否在周期末取消"`               // 标记是否设置了在周期末取消
	CancelAt               *time.Time `json:"cancel_at" gorm:"column:cancel_at;type:timestamptz"`                                                   // 订阅将在何时被取消（或已经取消）
	CancelReason           string     `json:"cancel_reason" gorm:"column:cancel_reason;type:varchar(255)"`                                          // 订阅取消原因
}

// TableName returns the database table name for the PaymentSubscription model.
func (PaymentSubscription) TableName() string {
	return "t_pay_payment_subscription"
}

// Constants for PaymentSubscription Status
const (
	PaymentSubscription_StatusActive   = "ACTIVE"   // 激活
	PaymentSubscription_StatusCanceled = "CANCELED" // 已取消
)
