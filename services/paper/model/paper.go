package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// Paper 论文结构体   TODO: status Bibtex字段需要参考原Java代码，了解这个状态到底是什么
type Paper struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	PaperId         string `json:"paperId" gorm:"column:paper_id;index;size:36;comment:论文ID"`              // 论文ID
	OwnerId         string `json:"ownerId" gorm:"column:owner_id;index;not null;size:36;comment:拥有者ID"`    // 拥有者ID，pdf的上传人不一定是paper的拥有人
	Title           string `json:"title" gorm:"column:title;type:varchar(255);comment:论文标题"`               // 论文标题
	Abstract        string `json:"abstract" gorm:"column:abstract;type:text;comment:论文摘要"`                 // 论文摘要
	Authors         string `json:"authors" gorm:"column:authors;type:text;comment:论文作者"`                   // 论文作者
	Bibtex          string `json:"bibtex" gorm:"column:bibtex;type:text;comment:论文bibtex"`                 // 论文bibtex
	Status          string `json:"status" gorm:"column:status;type:varchar(10);comment:状态"`                // 状态：公开或非公开
	PublishDate     string `json:"publishDate" gorm:"column:publish_date;type:varchar(20);comment:论文发布日期"` // 论文发布日期
	ParseStatus     int    `json:"parseStatus" gorm:"column:parse_status;type:int;comment:解析状态"`           // 解析状态
}

// TableName 返回表名
func (Paper) TableName() string {
	return "t_paper"
}
