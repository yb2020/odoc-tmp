package model

import (
	"github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/services/paper/constants"
)

// PaperComment 论文评论结构体
type PaperComment struct {
	model.BaseModel                // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	CommentLevel         string    `json:"commentLevel" gorm:"column:comment_level"`                    // 评审等级
	PaperCommentContent  string    `json:"paperCommentContent" gorm:"column:paper_comment_content;type:text"` // 评审内容
	Anonymous            bool      `json:"anonymous" gorm:"column:anonymous"`                           // 是否匿名
	PeerReviewScoreJson  string    `json:"peerReviewScoreJson" gorm:"column:peer_review_score_json"`   // 便于扩展，增加或者减少维度，不需要修改代码
	PaperId              int64     `json:"paperId" gorm:"column:paper_id;index;not null"`              // 论文ID
	Sort                 int32     `json:"sort" gorm:"column:sort"`                                     // 排序
	ExtParam             string    `json:"extParam" gorm:"column:ext_param"`                            // 扩展字段
}

// TableName 返回表名
func (PaperComment) TableName() string {
	return "t_paper_comment"
}

// GetExtParam 获取扩展参数值
func (p *PaperComment) GetExtParam(key constants.ExtParmKeys) string {
	// 这里可以实现从JSON字符串中解析特定键的值
	// 简化实现，实际使用时应该使用JSON解析库
	// 例如：
	// var extParams map[string]string
	// json.Unmarshal([]byte(p.ExtParam), &extParams)
	// return extParams[key.String()]
	return ""
}

// SetExtParam 设置扩展参数值
func (p *PaperComment) SetExtParam(key constants.ExtParmKeys, value string) {
	// 这里可以实现向JSON字符串中添加或更新特定键的值
	// 简化实现，实际使用时应该使用JSON解析库
	// 例如：
	// var extParams map[string]string
	// if p.ExtParam != "" {
	//     json.Unmarshal([]byte(p.ExtParam), &extParams)
	// } else {
	//     extParams = make(map[string]string)
	// }
	// extParams[key.String()] = value
	// jsonBytes, _ := json.Marshal(extParams)
	// p.ExtParam = string(jsonBytes)
}
