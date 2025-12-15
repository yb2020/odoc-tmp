package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// FullTextTranslateFix 全文翻译修复记录
type FullTextTranslateFix struct {
	model.BaseModel
	NoteId                  string `json:"noteId" gorm:"column:note_id;index;not null"`
	ErrorFileObjectId       string `json:"errorFileObjectId" gorm:"column:error_file_object_id;type:varchar(200);not null"`
	ErrorFileBucketId       string `json:"errorFileBucketId" gorm:"column:error_file_bucket_id;type:varchar(200);not null"`
	TranslationFileObjectId string `json:"translationFileObjectId" gorm:"column:translation_file_object_id;type:varchar(200)"`
	TranslationFileBucketId string `json:"translationFileBucketId" gorm:"column:translation_file_bucket_id;type:varchar(200)"`
	Progress                int    `json:"progress" gorm:"column:progress;not null;default:0"`
}

// TableName 指定表名
func (FullTextTranslateFix) TableName() string {
	return "t_full_text_translate_fix"
}
