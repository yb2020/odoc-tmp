package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PdfComment PDF评论实体
type PdfComment struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId           int64   `json:"paperId" gorm:"column:paper_id;index"`  // 论文ID
	Version           int64   `json:"version" gorm:"column:version"`         // 版本号
	VersionName       string  `json:"versionName" gorm:"column:version_name"` // 版本名称
	PdfUrl            string  `json:"pdfUrl" gorm:"column:pdf_url"`          // PDF URL
	Comment           string  `json:"comment" gorm:"column:comment"`         // 评论内容
	HtmlComment       string  `json:"htmlComment" gorm:"column:html_comment"` // HTML格式的评论内容
	Attr              string  `json:"attr" gorm:"column:attr"`               // 属性
	ParentId          int64   `json:"parentId" gorm:"column:parent_id;index"` // 父评论ID
}

// TableName 返回表名
func (PdfComment) TableName() string {
	return "t_pdf_comment"
}
