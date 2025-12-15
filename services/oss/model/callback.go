package model

// OsscallbackLog OSS回调日志
type OsscallbackLog struct {
	BucketName   string `gorm:"column:bucket_name;type:varchar(100);not null;index;comment:存储桶名称" json:"bucketName"` // 桶名称
	ObjectKey    string `gorm:"column:object_key;type:varchar(255);not null;index;comment:对象键名" json:"objectKey"`    // 对象键名
	EventType    string `gorm:"column:event_type;type:varchar(50);not null;index;comment:事件类型" json:"-"`             // 事件类型，JSON序列化时忽略
	Size         int64  `gorm:"column:size;type:bigint;not null;comment:对象大小" json:"size"`                           // 对象大小
	UserMetadata string `gorm:"column:user_metadata;type:text;comment:用户元数据" json:"userMetadata"`                    // 用户元数据(JSON)
	RecordId     string `gorm:"column:record_id;type:varchar(100);index;comment:OSS记录ID" json:"recordId"`            // OSS记录ID
	UploadUserId string `gorm:"column:upload_user_id;type:varchar(100);index;comment:上传用户ID" json:"uploadUserId"`    // 上传用户ID
	TopicName    string `gorm:"column:topic_name;type:varchar(100);comment:主题名称" json:"-"`                           // 通知主题名称，JSON序列化时忽略
}
