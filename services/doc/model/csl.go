package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// Csl 引用样式结构体
type Csl struct {
	model.BaseModel          // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	ShortTitle       string  `json:"shortTitle" gorm:"column:short_title"`             // 短标题
	CustomShortTitle string  `json:"customShortTitle" gorm:"column:custom_short_title"` // 自定义短标题
	CustomDefineTitle string `json:"customDefineTitle" gorm:"column:custom_define_title"` // 自定义定义标题
	Title           string  `json:"title" gorm:"column:title"`                         // 标题
	FileUrl         string  `json:"fileUrl" gorm:"column:file_url"`                    // 文件URL
	Updated         string  `json:"updated" gorm:"column:updated"`                     // 更新时间
	IsDefault       bool    `json:"isDefault" gorm:"column:is_default"`                // 是否默认
	UseCount        int32   `json:"useCount" gorm:"column:use_count"`                  // 使用次数
}

// TableName 返回表名
func (Csl) TableName() string {
	return "t_csl"
}
