package model

import (
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// Order 会员订阅订单实体 ER关系:会员1:N订阅; 订单1:1订阅
type Order struct {
	model.BaseModel                // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId               string    `json:"userId" gorm:"column:user_id;index;not null"`                              // 用户ID
	MembershipId         string    `json:"membershipId" gorm:"column:membership_id;index;not null"`                  // 会员ID
	PayOrderId           string    `json:"payOrderId" gorm:"column:pay_order_id;index"`                              // 支付订单ID, 与支付订单表关联
	OrderStatus          int32     `json:"orderStatus" gorm:"column:order_status;index;not null"`                    // 订单状态（1:待支付，2:已支付，3:已完成 4:已取消）
	OrderType            int32     `json:"orderType" gorm:"column:order_type;index;not null"`                        // 订单类型（1:Free会员，2:PRO会员，3:PRO会员附加积分包）
	SubName              string    `json:"subName" gorm:"column:sub_name;index"`                                     // 订阅名称
	SubCredit            int64     `json:"subCredit" gorm:"column:sub_credit;index; default:0"`                      // 订阅积分（单位：0.01个）
	SubAddOnCredit       int64     `json:"subAddOnCredit" gorm:"column:sub_add_on_credit;index; default:0"`          // 订阅附加积分（单位：0.01个）
	SubStartDate         time.Time `json:"subStartDate" gorm:"column:sub_start_date;index"`                          // 订阅开始时间
	SubEndDate           time.Time `json:"subEndDate" gorm:"column:sub_end_date;index"`                              // 订阅结束时间
	Price                int64     `json:"price" gorm:"column:price;index; default:0"`                               // 产品单价（单位：分）
	NumberCount          int32     `json:"numberCount" gorm:"column:number_count;index; default:0"`                  // 订单数量
	TotalAmount          int64     `json:"totalAmount" gorm:"column:total_amount;index; default:0"`                  // 订单总金额（单位：分）
	PayAmount            int64     `json:"payAmount" gorm:"column:pay_amount;index; default:0"`                      // 实际支付金额（单位：分）
	IsDiscount           bool      `json:"isDiscount" gorm:"column:is_discount;index; default:false"`                // 是否使用折扣
	TotalDiscountAmount  int64     `json:"totalDiscountAmount" gorm:"column:total_discount_amount;index; default:0"` // 折扣总金额（单位：分）总扣减金额
	DiscountPercent      int64     `json:"discountPercent" gorm:"column:discount_percent;index; default:0"`          // 折扣率 0-100 8折=80
	Currency             string    `json:"currency" gorm:"column:currency;index"`                                    // 结算货币
	PayExepiredAt        time.Time `json:"payExepiredAt" gorm:"column:pay_exepired_at;index"`                        // 支付过期时间, 订阅状态为待支付时,支付超时,过期后订阅状态更新为已取消
	PayTime              time.Time `json:"payTime" gorm:"column:pay_time;index"`                                     // 支付时间
	StripePayMode        string    `json:"stripePayMode" gorm:"column:stripe_pay_mode;index"`                        // 支付模式 payment: 一次性付款, subscription: 订阅
	StripePriceId        string    `json:"stripePriceId" gorm:"column:stripe_price_id;index"`                        // 价格ID, subscription模式下必填
	StripeSubscriptionId string    `json:"stripeSubscriptionId" gorm:"column:stripe_subscription_id;index"`          // Stripe订阅ID, subscription模式下通知返回
}

// TableName 返回表名
func (Order) TableName() string {
	return "t_order"
}
