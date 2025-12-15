package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperQuestion 论文问题结构体
type PaperQuestion struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId         int64     `json:"paperId" gorm:"column:paper_id;index;not null"` // 论文ID
	Content         string    `json:"content" gorm:"column:content;type:text"`       // 问题内容
	ViewCount       int64     `json:"viewCount" gorm:"column:view_count"`            // 查看次数
	PdfId           int64     `json:"pdfId" gorm:"column:pdf_id;index"`              // PDF文件ID
}

// TableName 返回表名
func (PaperQuestion) TableName() string {
	return "t_paper_question"
}
