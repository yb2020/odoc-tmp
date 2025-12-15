package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// NoteExportHistory 笔记导出历史实体
type NoteExportHistory struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	NoteId          string `json:"noteId" gorm:"column:note_id;index"`  // 笔记ID
	Version         int64  `json:"version" gorm:"column:version;index"` // 版本
	PdfUrl          string `json:"pdfUrl" gorm:"column:pdf_url"`        // PDF地址
	AppId           string `json:"appId" gorm:"column:app_id"`          // 应用ID
}

// TableName 返回表名
func (NoteExportHistory) TableName() string {
	return "t_note_export_history"
}
