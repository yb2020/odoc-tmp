package dao

import (
	"context"

	"gorm.io/gorm"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/oss/constant"
	"github.com/yb2020/odoc/services/oss/model"
)

// OssDAO oss访问对象
type OssDAO struct {
	*baseDao.GormBaseDAO[model.OssRecord]
	logger logging.Logger
}

// NewOssDAO NewOssDAO
func NewOssDAO(db *gorm.DB, logger logging.Logger) *OssDAO {
	return &OssDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.OssRecord](db, logger),
		logger:      logger,
	}
}

// GetPendingRecordByObjectKey 根据对象键获取状态为上传中的文件记录
func (d *OssDAO) GetPendingRecordByObjectKey(ctx context.Context, objectKey string) (*model.OssRecord, error) {
	var record model.OssRecord
	result := d.GetDB(ctx).Where("object_key = ? AND status = ?", objectKey, constant.FileStatusPending).First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 记录不存在
		}
		d.logger.Error("msg", "根据对象键获取上传中文件记录失败", "objectKey", objectKey, "error", result.Error.Error())
		return nil, result.Error
	}
	return &record, nil
}

// UpdateFileStatus 更新文件记录状态
func (d *OssDAO) UpdateFileStatus(ctx context.Context, id string, status string, fileSize int64) error {
	updates := map[string]interface{}{
		"status":    status,
		"file_size": fileSize,
	}

	result := d.GetDB(ctx).Model(&model.OssRecord{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		d.logger.Error("msg", "更新文件记录状态失败", "id", id, "status", status, "error", result.Error.Error())
		return result.Error
	}

	if result.RowsAffected == 0 {
		d.logger.Warn("msg", "未找到要更新的文件记录", "id", id)
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetSuccessRecordByFileSHA256 根据 fileSHA256 和成功状态获取文件记录
func (d *OssDAO) GetSuccessRecordByFileSHA256(ctx context.Context, fileSHA256 string) (*model.OssRecord, error) {
	var record model.OssRecord
	result := d.GetDB(ctx).Where("file_sha256 = ? AND status = ?", fileSHA256, constant.FileStatusSuccess).First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 记录不存在
		}
		d.logger.Error("msg", "根据 fileSHA256 获取成功文件记录失败", "fileSHA256", fileSHA256, "error", result.Error.Error())
		return nil, result.Error
	}
	return &record, nil
}

// GetRecordByID 根据 ID 获取文件记录
func (d *OssDAO) GetRecordByID(ctx context.Context, id string) (*model.OssRecord, error) {
	var record model.OssRecord
	result := d.GetDB(ctx).Where("id = ?", id).First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 记录不存在
		}
		d.logger.Error("msg", "根据 ID 获取文件记录失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &record, nil
}

// GetByObjectKey 根据对象键获取文件记录
func (d *OssDAO) GetByBucketNameAndObjectKey(ctx context.Context, bucketName, objectKey string) (*model.OssRecord, error) {
	var record model.OssRecord
	result := d.GetDB(ctx).Where("bucket_name = ? AND object_key = ?", bucketName, objectKey).First(&record)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 记录不存在
		}
		d.logger.Error("msg", "根据对象键获取文件记录失败", "objectKey", objectKey, "error", result.Error.Error())
		return nil, result.Error
	}
	return &record, nil
}
