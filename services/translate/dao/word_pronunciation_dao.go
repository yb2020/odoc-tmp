package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/translate/model"
	"gorm.io/gorm"
)

// WordPronunciationDAO GORM实现的单词发音DAO
type WordPronunciationDAO struct {
	*baseDao.GormBaseDAO[model.WordPronunciation]
	logger logging.Logger
}

// NewWordPronunciationDAO 创建一个新的GORM单词发音DAO
func NewWordPronunciationDAO(db *gorm.DB, logger logging.Logger) WordPronunciationDAO {
	return WordPronunciationDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.WordPronunciation](db, logger),
		logger:      logger,
	}
}

// FindByTargetContentAndSource 根据目标内容和源查找单词发音
func (d *WordPronunciationDAO) FindByTargetContentAndSource(ctx context.Context, targetContent, source string) ([]model.WordPronunciation, error) {
	var wordTranslates []model.WordPronunciation
	result := d.GetDB(ctx).
		Where("target_content = ? AND source = ?", targetContent, source).
		Find(&wordTranslates)

	if result.Error != nil {
		d.logger.Error("msg", "根据目标内容和源查找单词发音失败", "targetContent", targetContent, "source", source, "error", result.Error.Error())
		return nil, result.Error
	}
	return wordTranslates, nil
}
