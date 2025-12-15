package model

import (
	"github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/services/paper/constants"
)

// PaperCommentApproval 论文评论赞同状态结构体
type PaperCommentApproval struct {
	model.BaseModel                // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperCommentId     int64       `json:"paperCommentId" gorm:"column:paper_comment_id;index;not null"` // 同行评审的ID
	ApprovalStatus     string      `json:"approvalStatus" gorm:"column:approval_status"`                 // 赞同状态
}

// TableName 返回表名
func (PaperCommentApproval) TableName() string {
	return "t_paper_comment_approval"
}

// GetApprovalStatus 获取赞同状态枚举
func (p *PaperCommentApproval) GetApprovalStatus() constants.ApprovalStatus {
	return constants.ApprovalStatus(p.ApprovalStatus)
}

// SetApprovalStatus 设置赞同状态枚举
func (p *PaperCommentApproval) SetApprovalStatus(status constants.ApprovalStatus) {
	p.ApprovalStatus = status.String()
}

// NewPaperCommentApproval 创建新的论文评论赞同状态实例
func NewPaperCommentApproval(paperCommentId int64, approvalStatus constants.ApprovalStatus) *PaperCommentApproval {
	return &PaperCommentApproval{
		PaperCommentId: paperCommentId,
		ApprovalStatus: approvalStatus.String(),
	}
}
