package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperResources 论文资源结构体
type PaperResources struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId         int64     `json:"paperId" gorm:"column:paper_id;index;not null"` // 论文ID
	PaperTitle      string    `json:"paperTitle" gorm:"column:paper_title"`          // 论文标题
	ResourceUrl     string    `json:"resourceUrl" gorm:"column:resource_url"`        // 资源URL
	ResourceTitle   string    `json:"resourceTitle" gorm:"column:resource_title"`    // 资源标题
}

// TableName 返回表名
func (PaperResources) TableName() string {
	return "t_paper_resources"
}
