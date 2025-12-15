package model

import (
	"github.com/yb2020/odoc/pkg/model"
)

// PaperPdfParsed 存储PDF文件在OSS中的记录信息
type PaperPdfParsed struct {
	model.BaseModel // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	// 原始PDF文件SHA256
	SourcePdfSHA256 string `json:"sourcePdfSHA256" gorm:"column:source_pdf_sha256;type:varchar(64);not null;index;comment:原始PDF文件SHA256"`
	// 文件SHA256
	FileSHA256 string `json:"fileSHA256" gorm:"column:file_sha256;type:varchar(64);not null;comment:文件SHA256"`
	// 文件类型
	FileType string `json:"fileType" gorm:"column:file_type;type:varchar(32);not null;comment:文件类型"`
	// 文件名
	FileName string `json:"fileName" gorm:"column:file_name;type:varchar(255);not null;comment:文件名"`
	// 文件大小(字节)
	FileSize int64 `json:"fileSize" gorm:"column:file_size;type:bigint;not null;comment:文件大小(字节)"`
	// 文件在OSS中的对象键
	ObjectKey string `json:"objectKey" gorm:"column:object_key;type:varchar(255);not null;comment:文件在OSS中的对象键"`
	// 文件所在的OSS桶名称
	BucketName string `json:"bucketName" gorm:"column:bucket_name;type:varchar(64);not null;comment:文件所在的OSS桶名称"`
	// 图表的oss信息记录json字符串
	RecordsJson string `json:"recordsJson" gorm:"column:records_json;type:text;comment:图表的oss信息记录json字符串"`
	// 版本
	Version string `json:"version" gorm:"column:version;type:varchar(64);not null;comment:版本"`
}

// TableName 设置表名
func (PaperPdfParsed) TableName() string {
	return "t_paper_pdf_parsed"
}
