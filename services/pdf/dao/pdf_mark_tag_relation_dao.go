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

// PdfMarkTagRelationDAO GORM实现的PDF标记标签关系DAO
type PdfMarkTagRelationDAO struct {
	*baseDao.GormBaseDAO[model.PdfMarkTagRelation]
	logger logging.Logger
}

// NewPdfMarkTagRelationDAO 创建一个新的PDF标记标签关系DAO
func NewPdfMarkTagRelationDAO(db *gorm.DB, logger logging.Logger) *PdfMarkTagRelationDAO {
	return &PdfMarkTagRelationDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfMarkTagRelation](db, logger),
		logger:      logger,
	}
}

// Create 创建PDF标记标签关系
func (d *PdfMarkTagRelationDAO) Create(ctx context.Context, relation *model.PdfMarkTagRelation) error {
	return d.GetDB(ctx).Create(relation).Error
}

// GetPdfMarkTagRelationByID 根据ID获取PDF标记标签关系
func (d *PdfMarkTagRelationDAO) GetPdfMarkTagRelationByID(ctx context.Context, id string) (*model.PdfMarkTagRelation, error) {
	var relation model.PdfMarkTagRelation
	result := d.GetDB(ctx).Where("id = ?", id).First(&relation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF标记标签关系失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &relation, nil
}

// GetPdfMarkTagRelationsByTagID 根据标签ID获取PDF标记标签关系列表
func (d *PdfMarkTagRelationDAO) GetPdfMarkTagRelationsByTagID(ctx context.Context, tagID string) ([]model.PdfMarkTagRelation, error) {
	var relations []model.PdfMarkTagRelation
	result := d.GetDB(ctx).Where("tag_id = ?", tagID).Find(&relations)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记标签关系列表失败", "tag_id", tagID, "error", result.Error.Error())
		return nil, result.Error
	}
	return relations, nil
}

// GetByMarkId 根据标记ID获取PDF标记标签关系列表
func (d *PdfMarkTagRelationDAO) GetByMarkId(ctx context.Context, markID string) ([]model.PdfMarkTagRelation, error) {
	var relations []model.PdfMarkTagRelation
	result := d.GetDB(ctx).Where("mark_id = ? and is_deleted = false", markID).Find(&relations)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记标签关系列表失败", "mark_id", markID, "error", result.Error.Error())
		return nil, result.Error
	}
	return relations, nil
}

// GetByMarkIds 根据标记ID列表获取PDF标记标签关系列表
func (d *PdfMarkTagRelationDAO) GetByMarkIds(ctx context.Context, markIds []string) ([]model.PdfMarkTagRelation, error) {
	var relations []model.PdfMarkTagRelation
	result := d.GetDB(ctx).Where("mark_id in ? and is_deleted = false", markIds).Find(&relations)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记标签关系列表失败", "mark_id", markIds, "error", result.Error.Error())
		return nil, result.Error
	}
	return relations, nil
}

// RemoveByTagIDAndMarkID 根据标签ID和标记ID删除PDF标记标签关系(物理删除)
func (d *PdfMarkTagRelationDAO) RemoveByTagIDAndMarkID(ctx context.Context, tagID, markID string) error {
	result := d.GetDB(ctx).Where("tag_id = ? AND mark_id = ?", tagID, markID).Delete(&model.PdfMarkTagRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "删除PDF标记标签关系失败", "tag_id", tagID, "mark_id", markID, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// ListPdfMarkTagRelations 列出PDF标记标签关系
func (d *PdfMarkTagRelationDAO) ListPdfMarkTagRelations(ctx context.Context, limit, offset int) ([]model.PdfMarkTagRelation, error) {
	var relations []model.PdfMarkTagRelation
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&relations)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF标记标签关系失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return relations, nil
}

// CountPdfMarkTagRelations 获取PDF标记标签关系总数
func (d *PdfMarkTagRelationDAO) CountPdfMarkTagRelations(ctx context.Context) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfMarkTagRelation{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记标签关系总数失败", "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}

// GetPdfMarkTagRelationsByMarkId 根据标记ID获取PDF标记标签关系列表
func (d *PdfMarkTagRelationDAO) GetByMarkIdAndTagId(ctx context.Context, markId string, tagId string) (*model.PdfMarkTagRelation, error) {
	var relation model.PdfMarkTagRelation
	result := d.GetDB(ctx).Where("mark_id = ? AND tag_id = ? AND is_deleted = false", markId, tagId).First(&relation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF标记标签关系失败", "markId", markId, "tagId", tagId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &relation, nil
}

// DeleteByTagId 根据标签ID删除PDF标记标签关系
func (d *PdfMarkTagRelationDAO) DeleteByTagId(ctx context.Context, tagId string) error {
	result := d.GetDB(ctx).Model(&model.PdfMarkTagRelation{}).Where("tag_id = ? AND is_deleted = false", tagId).Update("is_deleted", true)
	if result.Error != nil {
		d.logger.Error("msg", "删除PDF标记标签关系失败", "tag_id", tagId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
