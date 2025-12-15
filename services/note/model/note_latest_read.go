package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// NoteLatestRead 笔记最近阅读实体
type NoteLatestRead struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	NoteId          string `json:"noteId" gorm:"column:note_id;index"` // 笔记ID
}

// TableName 返回表名
func (NoteLatestRead) TableName() string {
	return "t_note_latest_read"
}
