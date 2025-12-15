package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// NoteReadLocation 笔记阅读位置实体
type NoteReadLocation struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	NoteId          string `json:"noteId" gorm:"column:note_id;index"` // 笔记ID
	Location        string `json:"location" gorm:"column:location"`    // 位置信息
}

// TableName 返回表名
func (NoteReadLocation) TableName() string {
	return "t_note_read_location"
}
