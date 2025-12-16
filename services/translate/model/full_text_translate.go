package model

import (
	"github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/proto/gen/go/translate"
)

// FullTextTranslate 全文翻译历史记录
type FullTextTranslate struct {
	model.BaseModel
	UserId           string                            `json:"userId" gorm:"column:user_id;index;not null"`
	FileSHA256       string                            `json:"fileSHA256" gorm:"column:file_sha256;index;not null"`
	SourcePdfId      string                            `json:"sourcePdfId" gorm:"column:source_pdf_id;index;not null"`
	DocName          string                            `json:"docName" gorm:"column:doc_name;varchar(4000);default:''"`
	Message          string                            `json:"message" gorm:"column:message;varchar(2000);default:''"`
	SourceLanguage   string                            `json:"sourceLanguage" gorm:"column:source_language;type:varchar(200);not null"`
	TargetLanguage   string                            `json:"targetLanguage" gorm:"column:target_language;type:varchar(200);not null"`
	TargetBucketName string                            `json:"targetBucketName" gorm:"column:target_bucket_name;type:varchar(1000);default:''"`
	TargetObjectKey  string                            `json:"targetObjectKey" gorm:"column:target_object_key;type:varchar(1000);default:''"`
	FlowNumber       string                            `json:"flowNumber" gorm:"column:flow_number;varchar(200);default:''"`
	Alignment        string                            `json:"alignment" gorm:"column:alignment;text;default:''"`
	Status           translate.FullTranslateFlowStatus `json:"status" gorm:"column:status;numeric;index;default:0"`
	SessionId        string                            `json:"sessionId" gorm:"column:session_id;index;default:0"`
}

// TableName 指定表名
func (FullTextTranslate) TableName() string {
	return "t_full_text_translate"
}
