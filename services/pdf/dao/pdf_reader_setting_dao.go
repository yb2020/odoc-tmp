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

// PdfReaderSettingDAO GORM实现的PDF阅读器设置DAO
type PdfReaderSettingDAO struct {
	*baseDao.GormBaseDAO[model.PdfReaderSetting]
	logger logging.Logger
}

// NewPdfReaderSettingDAO 创建一个新的PDF阅读器设置DAO
func NewPdfReaderSettingDAO(db *gorm.DB, logger logging.Logger) *PdfReaderSettingDAO {
	return &PdfReaderSettingDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfReaderSetting](db, logger),
		logger:      logger,
	}
}

// Create 创建PDF阅读器设置
func (d *PdfReaderSettingDAO) Create(ctx context.Context, setting *model.PdfReaderSetting) error {
	return d.GetDB(ctx).Create(setting).Error
}

// GetPdfReaderSettingByID 根据ID获取PDF阅读器设置
func (d *PdfReaderSettingDAO) GetPdfReaderSettingByID(ctx context.Context, id string) (*model.PdfReaderSetting, error) {
	var setting model.PdfReaderSetting
	result := d.GetDB(ctx).Where("id = ?", id).First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF阅读器设置失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &setting, nil
}

// GetPdfReaderSettingByClientType 根据客户端类型获取PDF阅读器设置
func (d *PdfReaderSettingDAO) GetPdfReaderSettingByClientType(ctx context.Context, clientType int, creatorId string) (*model.PdfReaderSetting, error) {
	var setting model.PdfReaderSetting
	result := d.GetDB(ctx).Where("client_type = ? AND creator_id = ?", clientType, creatorId).First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF阅读器设置失败", "client_type", clientType, "creator_id", creatorId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &setting, nil
}

// GetPdfReaderSettingByClientType 根据客户端类型获取PDF阅读器设置
func (d *PdfReaderSettingDAO) GetPdfReaderSettingUserIdByClientType(ctx context.Context, creatorId string, clientType int) (*model.PdfReaderSetting, error) {
	var setting model.PdfReaderSetting
	result := d.GetDB(ctx).Where("client_type = ? AND creator_id = ?", clientType, creatorId).First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF阅读器设置失败", "client_type", clientType, "creator_id", creatorId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &setting, nil
}

// UpdatePdfReaderSetting 更新PDF阅读器设置
func (d *PdfReaderSettingDAO) UpdatePdfReaderSetting(ctx context.Context, setting *model.PdfReaderSetting) error {
	return d.Modify(ctx, setting)
}

// DeletePdfReaderSetting 删除PDF阅读器设置
func (d *PdfReaderSettingDAO) DeletePdfReaderSetting(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// ListPdfReaderSettings 列出PDF阅读器设置
func (d *PdfReaderSettingDAO) ListPdfReaderSettings(ctx context.Context, creatorId string, limit, offset int) ([]model.PdfReaderSetting, error) {
	var settings []model.PdfReaderSetting
	result := d.GetDB(ctx).Where("creator_id = ?", creatorId).Offset(offset).Limit(limit).Order("id DESC").Find(&settings)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF阅读器设置失败", "creator_id", creatorId, "error", result.Error.Error())
		return nil, result.Error
	}
	return settings, nil
}

// CountPdfReaderSettings 获取PDF阅读器设置总数
func (d *PdfReaderSettingDAO) CountPdfReaderSettings(ctx context.Context, creatorId string) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfReaderSetting{}).Where("creator_id = ?", creatorId).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF阅读器设置总数失败", "creator_id", creatorId, "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}
