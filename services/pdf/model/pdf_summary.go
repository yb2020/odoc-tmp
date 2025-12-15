package model

import (
	"encoding/json"

	"github.com/yb2020/odoc/pkg/model"
)

// PdfSummary PDF 摘要实体
type PdfSummary struct {
	model.BaseModel              // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	SourceFrom          string   `json:"from" gorm:"column:source_from;varchar(20);not null;uniqueIndex:idx_from_sha256_version_lang"` // 来源, copilot, user
	ConversationId      string   `json:"conversationId" gorm:"column:conversation_id;varchar(100);"`                                   // 对话ID
	FileSHA256          string   `json:"file_sha256" gorm:"column:file_sha256;varchar(64);uniqueIndex:idx_from_sha256_version_lang"`   // PDF的SHA256
	Lang                string   `json:"lang" gorm:"column:lang;varchar(30);uniqueIndex:idx_from_sha256_version_lang"`                 // 语言
	Summary             string   `json:"summary" gorm:"column:summary;varchar(4000);"`                                                 // 总结
	Version             string   `json:"version" gorm:"column:version;varchar(20);uniqueIndex:idx_from_sha256_version_lang"`           // 版本
	RelatedQuestion     string   `json:"relatedQuestion" gorm:"column:related_question;varchar(2000);"`                                // 关联问题
	RelatedQuestionList []string `json:"relatedQuestionList" gorm:"-"`
}

// TableName 返回表名
func (PdfSummary) TableName() string {
	return "t_pdf_summary"
}

// SetRelatedQuestions set related questions
func (p *PdfSummary) SetRelatedQuestions(questions []string) error {
	if len(questions) == 0 {
		p.RelatedQuestion = "[]"
		return nil
	}
	data, err := json.Marshal(questions)
	if err != nil {
		return err
	}
	p.RelatedQuestion = string(data)
	return nil
}

func (p *PdfSummary) GetRelatedQuestions() ([]string, error) {
	var questions []string
	if p.RelatedQuestion == "" {
		return questions, nil
	}
	err := json.Unmarshal([]byte(p.RelatedQuestion), &questions)
	return questions, err
}
