package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// DocClassifyRelationDAO GORM实现的文档分类关系DAO
type DocClassifyRelationDAO struct {
	*baseDao.GormBaseDAO[model.DocClassifyRelation]
	logger logging.Logger
}

// NewDocClassifyRelationDAO 创建一个新的文档分类关系DAO
func NewDocClassifyRelationDAO(db *gorm.DB, logger logging.Logger) *DocClassifyRelationDAO {
	return &DocClassifyRelationDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.DocClassifyRelation](db, logger),
		logger:      logger,
	}
}

// GetDocClassifyRelationsByUserID 根据用户ID获取文档分类关系
// 只返回未删除的分类关系
func (d *DocClassifyRelationDAO) GetDocClassifyRelationsByUserID(ctx context.Context, userID string) ([]model.DocClassifyRelation, error) {
	var docClassifyRelations []model.DocClassifyRelation
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false", userID).Find(&docClassifyRelations)
	if result.Error != nil {
		d.logger.Error("msg", "获取文档分类关系失败", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return docClassifyRelations, nil
}

// DeleteByClassifyId 根据分类ID删除文档分类关系
func (d *DocClassifyRelationDAO) DeleteByClassifyId(ctx context.Context, classifyId string) error {
	result := d.GetDB(ctx).Where("classify_id = ?", classifyId).Delete(&model.DocClassifyRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "删除文档分类关系失败", "classifyId", classifyId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

func (d *DocClassifyRelationDAO) DeleteByUserIdDocIdClassifyId(ctx context.Context, userId string, docId string, classifyId string) error {
	result := d.GetDB(ctx).Where("user_id = ? and doc_id = ? and classify_id = ?", userId, docId, classifyId).Delete(&model.DocClassifyRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "删除文档分类关系失败", "userId", userId, "docId", docId, "classifyId", classifyId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

func (d *DocClassifyRelationDAO) GetByClassifyIdAndDocId(ctx context.Context, classifyId string, docId string) ([]model.DocClassifyRelation, error) {
	var docClassifyRelations []model.DocClassifyRelation
	result := d.GetDB(ctx).Where("classify_id = ? and doc_id = ?", classifyId, docId).Find(&docClassifyRelations)
	if result.Error != nil {
		d.logger.Error("msg", "获取文档分类关系失败", "classifyId", classifyId, "docId", docId, "error", result.Error.Error())
		return nil, result.Error
	}
	return docClassifyRelations, nil
}
