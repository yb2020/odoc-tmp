package oss

import (
	"github.com/yb2020/odoc/config"
)

// BucketConfig 存储桶配置
type BucketConfig struct {
	Name       string
	Public     bool
	Versioning bool
	IsTemp     bool
	Expiration int64 // 过期时间（秒）
}

// NewStorageInterface 创建存储接口实例
func NewStorageInterface(cfg *config.Config) (StorageInterface, error) {
	return NewS3Storage(cfg)
}

// NewStorage 创建存储实例（向后兼容）
// Deprecated: 使用 NewStorageInterface 代替
func NewStorage(cfg *config.Config) (StorageInterface, error) {
	return NewStorageInterface(cfg)
}
