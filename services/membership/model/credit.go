package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// Credit 会员积分实体 ER关系:会员1:1积分
type Credit struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	MembershipId    string `json:"membershipId" gorm:"column:membership_id;index;not null"`  // 会员ID
	UserId          string `json:"userId" gorm:"column:user_id;index;size:36;not null"`      // 用户ID
	Credit          int64  `json:"credit" gorm:"column:credit;index; default:0"`             // 积分（单位：0.01个）
	AddOnCredit     int64  `json:"addOnCredit" gorm:"column:add_on_credit;index; default:0"` // 附加积分（单位：0.01个）
}

// TableName 返回表名
func (Credit) TableName() string {
	return "t_membership_credit"
}
