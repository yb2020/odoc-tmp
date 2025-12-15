package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// UserCslRelation 用户引用样式关联结构体
type UserCslRelation struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId          string `json:"userId" gorm:"column:user_id"` // 用户ID
	CslId           string `json:"cslId" gorm:"column:csl_id"`   // 引用样式ID
	Sort            int32  `json:"sort" gorm:"column:sort"`      // 排序顺序
}

// TableName 返回表名
func (UserCslRelation) TableName() string {
	return "t_user_csl_relation"
}
