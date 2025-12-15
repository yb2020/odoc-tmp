package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// NoteSummary 笔记摘要实体
type NoteSummary struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	NoteId          string `json:"noteId" gorm:"column:note_id;index"` // 笔记ID
	UserId          string `json:"userId" gorm:"column:user_id;index"` // 用户ID
	Content         string `json:"content" gorm:"column:content"`      // 内容
}

// TableName 返回表名
func (NoteSummary) TableName() string {
	return "t_note_summary"
}
