package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// UserDocClassify 用户文档分类结构体
type UserDocClassify struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	AppId           string `json:"appId" gorm:"column:app_id"`                  // 应用ID
	UserId          string `json:"userId" gorm:"column:user_id;index;not null"` // 用户ID
	Sort            int32  `json:"sort" gorm:"column:sort"`                     // 排序
	Name            string `json:"name" gorm:"column:name"`                     // 分类名称
	Remark          string `json:"remark" gorm:"column:remark"`                 // 备注
}

// TableName 返回表名
func (UserDocClassify) TableName() string {
	return "t_user_doc_classify"
}
