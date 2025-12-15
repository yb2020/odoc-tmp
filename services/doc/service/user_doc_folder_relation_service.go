package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	// 使用 ClientDoc.pb.go 中的定义
	clientdocpb "github.com/yb2020/odoc-proto/gen/go/doc"
	userContext "github.com/yb2020/odoc/pkg/context"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/dao"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// UserDocFolderRelationService 用户文档文件夹关系服务实现
type UserDocFolderRelationService struct {
	userDocFolderRelationDAO *dao.UserDocFolderRelationDAO
	userDocFolderDAO         *dao.UserDocFolderDAO
	logger                   logging.Logger
	tracer                   opentracing.Tracer
}

// NewUserDocFolderRelationService 创建新的用户文档文件夹关系服务
func NewUserDocFolderRelationService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	userDocFolderRelationDAO *dao.UserDocFolderRelationDAO,
	userDocFolderDAO *dao.UserDocFolderDAO,
) *UserDocFolderRelationService {
	return &UserDocFolderRelationService{
		logger:                   logger,
		tracer:                   tracer,
		userDocFolderRelationDAO: userDocFolderRelationDAO,
		userDocFolderDAO:         userDocFolderDAO,
	}
}

func (s *UserDocFolderRelationService) Save(ctx context.Context, userDocFolderRelation *model.UserDocFolderRelation) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.Save")
	defer span.Finish()

	return s.userDocFolderRelationDAO.Save(ctx, userDocFolderRelation)
}

func (s *UserDocFolderRelationService) SaveBatch(ctx context.Context, userDocFolderRelations []model.UserDocFolderRelation) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.SaveBatch")
	defer span.Finish()

	return s.userDocFolderRelationDAO.SaveBatch(ctx, userDocFolderRelations)
}

// GetUserDocFolderRelationsByUserID 根据用户ID获取用户文档文件夹关系列表
func (s *UserDocFolderRelationService) GetUserDocFolderRelationsByUserID(ctx context.Context, userId string) ([]model.UserDocFolderRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.GetUserDocFolderRelationsByUserID")
	defer span.Finish()

	relations, err := s.userDocFolderRelationDAO.GetUserDocFolderRelationsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取用户文档文件夹关系列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Wrap(err, "获取用户文档文件夹关系列表失败")
	}
	// 直接返回模型对象
	return relations, nil
}

// GetUserDocFolderRelationsByFolderId 根据用户ID和文件夹ID获取文件夹-文献关系列表
func (s *UserDocFolderRelationService) GetUserDocFolderRelationsByFolderId(ctx context.Context, userId string, folderId string) ([]model.UserDocFolderRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.GetUserDocFolderRelationsByFolderId")
	defer span.Finish()

	relations, err := s.userDocFolderRelationDAO.GetUserDocFolderRelationsByFolderId(ctx, userId, folderId)
	if err != nil {
		s.logger.Error("msg", "获取文件夹-文献关系列表失败", "userId", userId, "folderId", folderId, "error", err.Error())
		return nil, errors.Wrap(err, "获取文件夹-文献关系列表失败")
	}
	// 直接返回模型对象
	return relations, nil
}

// GetUserDocFolderRelationsByFolderIdAndDocId 根据用户ID和文件夹ID和文档id获取文献关系
func (s *UserDocFolderRelationService) GetUserDocFolderRelationsByFolderIdAndDocId(ctx context.Context, userId string, folderId string, docId string) (*model.UserDocFolderRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.GetUserDocFolderRelationsByFolderIdAndDocId")
	defer span.Finish()

	relations, err := s.userDocFolderRelationDAO.GetUserDocFolderRelationsByFolderIdAndDocId(ctx, userId, folderId, docId)
	if err != nil {
		return nil, errors.Wrap(err, "get user doc folder relations by folder id and doc id failed")
	}
	// 直接返回模型对象
	return relations, nil
}

// GetUserDocFolderRelationsByFolderIDs 根据文件夹ID列表获取用户文档文件夹关系列表
func (s *UserDocFolderRelationService) GetUserDocFolderRelationsByFolderIDs(ctx context.Context, folderIDs []string) ([]model.UserDocFolderRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.GetUserDocFolderRelationsByFolderIDs")
	defer span.Finish()

	// 使用 DAO 层方法获取关系列表
	relations, err := s.userDocFolderRelationDAO.GetUserDocFolderRelationsByFolderIDs(ctx, folderIDs)
	if err != nil {
		return nil, errors.Wrap(err, "根据文件夹ID列表获取文档关系失败")
	}
	// 直接返回模型对象
	return relations, nil
}

// GetRelationsByDocIds 根据文献ID列表物理查询用户文档文件夹关系
func (s *UserDocFolderRelationService) GetRelationsByUserIdAndDocIds(ctx context.Context, userId string, docIds []string) ([]model.UserDocFolderRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.DeleteRelationsByDocIds")
	defer span.Finish()

	if len(docIds) == 0 {
		return nil, nil
	}

	// 查询数据列表
	relations, err := s.userDocFolderRelationDAO.GetRelationsByUserIdAndDocIds(ctx, userId, docIds)
	if err != nil {
		s.logger.Error("msg", "批量删除文献关系失败", "docIds", docIds, "error", err.Error())
		return nil, errors.Wrap(err, "批量删除文献关系失败")
	}

	return relations, nil
}

// GetMaxSort 获取指定用户和文件夹下的最大排序值
func (s *UserDocFolderRelationService) GetMaxSort(ctx context.Context, userID string, folderID string) (int32, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.GetMaxSort")
	defer span.Finish()

	// 使用 DAO 层方法获取最大排序值
	maxSort, err := s.userDocFolderRelationDAO.GetMaxSort(ctx, userID, folderID)
	if err != nil {
		return 0, errors.Wrap(err, "获取最大排序值失败")
	}

	s.logger.Info("msg", "成功获取最大排序值", "userID", userID, "folderID", folderID, "maxSort", maxSort)
	return maxSort, nil
}

// CreateUserDocFolderRelation 创建用户文档文件夹关系
func (s *UserDocFolderRelationService) CreateUserDocFolderRelation(ctx context.Context, relation *model.UserDocFolderRelation) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.CreateUserDocFolderRelation")
	defer span.Finish()

	// 如果没有指定排序值，自动获取最大排序值+1
	if relation.Sort == 0 {
		maxSort, err := s.GetMaxSort(ctx, relation.UserId, relation.FolderId)
		if err != nil {
			return errors.Wrap(err, "获取最大排序值失败")
		}
		relation.Sort = maxSort + 1
	}

	// 使用 DAO 层方法创建关系
	err := s.userDocFolderRelationDAO.Save(ctx, relation)
	if err != nil {
		return errors.Wrap(err, "创建用户文档文件夹关系失败")
	}

	s.logger.Info("msg", "成功创建用户文档文件夹关系", "relation", relation)
	return nil
}

// DeleteRelationsByFolderIds 根据文件夹ID列表物理删除用户文档文件夹关系
func (s *UserDocFolderRelationService) DeleteRelationsByFolderIds(ctx context.Context, folderIds []string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.DeleteRelationsByFolderIds")
	defer span.Finish()

	if len(folderIds) == 0 {
		return nil
	}

	// 使用 DAO 层方法批量物理删除关系
	err := s.userDocFolderRelationDAO.BatchRemoveByFolderIds(ctx, folderIds)
	if err != nil {
		return errors.Wrap(err, "根据文件夹ID列表物理删除用户文档文件夹关系失败")
	}

	return nil
}

// 根据用户id、文件夹id、文献id物理删除
func (s *UserDocFolderRelationService) DeleteRelationsByUserIdAndFolderIdAndDocId(ctx context.Context, userId string, folderId string, docId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.DeleteRelationsByUserIdAndFolderIdAndDocId")
	defer span.Finish()

	// 批量删除关系
	err := s.userDocFolderRelationDAO.RemoveByUserIdAndFolderIdAndDocId(ctx, userId, folderId, docId)
	if err != nil {
		return errors.Wrap(err, "delete user doc folder relation failed")
	}
	return nil
}

// 根据用户id、文件夹id列表、文献id物理删除
func (s *UserDocFolderRelationService) DeleteRelationsByUserIdAndFolderIdsAndDocId(ctx context.Context, userId string, folderIds []string, docId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.DeleteRelationsByUserIdAndFolderIdsAndDocId")
	defer span.Finish()

	// 批量删除关系
	err := s.userDocFolderRelationDAO.RemoveByUserIdAndFolderIdsAndDocId(ctx, userId, folderIds, docId)
	if err != nil {
		return errors.Wrap(err, "delete user doc folder relation failed")
	}
	return nil
}

// 根据用户id、文献id物理删除
func (s *UserDocFolderRelationService) DeleteRelationsByUserIdAndDocIds(ctx context.Context, userId string, docId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.DeleteRelationsByUserIdAndDocIds")
	defer span.Finish()

	// 批量删除关系
	err := s.userDocFolderRelationDAO.RemoveByUserIdAndDocId(ctx, userId, docId)
	if err != nil {
		return errors.Wrap(err, "delete user doc folder relation failed")
	}
	return nil
}

// DeleteRelationsByUserIdAndFolderIds 根据用户ID与文件夹ID列表删除关系（物理删除，仅限该用户）
func (s *UserDocFolderRelationService) DeleteRelationsByUserIdAndFolderIds(ctx context.Context, userId string, folderIds []string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.DeleteRelationsByUserIdAndFolderIds")
	defer span.Finish()

	if len(folderIds) == 0 {
		return nil
	}
	return s.userDocFolderRelationDAO.RemoveByUserIdAndFolderIds(ctx, userId, folderIds)
}

// MoveDoc 移动文献到另一个文件夹
func (s *UserDocFolderRelationService) MoveDoc(ctx context.Context, req *clientdocpb.MoveFolderOrDocReq, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.MoveDoc")
	defer span.Finish()
	// 验证请求参数
	if len(req.TargetFolderItems) == 0 {
		return nil
	}
	if req.TargetFolderId == "0" {
		return errors.Biz("personal.doc.move.not.support")
	}
	// 不能从已有文件夹移到未分类
	if req.TargetFolderId != req.SourceFolderId && req.TargetFolderId == "0" {
		return errors.Biz("personal.doc.move.not.support")
	}
	// 验证文献存在并属于当前用户
	if len(req.MovedIds) == 0 {
		return errors.Biz("param.error")
	}
	// 将uint64类型的docIds转换为int64类型
	docIds := make([]string, 0, len(req.MovedIds))
	for _, id := range req.MovedIds {
		docIds = append(docIds, id)
	}
	// // 查询用户文献关系
	// dbDocFolderRelations, err := s.userDocFolderRelationDAO.FindByUserIdAndDocIds(ctx, userId, docIds)
	// if err != nil {
	// 	return errors.Wrap(err, "获取文献关系失败")
	// }
	// if len(dbDocFolderRelations) == 0 {
	// 	return errors.Biz("personal.doc.not.right")
	// }
	// 查询目标文件夹中的文献关系
	targetFolderRelations, err := s.userDocFolderRelationDAO.FindByUserIdAndFolderId(ctx, userId, req.TargetFolderId)
	if err != nil {
		return errors.Wrap(err, "获取目标文件夹关系失败")
	}

	// 检查目标文件夹中是否已存在相同的文献
	if req.TargetFolderId != req.SourceFolderId {
		existSameNameDoc := false
		for _, relation := range targetFolderRelations {
			for _, docId := range docIds {
				if relation.DocId == docId {
					existSameNameDoc = true
					break
				}
			}
			if existSameNameDoc {
				break
			}
		}
		if existSameNameDoc {
			return errors.Biz("personal.doc.folder.exist.same")
		}
	}
	// 使用事务执行移动操作
	return s.userDocFolderRelationDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务对象放入 context，确保内部操作使用同一个事务连接
		// 这对于 SQLite 尤其重要，因为 SQLite 使用数据库级锁
		txCtx := context.WithValue(ctx, baseDao.TransactionContextKey, tx)

		if req.TargetFolderId != req.SourceFolderId {
			// 删除旧的关系
			err := s.userDocFolderRelationDAO.DeleteByUserIdAndFolderIdAndDocIds(txCtx, userId, req.SourceFolderId, docIds)
			if err != nil {
				return errors.Wrap(err, "删除旧关系失败")
			}
			// 新增关系
			for _, docId := range docIds {
				relation := model.UserDocFolderRelation{
					UserId:   userId,
					FolderId: req.TargetFolderId,
					DocId:    docId,
				}
				// 保存新关系
				s.userDocFolderRelationDAO.Save(txCtx, &relation)
				targetFolderRelations = append(targetFolderRelations, relation)
			}
		}

		// 更新排序
		if len(req.TargetFolderItems) > 0 {
			for i, item := range req.TargetFolderItems {
				// 查找对应的关系
				var relationList []model.UserDocFolderRelation
				for _, relation := range targetFolderRelations {
					if relation.DocId == item.Id {
						relationList = append(relationList, relation)
					}
				}
				if len(relationList) == 0 {
					return errors.Biz("param.error")
				}
				// 更新排序
				relation := relationList[0]
				relation.Sort = int32(len(req.TargetFolderItems) - 1 - i)
				s.userDocFolderRelationDAO.ModifyExcludeNull(txCtx, &relation)
			}
		}

		return nil
	})
}

// MoveDocToAnotherFolder 将文档移动到另一个文件夹
func (s *UserDocFolderRelationService) MoveDocToAnotherFolder(ctx context.Context, req *clientdocpb.MoveDocOrFolderToAnotherFolderReq) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.MoveDocToAnotherFolder")
	defer span.Finish()

	if len(req.MovedDocItems) == 0 {
		return nil
	}

	// 获取用户上下文
	uc := userContext.GetUserContext(ctx)
	if uc == nil || uc.UserId == "0" {
		return errors.Biz("doc.user_doc_folder_relation.errors.user_context_not_found")
	}

	// 使用事务执行移动操作
	return s.userDocFolderRelationDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务对象放入 context，确保内部操作使用同一个事务连接
		// 这对于 SQLite 尤其重要，因为 SQLite 使用数据库级锁
		txCtx := context.WithValue(ctx, baseDao.TransactionContextKey, tx)

		for _, item := range req.MovedDocItems {
			docId := item.DocId
			sourceFolderId := item.SourceFolderId
			targetFolderId := req.TargetFolderId

			// 如果目标文件夹与源文件夹相同，则跳过
			if sourceFolderId == targetFolderId {
				continue
			}
			// 如果源文件夹ID不为0，则删除原关系
			if sourceFolderId != "0" {
				// 删除原关系（物理删除）
				s.userDocFolderRelationDAO.RemoveById(txCtx, sourceFolderId)
			}
			// 检查目标文件夹中是否已存在该文档的关系
			userDocFolderRelationObj, err := s.userDocFolderRelationDAO.GetUserDocFolderRelationsByFolderIdAndDocId(txCtx, uc.UserId, sourceFolderId, docId)
			if err != nil {
				return errors.Wrap(err, "update user doc folder relation failed")
			}
			if userDocFolderRelationObj == nil {
				return errors.Wrap(err, "update user doc folder relation failed")
			}
			// 更新目标文件夹中的文档关系
			userDocFolderRelationObj.FolderId = targetFolderId
			s.userDocFolderRelationDAO.ModifyExcludeNull(txCtx, userDocFolderRelationObj)
		}
		return nil
	})
}

// MoveRelationsToUnclassifiedByFolderIds 根据文件夹ID列表将关系移动到未分类（folder_id = 0）
func (s *UserDocFolderRelationService) MoveRelationsToUnclassifiedByFolderIds(ctx context.Context, folderIds []string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderRelationService.MoveRelationsToUnclassifiedByFolderIds")
	defer span.Finish()

	if len(folderIds) == 0 {
		return nil
	}
	return s.userDocFolderRelationDAO.MoveRelationsToUnclassifiedByFolderIds(ctx, folderIds)
}
