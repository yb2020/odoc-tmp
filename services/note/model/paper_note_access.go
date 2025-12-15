package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperNoteAccess 论文笔记访问记录实体
type PaperNoteAccess struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	NoteId          string `json:"noteId" gorm:"column:note_id;index;not null"` // 笔记ID
	OpenStatus      bool   `json:"openStatus" gorm:"column:open_status"`        // 开放状态
}

// TableName 返回表名
func (PaperNoteAccess) TableName() string {
	return "t_paper_note_access"
}
