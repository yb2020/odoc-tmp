package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/model"
	"gorm.io/gorm"
)

// PdfSummaryDAO GORM实现的PDF 摘要DAO
type PdfSummaryDAO struct {
	*baseDao.GormBaseDAO[model.PdfSummary]
	logger logging.Logger
}

// NewPdfSummaryDAO 创建一个新的PDF 摘要DAO
func NewPdfSummaryDAO(db *gorm.DB, logger logging.Logger) *PdfSummaryDAO {
	return &PdfSummaryDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfSummary](db, logger),
		logger:      logger,
	}
}

// FindExistBySourceFromAndFileSHA256AndVersionAndLang 根据source_from,file_sha256,version,lang查找数据
func (d *PdfSummaryDAO) FindExistBySourceFromAndFileSHA256AndVersionAndLang(ctx context.Context, sourceFrom, fileSHA256, version, lang string) (*model.PdfSummary, error) {
	var entity model.PdfSummary
	result := d.GetDB(ctx).Where("source_from = ? and file_sha256 = ? and version = ? and lang = ? and is_deleted = false", sourceFrom, fileSHA256, version, lang).First(&entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据source_from,file_sha256,version,lang查找数据失败", "source_from", sourceFrom, "file_sha256", fileSHA256, "version", version, "lang", lang, "error", result.Error.Error())
		return nil, result.Error
	}
	return &entity, nil
}
