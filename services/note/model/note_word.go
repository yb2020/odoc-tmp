package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// NoteWord 笔记单词实体
type NoteWord struct {
	model.BaseModel             // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	NoteId               string `json:"noteId" gorm:"column:note_id;index"`                       // 笔记ID
	UserId               string `json:"userId" gorm:"column:user_id;index"`                       // 用户ID
	Word                 string `json:"word" gorm:"column:word"`                                  // 单词
	Rectangle            string `json:"rectangle" gorm:"column:rectangle"`                        // 矩形区域
	TargetContent        string `json:"targetContent" gorm:"column:target_content"`               // 目标内容
	BritishSymbol        string `json:"britishSymbol" gorm:"column:british_symbol"`               // 英式音标
	AmericaSymbol        string `json:"americaSymbol" gorm:"column:america_symbol"`               // 美式音标
	BritishFormat        string `json:"britishFormat" gorm:"column:british_format"`               // 英式格式
	BritishPronunciation string `json:"britishPronunciation" gorm:"column:british_pronunciation"` // 英式发音
	AmericaFormat        string `json:"americaFormat" gorm:"column:america_format"`               // 美式格式
	AmericaPronunciation string `json:"americaPronunciation" gorm:"column:america_pronunciation"` // 美式发音
	TargetResp           string `json:"targetResp" gorm:"column:target_resp"`                     // 目标响应
}

// TableName 返回表名
func (NoteWord) TableName() string {
	return "t_note_word"
}
