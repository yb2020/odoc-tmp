package dto

import (
	pb "github.com/yb2020/odoc-proto/gen/go/membership"
)

// CreditPayIntent 支付积分DTO
type CreditPayIntent struct {
	Type        pb.CreditPayType     `json:"type"`        // 流水类型
	CreditType  pb.CreditType        `json:"creditType"`  // 积分类型
	ServiceType pb.CreditServiceType `json:"serviceType"` // 服务功能类型
	Credit      int64                `json:"credit"`      // 变动积分（单位：0.01个）
	AddOnCredit int64                `json:"addOnCredit"` // 变动附加积分（单位：0.01个）
	Content     string               `json:"content"`     // 内容
	Remark      string               `json:"remark"`      // 备注
}
