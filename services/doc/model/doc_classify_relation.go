package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// DocClassifyRelation 文档分类关系结构体
type DocClassifyRelation struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId          string `json:"userId" gorm:"column:user_id;index;not null"` // 用户ID
	DocId           string `json:"docId" gorm:"column:doc_id;index;not null"`   // 文档ID
	ClassifyId      string `json:"classifyId" gorm:"column:classify_id;index"`  // 分类ID
	Sort            int32  `json:"sort" gorm:"column:sort"`                     // 排序
}

// TableName 返回表名
func (DocClassifyRelation) TableName() string {
	return "t_doc_classify_relation"
}
