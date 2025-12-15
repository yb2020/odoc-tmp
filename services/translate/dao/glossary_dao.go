package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/translate/model"
	"gorm.io/gorm"
)

// GlossaryDAO 术语库数据访问实现
type GlossaryDAO struct {
	*baseDao.GormBaseDAO[model.Glossary]
	logger logging.Logger
}

// NewGlossaryDAO 创建术语库数据访问对象
func NewGlossaryDAO(db *gorm.DB, logger logging.Logger) GlossaryDAO {
	return GlossaryDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.Glossary](db, logger),
		logger:      logger,
	}
}

// GetGlossariesByUserID 获取用户的术语条目列表
func (d *GlossaryDAO) GetGlossariesByUserID(ctx context.Context, userID string) ([]model.Glossary, error) {
	var glossaries []model.Glossary
	result := d.GetDB(ctx).
		Where("user_id = ? AND is_deleted = false", userID).
		Find(&glossaries)

	if result.Error != nil {
		d.logger.Error("msg", "获取用户术语条目列表失败", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return glossaries, nil
}

// GetPublicGlossaries 获取公共术语条目列表
func (d *GlossaryDAO) GetPublicGlossaries(ctx context.Context) ([]model.Glossary, error) {
	var glossaries []model.Glossary
	result := d.GetDB(ctx).
		Where("is_public = ? AND is_deleted = false", true).
		Find(&glossaries)

	if result.Error != nil {
		d.logger.Error("msg", "获取公共术语条目列表失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return glossaries, nil
}

// GetGlossaryCountByUserID 获取用户的术语条目数量
func (d *GlossaryDAO) GetGlossaryCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	result := d.GetDB(ctx).
		Model(&model.Glossary{}).
		Where("user_id = ? AND is_deleted = false", userID).
		Count(&count)

	if result.Error != nil {
		d.logger.Error("msg", "获取用户术语条目数量失败", "userID", userID, "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// GetGlossariesByIDs 根据ID列表获取术语条目
func (d *GlossaryDAO) GetGlossariesByIDs(ctx context.Context, ids []int64) ([]model.Glossary, error) {
	var glossaries []model.Glossary
	result := d.GetDB(ctx).
		Where("id IN ? AND is_deleted = false", ids).
		Find(&glossaries)

	if result.Error != nil {
		d.logger.Error("msg", "根据ID列表获取术语条目失败", "ids", ids, "error", result.Error.Error())
		return nil, result.Error
	}
	return glossaries, nil
}

// SearchGlossaries 搜索术语条目
func (d *GlossaryDAO) SearchGlossaries(ctx context.Context, userID string, searchText string) ([]model.Glossary, error) {
	var glossaries []model.Glossary

	// 构建模糊查询条件
	likeQuery := "%" + searchText + "%"

	// 构建查询条件：用户自己的术语条目或公共术语条目，同时匹配搜索词
	result := d.GetDB(ctx).
		Where("(user_id = ? OR is_public = ?) AND is_deleted = false AND original_text LIKE ?", userID, true, likeQuery).
		Find(&glossaries)

	if result.Error != nil {
		d.logger.Error("msg", "搜索术语条目失败", "userID", userID, "searchText", searchText, "error", result.Error.Error())
		return nil, result.Error
	}
	return glossaries, nil
}
