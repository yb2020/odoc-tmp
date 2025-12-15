package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PdfMark PDF标记实体
type PdfMark struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	Type            int    `json:"type" gorm:"column:type"`                      // IDEAAnnotateType
	PaperId         string `json:"paperId" gorm:"column:paper_id;size:36;index"` // 论文ID
	NoteId          string `json:"noteId" gorm:"column:note_id;size:36;index"`   // 笔记ID
	PdfId           string `json:"pdfId" gorm:"column:pdf_id;size:36;index"`     // PDFID(新字段)
	Content         string `json:"content" gorm:"column:content"`                // 内容:原存储标注JSON字段，作废
	Idea            string `json:"idea" gorm:"column:idea"`                      // IDEA
	HtmlIdea        string `json:"htmlIdea" gorm:"column:html_idea"`             // HTML格式的IDEA
	IsHighlight     bool   `json:"isHighlight" gorm:"column:is_highlight"`       // 是否高亮
	KeyContent      string `json:"keyContent" gorm:"column:key_content"`         // 重点内容
	StyleId         int    `json:"styleId" gorm:"column:style_id"`               // 样式ID
	PicUrl          string `json:"picUrl" gorm:"column:pic_url"`                 // 图片地址
	Sort            int    `json:"sort" gorm:"column:sort"`                      // 排序
	QuadMd5         string `json:"quadMd5" gorm:"column:quad_md5"`               // 四边形MD5
	Page            int    `json:"page" gorm:"column:page"`                      // 页码
	CommentContent  string `json:"commentContent" gorm:"column:comment_content"` // 注释标注JSON(新字段)
	RectContent     string `json:"rectContent" gorm:"column:rect_content"`       // 矩形标注JSON(新字段)
	TextBoxContent  string `json:"textBoxContent" gorm:"column:txt_box_content"` // 文本框标注JSON(新字段)
}

// TableName 返回表名
func (PdfMark) TableName() string {
	return "t_pdf_mark"
}
