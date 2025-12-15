package model

import (
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// UserMembership 用户会员实体 ER关系:用户1:1会员
type UserMembership struct {
	model.BaseModel                // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId               string    `json:"userId" gorm:"column:user_id;index;not null"`                     // 用户ID
	Type                 int32     `json:"type" gorm:"column:type;index;not null"`                          // 会员类型 1-Free会员 2-PRO会员
	StripeSubscriptionId string    `json:"stripeSubscriptionId" gorm:"column:stripe_subscription_id;index"` // Stripe订阅ID,
	StartAt              time.Time `json:"startAt" gorm:"index"`                                            // 订阅开始时间
	EndAt                time.Time `json:"endAt" gorm:"index"`                                              // 订阅结束时间
}

// TableName 返回表名
func (UserMembership) TableName() string {
	return "t_user_membership"
}
