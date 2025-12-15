package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// UserDocFolderRelationDAO GORM实现的用户文档文件夹关系DAO
type UserDocFolderRelationDAO struct {
	*baseDao.GormBaseDAO[model.UserDocFolderRelation]
	logger logging.Logger
}

// NewUserDocFolderRelationDAO 创建一个新的用户文档文件夹关系DAO
func NewUserDocFolderRelationDAO(db *gorm.DB, logger logging.Logger) *UserDocFolderRelationDAO {
	return &UserDocFolderRelationDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.UserDocFolderRelation](db, logger),
		logger:      logger,
	}
}

// SaveBatch 批量保存用户文档文件夹关系
func (d *UserDocFolderRelationDAO) SaveBatch(ctx context.Context, relations []model.UserDocFolderRelation) error {
	if len(relations) == 0 {
		return nil
	}
	result := d.GetDB(ctx).Create(&relations)
	if result.Error != nil {
		d.logger.Error("msg", "batch save user doc folder relation failed", "count", len(relations), "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// GetUserDocFolderRelationsByUserID 根据用户ID获取用户文档分类
// 只返回未删除的分类
func (d *UserDocFolderRelationDAO) GetUserDocFolderRelationsByUserID(ctx context.Context, userID string) ([]model.UserDocFolderRelation, error) {
	var docFolderRelations []model.UserDocFolderRelation
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false", userID).Find(&docFolderRelations)
	if result.Error != nil {
		d.logger.Error("msg", "get user doc folder relations by user id failed", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return docFolderRelations, nil
}

// GetUserDocFolderRelationsByFolderId 根据用户ID和文件夹ID获取文件夹-文献关系列表
// 只返回未删除的分类
func (d *UserDocFolderRelationDAO) GetUserDocFolderRelationsByFolderId(ctx context.Context, userId string, folderId string) ([]model.UserDocFolderRelation, error) {
	var docFolderRelations []model.UserDocFolderRelation
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false and folder_id = ?", userId, folderId).Find(&docFolderRelations)
	if result.Error != nil {
		d.logger.Error("msg", "get user doc folder relation failed", "userId", userId, "folderId", folderId, "error", result.Error.Error())
		return nil, result.Error
	}
	return docFolderRelations, nil
}

func (d *UserDocFolderRelationDAO) GetUserDocFolderRelationsByFolderIdAndDocId(ctx context.Context, userId string, folderId string, docId string) (*model.UserDocFolderRelation, error) {
	var docFolderRelations []model.UserDocFolderRelation
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false and folder_id = ? and doc_id = ?", userId, folderId, docId).Find(&docFolderRelations)
	if result.Error != nil {
		d.logger.Error("msg", "get user doc folder relation failed", "userId", userId, "folderId", folderId, "docId", docId, "error", result.Error.Error())
		return nil, result.Error
	}
	if len(docFolderRelations) == 0 {
		return nil, nil // 没有找到记录，返回nil而不是错误
	}
	return &docFolderRelations[0], nil
}

// GetUserDocFolderRelationsByFolderIDs 根据文件夹ID列表获取用户文档文件夹关系列表
func (d *UserDocFolderRelationDAO) GetUserDocFolderRelationsByFolderIDs(ctx context.Context, folderIDs []string) ([]model.UserDocFolderRelation, error) {
	if len(folderIDs) == 0 {
		return []model.UserDocFolderRelation{}, nil
	}

	var relations []model.UserDocFolderRelation
	result := d.GetDB(ctx).Where("folder_id IN ? and is_deleted = false", folderIDs).Find(&relations)
	if result.Error != nil {
		d.logger.Error("msg", "get user doc folder relations by folder ids failed", "folderIDs", folderIDs, "error", result.Error.Error())
		return nil, result.Error
	}

	return relations, nil
}

// GetMaxSort 获取指定用户和文件夹下的最大排序值
func (d *UserDocFolderRelationDAO) GetMaxSort(ctx context.Context, userID string, folderID string) (int32, error) {
	var relation model.UserDocFolderRelation
	result := d.GetDB(ctx).Where("user_id = ? AND folder_id = ? AND is_deleted = false", userID, folderID).
		Order("sort DESC").Limit(1).First(&relation)

	// 如果没有找到记录，返回0作为起始排序值
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, nil
		}
		d.logger.Error("msg", "get max sort failed", "userID", userID, "folderID", folderID, "error", result.Error.Error())
		return 0, result.Error
	}

	return relation.Sort, nil
}

// BatchRemoveByFolderIds 根据文件夹ID列表物理删除用户文档文件夹关系
func (d *UserDocFolderRelationDAO) BatchRemoveByFolderIds(ctx context.Context, folderIds []string) error {
	if len(folderIds) == 0 {
		return nil
	}

	result := d.GetDB(ctx).Where("folder_id IN ?", folderIds).Delete(&model.UserDocFolderRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "batch remove user doc folder relation failed", "folderIds", folderIds, "error", result.Error.Error())
		return result.Error
	}

	return nil
}

func (d *UserDocFolderRelationDAO) GetRelationsByUserIdAndDocIds(ctx context.Context, userId string, docIds []string) ([]model.UserDocFolderRelation, error) {
	if len(docIds) == 0 {
		return nil, nil
	}

	var relations []model.UserDocFolderRelation
	result := d.GetDB(ctx).Where("user_id = ? AND doc_id IN ?", userId, docIds).Find(&relations)
	if result.Error != nil {
		d.logger.Error("msg", "get relations by doc ids failed", "docIds", docIds, "error", result.Error.Error())
		return nil, result.Error
	}

	return relations, nil
}

// RemoveByUserIdAndFolderIdAndDocId 根据用户ID、文件夹ID和文献ID删除关系
func (d *UserDocFolderRelationDAO) RemoveByUserIdAndFolderIdAndDocId(ctx context.Context, userId string, folderId string, docId string) error {
	result := d.GetDB(ctx).Where("user_id = ? AND folder_id = ? AND doc_id = ?", userId, folderId, docId).Delete(&model.UserDocFolderRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "remove user doc folder relation failed", "userId", userId, "folderId", folderId, "docId", docId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// RemoveByUserIdAndFolderIdsAndDocId 根据用户ID、文件夹ID列表和文献ID删除关系
func (d *UserDocFolderRelationDAO) RemoveByUserIdAndFolderIdsAndDocId(ctx context.Context, userId string, folderIds []string, docId string) error {
	result := d.GetDB(ctx).Where("user_id = ? AND folder_id IN ? AND doc_id = ?", userId, folderIds, docId).Delete(&model.UserDocFolderRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "remove user doc folder relation failed", "userId", userId, "folderIds", folderIds, "docId", docId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

func (d *UserDocFolderRelationDAO) RemoveByUserIdAndDocId(ctx context.Context, userId string, docId string) error {
	result := d.GetDB(ctx).Where("user_id = ? AND doc_id = ?", userId, docId).Delete(&model.UserDocFolderRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "remove user doc folder relation failed", "userId", userId, "docId", docId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// FindByUserIdAndFolderId 根据用户ID和文件夹ID查询关系
func (d *UserDocFolderRelationDAO) FindByUserIdAndFolderId(ctx context.Context, userId string, folderId string) ([]model.UserDocFolderRelation, error) {
	var relations []model.UserDocFolderRelation
	result := d.GetDB(ctx).Where("user_id = ? AND folder_id = ? AND is_deleted = false", userId, folderId).Find(&relations)
	if result.Error != nil {
		d.logger.Error("msg", "find by user id and folder id failed", "userId", userId, "folderId", folderId, "error", result.Error.Error())
		return nil, result.Error
	}

	return relations, nil
}

// DeleteByUserIdAndFolderIdAndDocIds 根据用户ID、文件夹ID和文献ID列表删除关系
func (d *UserDocFolderRelationDAO) DeleteByUserIdAndFolderIdAndDocIds(ctx context.Context, userId string, folderId string, docIds []string) error {
	if len(docIds) == 0 {
		return nil
	}

	result := d.GetDB(ctx).Where("user_id = ? AND folder_id = ? AND doc_id IN ?", userId, folderId, docIds).Delete(&model.UserDocFolderRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "delete by user id and folder id and doc ids failed", "userId", userId, "folderId", folderId, "docIds", docIds, "error", result.Error.Error())
		return result.Error
	}

	return nil
}

// Save 保存单个用户文档文件夹关系
func (d *UserDocFolderRelationDAO) Save(ctx context.Context, relation *model.UserDocFolderRelation) error {
	result := d.GetDB(ctx).Create(relation)
	if result.Error != nil {
		d.logger.Error("msg", "save user doc folder relation failed", "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// MoveRelationsToUnclassifiedByFolderIds 根据文件夹ID列表移动关系到未分类
func (d *UserDocFolderRelationDAO) MoveRelationsToUnclassifiedByFolderIds(ctx context.Context, folderIds []string) error {
	if len(folderIds) == 0 {
		return nil
	}
	result := d.GetDB(ctx).
		Model(&model.UserDocFolderRelation{}).
		Where("folder_id IN ?", folderIds).
		Update("folder_id", 0)
	if result.Error != nil {
		d.logger.Error("msg", "move relations to unclassified by folder ids failed", "folderIds", folderIds, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// RemoveByUserIdAndFolderIds 根据用户ID和文件夹ID列表删除关系（物理删除，仅限该用户）
func (d *UserDocFolderRelationDAO) RemoveByUserIdAndFolderIds(ctx context.Context, userId string, folderIds []string) error {
	if len(folderIds) == 0 {
		return nil
	}
	result := d.GetDB(ctx).Where("user_id = ? AND folder_id IN ?", userId, folderIds).Delete(&model.UserDocFolderRelation{})
	if result.Error != nil {
		d.logger.Error("msg", "remove by user id and folder ids failed", "userId", userId, "folderIds", folderIds, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
