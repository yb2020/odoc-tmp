package dao

import (
	"context"
	"errors"
	"time"

	userContext "github.com/yb2020/odoc/pkg/context"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// UserDocFolderDAO GORM实现的用户文档文件夹DAO
type UserDocFolderDAO struct {
	*baseDao.GormBaseDAO[model.UserDocFolder]
	logger logging.Logger
}

// NewUserDocFolderDAO 创建一个新的用户文档文件夹DAO
func NewUserDocFolderDAO(db *gorm.DB, logger logging.Logger) *UserDocFolderDAO {
	return &UserDocFolderDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.UserDocFolder](db, logger),
		logger:      logger,
	}
}

// GetUserDocFolderByID 根据ID获取用户文档文件夹
func (d *UserDocFolderDAO) GetUserDocFolderByID(ctx context.Context, id string) (*model.UserDocFolder, error) {
	var folder model.UserDocFolder
	result := d.GetDB(ctx).Where("id = ? and is_deleted = false", id).First(&folder)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户文档文件夹失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &folder, nil
}

// GetByIdAndUserId 根据ID获取用户文档文件夹
func (d *UserDocFolderDAO) GetByIdAndUserId(ctx context.Context, id string, userId string) (*model.UserDocFolder, error) {
	var folder model.UserDocFolder
	result := d.GetDB(ctx).Where("id = ? and user_id = ? and is_deleted = false", id, userId).First(&folder)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户文档文件夹失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &folder, nil
}

// GetUserDocFoldersByUserID 根据用户ID获取用户文档文件夹
func (d *UserDocFolderDAO) GetUserDocFoldersByUserID(ctx context.Context, userID string) ([]model.UserDocFolder, error) {
	var folders []model.UserDocFolder
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false", userID).Find(&folders)
	if result.Error != nil {
		d.logger.Error("msg", "获取用户文档文件夹失败", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return folders, nil
}

// GetUserDocFoldersByParentID 根据父文件夹ID获取用户文档文件夹
func (d *UserDocFolderDAO) GetUserDocFoldersByParentID(ctx context.Context, parentID string) ([]model.UserDocFolder, error) {
	var folders []model.UserDocFolder
	result := d.GetDB(ctx).Where("parent_id = ? and is_deleted = false", parentID).Find(&folders)
	if result.Error != nil {
		d.logger.Error("msg", "获取用户文档文件夹失败", "parentID", parentID, "error", result.Error.Error())
		return nil, result.Error
	}
	return folders, nil
}

// BatchDeleteByIds 批量逻辑删除文件夹
func (d *UserDocFolderDAO) BatchDeleteByIds(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	// 使用上下文中的用户信息更新修改者字段
	updates := make(map[string]interface{})
	updates["is_deleted"] = true
	updates["updated_at"] = time.Now()

	uc := userContext.GetUserContext(ctx)
	if uc != nil {
		if uc.UserId != "0" {
			updates["modifier_id"] = uc.UserId
		}
		if uc.Username != "" {
			updates["modifier"] = uc.Username
		}
	}

	// 执行批量更新
	result := d.GetDB(ctx).Model(&model.UserDocFolder{}).Where("id IN ?", ids).Updates(updates)
	if result.Error != nil {
		d.logger.Error("msg", "批量删除文件夹失败", "ids", ids, "error", result.Error.Error())
		return result.Error
	}

	d.logger.Info("msg", "批量删除文件夹成功", "ids", ids, "count", result.RowsAffected)
	return nil
}

// UpdateFolder 更新文件夹信息
func (d *UserDocFolderDAO) UpdateFolder(ctx context.Context, folderId string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	// 添加更新时间和修改者信息
	updates["updated_at"] = time.Now()

	uc := userContext.GetUserContext(ctx)
	if uc != nil {
		if uc.UserId != "0" {
			updates["modifier_id"] = uc.UserId
		}
		if uc.Username != "" {
			updates["modifier"] = uc.Username
		}
	}

	// 执行更新
	result := d.GetDB(ctx).Model(&model.UserDocFolder{}).Where("id = ? AND is_deleted = false", folderId).Updates(updates)
	if result.Error != nil {
		d.logger.Error("msg", "更新文件夹失败", "folderId", folderId, "updates", updates, "error", result.Error.Error())
		return result.Error
	}

	if result.RowsAffected == 0 {
		d.logger.Warn("msg", "未找到要更新的文件夹或文件夹已删除", "folderId", folderId)
		return nil
	}

	d.logger.Info("msg", "更新文件夹成功", "folderId", folderId)
	return nil
}

// FindExistById 根据ID查找存在的文件夹
func (d *UserDocFolderDAO) FindExistById(ctx context.Context, id string) (*model.UserDocFolder, error) {
	return d.GetUserDocFolderByID(ctx, id)
}

// FindAllByIds 根据ID列表查找所有存在的文件夹
func (d *UserDocFolderDAO) FindAllByIds(ctx context.Context, ids []string) ([]model.UserDocFolder, error) {
	if len(ids) == 0 {
		return []model.UserDocFolder{}, nil
	}

	var folders []model.UserDocFolder
	result := d.GetDB(ctx).Where("id IN ? AND is_deleted = false", ids).Find(&folders)
	if result.Error != nil {
		d.logger.Error("msg", "根据ID列表查找文件夹失败", "ids", ids, "error", result.Error.Error())
		return nil, result.Error
	}

	return folders, nil
}

func (d *UserDocFolderDAO) SaveBatch(ctx context.Context, folders []model.UserDocFolder) error {
	if len(folders) == 0 {
		return nil
	}

	result := d.GetDB(ctx).Create(&folders)
	if result.Error != nil {
		d.logger.Error("msg", "批量保存文件夹失败", "count", len(folders), "error", result.Error.Error())
		return result.Error
	}

	d.logger.Info("msg", "批量保存文件夹成功", "count", len(folders))
	return nil
}
