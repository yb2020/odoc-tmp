package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/translate/model"
	"gorm.io/gorm"
)

// FullTextTranslateFixDAO 翻译修复DAO接口
type FullTextTranslateFixDAO interface {
	// Save 保存翻译修复记录
	Save(ctx context.Context, fix *model.FullTextTranslateFix) error

	// FindById 根据ID查询翻译修复记录
	FindById(ctx context.Context, id string) (*model.FullTextTranslateFix, error)

	// FindByNoteId 根据笔记ID查询翻译修复记录
	FindByNoteId(ctx context.Context, noteId string) ([]*model.FullTextTranslateFix, error)

	// FindAll 查询所有翻译修复记录，支持分页
	FindAll(ctx context.Context, id, noteId string) ([]*model.FullTextTranslateFix, error)
}

// FullTextTranslateFixDAOImpl 翻译修复DAO实现
type FullTextTranslateFixDAOImpl struct {
	*baseDao.GormBaseDAO[model.FullTextTranslateFix]
	logger logging.Logger
}

// NewFullTextTranslateFixDAO 创建一个新的翻译修复DAO
func NewFullTextTranslateFixDAO(db *gorm.DB, logger logging.Logger) FullTextTranslateFixDAO {
	return &FullTextTranslateFixDAOImpl{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.FullTextTranslateFix](db, logger),
		logger:      logger,
	}
}

// Save 保存翻译修复记录
func (dao *FullTextTranslateFixDAOImpl) Save(ctx context.Context, fix *model.FullTextTranslateFix) error {
	return dao.GormBaseDAO.Save(ctx, fix)
}

// FindById 根据ID查询翻译修复记录
func (dao *FullTextTranslateFixDAOImpl) FindById(ctx context.Context, id string) (*model.FullTextTranslateFix, error) {
	var fix model.FullTextTranslateFix
	err := dao.GetDB(ctx).
		Where("id = ? AND is_delete = 0", id).
		First(&fix).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		dao.logger.Error("查询翻译修复记录失败", "error", err.Error(), "id", id)
		return nil, err
	}
	return &fix, nil
}

// FindByNoteId 根据笔记ID查询翻译修复记录
func (dao *FullTextTranslateFixDAOImpl) FindByNoteId(ctx context.Context, noteId string) ([]*model.FullTextTranslateFix, error) {
	var fixes []*model.FullTextTranslateFix
	err := dao.GetDB(ctx).
		Where("note_id = ? AND is_delete = 0", noteId).
		Order("create_time DESC").
		Find(&fixes).Error
	if err != nil {
		dao.logger.Error("查询翻译修复记录失败", "error", err.Error(), "noteId", noteId)
		return nil, err
	}
	return fixes, nil
}

// FindAll 查询所有翻译修复记录，支持分页
func (dao *FullTextTranslateFixDAOImpl) FindAll(ctx context.Context, id, noteId string) ([]*model.FullTextTranslateFix, error) {
	db := dao.GetDB(ctx).Where("is_delete = 0")

	if id != "" {
		db = db.Where("id = ?", id)
	}

	if noteId != "" {
		db = db.Where("note_id = ?", noteId)
	}

	var fixes []*model.FullTextTranslateFix
	err := db.Order("create_time DESC").Find(&fixes).Error
	if err != nil {
		dao.logger.Error("查询翻译修复记录失败", "error", err.Error())
		return nil, err
	}
	return fixes, nil
}
