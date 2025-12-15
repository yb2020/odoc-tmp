package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperAccess 描述论文是否可以被某个用户访问的结构体
type PaperAccess struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId         int64  `json:"paperId" gorm:"column:paper_id;index;not null"`  // 论文ID
	UserId          int64  `json:"userId" gorm:"column:user_id;index;not null"`    // 用户ID
}

// TableName 返回表名
func (PaperAccess) TableName() string {
	return "t_paper_access"
}
