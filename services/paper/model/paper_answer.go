package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperAnswer 论文问答答案结构体
type PaperAnswer struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	QuestionId      int64     `json:"questionId" gorm:"column:question_id;index;not null"` // 问题ID
	Raw             string    `json:"raw" gorm:"column:raw;type:text"`                     // 问答中答案的元数据
	RawType         string    `json:"rawType" gorm:"column:raw_type"`                      // 元数据类型，参考PaperConstant.AnswerRawType
	ReplyUserId     int64     `json:"replyUserId" gorm:"column:reply_user_id;index"`       // 回复用户ID
	ReplyAnswerId   int64     `json:"replyAnswerId" gorm:"column:reply_answer_id;index"`   // 回复答案ID
}

// TableName 返回表名
func (PaperAnswer) TableName() string {
	return "t_paper_answer"
}
