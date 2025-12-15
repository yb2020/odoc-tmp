package model

import (
	"github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/services/paper/constants"
)

// PaperAttachment 论文附件结构体
type PaperAttachment struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId         int64     `json:"paperId" gorm:"column:paper_id;index;not null"`  // 论文ID
	FileName        string    `json:"fileName" gorm:"column:file_name"`               // 文件名
	Type            int       `json:"type" gorm:"column:type"`                        // 附件类型，参考constants.PaperAttachType
	TopStamp        int64     `json:"topStamp" gorm:"column:top_stamp"`               // 置顶时间戳
	Url             string    `json:"url" gorm:"column:url"`                          // 文件URL
	FileSize        int64     `json:"fileSize" gorm:"column:file_size"`               // 文件大小（字节）
}

// TableName 返回表名
func (PaperAttachment) TableName() string {
	return "t_paper_attachment"
}

// GetAttachType 获取附件类型枚举
func (p *PaperAttachment) GetAttachType() constants.PaperAttachType {
	return constants.PaperAttachTypeFromValue(p.Type)
}

// SetAttachType 设置附件类型枚举
func (p *PaperAttachment) SetAttachType(attachType constants.PaperAttachType) {
	p.Type = attachType.Value()
}
