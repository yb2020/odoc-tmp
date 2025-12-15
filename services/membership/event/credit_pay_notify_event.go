package event

import "github.com/yb2020/odoc/pkg/eventbus"

// 支付模块下的事件类型
const (
	// 预扣成功
	CreditPayNotifyEvent_PrePaySuccess eventbus.EventType = "credit.pay.notify.pre_pay_success"
	// 预扣失败
	CreditPayNotifyEvent_PrePayFailed eventbus.EventType = "credit.pay.notify.pre_pay_failed"
	// 确认支付
	CreditPayNotifyEvent_PayConfirm eventbus.EventType = "credit.pay.notify.pay_confirm"
	// 回滚积分
	CreditPayNotifyEvent_PayRetrieve eventbus.EventType = "credit.pay.notify.pay_retrieve"
)

// CreditPayNotifyEvent 会员积分支付通知事件
type CreditPayNotifyEvent struct {
	EventType    eventbus.EventType // 事件类型
	RecordId     string             // 支付订单ID
	UserId       string             // 用户ID
	BillId       string             // 支付流水ID
	ErrorCode    string             // 错误码
	ErrorMessage string             // 错误信息
}
