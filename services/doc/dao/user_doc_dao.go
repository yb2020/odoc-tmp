package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// UserDocDAO GORM实现的用户文档DAO
type UserDocDAO struct {
	*baseDao.GormBaseDAO[model.UserDoc]
	logger logging.Logger
}

// NewUserDocDAO 创建一个新的用户文档DAO
func NewUserDocDAO(db *gorm.DB, logger logging.Logger) *UserDocDAO {
	return &UserDocDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.UserDoc](db, logger),
		logger:      logger,
	}
}

func (d *UserDocDAO) GetUserDocByUserIdAndPdfId(ctx context.Context, userId string, pdfId string) (*model.UserDoc, error) {
	var doc model.UserDoc
	result := d.GetDB(ctx).Where("user_id = ? AND pdf_id = ? and is_deleted = false", userId, pdfId).First(&doc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户文档失败", "userId", userId, "pdfId", pdfId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &doc, nil
}

// GetUserDocByUserIDAndPaperID 根据用户ID和论文ID获取用户文档
func (d *UserDocDAO) GetUserDocByUserIDAndPaperID(ctx context.Context, userID string, paperID string) (*model.UserDoc, error) {
	var doc model.UserDoc
	result := d.GetDB(ctx).Where("user_id = ? AND paper_id = ? and is_deleted = false", userID, paperID).First(&doc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户文档失败", "userID", userID, "paperID", paperID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &doc, nil
}

// 根据用户id和文献名称查询是否存在未删除的记录
func (d *UserDocDAO) GetUserDocByUserIdAndFileName(ctx context.Context, userId string, fileName string) (*model.UserDoc, error) {
	var doc model.UserDoc
	result := d.GetDB(ctx).Where("user_id = ? AND doc_name = ? and is_deleted = false", userId, fileName).First(&doc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户文档失败", "userId", userId, "fileName", fileName, "error", result.Error.Error())
		return nil, result.Error
	}
	return &doc, nil
}

// GetUserDocsByUserID 根据用户ID查询用户文档列表
func (d *UserDocDAO) GetUserDocsByUserID(ctx context.Context, userId string) ([]model.UserDoc, error) {
	var docs []model.UserDoc
	result := d.GetDB(ctx).Where("user_id = ? AND is_deleted = false", userId).Find(&docs)
	if result.Error != nil {
		d.logger.Error("msg", "根据用户ID查询用户文档列表失败", "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return docs, nil
}

func (d *UserDocDAO) GetUserDocByNoteId(ctx context.Context, noteId string) (*model.UserDoc, error) {
	var doc model.UserDoc
	result := d.GetDB(ctx).Where("note_id = ? and is_deleted = false", noteId).First(&doc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户文档失败", "noteId", noteId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &doc, nil
}

func (d *UserDocDAO) GetUserDocByPdfId(ctx context.Context, pdfId string) (*model.UserDoc, error) {
	var doc model.UserDoc
	result := d.GetDB(ctx).Where("pdf_id = ? and is_deleted = false", pdfId).First(&doc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户文档失败", "pdfId", pdfId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &doc, nil
}

// GetAllUserDocsByUserID 获取用户所有未删除的文档列表（不分页）
func (d *UserDocDAO) GetAllUserDocsByUserID(ctx context.Context, userID string) ([]model.UserDoc, error) {
	var docs []model.UserDoc
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false", userID).Find(&docs)
	if result.Error != nil {
		d.logger.Error("msg", "获取用户所有文档失败", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return docs, nil
}

// GetDocIdsByUserID 获取用户所有未删除的文档ID列表
func (d *UserDocDAO) GetDocIdsByUserID(ctx context.Context, userID string) ([]string, error) {
	var docIds []string
	result := d.GetDB(ctx).Model(&model.UserDoc{}).Where("user_id = ? and is_deleted = false", userID).Pluck("id", &docIds)
	if result.Error != nil {
		d.logger.Error("msg", "获取用户文档ID列表失败", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return docIds, nil
}

// GetUserDocsByUserID 根据用户ID查询用户文档列表(含Ids)
func (d *UserDocDAO) GetByUserIdAndWithIds(ctx context.Context, userId string, ids []string) ([]model.UserDoc, error) {
	var docs []model.UserDoc
	result := d.GetDB(ctx).Where("user_id = ? AND is_deleted = false AND id IN(?)", userId, ids).Find(&docs)
	if result.Error != nil {
		d.logger.Error("msg", "根据用户ID查询用户文档列表失败", "userId", userId, "ids", ids, "error", result.Error.Error())
		return nil, result.Error
	}
	return docs, nil
}

// GetUserDocsByUserID 根据用户IDS查询用户文档列表
func (d *UserDocDAO) GetUserDocsByIds(ctx context.Context, ids []string) ([]model.UserDoc, error) {
	var docs []model.UserDoc
	result := d.GetDB(ctx).Where("id IN(?) AND is_deleted = false", ids).Find(&docs)
	if result.Error != nil {
		d.logger.Error("msg", "根据用户ID查询用户文档列表失败", "ids", ids, "error", result.Error.Error())
		return nil, result.Error
	}
	return docs, nil
}

func (d *UserDocDAO) GetUserDocsByUserIDAndLatestReadTimeIsNotNull(ctx context.Context, userId string, latestReadSize int) ([]model.UserDoc, error) {
	var docs []model.UserDoc
	result := d.GetDB(ctx).Where("user_id = ? AND is_deleted = false AND last_read_time IS NOT NULL ", userId).Order("last_read_time DESC").Limit(latestReadSize).Find(&docs)
	if result.Error != nil {
		d.logger.Error("msg", "根据用户ID查询用户文档列表失败", "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return docs, nil
}

// GetUsersNotParsedByFileSHA256 根据文件SHA256查询未解析完成的用户ID列表
func (d *UserDocDAO) GetUsersNotParsedByFileSHA256(ctx context.Context, fileSHA256 string, parseStatus int32) ([]string, error) {
	var userIds []string

	// 通过fileSHA256查询t_paper_pdf的记录，然后left join t_user_doc表
	// 筛选parse_status不等于指定状态的记录，获取creator_id列表
	result := d.GetDB(ctx).Table("t_paper_pdf").
		Select("DISTINCT t_user_doc.creator_id").
		Joins("LEFT JOIN t_user_doc ON t_paper_pdf.id = t_user_doc.pdf_id").
		Where("t_paper_pdf.file_sha256 = ? AND t_user_doc.parse_status != ? AND t_user_doc.is_deleted = false",
			fileSHA256, parseStatus).
		Pluck("creator_id", &userIds)

	if result.Error != nil {
		d.logger.Error("msg", "根据文件SHA256查询未解析完成的用户列表失败", "fileSHA256", fileSHA256, "error", result.Error.Error())
		return nil, result.Error
	}

	return userIds, nil
}
