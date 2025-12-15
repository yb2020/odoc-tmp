package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

type DoiMetaInfoDAO struct {
	*baseDao.GormBaseDAO[model.DoiMetaInfo]
	logger logging.Logger
}

func NewDoiMetaInfoDAO(db *gorm.DB, logger logging.Logger) *DoiMetaInfoDAO {
	return &DoiMetaInfoDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.DoiMetaInfo](db, logger),
		logger:      logger,
	}
}

// GetByDoi 根据元数据By DOI
func (d *DoiMetaInfoDAO) GetByDoi(ctx context.Context, doi string) (*model.DoiMetaInfo, error) {
	var doiMetaInfo model.DoiMetaInfo
	result := d.GetDB(ctx).Where("doi = ? and is_deleted = false", doi).First(&doiMetaInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取doi失败", "doi", doi, "error", result.Error.Error())
		return nil, result.Error
	}
	return &doiMetaInfo, nil
}

// GetByPaperId 根据元数据By PaperId
func (d *DoiMetaInfoDAO) GetByPaperId(ctx context.Context, paperId string) (*model.DoiMetaInfo, error) {
	var doiMetaInfo model.DoiMetaInfo
	result := d.GetDB(ctx).Where("paper_id = ? and is_deleted = false", paperId).First(&doiMetaInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取doi失败", "paperId", paperId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &doiMetaInfo, nil
}

// SelectForInitMetaField selects a batch of DoiMetaInfo records for initializing meta fields.
func (d *DoiMetaInfoDAO) SelectForInitMetaField(ctx context.Context, batchSize int, latestId *string) ([]*model.DoiMetaInfo, error) {
	var doiMetaInfos []*model.DoiMetaInfo
	query := d.GetDB(ctx).Model(&model.DoiMetaInfo{}).Where("is_deleted = false")

	if latestId != nil {
		query = query.Where("id > ?", *latestId)
	}

	result := query.Order("id asc").Limit(batchSize).Find(&doiMetaInfos)
	if result.Error != nil {
		d.logger.Error("msg", "select for init meta field failed", "error", result.Error.Error())
		return nil, result.Error
	}
	return doiMetaInfos, nil
}
