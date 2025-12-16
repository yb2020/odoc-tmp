package dao

import (
	"context"
	"time"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	translate "github.com/yb2020/odoc/proto/gen/go/translate"
	"github.com/yb2020/odoc/services/translate/model"
	"gorm.io/gorm"
)

// FullTextTranslateDAO 实现的翻译DAO
type FullTextTranslateDAO struct {
	*baseDao.GormBaseDAO[model.FullTextTranslate]
	logger logging.Logger
}

// NewFullTextTranslateHistoryDAO 创建一个新的GORM翻译DAO
func NewFullTextTranslateHistoryDAO(db *gorm.DB, logger logging.Logger) FullTextTranslateDAO {
	return FullTextTranslateDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.FullTextTranslate](db, logger),
		logger:      logger,
	}
}

// FindByPdfId 根据PdfId查询翻译历史记录
func (dao *FullTextTranslateDAO) FindByPdfId(ctx context.Context, pdfId string) ([]*model.FullTextTranslate, error) {
	var histories []*model.FullTextTranslate
	err := dao.GetDB(ctx).
		Where("source_pdf_id = ? AND is_deleted = false", pdfId).
		Find(&histories).Error
	if err != nil {
		dao.logger.Error("查询翻译历史记录失败", "error", err.Error(), "pdfId", pdfId)
		return nil, err
	}
	return histories, nil
}

// FindLatestByFileSHA256AndStatus 根据FileSHA256和状态查询最新的一条翻译历史记录
func (dao *FullTextTranslateDAO) FindLatestByFileSHA256AndStatus(ctx context.Context, fileSHA256 string, status translate.FullTranslateFlowStatus) (*model.FullTextTranslate, error) {
	var history model.FullTextTranslate
	err := dao.GetDB(ctx).
		Where("file_sha256 = ? AND status = ? AND is_deleted = false", fileSHA256, status).
		Order("created_at DESC").
		Limit(1).
		First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		dao.logger.Error("查询最新翻译历史记录失败", "error", err.Error(), "fileSHA256", fileSHA256, "status", status)
		return nil, err
	}
	return &history, nil
}

// FindByPdfIdAndUserIdAndStatusWithinDays 根据PdfId、UserId和状态查询指定天数内的翻译历史记录
func (dao *FullTextTranslateDAO) FindByPdfIdAndUserIdAndStatusWithinDays(ctx context.Context, pdfId, userId string, status translate.FullTranslateFlowStatus, days int) ([]*model.FullTextTranslate, error) {
	var histories []*model.FullTextTranslate
	startDate := time.Now().AddDate(0, 0, -days)
	err := dao.GetDB(ctx).
		Where("source_pdf_id = ? AND user_id = ? AND status = ? AND is_deleted = false AND created_at >= ?", pdfId, userId, status, startDate).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		dao.logger.Error("查询指定天数内的翻译历史记录失败", "error", err.Error(), "pdfId", pdfId, "userId", userId, "status", status, "days", days)
		return nil, err
	}
	return histories, nil
}

// FindByPdfIdAndUserIdWithinDays 根据PdfId和UserId查询指定天数内的翻译历史记录
func (dao *FullTextTranslateDAO) FindByPdfIdAndUserIdWithinDays(ctx context.Context, pdfId, userId string, days int) ([]*model.FullTextTranslate, error) {
	var histories []*model.FullTextTranslate
	startDate := time.Now().AddDate(0, 0, -days)
	err := dao.GetDB(ctx).
		Where("source_pdf_id = ? AND user_id = ? AND is_deleted = false AND created_at >= ?", pdfId, userId, startDate).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		dao.logger.Error("查询指定天数内的翻译历史记录失败", "error", err.Error(), "pdfId", pdfId, "userId", userId, "days", days)
		return nil, err
	}
	return histories, nil
}

// FindByPdfIdAndStatusList 根据PdfId和状态列表查询翻译历史记录
func (dao *FullTextTranslateDAO) FindByPdfIdAndStatusList(ctx context.Context, pdfId string, statusList []translate.FullTranslateFlowStatus) ([]*model.FullTextTranslate, error) {
	var histories []*model.FullTextTranslate
	err := dao.GetDB(ctx).
		Where("source_pdf_id = ? AND status IN ? AND is_deleted = false", pdfId, statusList).
		Find(&histories).Error
	if err != nil {
		dao.logger.Error("查询翻译历史记录失败", "error", err.Error(), "pdfId", pdfId, "statusList", statusList)
		return nil, err
	}
	return histories, nil
}

// FindByUserIdWithinDays 根据UserId查询指定天数内的翻译历史记录
func (dao *FullTextTranslateDAO) FindByUserIdWithinDays(ctx context.Context, userId string, days int) ([]*model.FullTextTranslate, error) {
	var histories []*model.FullTextTranslate
	startDate := time.Now().AddDate(0, 0, -days)
	err := dao.GetDB(ctx).
		Where("user_id = ? AND is_deleted = false AND created_at >= ?", userId, startDate).
		Find(&histories).Error
	if err != nil {
		dao.logger.Error("查询指定天数内的翻译历史记录失败", "error", err.Error(), "userId", userId, "days", days)
		return nil, err
	}
	return histories, nil
}

func (dao *FullTextTranslateDAO) FindByFlowNumber(ctx context.Context, flowNumber string) ([]model.FullTextTranslate, error) {
	var histories []model.FullTextTranslate
	err := dao.GetDB(ctx).
		Where("flow_number = ? AND is_deleted = false", flowNumber).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		dao.logger.Error("查询指定天数内的翻译历史记录失败", "error", err.Error(), "flowNumber", flowNumber)
		return nil, err
	}
	return histories, nil
}

// FindByUserIdAndIntervalDays 根据UserId和天数查询翻译历史记录
func (dao *FullTextTranslateDAO) FindByUserIdAndIntervalDays(ctx context.Context, userId string, days int) ([]*model.FullTextTranslate, error) {
	var histories []*model.FullTextTranslate
	startDate := time.Now().AddDate(0, 0, -days)
	err := dao.GetDB(ctx).
		Where("user_id = ? AND is_deleted = false AND created_at >= ?", userId, startDate).
		Find(&histories).Error
	if err != nil {
		dao.logger.Error("查询指定天数内的翻译历史记录失败", "error", err.Error(), "userId", userId, "days", days)
		return nil, err
	}
	return histories, nil
}
