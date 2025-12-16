package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	docpb "github.com/yb2020/odoc/proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/dao"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// UserDocClassifyService 用户文档分类服务实现
type UserDocClassifyService struct {
	userDocClassifyDAO      *dao.UserDocClassifyDAO
	classifyRelationService *DocClassifyRelationService
	logger                  logging.Logger
	tracer                  opentracing.Tracer
}

// NewUserDocClassifyService 创建新的用户文档分类服务
func NewUserDocClassifyService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	userDocClassifyDAO *dao.UserDocClassifyDAO,
	classifyRelationService *DocClassifyRelationService,
) *UserDocClassifyService {
	return &UserDocClassifyService{
		logger:                  logger,
		tracer:                  tracer,
		userDocClassifyDAO:      userDocClassifyDAO,
		classifyRelationService: classifyRelationService,
	}
}

// GetUserDocClassifies 获取用户文档分类列表
func (s *UserDocClassifyService) GetUserDocClassifies(ctx context.Context, userId string) ([]model.UserDocClassify, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocClassifyService.GetUserDocClassifies")
	defer span.Finish()

	// 调用DAO层查询用户分类列表
	classifies, err := s.userDocClassifyDAO.GetUserDocClassifiesByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取用户文档分类列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Wrap(err, "获取用户文档分类列表失败")
	}

	// 直接返回模型对象
	return classifies, nil
}

func (s *UserDocClassifyService) Save(ctx context.Context, userDocClassify *model.UserDocClassify) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocClassifyService.Save")
	defer span.Finish()

	// 调用DAO层保存用户分类
	return s.userDocClassifyDAO.Save(ctx, userDocClassify)
}

func (s *UserDocClassifyService) FindById(ctx context.Context, classifyId string) (*model.UserDocClassify, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocClassifyService.FindById")
	defer span.Finish()

	// 调用DAO层查询用户分类
	return s.userDocClassifyDAO.FindById(ctx, classifyId)
}

func (s *UserDocClassifyService) DeleteById(ctx context.Context, classifyId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocClassifyService.DeleteById")
	defer span.Finish()

	// 调用DAO层删除用户分类
	return s.userDocClassifyDAO.DeleteById(ctx, classifyId)
}

// 添加用户文档分类
func (s *UserDocClassifyService) AddUserDocClassify(ctx context.Context, userId string, classifyName string, remark string) (*model.UserDocClassify, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocClassifyService.AddUserDocClassify")
	defer span.Finish()

	classifyId := idgen.GenerateUUID()
	UserDocClassify := model.UserDocClassify{
		UserId: userId,
		Name:   classifyName,
		Remark: remark,
	}
	UserDocClassify.Id = classifyId
	err := s.Save(ctx, &UserDocClassify)
	if err != nil {
		return nil, errors.Biz("doc.user_doc.errors.update_failed")
	}
	return &UserDocClassify, nil
}

// 删除用户文档分类
func (s *UserDocClassifyService) DeleteUserDocClassify(ctx context.Context, userId string, classifyId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocClassifyService.DeleteUserDocClassify")
	defer span.Finish()

	// 使用事务保证删除分类和关联关系的原子性
	err := s.userDocClassifyDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务对象放入 context，确保内部操作使用同一个事务连接
		// 这对于 SQLite 尤其重要，因为 SQLite 使用数据库级锁
		txCtx := context.WithValue(ctx, baseDao.TransactionContextKey, tx)

		// 1. 先删除对应关系
		err := s.classifyRelationService.DeleteByClassifyId(txCtx, classifyId)
		if err != nil {
			s.logger.Error("msg", "delete doc classify relations by classify id failed", "classifyId", classifyId, "error", err)
			return errors.Biz("doc.user_doc.errors.delete_relation_failed")
		}

		// 2. 然后删除分类
		err = s.DeleteById(txCtx, classifyId)
		if err != nil {
			s.logger.Error("msg", "delete classify id failed", "classifyId", classifyId, "error", err)
			return errors.Biz("doc.user_doc.errors.delete_failed")
		}

		return nil
	})

	if err != nil {
		s.logger.Error("msg", "transaction failed when deleting user doc classify", "classifyId", classifyId, "error", err)
		return err
	}

	return nil
}

// GetUserAllClassifyList 获取用户所有的文档分类列表
func (s *UserDocClassifyService) GetUserAllClassifyList(ctx context.Context, userId string) (*docpb.GetUserAllClassifyListResp, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocClassifyService.GetUserAllClassifyList")
	defer span.Finish()

	// 调用userDocClassifyService获取用户分类列表
	docClassifies, err := s.GetUserDocClassifies(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取用户文档分类列表失败", "userId", userId, "error", err.Error())
		return nil, err
	}

	// 如果分类列表为空，返回空列表
	if len(docClassifies) == 0 {
		return &docpb.GetUserAllClassifyListResp{
			Results: []*docpb.UserDocClassify{},
		}, nil
	}

	// 创建结果列表
	results := make([]*docpb.UserDocClassify, 0, len(docClassifies))

	// 遍历分类列表，转换为Proto对象
	for _, classify := range docClassifies {
		vo := &docpb.UserDocClassify{
			ClassifyId:   classify.Id,
			ClassifyName: classify.Name,
		}
		results = append(results, vo)
	}
	// 返回结果
	return &docpb.GetUserAllClassifyListResp{
		Results: results,
	}, nil
}
