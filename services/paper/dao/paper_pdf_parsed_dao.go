package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperPdfParsedDAO 提供PaperPdfParsed数据访问功能
type PaperPdfParsedDAO struct {
	*baseDao.GormBaseDAO[model.PaperPdfParsed]
	logger logging.Logger
}

// NewPaperPdfParsedDAO 创建一个新的PaperPdfParsedDAO
func NewPaperPdfParsedDAO(db *gorm.DB, logger logging.Logger) *PaperPdfParsedDAO {
	return &PaperPdfParsedDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperPdfParsed](db, logger),
		logger:      logger,
	}
}

// BatchSave 批量保存PaperPdfParsed（SQL级别的批量插入）
func (d *PaperPdfParsedDAO) BatchSave(ctx context.Context, records []*model.PaperPdfParsed) error {
	if len(records) == 0 {
		return nil
	}

	// 使用事务进行批量插入
	err := d.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 使用CreateInBatches进行批量插入，每批100条记录
		result := tx.CreateInBatches(records, 100)
		if result.Error != nil {
			d.logger.Error("批量保存PaperPdfParsed失败", "error", result.Error.Error())
			return result.Error
		}
		d.logger.Info("批量保存PaperPdfParsed成功", "count", len(records))
		return nil
	})

	if err != nil {
		d.logger.Error("批量保存PaperPdfParsed事务失败", "error", err.Error())
		return err
	}

	return nil
}

// GetBySourcePdfFileSHA256AndVersion 根据源PDF文件的fileSHA256查询记录，返回可能的多条记录
func (d *PaperPdfParsedDAO) GetBySourcePdfFileSHA256AndVersion(ctx context.Context, sha256, version string) ([]*model.PaperPdfParsed, error) {
	var records []*model.PaperPdfParsed
	result := d.GetDB(ctx).Where("source_pdf_sha256 = ? AND version = ? AND is_deleted = ?", sha256, version, false).Find(&records)
	if result.Error != nil {
		d.logger.Error("查询PaperPdfParsed记录失败", "error", result.Error.Error(), "source_pdf_sha256", sha256, "version", version)
		return nil, result.Error
	}

	if len(records) == 0 {
		d.logger.Info("未找到对应的PaperPdfParsed记录", "source_pdf_sha256", sha256, "version", version)
	} else {
		d.logger.Info("查询到PaperPdfParsed记录", "source_pdf_sha256", sha256, "version", version, "count", len(records))
	}

	return records, nil
}

// GetBySourcePdfFileSHA256AndTypeAndVersion 根据源PDF文件的fileSHA256和文件类型查询记录，必须是唯一的
func (d *PaperPdfParsedDAO) GetBySourcePdfFileSHA256AndTypeAndVersion(ctx context.Context, sha256, fileType, version string) (*model.PaperPdfParsed, error) {
	var records *model.PaperPdfParsed
	result := d.GetDB(ctx).Where("source_pdf_sha256 = ? AND file_type = ? AND version = ? AND is_deleted = ?", sha256, fileType, version, false).Find(&records)
	if result.Error != nil {
		d.logger.Error("查询PaperPdfParsed记录失败", "error", result.Error.Error(), "source_pdf_sha256", sha256, "file_type", fileType)
		return nil, result.Error
	}

	if records == nil {
		d.logger.Info("未找到对应的PaperPdfParsed记录", "source_pdf_sha256", sha256, "file_type", fileType)
	} else {
		d.logger.Info("查询到PaperPdfParsed记录", "source_pdf_sha256", sha256, "file_type", fileType)
	}

	return records, nil
}

func (d *PaperPdfParsedDAO) HasExistBySourcePdfFileSHA256AndVersion(ctx context.Context, sha256, version string) (bool, error) {
	var records []*model.PaperPdfParsed
	result := d.GetDB(ctx).Where("source_pdf_sha256 = ? AND version = ? AND is_deleted = ?", sha256, version, false).Find(&records)
	if result.Error != nil {
		d.logger.Error("查询PaperPdfParsed记录失败", "error", result.Error.Error(), "source_pdf_sha256", sha256, "version", version)
		return false, result.Error
	}

	if len(records) == 0 {
		d.logger.Info("未找到对应的PaperPdfParsed记录", "source_pdf_sha256", sha256, "version", version)
		return false, nil
	}

	return true, nil
}
