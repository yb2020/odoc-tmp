package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// Glossary 表示术语库
type Glossary struct {
	model.BaseModel
	UserId          string `json:"userId" gorm:"column:user_id;index;not null"`
	OriginalText    string `json:"name" gorm:"column:original_text;type:varchar(100);not null"`
	TranslationText string `json:"description" gorm:"column:translation_text;type:text"`
	MatchCase       bool   `json:"matchCase" gorm:"column:match_case;default:false"`
	Ignored         bool   `json:"ignored" gorm:"column:ignored;default:false"`
	IsPublic        bool   `json:"isPublic" gorm:"column:is_public;default:false"`
}

// TableName 指定表名
func (Glossary) TableName() string {
	return "t_glossary"
}

// GlossaryTranslateModel 术语库翻译处理Bean
type GlossaryTranslateModel struct {
	// 用户选择的文本
	UserSelectedText string `json:"userSelectedText"`
	// 替换后的文本
	UserSelectedTextAfterReplace string `json:"userSelectedTextAfterReplace"`
	// 原始翻译
	OriginalTranslation string `json:"originalTranslation"`
	// 替换后的翻译
	OriginalTranslationAfterReplace string `json:"originalTranslationAfterReplace"`
	// 关联信息
	RelInfos []RelInfo `json:"relInfos"`
}

// RelInfo 术语关联信息
type RelInfo struct {
	// 术语ID
	Id string `json:"id"`
	// 特殊替换文本的Key
	Key string `json:"key"`
	// 译文
	Translation string `json:"translation"`
	// 原文
	OriginalText string `json:"originalText"`
	// 是否区分大小写
	MatchCase bool `json:"matchCase"`
}
