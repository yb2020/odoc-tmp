package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// TextTranslateLog 表示翻译日志
type TextTranslateLog struct {
	model.BaseModel
	UserId         string `json:"userId" gorm:"column:user_id;index;"`
	RequestId      string `json:"requestId" gorm:"column:request_id;index;type:varchar(200);not null"`
	Channel        string `json:"channel" gorm:"column:channel;type:varchar(100);not null"`
	SourceLanguage string `json:"sourceLanguage" gorm:"column:source_language;type:varchar(100);not null"`
	TargetLanguage string `json:"targetLanguage" gorm:"column:target_language;type:varchar(100);not null"`
	SourceContent  string `json:"sourceContent" gorm:"column:source_content;type:text;not null"`
	TargetContent  string `json:"targetContent" gorm:"column:target_content;type:text"`
	NetworkContent string `json:"networkContent" gorm:"-"`
	Md5Hash        string `json:"md5Hash" gorm:"column:md5_hash;index;type:varchar(32);not null"`
	UseGlossary    bool   `json:"useGlossary" gorm:"column:use_glossary;default:false"`
	PdfId          string `json:"pdfId" gorm:"column:pdf_id;default:0"`
	FeedBackResult int    `json:"feedBackResult" gorm:"column:feed_back_result;type:int;default:0"` // 0: 未纠错, 1: 已纠错
	Status         int    `json:"status" gorm:"column:status;not null"`                             // 0: 失败, 1: 成功
	CostMs         int64  `json:"costMs" gorm:"column:cost_ms;not null"`                            // 耗时（毫秒）
	ExtParams      string `json:"extParams" gorm:"column:ext_params;type:varchar(4000)"`
}

// TableName 指定表名
func (TextTranslateLog) TableName() string {
	return "t_translate_log"
}
