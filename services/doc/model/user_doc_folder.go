package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// UserDocFolder 用户文档文件夹结构体
type UserDocFolder struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId          string `json:"userId" gorm:"column:user_id;index;not null"` // 用户ID
	ParentId        string `json:"parentId" gorm:"column:parent_id"`            // 父文件夹ID
	Name            string `json:"name" gorm:"column:name"`                     // 文件夹名称
	Remark          string `json:"remark" gorm:"column:remark"`                 // 备注
	Sort            int32  `json:"sort" gorm:"column:sort"`                     // 排序
}

// TableName 返回表名
func (UserDocFolder) TableName() string {
	return "t_user_doc_folder"
}
