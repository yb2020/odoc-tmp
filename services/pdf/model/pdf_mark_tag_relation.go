package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PdfMarkTagRelation PDF标记标签关系实体
type PdfMarkTagRelation struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	TagId           string `json:"tagId" gorm:"column:tag_id;index"`   // 标签ID
	MarkId          string `json:"markId" gorm:"column:mark_id;index"` // 标记ID
}

// TableName 返回表名
func (PdfMarkTagRelation) TableName() string {
	return "t_pdf_mark_tag_relation"
}
