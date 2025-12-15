package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// WordPronunciationLog 表示单词发音日志
type WordPronunciation struct {
	model.BaseModel
	TargetContent        string    `json:"target_content" gorm:"column:target_content;type:varchar(500);index:idx_content_source"`
	BritishSymbol        string    `json:"british_symbol" gorm:"column:british_symbol;type:varchar(500)"`
	AmericaSymbol        string    `json:"america_symbol" gorm:"column:america_symbol;type:varchar(500)"`
	BritishFormat        string    `json:"british_format" gorm:"column:british_format;type:varchar(500)"`
	BritishPronunciation string    `json:"british_pronunciation" gorm:"column:british_pronunciation;type:text"`
	AmericaFormat        string    `json:"america_format" gorm:"column:america_format;type:varchar(500)"`
	AmericaPronunciation string    `json:"america_pronunciation" gorm:"column:america_pronunciation;type:text"`
	Source               string    `json:"source" gorm:"column:source;type:varchar(50);index:idx_content_source;comment:翻译来源,baidu/youdao"`
	TargetResp           []WordExp `json:"target_resp" json_serialize:"true" json_field:"TargetRespJSON" gorm:"-"`
	TargetRespJSON       string    `json:"-" gorm:"column:target_resp;type:json"`
}

// WordExp 表示单词释义
type WordExp struct {
	Part          string   `json:"part"`
	TargetContent []string `json:"target_content"`
}

// TableName 指定表名
func (WordPronunciation) TableName() string {
	return "t_word_pronunciation"
}
