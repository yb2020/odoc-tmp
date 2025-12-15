package oss

import (
	"context"
	"io"
	"time"
)

// ObjectInfo 对象信息
type ObjectInfo struct {
	Name         string            `json:"name"`          // 对象名称
	Size         int64             `json:"size"`          // 对象大小
	LastModified time.Time         `json:"last_modified"` // 最后修改时间
	ContentType  string            `json:"content_type"`  // 内容类型
	ETag         string            `json:"etag"`          // ETag
	Metadata     map[string]string `json:"metadata"`      // 元数据
}

// PreSignedUploadResponse 预签名上传响应
type PreSignedUploadResponse struct {
	// 是否需要上传
	NeedUpload bool `json:"need_upload"`
	// 上传信息
	UploadInfo *PreSignedUploadInfo `json:"upload_info,omitempty"`
	// OSS信息
	OssInfo *OssInfo `json:"oss_info,omitempty"`
	//公共桶的永久地址
	PublicUrl string `json:"public_url,omitempty"`
}

// PreSignedUploadInfo 预签名上传信息
type PreSignedUploadInfo struct {
	// 预签名URL
	URL string `json:"url"` // 预签名URL
	// HTTP方法
	Method string `json:"method"` // HTTP方法（通常是 PUT）
	// HTTP头
	Headers map[string]string `json:"headers"` // 需要设置的HTTP头
}

// OssInfo OSS信息
type OssInfo struct {
	BucketName string `json:"bucket_name"` // 存储桶名称
	ObjectKey  string `json:"object_key"`  // 对象键（文件路径）
	FileName   string `json:"file_name"`   // 原始文件名
	FileSize   int64  `json:"file_size"`   // 文件大小（字节）
	FileSHA256 string `json:"file_sha256"` // 文件sha256值
}

// StorageInterface 存储接口
type StorageInterface interface {
	// Upload 上传对象
	Upload(ctx context.Context, bucket, objectName string, reader io.Reader, objectSize int64, opts map[string]string, userMetadata map[string]string) error

	// Download 下载对象
	Download(ctx context.Context, bucket, objectName string) (io.Reader, error)

	// Delete 删除对象
	Delete(ctx context.Context, bucket, objectName string) error

	// GetObjectURL 获取对象的临时URL
	GetObjectURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error)

	// GetPermanentURL 获取对象的永久URL（仅适用于公开桶）
	// 注意：只有公开桶中的对象才能使用此方法
	GetPermanentURL(ctx context.Context, bucket, objectName string) (string, error)

	// GetBucketConfig 获取桶配置
	GetBucketConfig(bucket string) (BucketConfig, error)

	// GeneratePreSignedUpload 生成预签名上传URL和相关信息
	// contentType: 文件的MIME类型
	// metadata: 用户自定义元数据，会在回调时返回
	GeneratePreSignedUpload(ctx context.Context, bucketType, objectName, contentType string, fileSize int64, metadata map[string]string) (*PreSignedUploadResponse, error)

	// Close 关闭存储客户端
	Close() error
}
