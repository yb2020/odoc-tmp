package model

import (
	"time"

	"github.com/yb2020/odoc/pkg/model"
)

// OssRecord 文件记录表
type OssRecord struct {
	model.BaseModel        // 嵌入基础模型
	BucketName      string `gorm:"column:bucket_name;type:varchar(100);not null;index;comment:存储桶名称"`
	ObjectKey       string `gorm:"column:object_key;type:varchar(255);not null;index;comment:对象键（文件路径）"`
	FileName        string `gorm:"column:file_name;type:varchar(255);comment:原始文件名"`
	FileSize        int64  `gorm:"column:file_size;type:bigint;comment:文件大小（字节）"`
	FileSHA256      string `gorm:"column:file_sha256;type:varchar(100);comment:文件sha256值"`
	ContentType     string `gorm:"column:content_type;type:varchar(100);comment:文件类型"`
	HashValue       string `gorm:"column:hash_value;type:varchar(100);comment:文件哈希值"`
	// 状态控制
	Status     string `gorm:"column:status;type:varchar(20);not null;default:'pending';index;comment:文件状态"`
	Visibility string `gorm:"column:visibility;type:varchar(20);not null;default:'private';comment:可见性"`

	// 临时文件
	IsTemp   bool       `gorm:"column:is_temp;not null;default:false;index;comment:是否为临时文件"`
	ExpireAt *time.Time `gorm:"column:expire_at;index;comment:过期时间"`

	// 上传信息
	UploadIP  string `gorm:"column:upload_ip;type:varchar(50);comment:上传者IP"`
	UserAgent string `gorm:"column:user_agent;type:varchar(255);comment:用户代理"`

	// Callback and metadata
	CallbackTopic string `gorm:"column:callback_topic;type:varchar(255);comment:回调主题"`
	BizMetadata   string `gorm:"column:biz_metadata;type:text;comment:业务元数据"`
}

// TableName 返回表名
func (OssRecord) TableName() string {
	return "t_oss_upload_record"
}
