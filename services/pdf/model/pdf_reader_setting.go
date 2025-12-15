package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PdfReaderSetting PDF阅读器设置实体
type PdfReaderSetting struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	ClientType        int     `json:"clientType" gorm:"column:client_type"`  // 客户端类型：0-网页端，1-Windows客户端，2-iPadOS客户端
	Setting           string  `json:"setting" gorm:"column:setting"`         // 设置信息
}

// TableName 返回表名
func (PdfReaderSetting) TableName() string {
	return "t_pdf_reader_setting"
}
