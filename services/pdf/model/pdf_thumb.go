package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PdfThumb PDF缩略图实体
type PdfThumb struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId           int64   `json:"paperId" gorm:"column:paper_id;index"`  // 论文ID
	SourceUrl         string  `json:"sourceUrl" gorm:"column:source_url"`     // 源URL
	ThumbUrl          string  `json:"thumbUrl" gorm:"column:thumb_url"`       // 缩略图URL
	PdfId             int64   `json:"pdfId" gorm:"column:pdf_id;index"`       // PDF ID
}

// TableName 返回表名
func (PdfThumb) TableName() string {
	return "t_pdf_thumb"
}
