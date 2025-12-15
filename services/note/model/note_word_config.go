package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// NoteWordConfig 笔记单词配置实体
type NoteWordConfig struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	AppId           string `json:"appId" gorm:"column:app_id"`             // 应用ID
	Color           string `json:"color" gorm:"column:color"`              // 颜色
	DisplayMode     string `json:"displayMode" gorm:"column:display_mode"` // 显示模式
	NoteId          string `json:"noteId" gorm:"column:note_id;index"`     // 笔记ID
}

// TableName 返回表名
func (NoteWordConfig) TableName() string {
	return "t_note_word_config"
}
