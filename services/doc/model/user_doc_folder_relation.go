package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// UserDocFolderRelation 用户文档文件夹关系结构体
type UserDocFolderRelation struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId          string `json:"userId" gorm:"column:user_id;index;not null"` // 用户ID
	DocId           string `json:"docId" gorm:"column:doc_id;index;not null"`   // 文档ID
	FolderId        string `json:"folderId" gorm:"column:folder_id;index"`      // 文件夹ID
	Sort            int32  `json:"sort" gorm:"column:sort"`                     // 排序
}

// TableName 返回表名
func (UserDocFolderRelation) TableName() string {
	return "t_user_doc_folder_relation"
}
