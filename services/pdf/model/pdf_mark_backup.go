package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PdfMarkBackup PDF标记备份实体
type PdfMarkBackup struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId           int64   `json:"paperId" gorm:"column:paper_id;index"`  // 论文ID
	NoteId            int64   `json:"noteId" gorm:"column:note_id;index"`    // 笔记ID
	Sort              int     `json:"sort" gorm:"column:sort"`               // 排序
}

// TableName 返回表名
func (PdfMarkBackup) TableName() string {
	return "t_pdf_mark_backup"
}
