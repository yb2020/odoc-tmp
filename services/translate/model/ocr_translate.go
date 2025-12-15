package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// OCRTranslateLog 表示OCR翻译日志
type OCRTranslateLog struct {
	model.BaseModel
	UserId         string `json:"userId" gorm:"column:user_id;index;not null"`
	RequestId      string `json:"requestId" gorm:"column:request_id;index;type:varchar(64);not null"`
	ImageBase64    string `json:"imageBase64" gorm:"column:image_base64;type:text;"`
	SourceLanguage string `json:"sourceLanguage" gorm:"column:source_language;type:varchar(100);not null"`
	TargetLanguage string `json:"targetLanguage" gorm:"column:target_language;type:varchar(100);not null"`
	OCRText        string `json:"ocrText" gorm:"column:ocr_text;type:text"`
	TranslatedText string `json:"translatedText" gorm:"column:translated_text;type:text"`
	Status         int    `json:"status" gorm:"column:status;not null"`  // 0: 失败, 1: 成功
	CostMS         int64  `json:"costMs" gorm:"column:cost_ms;not null"` // 耗时（毫秒）
}

// TableName 指定表名
func (OCRTranslateLog) TableName() string {
	return "t_ocr_translate_log"
}
