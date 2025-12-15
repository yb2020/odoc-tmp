package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// UserDocAttachment 用户文档附件结构体
type UserDocAttachment struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId          string `json:"userId" gorm:"column:user_id;index;not null"` // 用户ID
	DocId           string `json:"docId" gorm:"column:doc_id;index;not null"`   // 文档ID
	Size            int32  `json:"size" gorm:"column:size"`                     // 附件大小
	Name            string `json:"name" gorm:"column:name"`                     // 附件名称
	OssObjectName   string `json:"ossObjectName" gorm:"column:oss_object_name"` // OSS对象名称
	ContentType     string `json:"contentType" gorm:"column:content_type"`      // 内容类型
}

// TableName 返回表名
func (UserDocAttachment) TableName() string {
	return "t_user_doc_attachment"
}
