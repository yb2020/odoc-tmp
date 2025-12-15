package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperJcrDAO GORM实现的论文JCR DAO
type PaperJcrDAO struct {
	*baseDao.GormBaseDAO[model.PaperJcrEntity]
	logger logging.Logger
}

// NewPaperJcrDAO 创建一个新的论文JCR DAO
func NewPaperJcrDAO(db *gorm.DB, logger logging.Logger) *PaperJcrDAO {
	return &PaperJcrDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperJcrEntity](db, logger),
		logger:      logger,
	}
}

// GetPaperJcrEntityByVenue 根据 venue 查找 PaperJcrEntity
func (d *PaperJcrDAO) GetPaperJcrEntityByVenue(ctx context.Context, venue string) (*model.PaperJcrEntity, error) {
	var entity model.PaperJcrEntity
	err := d.GetDB(ctx).
		Where("venue = ?", venue).
		First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		d.logger.Error("GetPaperJcrEntityByVenue failed, venue=%s, err=%v", venue, err)
		return nil, err
	}
	return &entity, nil
}
