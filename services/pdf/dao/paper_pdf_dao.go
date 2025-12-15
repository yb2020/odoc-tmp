package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/model"
	"gorm.io/gorm"
)

// PaperPDFDAO 提供论文PDF数据访问功能

// PaperPDFDAO GORM实现的论文PDF DAO
type PaperPDFDAO struct {
	*baseDao.GormBaseDAO[model.PaperPdf]
	logger logging.Logger
}

// NewPaperPDFDAO 创建一个新的论文PDF DAO
func NewPaperPDFDAO(db *gorm.DB, logger logging.Logger) *PaperPDFDAO {
	return &PaperPDFDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperPdf](db, logger),
		logger:      logger,
	}
}

// GetPaperPDFByID 根据ID获取PDF记录
func (d *PaperPDFDAO) GetPaperPDFByID(ctx context.Context, id string) (*model.PaperPdf, error) {
	var pdf model.PaperPdf
	result := d.GetDB(ctx).Where("id = ? and is_deleted = false", id).First(&pdf)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文PDF失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &pdf, nil
}

// GetPaperPDFByPaperID 根据论文ID获取PDF记录
func (d *PaperPDFDAO) GetPaperPDFByPaperID(ctx context.Context, paperID string) (*model.PaperPdf, error) {
	var pdf model.PaperPdf
	result := d.GetDB(ctx).Where("paper_id = ? and is_deleted = false", paperID).First(&pdf)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文PDF失败", "paper_id", paperID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &pdf, nil
}

// GetPaperPDFByFileSHA256 根据FileSHA256获取PDF记录
func (d *PaperPDFDAO) GetPaperPDFByFileSHA256(ctx context.Context, fileSHA256 string) (*model.PaperPdf, error) {
	var pdf model.PaperPdf
	result := d.GetDB(ctx).Where("file_sha256 = ? and is_deleted = false", fileSHA256).Limit(1).First(&pdf)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文PDF失败", "file_sha256", fileSHA256, "error", result.Error.Error())
		return nil, result.Error
	}
	return &pdf, nil
}

// UpdatePaperPDF 更新PDF记录
func (d *PaperPDFDAO) UpdatePaperPDF(ctx context.Context, pdf *model.PaperPdf) error {
	return d.Modify(ctx, pdf)
}

// DeletePaperPDF 删除PDF记录
func (d *PaperPDFDAO) DeletePaperPDF(ctx context.Context, id string) error {
	// 由于 BaseDAO 中的 DeleteById 接收的是 int64 类型，这里需要转换
	return d.RemoveById(ctx, id)
}

// ListPaperPDFs 列出PDF记录
func (d *PaperPDFDAO) ListPaperPDFs(ctx context.Context, limit, offset int) ([]model.PaperPdf, error) {
	var pdfs []model.PaperPdf
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&pdfs)
	if result.Error != nil {
		d.logger.Error("msg", "列出论文PDF失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return pdfs, nil
}

// CountPaperPDFs 获取PDF记录总数
func (d *PaperPDFDAO) CountPaperPDFs(ctx context.Context) (int64, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PaperPdf{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文PDF总数失败", "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// GetDefaultPdfByPaperId 根据论文ID获取默认PDF记录
func (d *PaperPDFDAO) GetDefaultPdfByPaperId(ctx context.Context, paperId string) (*model.PaperPdf, error) {
	var pdf model.PaperPdf
	result := d.GetDB(ctx).Where("paper_id = ? AND is_deleted = false", paperId).First(&pdf)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文PDF失败", "paper_id", paperId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &pdf, nil
}

// GetVisitablePdfIdByPaperIdAndUserId 根据论文ID和用户ID获取用户可访问的PDF记录
func (d *PaperPDFDAO) GetVisitablePdfIdByPaperIdAndUserId(ctx context.Context, paperId string, userId string) (*model.PaperPdf, error) {
	var pdf model.PaperPdf
	// 由于GORM不直接支持这种复杂的SQL查询，我们可以使用原生SQL
	// todo:  这里的sql需要讨论一下是否按照这种模式实现
	sql := `
	SELECT pdf.*
	FROM (
		SELECT *
		FROM t_paper_pdf
		WHERE paper_id = ? AND is_deleted = false
	) AS pdf
	LEFT JOIN (
		SELECT *
		FROM t_user_relation_other
		WHERE user_id = ? AND sign = 7 AND type = 'own' AND is_deleted = false
	) AS relation
	ON pdf.id = relation.relation_id
	WHERE relation.user_id = ? OR pdf.creator_id = 0
	ORDER BY GREATEST(pdf.created_at, IFNULL(relation.updated_at, 0)) DESC
	LIMIT 1
	`

	result := d.GetDB(ctx).Raw(sql, paperId, userId, userId).Scan(&pdf)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户可访问的PDF记录失败", "paperId", paperId, "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &pdf, nil
}

// GetByUserIdAndFileSHA256 根据用户ID和FileSHA256获取PDF记录
func (d *PaperPDFDAO) GetByUserIdAndFileSHA256(ctx context.Context, userId string, fileSHA256 string) (*model.PaperPdf, error) {
	var pdf model.PaperPdf
	result := d.GetDB(ctx).
		Where("creator_id = ? AND file_sha256 = ? AND is_deleted = false", userId, fileSHA256).
		First(&pdf)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据用户ID和FileSHA256获取PDF记录失败", "userId", userId, "file_sha256", fileSHA256, "error", result.Error.Error())
		return nil, result.Error
	}

	return &pdf, nil
}

// 根据用户id获取用户的所有paperpdf记录列表
func (d *PaperPDFDAO) GetUserPaperPdfList(ctx context.Context, userId string) ([]model.PaperPdf, error) {
	var pdfs []model.PaperPdf
	result := d.GetDB(ctx).Where("creator_id = ? AND is_deleted = false", userId).Find(&pdfs)
	if result.Error != nil {
		d.logger.Error("msg", "根据用户id获取用户的所有paperpdf记录列表失败", "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return pdfs, nil
}
