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

// PdfThumbDAO GORM实现的PDF缩略图DAO
type PdfThumbDAO struct {
	*baseDao.GormBaseDAO[model.PdfThumb]
	logger logging.Logger
}

// NewPdfThumbDAO 创建一个新的PDF缩略图DAO
func NewPdfThumbDAO(db *gorm.DB, logger logging.Logger) *PdfThumbDAO {
	return &PdfThumbDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfThumb](db, logger),
		logger:      logger,
	}
}

// Create 创建PDF缩略图
func (d *PdfThumbDAO) Create(ctx context.Context, thumb *model.PdfThumb) error {
	return d.GetDB(ctx).Create(thumb).Error
}

// GetPdfThumbByID 根据ID获取PDF缩略图
func (d *PdfThumbDAO) GetPdfThumbByID(ctx context.Context, id string) (*model.PdfThumb, error) {
	var thumb model.PdfThumb
	result := d.GetDB(ctx).Where("id = ?", id).First(&thumb)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF缩略图失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &thumb, nil
}

// GetPdfThumbByPaperID 根据论文ID获取PDF缩略图
func (d *PdfThumbDAO) GetPdfThumbByPaperID(ctx context.Context, paperID string) (*model.PdfThumb, error) {
	var thumb model.PdfThumb
	result := d.GetDB(ctx).Where("paper_id = ?", paperID).First(&thumb)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF缩略图失败", "paper_id", paperID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &thumb, nil
}

// GetPdfThumbByPdfID 根据PDF ID获取PDF缩略图
func (d *PdfThumbDAO) GetPdfThumbByPdfID(ctx context.Context, pdfID string) (*model.PdfThumb, error) {
	var thumb model.PdfThumb
	result := d.GetDB(ctx).Where("pdf_id = ?", pdfID).First(&thumb)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF缩略图失败", "pdf_id", pdfID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &thumb, nil
}

// UpdatePdfThumb 更新PDF缩略图
func (d *PdfThumbDAO) UpdatePdfThumb(ctx context.Context, thumb *model.PdfThumb) error {
	return d.Modify(ctx, thumb)
}

// DeletePdfThumb 删除PDF缩略图
func (d *PdfThumbDAO) DeletePdfThumb(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// ListPdfThumbs 列出PDF缩略图
func (d *PdfThumbDAO) ListPdfThumbs(ctx context.Context, limit, offset int) ([]model.PdfThumb, error) {
	var thumbs []model.PdfThumb
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&thumbs)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF缩略图失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return thumbs, nil
}

// CountPdfThumbs 获取PDF缩略图总数
func (d *PdfThumbDAO) CountPdfThumbs(ctx context.Context) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfThumb{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF缩略图总数失败", "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}
