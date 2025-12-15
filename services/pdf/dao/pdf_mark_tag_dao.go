package dao

import (
	"context"
	"errors"
	"fmt"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/model"
	"gorm.io/gorm"
)

// PdfMarkTagDAO GORM实现的PDF标记标签DAO
type PdfMarkTagDAO struct {
	*baseDao.GormBaseDAO[model.PdfMarkTag]
	logger logging.Logger
}

// NewPdfMarkTagDAO 创建一个新的PDF标记标签DAO
func NewPdfMarkTagDAO(db *gorm.DB, logger logging.Logger) *PdfMarkTagDAO {
	return &PdfMarkTagDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfMarkTag](db, logger),
		logger:      logger,
	}
}

// GetPdfMarkTagByID 根据ID获取PDF标记标签
func (d *PdfMarkTagDAO) GetPdfMarkTagByID(ctx context.Context, id string) (*model.PdfMarkTag, error) {
	var tag model.PdfMarkTag
	result := d.GetDB(ctx).Where("id = ?", id).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF标记标签失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &tag, nil
}

// GetPdfMarkTagByName 根据标签名称获取PDF标记标签
func (d *PdfMarkTagDAO) GetPdfMarkTagByName(ctx context.Context, tagName string, creatorId string) (*model.PdfMarkTag, error) {
	var tag model.PdfMarkTag
	result := d.GetDB(ctx).Where("tag_name = ? AND creator_id = ? AND idempotent = 0", tagName, creatorId).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF标记标签失败", "tag_name", tagName, "creator_id", creatorId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &tag, nil
}

// UpdatePdfMarkTag 更新PDF标记标签
func (d *PdfMarkTagDAO) UpdatePdfMarkTag(ctx context.Context, tag *model.PdfMarkTag) error {
	return d.Modify(ctx, tag)
}

// DeletePdfMarkTag 删除PDF标记标签
func (d *PdfMarkTagDAO) DeletePdfMarkTag(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// LogicalDeletePdfMarkTag 逻辑删除PDF标记标签（设置idempotent为当前时间戳）
func (d *PdfMarkTagDAO) LogicalDeletePdfMarkTag(ctx context.Context, id string) error {
	result := d.GetDB(ctx).Model(&model.PdfMarkTag{}).Where("id = ?", id).Update("idempotent", gorm.Expr("UNIX_TIMESTAMP()"))
	if result.Error != nil {
		d.logger.Error("msg", "逻辑删除PDF标记标签失败", "id", id, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// ListPdfMarkTags 列出PDF标记标签
func (d *PdfMarkTagDAO) ListPdfMarkTags(ctx context.Context, creatorId string, limit, offset int) ([]model.PdfMarkTag, error) {
	var tags []model.PdfMarkTag
	result := d.GetDB(ctx).Where("creator_id = ? AND idempotent = 0 AND is_deleted = false", creatorId).Offset(offset).Limit(limit).Order("id DESC").Find(&tags)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF标记标签失败", "creator_id", creatorId, "error", result.Error.Error())
		return nil, result.Error
	}
	return tags, nil
}

// GetTagsByUserId 根据用户ID获取PDF标记标签
func (d *PdfMarkTagDAO) GetTagsByUserId(ctx context.Context, userId string) ([]model.PdfMarkTag, error) {
	var tags []model.PdfMarkTag
	result := d.GetDB(ctx).Where("creator_id = ? AND idempotent = 0 AND is_deleted = false", userId).Order("id DESC").Find(&tags)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF标记标签失败", "creator_id", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return tags, nil
}

// CountPdfMarkTags 获取PDF标记标签总数
func (d *PdfMarkTagDAO) CountPdfMarkTags(ctx context.Context, creatorId string) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfMarkTag{}).Where("creator_id = ? AND idempotent = 0 AND is_deleted = false", creatorId).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记标签总数失败", "creator_id", creatorId, "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}

// GetByIds 根据多个ID查询PDF标记标签列表
func (d *PdfMarkTagDAO) GetByIds(ctx context.Context, ids []string) ([]model.PdfMarkTag, error) {
	var tags []model.PdfMarkTag
	result := d.GetDB(ctx).Where("id IN (?) AND is_deleted = false", ids).Find(&tags)
	if result.Error != nil {
		d.logger.Error("msg", "根据多个ID查询PDF标记标签列表", "ids", ids, "error", result.Error.Error())
		return nil, result.Error
	}
	return tags, nil
}
