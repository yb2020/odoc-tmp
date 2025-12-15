package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PdfAnnotation PDF注释实体
type PdfAnnotation struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PdfId           int64  `json:"pdfId" gorm:"column:pdf_id;index,type:bigint"`                      // PDF ID
	Page            int    `json:"page" gorm:"column:page,type:int"`                                  // 页码
	StyleId         int    `json:"styleId" gorm:"column:style_id,type:int"`                           // 样式ID
	Content         string `json:"content" gorm:"column:content,type:text"`                           // 内容
	Rectangles      string `json:"rectangles" gorm:"column:rectangles,type:text"`                     // 矩形区域
	AnnotationMd5   string `json:"annotationMd5" gorm:"column:annotation_md5;index,type:varchar(64)"` // 注释MD5
	Type            int    `json:"type" gorm:"column:type,type:int"`                                  // 类型
	PicUrl          string `json:"picUrl" gorm:"column:pic_url,type:varchar(255)"`                    // 图片URL
	RectStr         string `json:"rectStr" gorm:"column:rect_str,type:varchar(255)"`                  // 矩形字符串
}

// TableName 返回表名
func (PdfAnnotation) TableName() string {
	return "t_pdf_annotation"
}
