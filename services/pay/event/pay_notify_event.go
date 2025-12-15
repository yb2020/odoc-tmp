package event

import (
	"time"

	"github.com/yb2020/odoc/pkg/eventbus"
)

// 支付模块下的事件类型
const (
	PayNotifyEvent_PaySuccess                   eventbus.EventType = "pay.notify.pay_success"
	PayNotifyEvent_PayFailed                    eventbus.EventType = "pay.notify.pay_failed"
	PayNotifyEvent_PayExpire                    eventbus.EventType = "pay.notify.pay_expire"
	PayNotifyEvent_CustomerSubscriptionsCreated eventbus.EventType = "pay.notify.customer.subscriptions.created"
	PayNotifyEvent_CustomerSubscriptionsUpdated eventbus.EventType = "pay.notify.customer.subscriptions.updated"
	PayNotifyEvent_CustomerSubscriptionsDeleted eventbus.EventType = "pay.notify.customer.subscriptions.deleted"
	PayNotifyEvent_InvoicePaymentSucceeded      eventbus.EventType = "pay.notify.invoice.payment_succeeded"
	PayNotifyEvent_InvoicePaymentFailed         eventbus.EventType = "pay.notify.invoice.payment_failed"
)

// 支付通知事件
type PayNotifyEvent struct {
	OrderId        string    `json:"order_id"`
	PayRecordId    string    `json:"pay_record_id"`
	SubscriptionId string    `json:"subscription_id"`
	InvoiceId      string    `json:"invoice_id"`
	PayMode        string    `json:"pay_mode"`
	UserId         string    `json:"user_id"`
	SubStartAt     time.Time `json:"sub_start_at"`
	SubEndAt       time.Time `json:"sub_end_at"`
}
