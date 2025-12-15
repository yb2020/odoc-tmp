package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperAttachmentDAO 提供论文附件数据访问功能
type PaperAttachmentDAO struct {
	*baseDao.GormBaseDAO[model.PaperAttachment]
	logger logging.Logger
}

// NewPaperAttachmentDAO 创建一个新的论文附件DAO
func NewPaperAttachmentDAO(db *gorm.DB, logger logging.Logger) *PaperAttachmentDAO {
	return &PaperAttachmentDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperAttachment](db, logger),
		logger:      logger,
	}
}

// Create 创建论文附件
func (d *PaperAttachmentDAO) Create(ctx context.Context, attachment *model.PaperAttachment) error {
	return d.GetDB(ctx).Create(attachment).Error
}

// FindById 根据ID获取论文附件
func (d *PaperAttachmentDAO) FindById(ctx context.Context, id string) (*model.PaperAttachment, error) {
	var attachment model.PaperAttachment
	result := d.GetDB(ctx).Where("id = ?", id).First(&attachment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文附件失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &attachment, nil
}

// FindByPaperId 根据论文ID获取论文附件列表
func (d *PaperAttachmentDAO) FindByPaperId(ctx context.Context, paperId string) ([]model.PaperAttachment, error) {
	var attachments []model.PaperAttachment
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Find(&attachments)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文附件列表失败", "paper_id", paperId, "error", result.Error.Error())
		return nil, result.Error
	}
	return attachments, nil
}

// FindByType 根据附件类型获取论文附件列表
func (d *PaperAttachmentDAO) FindByType(ctx context.Context, attachType int) ([]model.PaperAttachment, error) {
	var attachments []model.PaperAttachment
	result := d.GetDB(ctx).Where("type = ?", attachType).Find(&attachments)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文附件列表失败", "type", attachType, "error", result.Error.Error())
		return nil, result.Error
	}
	return attachments, nil
}

// UpdateById 更新论文附件
func (d *PaperAttachmentDAO) UpdateById(ctx context.Context, attachment *model.PaperAttachment) error {
	return d.Modify(ctx, attachment)
}

// DeleteById 删除论文附件
func (d *PaperAttachmentDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// DeleteByPaperId 根据论文ID删除论文附件
func (d *PaperAttachmentDAO) DeleteByPaperId(ctx context.Context, paperId string) error {
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Delete(&model.PaperAttachment{})
	if result.Error != nil {
		d.logger.Error("msg", "删除论文附件失败", "paper_id", paperId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
