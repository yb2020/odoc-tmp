package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperPdfSelectRecord 论文PDF选择记录实体
type PaperPdfSelectRecord struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId           int64   `json:"paperId" gorm:"column:paper_id;index"`  // 论文ID
	UserId            int64   `json:"userId" gorm:"column:user_id;index"`    // 用户ID
	PdfId             int64   `json:"pdfId" gorm:"column:pdf_id;index"`      // PDF ID
}

// TableName 返回表名
func (PaperPdfSelectRecord) TableName() string {
	return "t_paper_pdf_select_record"
}
