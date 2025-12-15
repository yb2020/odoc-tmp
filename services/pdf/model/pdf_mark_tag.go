package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PdfMarkTag PDF标记标签实体
type PdfMarkTag struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	TagName           string  `json:"tagName" gorm:"column:tag_name;uniqueIndex:idx_creator_tag_idempotent"`  // 标签名称
	Idempotent        int64   `json:"idempotent" gorm:"column:idempotent;uniqueIndex:idx_creator_tag_idempotent"`  // 幂等标记：实体有效且存在时，为0。实体被删除时被赋值为删除时的时间戳
}

// TableName 返回表名
func (PdfMarkTag) TableName() string {
	return "t_pdf_mark_tag"
}
