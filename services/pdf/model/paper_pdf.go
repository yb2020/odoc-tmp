package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperPdf 论文PDF实体
type PaperPdf struct {
	model.BaseModel        // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	AppId           string `json:"appId" gorm:"column:app_id;type:varchar(36);comment:应用ID"`                       // 应用ID
	PaperId         string `json:"paperId" gorm:"column:paper_id;size:36;index;comment:论文ID"`                      // 论文ID
	FileSHA256      string `json:"fileSHA256" gorm:"column:file_sha256;index;type:varchar(64);comment:SHA256值"`    // SHA256值
	ParseCount      int    `json:"parseCount" gorm:"column:parse_count;type:int;comment:解析次数"`                     // 解析次数
	Size            int64  `json:"size" gorm:"column:size;type:bigint;comment:PDF文件大小(bit)"`                       // PDF文件大小(bit)
	PageCount       int    `json:"pageCount" gorm:"column:page_count;type:int;comment:页数"`                         // 页数
	Language        string `json:"language" gorm:"column:language;type:varchar(10);comment:语言"`                    // 语言
	OssBucketName   string `json:"ossBucketName" gorm:"column:oss_bucket_name;type:varchar(100);comment:OSS存储桶名称"` // OSS存储桶名称
	OssObjectKey    string `json:"ossObjectKey" gorm:"column:oss_object_key;type:varchar(255);comment:OSS对象名称"`    // OSS对象名称
}

// TableName 返回表名
func (PaperPdf) TableName() string {
	return "t_paper_pdf"
}
