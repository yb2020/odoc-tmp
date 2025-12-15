package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// CreditBill 会员积分流水账单
type CreditBill struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId          string `json:"userId" gorm:"column:user_id;index"`             // 用户ID
	MembershipId    string `json:"membershipId" gorm:"column:membership_id;index"` // 会员ID
	CreditId        string `json:"creditId" gorm:"column:credit_id;index"`         // 账号积分ID

	Type       int32 `json:"type" gorm:"column:type;index"`              // 流水类型 constant.CreditBillType
	CreditType int32 `json:"creditType" gorm:"column:credit_type;index"` // 积分类型 constant.CreditType
	InOutType  int32 `json:"inOutType" gorm:"column:in_out_type;index"`  // 收入/支出类型 1:收入 2:支出 constant.CreditBillType_InOutType

	Credit       int64 `json:"credit" gorm:"column:credit;index"`              // 变动积分（单位：0.01个）正数为增加，负数为减少
	BeforeCredit int64 `json:"beforeCredit" gorm:"column:before_credit;index"` // 变动前积分（单位：0.01个）
	AfterCredit  int64 `json:"afterCredit" gorm:"column:after_credit;index"`   // 变动后积分（单位：0.01个）

	AddOnCredit       int64 `json:"addOnCredit" gorm:"column:add_on_credit;index"`              // 变动附加积分（单位：0.01个）正数为增加，负数为减少
	BeforeAddOnCredit int64 `json:"beforeAddOnCredit" gorm:"column:before_add_on_credit;index"` // 变动前附加积分（单位：0.01个）
	AfterAddOnCredit  int64 `json:"afterAddOnCredit" gorm:"column:after_add_on_credit;index"`   // 变动后附加积分（单位：0.01个）

	Content string `json:"content" gorm:"column:content;index"` // 内容
	Remark  string `json:"remark" gorm:"column:remark;index"`   // 备注
}

// TableName 返回表名
func (CreditBill) TableName() string {
	return "t_membership_credit_bill"
}
