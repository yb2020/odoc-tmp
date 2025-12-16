package service

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	// 使用 ClientDoc.pb.go 中的定义
	baseDao "github.com/yb2020/odoc/pkg/dao"
	errors "github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	clientdocpb "github.com/yb2020/odoc/proto/gen/go/doc"
	docpb "github.com/yb2020/odoc/proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/bean"
	"github.com/yb2020/odoc/services/doc/dao"
	"github.com/yb2020/odoc/services/doc/model"
)

// UserDocFolderService 用户文档文件夹服务实现
type UserDocFolderService struct {
	userDocFolderDAO             *dao.UserDocFolderDAO
	userDocFolderRelationService *UserDocFolderRelationService
	logger                       logging.Logger
	tracer                       opentracing.Tracer
	transactionManager           *baseDao.TransactionManager
}

// NewUserDocFolderService 创建新的用户文档文件夹服务
func NewUserDocFolderService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	userDocFolderDAO *dao.UserDocFolderDAO,
	transactionManager *baseDao.TransactionManager,
) *UserDocFolderService {
	return &UserDocFolderService{
		logger:             logger,
		tracer:             tracer,
		userDocFolderDAO:   userDocFolderDAO,
		transactionManager: transactionManager,
	}
}

// CreateUserDocFolder 创建用户文档文件夹
func (s *UserDocFolderService) CreateUserDocFolder(ctx context.Context, userId string, req *clientdocpb.CreateUserDocFolderRequest) (*clientdocpb.CreateUserDocFolderResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.CreateUserDocFolder")
	defer span.Finish()
	folderId := idgen.GenerateUUID()
	// 创建文件夹对象
	folder := &model.UserDocFolder{
		UserId:   userId,
		ParentId: req.GetParentId(),
		Name:     req.GetName(),
		Sort:     int32(req.GetSort()),
	}

	// 验证父文件夹ID
	if folder.ParentId != "0" {
		parentFolder, err := s.userDocFolderDAO.GetUserDocFolderByID(ctx, folder.ParentId)
		if err != nil {
			return nil, errors.Wrap(err, "get parent folder failed")
		}
		if parentFolder == nil {
			folder.ParentId = "0" // 如果父文件夹不存在，则设置为根目录
		} else if parentFolder.UserId != folder.UserId {
			return nil, errors.Biz("doc.user_doc_folder.errors.not_belong_to_current_user")
		}
	}

	if err := s.userDocFolderDAO.Save(ctx, folder); err != nil {
		return nil, errors.Wrap(err, "create user doc folder failed")
	}

	// 创建响应对象并设置属性
	// 将 folderId 从 int64 转换为 uint64
	resp := &clientdocpb.CreateUserDocFolderResponse{
		Id:   folderId,
		Name: folder.Name,
	}

	return resp, nil
}

// GetUserDocFolders 获取用户文档文件夹列表
func (s *UserDocFolderService) GetUserDocFolders(ctx context.Context, userId string) ([]model.UserDocFolder, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.GetUserDocFolders")
	defer span.Finish()

	// 调用DAO层查询用户文件夹列表
	folders, err := s.userDocFolderDAO.GetUserDocFoldersByUserID(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "get user doc folders failed")
	}
	return folders, nil
}

// GetAllAuthDocFolderByIds 获取指定文件夹ID列表下的所有子孙级文件夹
func (s *UserDocFolderService) GetAllAuthDocFolderByIds(ctx context.Context, docFolderIds []string, userId string) ([]model.UserDocFolder, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.GetAllAuthDocFolderByIds")
	defer span.Finish()

	if len(docFolderIds) == 0 {
		return []model.UserDocFolder{}, nil
	}
	// 获取该用户下的所有文件目录
	folders, err := s.userDocFolderDAO.GetUserDocFoldersByUserID(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "get user doc folders failed")
	}
	// 创建文件夹ID到文件夹对象的映射，方便查找
	folderMap := make(map[string]model.UserDocFolder)
	for _, folder := range folders {
		folderMap[folder.Id] = folder
	}
	// 创建父子关系映射
	childrenMap := make(map[string][]string)
	for _, folder := range folders {
		if folder.ParentId != "0" {
			childrenMap[folder.ParentId] = append(childrenMap[folder.ParentId], folder.Id)
		}
	}
	// 用于存储结果的集合
	resultSet := make(map[string]model.UserDocFolder)
	// 递归查找所有子文件夹
	var findAllDescendants func(folderId string)
	findAllDescendants = func(folderId string) {
		// 如果当前文件夹已经在结果集中，则跳过
		if _, exists := resultSet[folderId]; exists {
			return
		}
		// 将当前文件夹添加到结果集
		if folder, exists := folderMap[folderId]; exists {
			resultSet[folderId] = folder
			// 递归处理所有子文件夹
			for _, childId := range childrenMap[folderId] {
				findAllDescendants(childId)
			}
		}
	}
	// 对每个指定的文件夹ID，查找其所有子孙文件夹
	for _, folderId := range docFolderIds {
		findAllDescendants(folderId)
	}
	result := make([]model.UserDocFolder, 0, len(resultSet))
	for _, folder := range resultSet {
		result = append(result, folder)
	}
	return result, nil
}

// DeleteFoldersByIds 根据文件夹ID列表删除文件夹及其所有子文件夹，同时删除关联的文件夹关系
func (s *UserDocFolderService) DeleteFoldersByIds(ctx context.Context, docFolderIds []string, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.DeleteFoldersByIds")
	defer span.Finish()

	if len(docFolderIds) == 0 {
		return nil
	}

	// 获取所有需要删除的文件夹（包括子文件夹）
	allFolders, err := s.GetAllAuthDocFolderByIds(ctx, docFolderIds, userId)
	if err != nil {
		return errors.Wrap(err, "get all auth doc folder by ids failed")
	}

	if len(allFolders) == 0 {
		return nil
	}

	// 提取所有文件夹ID
	folderIds := make([]string, 0, len(allFolders))
	for _, folder := range allFolders {
		folderIds = append(folderIds, folder.Id)
	}

	// 使用事务执行删除操作，确保原子性
	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		if s.userDocFolderRelationService != nil {
			// Step 1: 查询删除目录中该用户所有的 doc_id
			relations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByFolderIDs(txCtx, folderIds)
			if err != nil {
				return errors.Wrap(err, "get relations by folder ids failed")
			}
			docIdSet := make(map[string]struct{})
			for _, r := range relations {
				// 保险起见，只收集当前用户的关系
				if r.UserId == userId {
					docIdSet[r.DocId] = struct{}{}
				}
			}
			docIds := make([]string, 0, len(docIdSet))
			for id := range docIdSet {
				docIds = append(docIds, id)
			}

			// Step 2: 直接删除该用户文件夹关系（按文件夹ID删除关系）
			if err := s.userDocFolderRelationService.DeleteRelationsByFolderIds(txCtx, folderIds); err != nil {
				return errors.Wrap(err, "delete relations by folder ids failed")
			}

			// Step 3: 对于这些 doc_id，判断当前用户是否还有其他关系
			// 查询该用户与这些 doc 的剩余关系
			if len(docIds) > 0 {
				leftRelations, err := s.userDocFolderRelationService.GetRelationsByUserIdAndDocIds(txCtx, userId, docIds)
				if err != nil {
					return errors.Wrap(err, "get relations by user and doc ids failed")
				}
				stillRelated := make(map[string]struct{})
				for _, lr := range leftRelations {
					stillRelated[lr.DocId] = struct{}{}
				}
				// 对没有任何关系的 doc，新增一条 folderId=0 的关系
				for _, docId := range docIds {
					if _, ok := stillRelated[docId]; ok {
						continue
					}
					newRel := &model.UserDocFolderRelation{
						UserId:   userId,
						FolderId: "0",
						DocId:    docId,
						// Sort 由服务层在 CreateUserDocFolderRelation 中自动补全（max+1）
					}
					if err := s.userDocFolderRelationService.CreateUserDocFolderRelation(txCtx, newRel); err != nil {
						return errors.Wrap(err, "create unclassified relation failed")
					}
				}
			}
		}

		// 2. 使用批量删除方法（逻辑删除文件夹）
		err := s.userDocFolderDAO.BatchDeleteByIds(txCtx, folderIds)
		if err != nil {
			return errors.Wrap(err, "batch delete user doc folder failed")
		}
		return nil
	})
}

// DeleteUserDocFolder 根据DeleteUserDocFolderReq请求删除文件夹
func (s *UserDocFolderService) DeleteUserDocFolder(ctx context.Context, req *clientdocpb.DeleteUserDocFolderReq, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.DeleteUserDocFolder")
	defer span.Finish()

	// 将uint64类型的folderIds转换为int64类型
	folderIds := make([]string, 0, len(req.FolderIds))
	for _, id := range req.FolderIds {
		folderIds = append(folderIds, id)
	}

	// 调用DeleteFoldersByIds方法执行删除操作
	return s.DeleteFoldersByIds(ctx, folderIds, userId)
}

// UpdateUserDocFolder 根据UpdateUserDocFolderReq请求更新文件夹信息
func (s *UserDocFolderService) UpdateUserDocFolder(ctx context.Context, req *clientdocpb.UpdateUserDocFolderReq, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.UpdateUserDocFolder")
	defer span.Finish()

	// 将uint64类型的folderId转换为int64类型
	folderId := req.FolderId

	// 验证文件夹存在并属于当前用户
	folder, err := s.userDocFolderDAO.FindExistById(ctx, folderId)
	if err != nil {
		return errors.Wrap(err, "get folder failed")
	}

	if folder == nil {
		return errors.Biz("doc.user_doc_folder.errors.folder_not_found")
	}

	if folder.UserId != userId {
		return errors.Biz("doc.user_doc_folder.errors.not_belong_to_current_user")
	}

	// 创建更新映射
	updates := make(map[string]interface{})

	// 只更新非空字段
	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 调用DAO层执行更新
	err = s.userDocFolderDAO.UpdateFolder(ctx, folderId, updates)
	if err != nil {
		return errors.Wrap(err, "update folder failed")
	}

	return nil
}

// SetUserDocFolderRelationService 设置用户文档文件夹关系服务，用于解决循环依赖问题
func (s *UserDocFolderService) SetUserDocFolderRelationService(relationService *UserDocFolderRelationService) error {
	if relationService == nil {
		return errors.Biz("doc.user_doc_folder.errors.relation_service_nil")
	}
	s.userDocFolderRelationService = relationService
	return nil
}

// MoveDocOrFolderToAnotherFolder 将文档或文件夹移动到另一个文件夹
func (s *UserDocFolderService) MoveDocOrFolderToAnotherFolder(ctx context.Context, req *clientdocpb.MoveDocOrFolderToAnotherFolderReq, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.MoveDocOrFolderToAnotherFolder")
	defer span.Finish()

	// 验证目标文件夹
	targetFolderId := req.TargetFolderId
	if targetFolderId != "0" { // 如果目标文件夹不是根目录
		targetFolder, err := s.userDocFolderDAO.FindExistById(ctx, targetFolderId)
		if err != nil {
			return errors.Wrap(err, "get target folder failed")
		}

		if targetFolder == nil {
			return errors.Biz("doc.user_doc_folder.errors.target_folder_not_found")
		}

		if targetFolder.UserId != userId {
			return errors.Biz("doc.user_doc_folder.errors.target_folder_not_belong_to_current_user")
		}
	}

	// 验证请求参数
	if len(req.MovedFolderIds) > 0 && len(req.MovedDocItems) > 0 {
		return errors.Biz("doc.user_doc_folder.errors.cannot_move_docs_and_folders_at_same_time")
	}

	// 移动文献
	if len(req.MovedDocItems) > 0 {
		// 验证文献参数
		for _, item := range req.MovedDocItems {
			if item.DocId == "0" {
				return errors.Biz("doc.user_doc_folder.errors.invalid_doc_params")
			}
		}
		err := s.userDocFolderRelationService.MoveDocToAnotherFolder(ctx, req)
		if err != nil {
			return errors.Wrap(err, "move doc to another folder failed")
		}
	}
	// 移动文件夹
	if len(req.MovedFolderIds) > 0 {
		// 将uint64类型的folderIds转换为int64类型
		folderIds := make([]string, 0, len(req.MovedFolderIds))
		for _, id := range req.MovedFolderIds {
			folderIds = append(folderIds, id)
		}
		// 验证要移动的文件夹存在并属于当前用户
		folders, err := s.userDocFolderDAO.FindAllByIds(ctx, folderIds)
		if err != nil {
			return errors.Wrap(err, "get folders failed")
		}
		if len(folders) == 0 {
			return errors.Biz("doc.user_doc_folder.errors.folders_not_found")
		}
		// 验证文件夹属于当前用户
		for _, folder := range folders {
			if folder.UserId != userId {
				return errors.Biz("doc.user_doc_folder.errors.folder_not_belong_to_current_user")
			}
		}
		// 使用事务执行移动操作
		err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
			for _, folderId := range folderIds {
				// 如果要移动的文件夹就是目标文件夹，则移动到根目录
				if folderId == targetFolderId {
					updates := map[string]interface{}{
						"parent_id": 0,
					}
					err := s.userDocFolderDAO.UpdateFolder(txCtx, folderId, updates)
					if err != nil {
						return errors.Wrap(err, "update folder parent id failed")
					}
				} else {
					updates := map[string]interface{}{
						"parent_id": targetFolderId,
					}
					err := s.userDocFolderDAO.UpdateFolder(txCtx, folderId, updates)
					if err != nil {
						return errors.Wrap(err, "update folder parent id failed")
					}
				}
			}
			// 获取用户的所有文件夹，检查是否存在循环依赖
			folders, err := s.userDocFolderDAO.GetUserDocFoldersByUserID(txCtx, userId)
			if err != nil {
				return errors.Wrap(err, "get user folders failed")
			}
			// 检查是否存在循环依赖
			hasCircular := s.hasCircularDependency(folders)
			if hasCircular {
				return errors.Biz("doc.user_doc_folder.errors.circular_dependency")
			}

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// hasCircularDependency 检查文件夹列表是否存在循环依赖
func (s *UserDocFolderService) hasCircularDependency(folders []model.UserDocFolder) bool {
	// 构建文件夹的父子关系图
	folderMap := make(map[string]*model.UserDocFolder)
	for i := range folders {
		folderMap[folders[i].Id] = &folders[i]
	}
	// 使用深度优先搜索检测是否有环
	visited := make(map[string]bool)  // 记录已访问过的节点
	recStack := make(map[string]bool) // 记录当前递归调用栈中的节点
	// 对每个节点进行检查
	for _, folder := range folders {
		if !visited[folder.Id] {
			if s.isCyclicUtil(folder.Id, folderMap, visited, recStack) {
				return true
			}
		}
	}
	return false
}

// isCyclicUtil 递归检查是否存在环
func (s *UserDocFolderService) isCyclicUtil(folderId string, folderMap map[string]*model.UserDocFolder, visited map[string]bool, recStack map[string]bool) bool {
	// 标记当前节点为已访问
	visited[folderId] = true
	// 将当前节点加入递归栈
	recStack[folderId] = true
	// 获取当前文件夹
	folder, exists := folderMap[folderId]
	if !exists {
		// 如果文件夹不存在，则不可能有环
		recStack[folderId] = false
		return false
	}
	// 获取父文件夹ID
	parentId := folder.ParentId
	// 如果父文件夹为0，表示这是根目录，不可能有环
	if parentId == "0" {
		recStack[folderId] = false
		return false
	}
	// 检查父文件夹是否在递归栈中，如果在，则存在环
	if recStack[parentId] {
		return true
	}
	// 如果父文件夹已经被访问过，但不在递归栈中，则不存在环
	if visited[parentId] {
		recStack[folderId] = false
		return false
	}
	// 递归检查父文件夹
	if s.isCyclicUtil(parentId, folderMap, visited, recStack) {
		return true
	}
	// 当前节点的所有路径都检查完毕，从递归栈中移除
	recStack[folderId] = false
	return false
}

// MoveFolder 移动文件夹到另一个文件夹
func (s *UserDocFolderService) MoveFolder(ctx context.Context, req *clientdocpb.MoveFolderOrDocReq, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.MoveFolder")
	defer span.Finish()
	// 验证请求参数
	if len(req.TargetFolderItems) == 0 {
		return nil
	}

	// 验证文件夹存在并属于当前用户
	if len(req.MovedIds) == 0 {
		return errors.Biz("param.error")
	}

	// 获取要移动的文件夹
	movedFolder, err := s.userDocFolderDAO.FindById(ctx, req.MovedIds[0])
	if err != nil {
		return errors.Wrap(err, "get moved folder failed")
	}

	if movedFolder == nil || movedFolder.UserId != userId {
		return errors.Biz("personal.doc.folder.sort.status.error")
	}

	// 获取用户的所有文件夹
	userDocFolders, err := s.userDocFolderDAO.GetUserDocFoldersByUserID(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "get user folders failed")
	}

	movedFolderIds := make([]string, 0, len(req.MovedIds))
	for _, id := range req.MovedIds {
		movedFolderIds = append(movedFolderIds, id)
	}

	// 检查目标文件夹中是否已存在相同名称的文件夹
	if req.TargetFolderId != req.SourceFolderId {
		existSameNameFolder := false
		targetFolderId := req.TargetFolderId

		for _, folder := range userDocFolders {
			if folder.ParentId == targetFolderId && folder.Name == movedFolder.Name && folder.Id != movedFolder.Id {
				existSameNameFolder = true
				break
			}
		}

		if existSameNameFolder {
			return errors.Biz("personal.doc.folder.exist.same")
		}

		// 更新文件夹的父目录
		for _, movedFolderId := range movedFolderIds {
			var newParentId string = req.TargetFolderId

			// 如果要移动的文件夹就是目标文件夹，则移动到根目录
			if movedFolderId == req.TargetFolderId {
				newParentId = "0"
			}

			// 更新父目录ID
			updates := map[string]interface{}{
				"parent_id": newParentId,
			}

			err := s.userDocFolderDAO.UpdateFolder(ctx, movedFolderId, updates)
			if err != nil {
				return errors.Wrap(err, "update folder parent id failed")
			}
		}
	}

	// 更新文件夹排序
	if len(req.TargetFolderItems) > 0 {
		for i, item := range req.TargetFolderItems {
			// 计算新的排序值
			newSort := int32(len(req.TargetFolderItems) - 1 - i)

			// 更新排序
			updates := map[string]interface{}{
				"sort": newSort,
			}

			err := s.userDocFolderDAO.UpdateFolder(ctx, item.Id, updates)
			if err != nil {
				return errors.Wrap(err, "update folder sort failed")
			}
		}
	}

	// 获取用户的所有文件夹，检查是否存在循环依赖
	folders, err := s.userDocFolderDAO.GetUserDocFoldersByUserID(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "get user folders failed")
	}

	// 检查是否存在循环依赖
	hasCircular := s.hasCircularDependency(folders)
	if hasCircular {
		return errors.Biz("doc.user_doc_folder.errors.circular_dependency")
	}
	return nil
}

// GetByIdAndUserId 获取用户文档文件夹，如果id为0，返回根文件夹
func (s *UserDocFolderService) GetByIdAndUserId(ctx context.Context, id string, userId string) (*model.UserDocFolder, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.GetDocFolderByIdAndUserId")
	defer span.Finish()

	if id == "0" {
		rootFolder := &model.UserDocFolder{
			UserId:   userId,
			Name:     "root",
			ParentId: "0",
			Sort:     0,
		}
		rootFolder.Id = "0"
		return rootFolder, nil
	}

	// 调用DAO层查询用户文件夹列表
	folder, err := s.userDocFolderDAO.GetByIdAndUserId(ctx, id, userId)
	if err != nil {
		return nil, errors.Wrap(err, "get user doc folder failed")
	}
	return folder, nil
}

func (s *UserDocFolderService) GetChildrenFoldersByFolderId(ctx context.Context, folderId string) ([]model.UserDocFolder, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.GetChildrenFoldersByFolderId")
	defer span.Finish()

	// 调用DAO层查询用户文件夹列表
	folders, err := s.userDocFolderDAO.GetUserDocFoldersByParentID(ctx, folderId)
	if err != nil {
		return nil, errors.Wrap(err, "get user doc folders by parent id failed")
	}
	return folders, nil
}

// GetAllByIdAndUserId 获取用户文档文件夹及其下所有子文件夹
func (s *UserDocFolderService) GetAllByIdAndUserId(ctx context.Context, id string, userId string) ([]model.UserDocFolder, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.GetAllByIdAndUserId")
	defer span.Finish()

	var folders []model.UserDocFolder

	// 调用DAO层查询用户文件夹列表
	folder, err := s.GetByIdAndUserId(ctx, id, userId)
	if err != nil {
		return nil, errors.Wrap(err, "get user doc folder failed")
	}

	if folder == nil {
		return folders, nil
	}

	folders = append(folders, *folder)

	//通过s.GetChildrenFoldersByFolderId递归获取子文件夹及其子文件夹
	childrenFolders, err := s.GetChildrenFoldersByFolderId(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get children folders by folder id failed")
	}

	for _, childFolder := range childrenFolders {
		//folders = append(folders, childFolder)
		//查询子文件夹的子文件夹
		subFolders, err := s.GetAllByIdAndUserId(ctx, childFolder.Id, userId)
		if err != nil {
			return nil, errors.Wrap(err, "get children folders by folder id failed")
		}
		if subFolders != nil {
			folders = append(folders, subFolders...)
		}

	}

	return folders, nil
}

func (s *UserDocFolderService) GetById(ctx context.Context, id string) (*model.UserDocFolder, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.GetById")
	defer span.Finish()

	// 调用DAO层查询用户文件夹列表
	folder, err := s.userDocFolderDAO.FindById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get user doc folder failed")
	}
	return folder, nil
}

func (s *UserDocFolderService) CopyDocToAnotherFolder(ctx context.Context, userId string, folderId string, docIds []string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.CopyDocOrFolderToAnotherFolder")
	defer span.Finish()

	userDocFolderRelations := make([]model.UserDocFolderRelation, 0)
	maxSort, err := s.userDocFolderRelationService.GetMaxSort(ctx, userId, folderId)
	if err != nil {
		return errors.Wrap(err, "get max sort failed")
	}
	for _, docId := range docIds {
		userDocFolderRelation, err := s.GenDocFolderRelation(ctx, userId, folderId, docId)
		if err != nil {
			return errors.Wrap(err, "get user doc folder failed")
		}
		if userDocFolderRelation != nil {
			maxSort++
			userDocFolderRelation.Sort = maxSort
			userDocFolderRelations = append(userDocFolderRelations, *userDocFolderRelation)
		}
	}
	if len(userDocFolderRelations) > 0 {
		err := s.userDocFolderRelationService.SaveBatch(ctx, userDocFolderRelations)
		if err != nil {
			return errors.Wrap(err, "save user doc folder failed")
		}
	}
	return nil
}

func (s *UserDocFolderService) GenDocFolderRelation(ctx context.Context, userId string, folderId string, docId string) (*model.UserDocFolderRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.GenDocFolderRelation")
	defer span.Finish()

	userDocFolderRelation := &model.UserDocFolderRelation{
		UserId:   userId,
		FolderId: folderId,
		DocId:    docId,
	}
	folderRelationId := idgen.GenerateUUID()
	userDocFolderRelation.Id = folderRelationId
	userDocFolderRelation.CreatorId = userId
	userDocFolderRelation.ModifierId = userId
	userDocFolderRelation.CreatedAt = time.Now()
	//
	sourceRelation, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByFolderIdAndDocId(ctx, userId, folderId, docId)
	if err != nil {
		return nil, errors.Wrap(err, "get user doc folder failed")
	}
	if sourceRelation == nil {
		return userDocFolderRelation, nil
	}
	return nil, nil
}

func (s *UserDocFolderService) CopyFolderToAnotherFolder(ctx context.Context, userId string, targetFolderId string, folderIds []string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.CopyFolderToAnotherFolder")
	defer span.Finish()

	folders, err := s.userDocFolderDAO.GetUserDocFoldersByUserID(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "get user doc folders failed")
	}
	// 创建文件夹ID到文件夹对象的映射
	folderMap := make(map[string]*model.UserDocFolder)
	for i := range folders {
		folderMap[folders[i].Id] = &folders[i]
	}
	// 获取目标文件夹下所有子文件夹的名称
	targetSubFolderNames := make(map[string]bool)
	for _, folder := range folders {
		if folder.ParentId == targetFolderId {
			targetSubFolderNames[folder.Name] = true
		}
	}
	// 检查是否存在同名文件夹
	for _, folderId := range folderIds {
		if folder, exists := folderMap[folderId]; exists {
			if _, nameExists := targetSubFolderNames[folder.Name]; nameExists {
				return errors.Biz("doc.user_doc_folder.errors.folder_name_exists")
			}
		}
	}
	relations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByUserID(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "get user doc folders failed")
	}
	relationMapByFolderId := make(map[string]*model.UserDocFolderRelation)
	for i := range relations {
		relationMapByFolderId[relations[i].FolderId] = &relations[i]
	}
	//需要保存的文件夹列表
	foldersToSave := []model.UserDocFolder{}
	//需要保存的文献列表
	documentsToSave := []model.UserDocFolderRelation{}
	//需要修改的文件夹列表
	foldersToUpdate := []model.UserDocFolder{}
	for _, folderId := range folderIds {
		folder, err := s.userDocFolderDAO.FindById(ctx, folderId)
		if err != nil || folder == nil {
			continue
		}
		//找到所有的子文件夹
		descendantFolderIds := []string{}
		s.getDescendantFolderIds(folderId, folders, &descendantFolderIds)
		descendantFolderIds = append(descendantFolderIds, folderId)
		currentFolderTempInfos := []*bean.CopyFolderTempInfo{}
		oldNewFolderIdMap := make(map[string]string)
		for _, descendantFolderId := range descendantFolderIds {
			docFolder := folderMap[descendantFolderId]
			if docFolder == nil {
				continue
			}
			userDocFolder := &model.UserDocFolder{
				UserId:   userId,
				Name:     docFolder.Name,
				ParentId: targetFolderId,
			}
			userDocFolder.CreatorId = userId
			userDocFolder.ModifierId = userId
			// 生成ID并保存
			folderId := idgen.GenerateUUID()
			userDocFolder.Id = folderId
			userDocFolder.CreatedAt = time.Now()
			foldersToSave = append(foldersToSave, *userDocFolder)
			tempInfo := &bean.CopyFolderTempInfo{
				NewFolderId:   userDocFolder.Id,
				OldFolderId:   descendantFolderId,
				OldParentId:   docFolder.ParentId,
				UserDocFolder: userDocFolder,
			}
			oldNewFolderIdMap[descendantFolderId] = userDocFolder.Id
			// 获取关联的文档关系
			var docFolderRelations []model.UserDocFolderRelation
			for _, relation := range relations {
				if relation.FolderId == descendantFolderId {
					docFolderRelations = append(docFolderRelations, relation)
				}
			}
			tempInfo.Relations = docFolderRelations
			currentFolderTempInfos = append(currentFolderTempInfos, tempInfo)
		}
		//更新文件夹的父ID
		for _, tempInfo := range currentFolderTempInfos {
			newParentId, exists := oldNewFolderIdMap[tempInfo.OldParentId]
			if !exists {
				newParentId = targetFolderId
			}
			tempInfo.NewParentId = newParentId
			tempInfo.UserDocFolder.ParentId = newParentId
			// 更新文件夹
			foldersToUpdate = append(foldersToUpdate, *tempInfo.UserDocFolder)
			// 处理文件夹文献关系
			if len(tempInfo.Relations) > 0 {
				maxSort, err := s.userDocFolderRelationService.GetMaxSort(ctx, userId, tempInfo.NewFolderId)
				if err != nil {
					return errors.Wrap(err, "get max sort failed")
				}
				for _, relation := range tempInfo.Relations {
					// 创建新的文档-文件夹关系
					newRelation, err := s.GenDocFolderRelation(ctx, relation.UserId, tempInfo.NewFolderId, relation.DocId)
					if err != nil {
						return errors.Wrap(err, "get user doc folder relation failed")
					}
					if newRelation != nil {
						maxSort++
						newRelation.Sort = maxSort
						documentsToSave = append(documentsToSave, *newRelation)
					}
				}
			}
		}
	}
	// 使用事务确保数据一致性
	err = s.userDocFolderDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务对象放入 context，确保内部操作使用同一个事务连接
		// 这对于 SQLite 尤其重要，因为 SQLite 使用数据库级锁
		txCtx := context.WithValue(ctx, baseDao.TransactionContextKey, tx)

		if len(foldersToUpdate) > 0 {
			for _, folder := range foldersToUpdate {
				s.userDocFolderDAO.ModifyExcludeNull(txCtx, &folder)
			}
		}
		if len(documentsToSave) > 0 {
			s.userDocFolderRelationService.SaveBatch(txCtx, documentsToSave)
		}
		if len(foldersToSave) > 0 {
			s.userDocFolderDAO.SaveBatch(txCtx, foldersToSave)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}

func (s *UserDocFolderService) getDescendantFolderIds(pid string, allFolders []model.UserDocFolder, descendantFolderIds *[]string) {
	for _, folder := range allFolders {
		if folder.ParentId == pid {
			*descendantFolderIds = append(*descendantFolderIds, folder.Id)
			s.getDescendantFolderIds(folder.Id, allFolders, descendantFolderIds)
		}
	}
}

func (s *UserDocFolderService) CopyDocOrFolderToAnotherFolder(ctx context.Context, userId string, req *docpb.CopyDocOrFolderToAnotherFolderReq) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocFolderService.CopyDocOrFolderToAnotherFolder")
	defer span.Finish()

	folder, err := s.userDocFolderDAO.FindById(ctx, req.TargetFolderId)
	if err != nil {
		return errors.Wrap(err, "get user doc folder failed")
	}
	if folder == nil {
		return errors.Biz("doc.user_doc_folder.errors.folder_not_found")
	}
	if folder.UserId != userId {
		return errors.Biz("doc.user_doc_folder.errors.not_belong_to_current_user")
	}
	//复制文献
	if len(req.DocIds) > 0 {
		s.CopyDocToAnotherFolder(ctx, userId, folder.Id, req.DocIds)
	}
	// 复制文件夹
	if len(req.FolderIds) > 0 {
		s.CopyFolderToAnotherFolder(ctx, userId, folder.Id, req.FolderIds)
	}
	return nil
}

func (s *UserDocFolderService) RemoveDocFromFolder(ctx context.Context, userId string, req *docpb.RemoveDocFromFolderReq) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.RemoveDocFromFolder")
	defer span.Finish()

	//参数转list
	removeDocIdList := make([]string, 0)
	for i := range req.RemovedDocItems {
		if req.RemovedDocItems[i].DocId == "0" {
			continue
		}
		removeDocIdList = append(removeDocIdList, req.RemovedDocItems[i].DocId)
	}
	//获取所有的对应关系
	relations, err := s.userDocFolderRelationService.GetRelationsByUserIdAndDocIds(ctx, userId, removeDocIdList)
	if err != nil {
		return errors.Wrap(err, "get user doc folder relation failed")
	}
	if len(relations) == 0 {
		return nil
	}
	if req.IsHierarchicallyRemove {
		//使用事务
		err = s.userDocFolderDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
			// 将事务对象放入 context，确保内部操作使用同一个事务连接
			txCtx := context.WithValue(ctx, baseDao.TransactionContextKey, tx)

			for _, item := range req.RemovedDocItems {
				if item.FolderId == "0" {
					continue
				}
				err := s.userDocFolderRelationService.DeleteRelationsByUserIdAndFolderIdAndDocId(txCtx, userId, item.FolderId, item.DocId)
				if err != nil {
					return errors.Wrap(err, "delete user doc folder relation failed")
				}
			}
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "delete user doc folder relation failed")
		}
	} else {
		folders, err := s.GetUserDocFolders(ctx, userId)
		if err != nil {
			return errors.Wrap(err, "get user doc folders failed")
		}
		//使用事务
		err = s.userDocFolderDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
			// 将事务对象放入 context，确保内部操作使用同一个事务连接
			txCtx := context.WithValue(ctx, baseDao.TransactionContextKey, tx)

			for _, item := range req.RemovedDocItems {
				if item.FolderId == "0" {
					continue
				}
				descendantFolderIds := []string{}
				s.getDescendantFolderIds(item.FolderId, folders, &descendantFolderIds)
				descendantFolderIds = append(descendantFolderIds, item.FolderId)
				err := s.userDocFolderRelationService.DeleteRelationsByUserIdAndFolderIdsAndDocId(txCtx, userId, descendantFolderIds, item.DocId)
				if err != nil {
					return errors.Wrap(err, "delete user doc folder relation failed")
				}
			}
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "delete user doc folder relation failed")
		}
	}
	return nil
}
