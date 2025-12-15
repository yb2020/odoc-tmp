package dao

import (
	"context"
	"errors"

	pb "github.com/yb2020/odoc-proto/gen/go/nav"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/nav/model"
	"gorm.io/gorm"
)

// WebsiteDAO GORM实现的网站DAO
type WebsiteDAO struct {
	*baseDao.GormBaseDAO[model.Website]
	logger logging.Logger
}

// NewWebsiteDAO 创建一个新的网站DAO
func NewWebsiteDAO(db *gorm.DB, logger logging.Logger) *WebsiteDAO {
	return &WebsiteDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.Website](db, logger),
		logger:      logger,
	}
}

// GetUserWebsiteListBySortOrder 获取用户学术网站列表(按SortOrder升序排序)
func (s *WebsiteDAO) FindExistById(ctx context.Context, id string) (*model.Website, error) {
	return s.GormBaseDAO.FindExistById(ctx, id)
}

// GetUserWebsiteListBySortOrder 获取用户学术网站列表(按SortOrder升序排序)
func (s *WebsiteDAO) GetUserWebsiteListBySortOrder(ctx context.Context, userId string) ([]model.Website, error) {
	var entities []model.Website
	// 使用 sort_order, id 联合排序，确保在 sort_order 相同时获得一个稳定的排序结果
	result := s.GetDB(ctx).Where("is_deleted = false and user_id = ?", userId).Order("sort_order asc, id asc").Find(&entities)
	if result.Error != nil {
		s.logger.Error("msg", "获取用户学术网站失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return entities, nil
}

func (s *WebsiteDAO) FindOne(ctx context.Context, conditions map[string]interface{}) (*model.Website, error) {
	var entity model.Website
	result := s.GetDB(ctx).Where(conditions).First(&entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Record not found is not an error
		}
		return nil, result.Error
	}
	return &entity, nil
}

// BatchUpdateSortOrder 批量更新排序
func (s *WebsiteDAO) BatchUpdateSortOrder(ctx context.Context, updates map[string]int32) error {
	return s.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		for id, sortOrder := range updates {
			result := tx.Model(&model.Website{}).Where("id = ?", id).Update("sort_order", sortOrder)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

// UpdateSortOrder updates the sort order for a single website.
func (s *WebsiteDAO) UpdateSortOrder(ctx context.Context, id string, sortOrder int32) error {
	result := s.GetDB(ctx).Model(&model.Website{}).Where("id = ?", id).Update("sort_order", sortOrder)
	return result.Error
}

// GetMaxSortOrder gets the maximum sort order for a user.
func (s *WebsiteDAO) GetMaxSortOrder(ctx context.Context, userId string) (int32, error) {
	var website model.Website
	result := s.GetDB(ctx).
		Where("user_id = ? AND is_deleted = false", userId).
		Order("sort_order DESC").
		First(&website)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// This is the first website for the user, which is a valid case.
			return 0, nil
		}
		// For other errors, return them.
		return 0, result.Error
	}

	return website.SortOrder, nil
}

// CheckUserInitSystemWebsite 检查用户是否初始化了系统网站
func (s *WebsiteDAO) CheckUserInitSystemWebsite(ctx context.Context, userId string) (bool, error) {
	var websites []model.Website
	result := s.GetDB(ctx).Where("user_id = ? AND source = ?", userId, int32(pb.WebsiteSource_WebsiteSource_System)).Find(&websites)
	if result.Error != nil {
		return false, result.Error
	}

	if len(websites) > 0 {
		return true, nil
	}
	return false, nil
}
