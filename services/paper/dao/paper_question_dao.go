package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperQuestionDAO 提供论文问题数据访问功能
type PaperQuestionDAO struct {
	*baseDao.GormBaseDAO[model.PaperQuestion]
	logger logging.Logger
}

// NewPaperQuestionDAO 创建一个新的论文问题DAO
func NewPaperQuestionDAO(db *gorm.DB, logger logging.Logger) *PaperQuestionDAO {
	return &PaperQuestionDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperQuestion](db, logger),
		logger:      logger,
	}
}

// Create 创建论文问题
func (d *PaperQuestionDAO) Create(ctx context.Context, question *model.PaperQuestion) error {
	return d.GetDB(ctx).Create(question).Error
}

// FindById 根据ID获取论文问题
func (d *PaperQuestionDAO) FindById(ctx context.Context, id string) (*model.PaperQuestion, error) {
	var question model.PaperQuestion
	result := d.GetDB(ctx).Where("id = ?", id).First(&question)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文问题失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &question, nil
}

// FindByPaperId 根据论文ID获取论文问题列表
func (d *PaperQuestionDAO) FindByPaperId(ctx context.Context, paperId string) ([]model.PaperQuestion, error) {
	var questions []model.PaperQuestion
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Find(&questions)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文问题列表失败", "paper_id", paperId, "error", result.Error.Error())
		return nil, result.Error
	}
	return questions, nil
}

// FindByPdfId 根据PDF ID获取论文问题列表
func (d *PaperQuestionDAO) FindByPdfId(ctx context.Context, pdfId string) ([]model.PaperQuestion, error) {
	var questions []model.PaperQuestion
	result := d.GetDB(ctx).Where("pdf_id = ?", pdfId).Find(&questions)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文问题列表失败", "pdf_id", pdfId, "error", result.Error.Error())
		return nil, result.Error
	}
	return questions, nil
}

// UpdateById 更新论文问题
func (d *PaperQuestionDAO) UpdateById(ctx context.Context, question *model.PaperQuestion) error {
	return d.Modify(ctx, question)
}

// IncrementViewCount 增加查看次数
func (d *PaperQuestionDAO) IncrementViewCount(ctx context.Context, id string) error {
	result := d.GetDB(ctx).Model(&model.PaperQuestion{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	if result.Error != nil {
		d.logger.Error("msg", "增加论文问题查看次数失败", "id", id, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// DeleteById 删除论文问题
func (d *PaperQuestionDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// DeleteByPaperId 根据论文ID删除论文问题
func (d *PaperQuestionDAO) DeleteByPaperId(ctx context.Context, paperId string) error {
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Delete(&model.PaperQuestion{})
	if result.Error != nil {
		d.logger.Error("msg", "删除论文问题失败", "paper_id", paperId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
