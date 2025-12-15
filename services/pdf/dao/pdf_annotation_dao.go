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

// PdfAnnotationDAO GORM实现的PDF注释DAO
type PdfAnnotationDAO struct {
	*baseDao.GormBaseDAO[model.PdfAnnotation]
	logger logging.Logger
}

// NewPdfAnnotationDAO 创建一个新的PDF注释DAO
func NewPdfAnnotationDAO(db *gorm.DB, logger logging.Logger) *PdfAnnotationDAO {
	return &PdfAnnotationDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfAnnotation](db, logger),
		logger:      logger,
	}
}

// Create 创建PDF注释
func (d *PdfAnnotationDAO) Create(ctx context.Context, annotation *model.PdfAnnotation) error {
	return d.GetDB(ctx).Create(annotation).Error
}

// GetPdfAnnotationByID 根据ID获取PDF注释
func (d *PdfAnnotationDAO) GetPdfAnnotationByID(ctx context.Context, id string) (*model.PdfAnnotation, error) {
	var annotation model.PdfAnnotation
	result := d.GetDB(ctx).Where("id = ?", id).First(&annotation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF注释失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &annotation, nil
}

// GetPdfAnnotationsByPdfID 根据PDF ID获取PDF注释列表
func (d *PdfAnnotationDAO) GetPdfAnnotationsByPdfID(ctx context.Context, pdfID string) ([]model.PdfAnnotation, error) {
	var annotations []model.PdfAnnotation
	result := d.GetDB(ctx).Where("pdf_id = ?", pdfID).Find(&annotations)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF注释列表失败", "pdf_id", pdfID, "error", result.Error.Error())
		return nil, result.Error
	}
	return annotations, nil
}

// GetPdfAnnotationsByPage 根据PDF ID和页码获取PDF注释列表
func (d *PdfAnnotationDAO) GetPdfAnnotationsByPage(ctx context.Context, pdfID string, page int) ([]model.PdfAnnotation, error) {
	var annotations []model.PdfAnnotation
	result := d.GetDB(ctx).Where("pdf_id = ? AND page = ?", pdfID, page).Find(&annotations)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF注释列表失败", "pdf_id", pdfID, "page", page, "error", result.Error.Error())
		return nil, result.Error
	}
	return annotations, nil
}

// GetPdfAnnotationByMD5 根据注释MD5获取PDF注释
func (d *PdfAnnotationDAO) GetPdfAnnotationByMD5(ctx context.Context, md5 string) (*model.PdfAnnotation, error) {
	var annotation model.PdfAnnotation
	result := d.GetDB(ctx).Where("annotation_md5 = ?", md5).First(&annotation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF注释失败", "annotation_md5", md5, "error", result.Error.Error())
		return nil, result.Error
	}
	return &annotation, nil
}

// UpdatePdfAnnotation 更新PDF注释
func (d *PdfAnnotationDAO) UpdatePdfAnnotation(ctx context.Context, annotation *model.PdfAnnotation) error {
	return d.Modify(ctx, annotation)
}

// DeletePdfAnnotation 删除PDF注释
func (d *PdfAnnotationDAO) DeletePdfAnnotation(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// ListPdfAnnotations 列出PDF注释
func (d *PdfAnnotationDAO) ListPdfAnnotations(ctx context.Context, limit, offset int) ([]model.PdfAnnotation, error) {
	var annotations []model.PdfAnnotation
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&annotations)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF注释失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return annotations, nil
}

// CountPdfAnnotations 获取PDF注释总数
func (d *PdfAnnotationDAO) CountPdfAnnotations(ctx context.Context) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfAnnotation{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF注释总数失败", "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}
