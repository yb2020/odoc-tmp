package model

import (
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// DoiMetaInfo DOI元信息结构体
type DoiMetaInfo struct {
	model.BaseModel           // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	Language        string    `json:"language" gorm:"column:language"`
	Content         string    `json:"content" gorm:"column:content"`
	Title           string    `json:"title" gorm:"column:title"`
	DocType         string    `json:"docType" gorm:"column:doc_type"`
	AuthorList      string    `json:"authorList" gorm:"column:author_list"`
	PublishTime     time.Time `json:"publishTime" gorm:"column:publish_time;type:timestamptz;index"`
	IsParse         bool      `json:"isParse" gorm:"column:is_parse"`
	Doi             string    `json:"doi" gorm:"column:doi"`
	PaperId         string    `json:"paperId" gorm:"column:paper_id"`
	EventTitle      string    `json:"eventTitle" gorm:"column:event_title"`
	EventPlace      string    `json:"eventPlace" gorm:"column:event_place"`
	EventDate       string    `json:"eventDate" gorm:"column:event_date"`
	ContainerTitle  string    `json:"containerTitle" gorm:"column:container_title"`
	Page            string    `json:"page" gorm:"column:page"`
	Source          string    `json:"source" gorm:"column:source"`
	Volume          string    `json:"volume" gorm:"column:volume"`
	Issue           string    `json:"issue" gorm:"column:issue"`
	Url             string    `json:"url" gorm:"column:url"`
	IsCrawlMeta     bool      `json:"isCrawlMeta" gorm:"column:is_crawl_meta"`
}

func (DoiMetaInfo) TableName() string {
	return "t_doi_meta_info"
}
