package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperAnswerDAO 提供论文问答答案数据访问功能
type PaperAnswerDAO struct {
	*baseDao.GormBaseDAO[model.PaperAnswer]
	logger logging.Logger
}

// NewPaperAnswerDAO 创建一个新的论文问答答案DAO
func NewPaperAnswerDAO(db *gorm.DB, logger logging.Logger) *PaperAnswerDAO {
	return &PaperAnswerDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperAnswer](db, logger),
		logger:      logger,
	}
}

// Create 创建论文问答答案
func (d *PaperAnswerDAO) Create(ctx context.Context, answer *model.PaperAnswer) error {
	return d.GetDB(ctx).Create(answer).Error
}

// FindById 根据ID获取论文问答答案
func (d *PaperAnswerDAO) FindById(ctx context.Context, id string) (*model.PaperAnswer, error) {
	var answer model.PaperAnswer
	result := d.GetDB(ctx).Where("id = ?", id).First(&answer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文问答答案失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &answer, nil
}

// FindByQuestionId 根据问题ID获取论文问答答案列表
func (d *PaperAnswerDAO) FindByQuestionId(ctx context.Context, questionId string) ([]model.PaperAnswer, error) {
	var answers []model.PaperAnswer
	result := d.GetDB(ctx).Where("question_id = ?", questionId).Find(&answers)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文问答答案列表失败", "question_id", questionId, "error", result.Error.Error())
		return nil, result.Error
	}
	return answers, nil
}

// FindByReplyUserId 根据回复用户ID获取论文问答答案列表
func (d *PaperAnswerDAO) FindByReplyUserId(ctx context.Context, replyUserId string) ([]model.PaperAnswer, error) {
	var answers []model.PaperAnswer
	result := d.GetDB(ctx).Where("reply_user_id = ?", replyUserId).Find(&answers)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文问答答案列表失败", "reply_user_id", replyUserId, "error", result.Error.Error())
		return nil, result.Error
	}
	return answers, nil
}

// FindByReplyAnswerId 根据回复答案ID获取论文问答答案列表
func (d *PaperAnswerDAO) FindByReplyAnswerId(ctx context.Context, replyAnswerId string) ([]model.PaperAnswer, error) {
	var answers []model.PaperAnswer
	result := d.GetDB(ctx).Where("reply_answer_id = ?", replyAnswerId).Find(&answers)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文问答答案列表失败", "reply_answer_id", replyAnswerId, "error", result.Error.Error())
		return nil, result.Error
	}
	return answers, nil
}

// UpdateById 更新论文问答答案
func (d *PaperAnswerDAO) UpdateById(ctx context.Context, answer *model.PaperAnswer) error {
	return d.Modify(ctx, answer)
}

// DeleteById 删除论文问答答案
func (d *PaperAnswerDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}
