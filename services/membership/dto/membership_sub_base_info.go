package dto

import pb "github.com/yb2020/odoc/proto/gen/go/order"

// MembershipSubBaseInfo 会员订阅计划
type MembershipSubBaseInfo struct {
	Type          pb.OrderType `json:"type"`          // 订阅类型：1-Free会员订阅, 2-PRO会员订阅 3-PRO版附加积分包订阅
	Name          string       `json:"name"`          // 订阅名称
	IsFree        bool         `json:"isFree"`        // 是否为免费
	Credit        int64        `json:"credit"`        // 积分（单位：0.01个）
	AddOnCredit   int64        `json:"addOnCredit"`   // 附加积分（单位：0.01个）
	Duration      int          `json:"duration"`      // 有效期（单位：月）
	Price         int64        `json:"price"`         // 价格（单位：分）
	Currency      string       `json:"currency"`      // 货币单位
	StripePayMode string       `json:"stripePayMode"` // 支付模式
	StripePriceId string       `json:"stripePriceId"` // 价格ID
}
