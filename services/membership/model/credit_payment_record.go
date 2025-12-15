package model

import (
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// CreditPaymentRecord 记录用户使用积分支付服务的凭证
type CreditPaymentRecord struct {
	model.BaseModel // 嵌入基础模型 (Id, CreatedAt, UpdatedAt, IsDeleted etc.)

	UserId            string     `json:"userId" gorm:"column:user_id;index;not null;comment:用户ID"`                                              //用户ID
	MembershipId      string     `json:"membershipId" gorm:"column:membership_id;index;comment:会员ID"`                                           //会员ID
	ServiceType       int32      `json:"serviceType" gorm:"column:service_type;index;comment:服务类型"`                                             //服务类型dto.CreditFunType
	CreditType        int32      `json:"creditType" gorm:"column:credit_type;index;comment:积分类型"`                                               //积分类型dto.CreditType
	Credit            int64      `json:"credit" gorm:"column:credit;not null;comment:消耗的积分数量"`                                                  //消耗的积分数量
	Status            int32      `json:"status" gorm:"column:status;not null;index;comment:支付状态 (0:Unknown, 1:Pending, 2:Succeeded, 3:Failed)"` //支付状态
	Content           string     `json:"content" gorm:"column:content;varchar(255);comment:支付内容"`                                               //支付内容Content
	Remark            string     `json:"remark" gorm:"column:remark;varchar(255);comment:支付备注"`                                                 //支付备注Remark
	PayAt             *time.Time `json:"payAt" gorm:"column:pay_at;comment:支付时间"`                                                               //支付时间
	RelCreditBid      string     `json:"relCreditBid" gorm:"column:rel_credit_bid;index;comment:关联的积分流水ID"`                                     //关联到实际扣款的CreditBill
	ErrorCode         string     `json:"errorCode" gorm:"column:error_code;varchar(255);comment:支付失败错误码"`                                       //如果支付失败，记录错误码
	ErrorMessage      string     `json:"errorMessage" gorm:"column:error_message;varchar(255);comment:支付失败原因"`                                  //如果支付失败，记录原因
	TransactionAt     *time.Time `json:"transactionAt" gorm:"column:transaction_at;comment:交易发生时间（成功或失败的最终确认时间）"`                               //记录支付状态最终确定的时间
	ConfirmAt         *time.Time `json:"confirmAt" gorm:"column:confirm_at;comment:支付确认时间（支付成功或失败的确认时间）"`                                       //支付确认时间
	ConfirmExpiredAt  *time.Time `json:"confirmExpiredAt" gorm:"column:confirm_expired_at;comment:支付确认过期时间"`                                    //支付确认过期时间
	RetrieveCreditBid string     `json:"retrieveCreditBid" gorm:"column:retrieve_credit_bid;index;comment:回退积分关联的积分流水ID"`
}

// TableName 指定表名
func (CreditPaymentRecord) TableName() string {
	return "t_membership_credit_payment_record"
}
