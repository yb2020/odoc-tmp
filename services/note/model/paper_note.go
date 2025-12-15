package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperNote 论文笔记实体
type PaperNote struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId         string `json:"paperId" gorm:"column:paper_id;size:36;index;not null"`                             // 论文ID
	PdfId           string `json:"pdfId" gorm:"column:pdf_id;size:36;index;not null;default:''"`                      // PDF文件ID
	NoteCount       uint32 `json:"noteCount" gorm:"column:note_count;default:0"`                                      // 笔记数量
	AnnotationPdfId string `json:"annotationPdfId" gorm:"column:annotation_pdf_id;size:36;index;not null;default:''"` // 标注PDFID
}

// TableName 返回表名
func (PaperNote) TableName() string {
	return "t_paper_note"
}
