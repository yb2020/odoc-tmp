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

// PaperPdfSelectRecordDAO GORM实现的论文PDF选择记录DAO
type PaperPdfSelectRecordDAO struct {
	*baseDao.GormBaseDAO[model.PaperPdfSelectRecord]
	logger logging.Logger
}

// NewPaperPdfSelectRecordDAO 创建一个新的论文PDF选择记录DAO
func NewPaperPdfSelectRecordDAO(db *gorm.DB, logger logging.Logger) *PaperPdfSelectRecordDAO {
	return &PaperPdfSelectRecordDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperPdfSelectRecord](db, logger),
		logger:      logger,
	}
}

// CreatePaperPdfSelectRecord 创建论文PDF选择记录
func (d *PaperPdfSelectRecordDAO) Create(ctx context.Context, record *model.PaperPdfSelectRecord) error {
	return d.GetDB(ctx).Create(record).Error
}

// GetPaperPdfSelectRecordByID 根据ID获取论文PDF选择记录
func (d *PaperPdfSelectRecordDAO) GetPaperPdfSelectRecordByID(ctx context.Context, id string) (*model.PaperPdfSelectRecord, error) {
	var record model.PaperPdfSelectRecord
	result := d.GetDB(ctx).Where("id = ?", id).First(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文PDF选择记录失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &record, nil
}

// GetPaperPdfSelectRecordByPaperID 根据论文ID获取论文PDF选择记录
func (d *PaperPdfSelectRecordDAO) GetPaperPdfSelectRecordByPaperID(ctx context.Context, paperID string) (*model.PaperPdfSelectRecord, error) {
	var record model.PaperPdfSelectRecord
	result := d.GetDB(ctx).Where("paper_id = ?", paperID).First(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文PDF选择记录失败", "paper_id", paperID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &record, nil
}

// GetPaperPdfSelectRecordByUserID 根据用户ID获取论文PDF选择记录
func (d *PaperPdfSelectRecordDAO) GetPaperPdfSelectRecordByUserID(ctx context.Context, userID string) ([]model.PaperPdfSelectRecord, error) {
	var records []model.PaperPdfSelectRecord
	result := d.GetDB(ctx).Where("user_id = ?", userID).Find(&records)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文PDF选择记录失败", "user_id", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return records, nil
}

// UpdatePaperPdfSelectRecord 更新论文PDF选择记录
func (d *PaperPdfSelectRecordDAO) UpdatePaperPdfSelectRecord(ctx context.Context, record *model.PaperPdfSelectRecord) error {
	return d.Modify(ctx, record)
}

// DeletePaperPdfSelectRecord 删除论文PDF选择记录
func (d *PaperPdfSelectRecordDAO) DeletePaperPdfSelectRecord(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// ListPaperPdfSelectRecords 列出论文PDF选择记录
func (d *PaperPdfSelectRecordDAO) ListPaperPdfSelectRecords(ctx context.Context, limit, offset int) ([]model.PaperPdfSelectRecord, error) {
	var records []model.PaperPdfSelectRecord
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&records)
	if result.Error != nil {
		d.logger.Error("msg", "列出论文PDF选择记录失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return records, nil
}

// CountPaperPdfSelectRecords 获取论文PDF选择记录总数
func (d *PaperPdfSelectRecordDAO) CountPaperPdfSelectRecords(ctx context.Context) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PaperPdfSelectRecord{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文PDF选择记录总数失败", "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}
