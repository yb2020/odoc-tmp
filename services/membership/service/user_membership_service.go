package service

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/membership"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/internal/biz"
	userContext "github.com/yb2020/odoc/pkg/context"
	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/membership/dao"
	"github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/membership/model"
	userService "github.com/yb2020/odoc/services/user/service"
)

// UserMembershipService 用户会员服务实现
type UserMembershipService struct {
	logger             logging.Logger
	tracer             opentracing.Tracer
	transactionManager *baseDao.TransactionManager
	userService        *userService.UserService
	userMembershipDAO  *dao.UserMembershipDAO

	msConfigService *ConfigService
	creditService   interfaces.ICreditService
}

func NewUserMembershipService(logger logging.Logger, tracer opentracing.Tracer, userService *userService.UserService, userMembershipDAO *dao.UserMembershipDAO, msConfigService *ConfigService, creditService interfaces.ICreditService, transactionManager *baseDao.TransactionManager) *UserMembershipService {
	return &UserMembershipService{
		logger:             logger,
		tracer:             tracer,
		transactionManager: transactionManager,
		userService:        userService,
		userMembershipDAO:  userMembershipDAO,

		msConfigService: msConfigService,
		creditService:   creditService,
	}
}

// GetUserMembership 获取用户会员
func (s *UserMembershipService) GetById(ctx context.Context, id string) (*model.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetById")
	defer span.Finish()

	userMembership, err := s.userMembershipDAO.FindExistById(ctx, id)
	if err != nil {
		return nil, err
	}

	return userMembership, nil
}

// GetUserMembership 获取用户会员
func (s *UserMembershipService) GetByUserId(ctx context.Context, userId string) (*model.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetByUserId")
	defer span.Finish()

	userMembership, err := s.userMembershipDAO.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return userMembership, nil
}

// GetBaseInfo 获取用户会员基础信息
func (s *UserMembershipService) GetBaseInfo(ctx context.Context, userId string) (*pb.UserMembershipBaseInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetBaseInfo")
	defer span.Finish()

	userMembership, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	if userMembership == nil {
		return nil, errors.BizWithStatus(biz.Membership_Status_UserAccountNotFound, "user account not found")
	}
	creditAccount, err := s.creditService.GetByMembershipId(ctx, userMembership.Id)
	if err != nil {
		return nil, err
	}
	if creditAccount == nil {
		return nil, errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "user credit account not found")
	}

	msConfig := s.msConfigService.GetMemberConfigByType(ctx, pb.MembershipType(userMembership.Type))
	if msConfig == nil {
		return nil, errors.Biz("member config not found")
	}

	typeName := msConfig.Name

	expiredDay := int32(time.Until(userMembership.EndAt).Hours() / 24)
	if expiredDay < 0 {
		expiredDay = 0
	}
	return &pb.UserMembershipBaseInfo{
		UserId:      userMembership.UserId,
		Type:        uint32(userMembership.Type),
		TypeName:    typeName,
		StartAt:     uint64(userMembership.StartAt.UnixMilli()),
		EndAt:       uint64(userMembership.EndAt.UnixMilli()),
		IsExpired:   userMembership.EndAt.Before(time.Now()),
		ExpiredDay:  uint32(expiredDay),
		Credit:      uint64(creditAccount.Credit),
		AddOnCredit: uint64(creditAccount.AddOnCredit),
	}, nil
}

// GetInfo 获取用户会员信息
func (s *UserMembershipService) GetInfo(ctx context.Context, userId string) (*pb.UserMembershipInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetInfo")
	defer span.Finish()

	userBaseInfo, err := s.GetBaseInfo(ctx, userId)
	if err != nil {
		return nil, err
	}

	userPermisson, err := s.GetUserPermission(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &pb.UserMembershipInfo{
		BaseInfo:   userBaseInfo,
		Permission: userPermisson,
	}, nil
}

func (s *UserMembershipService) GetUserPermission(ctx context.Context, userId string) (*pb.UserPermission, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetUserPermission")
	defer span.Finish()
	userMembership, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	if userMembership == nil {
		return nil, errors.BizWithStatus(biz.Membership_Status_UserAccountNotFound, "user membership not found")
	}
	msConfig := s.msConfigService.GetMemberConfigByType(ctx, pb.MembershipType(userMembership.Type))
	if msConfig == nil {
		return nil, errors.Biz("member config not found")
	}

	basePermission := &pb.UserBasePermission{
		IsEnableAddOnCredit:    msConfig.Base.IsEnableAddOnCredit,
		IsEnableSubAddOnCredit: msConfig.Base.IsEnableSubAddOnCredit,
	}

	docsPermission := &pb.UserDocsPermission{
		MaxStorageCapacity:    uint64(msConfig.Docs.MaxStorageCapacity),
		DocUploadMaxSize:      uint64(msConfig.Docs.DocUploadMaxSize),
		DocUploadMaxPageCount: uint32(msConfig.Docs.DocUploadMaxPageCount),
		UseStorageCapacity:    0, // TODO: 计算用户已使用存储容量
	}
	notePermission := &pb.UserNotePermission{
		IsNoteSummary:     msConfig.Note.IsNoteSummary,
		IsNoteWord:        msConfig.Note.IsNoteWord,
		IsNoteExtract:     msConfig.Note.IsNoteExtract,
		IsNoteManage:      msConfig.Note.IsNoteManage,
		IsNotePdfDownload: msConfig.Note.IsNotePdfDownload,
	}

	translatePermission := &pb.UserTranslatePermission{
		IsOcr:                         msConfig.Translate.IsOcr,
		OcrCreditCost:                 uint64(msConfig.Translate.OcrCreditCost),
		IsWordTranslate:               msConfig.Translate.IsWordTranslate,
		WordTranslateCreditCost:       uint64(msConfig.Translate.WordTranslateCreditCost),
		IsFullTextTranslate:           msConfig.Translate.IsFullTextTranslate,
		FullTextTranslateCreditCost:   uint64(msConfig.Translate.FullTextTranslateCreditCost),
		FullTextTranslateMaxPageCount: uint64(msConfig.Translate.FullTextTranslateMaxPageCount),
		IsAiTranslation:               msConfig.Translate.IsAiTranslation,
		AiTranslationCreditCost:       uint64(msConfig.Translate.AiTranslationCreditCost),
	}

	var pbAICopilotModels []*pb.CopilotModel
	for _, model := range msConfig.AI.Copilot.Models {
		pbAICopilotModels = append(pbAICopilotModels, &pb.CopilotModel{
			Key:        model.Key,
			Name:       model.Name,
			IsEnable:   model.IsEnable,
			IsFree:     model.IsFree,
			CreditCost: uint64(model.CreditCost),
		})
	}
	aiPermission := &pb.UserAIPermission{
		Copilot: &pb.CopilotPermission{
			IsEnable: msConfig.AI.Copilot.IsEnable,
			Models:   pbAICopilotModels,
		},
	}

	userPermisson := &pb.UserPermission{
		Base:      basePermission,
		Docs:      docsPermission,
		Note:      notePermission,
		Translate: translatePermission,
		Ai:        aiPermission,
	}

	return userPermisson, nil
}

// NewAccount 创建用户会员账号
func (s *UserMembershipService) NewAccount(ctx context.Context, userId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.NewAccount")
	defer span.Finish()

	// 检查用户状态
	userActive, err := s.userService.CheckUserStatusActive(ctx, userId)
	if err != nil {
		return "0", err
	}
	if !userActive {
		return "0", err
	}

	// 检查用户会员是否存在
	userMembershipExit, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return "0", err
	}
	if userMembershipExit != nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserAccountAlreadyExists, "user membership already exists")
	}

	userMembership := &model.UserMembership{
		UserId:  userId,
		Type:    int32(pb.MembershipType_MEMBERSHIP_TYPE_FREE),
		StartAt: time.Now().Add(-10 * time.Minute), // 设置开始时间为10分钟前
		EndAt:   time.Now().Add(-5 * time.Minute),  // 设置过期时间为5分钟前
	}
	id := idgen.GenerateUUID()
	userMembership.Id = id

	// 创建用户会员账号事务
	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		// 1.创建用户会员
		err = s.userMembershipDAO.Save(txCtx, userMembership)
		if err != nil {
			return err
		}
		// 2.创建积分账户
		_, err = s.creditService.NewCreditAccount(txCtx, userMembership.Id, userMembership.UserId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "0", err
	}

	return id, nil
}

// DeleteAccount 删除用户会员账号
func (s *UserMembershipService) DeleteAccount(ctx context.Context, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.DeleteAccount")
	defer span.Finish()

	// 检查用户会员是否存在
	userMembership, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return err
	}
	if userMembership == nil {
		return errors.BizWithStatus(biz.Membership_Status_UserAccountNotFound, "user membership not found")
	}
	creditAccount, err := s.creditService.GetByMembershipId(ctx, userMembership.Id)
	if err != nil {
		return err
	}
	if creditAccount == nil {
		return errors.BizWithStatus(biz.Membership_Status_UserCreditAccountNotFound, "user credit account not found")
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		// 1.删除积分账户
		err = s.creditService.DeleteCreditAccount(txCtx, creditAccount.Id)
		if err != nil {
			return err
		}
		// 2.删除用户会员
		return s.userMembershipDAO.DeleteById(txCtx, userMembership.Id)
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateAccountType 更新用户会员类型
// 参数说明: userId: 用户ID, memberType: 会员类型, stripeSubscriptionId: Stripe订阅ID, startDate: 会员开始时间, endDate: 会员结束时间
func (s *UserMembershipService) UpdateAccountType(ctx context.Context, userId string, memberType int32, stripeSubscriptionId string, startDate time.Time, endDate time.Time) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.UpdateAccountType")
	defer span.Finish()

	userMembership, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return err
	}
	if userMembership == nil {
		return errors.Biz("user membership not found")
	}

	userMembership.Type = memberType
	userMembership.StripeSubscriptionId = stripeSubscriptionId
	userMembership.StartAt = startDate
	userMembership.EndAt = endDate
	return s.userMembershipDAO.Modify(ctx, userMembership)
}

// GetUserConfig 获取用户会员配置信息
func (s *UserMembershipService) GetUserConfig(ctx context.Context) (config.MembershipTypeConfig, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetUserConfig")
	defer span.Finish()

	userId, _ := userContext.GetUserID(ctx)
	userMembership, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return config.MembershipTypeConfig{}, err
	}
	if userMembership == nil {
		return config.MembershipTypeConfig{}, errors.Biz("user membership not found")
	}

	return *s.msConfigService.GetMemberConfigByType(ctx, pb.MembershipType(userMembership.Type)), nil
}

// GetUserConfigByUserId 获取用户会员配置信息
func (s *UserMembershipService) GetUserConfigByUserId(ctx context.Context, userId string) (config.MembershipTypeConfig, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetUserConfigByUserId")
	defer span.Finish()

	userMembership, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return config.MembershipTypeConfig{}, err
	}
	if userMembership == nil {
		return config.MembershipTypeConfig{}, errors.Biz("user membership not found")
	}

	return *s.msConfigService.GetMemberConfigByType(ctx, pb.MembershipType(userMembership.Type)), nil
}

// CheckExpired 检查用户会员是否过期 返回是否过期，会员类型，错误
func (s *UserMembershipService) CheckExpired(ctx context.Context, userId string) (bool, int32, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.CheckExpired")
	defer span.Finish()

	userMembership, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return false, 0, err
	}
	if userMembership == nil {
		return false, 0, errors.Biz("user membership not found")
	}

	return userMembership.EndAt.Before(time.Now()), userMembership.Type, nil
}

// HasPermUseAddOnCredit 检查用户会员是否可以使用附加积分
func (s *UserMembershipService) HasPermUseAddOnCredit(ctx context.Context, userId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.HasPermUseAddOnCredit")
	defer span.Finish()

	userMembership, err := s.GetByUserId(ctx, userId)
	if err != nil {
		return false, err
	}
	if userMembership == nil {
		return false, errors.Biz("user membership not found")
	}

	config, err := s.GetUserConfigByUserId(ctx, userId)
	if err != nil {
		return false, err
	}

	return config.Base.IsEnableAddOnCredit, nil
}

// GetExpiredList 获取过期的会员列表
func (s *UserMembershipService) GetExpiredList(ctx context.Context, size int) ([]model.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetExpiredList")
	defer span.Finish()

	return s.userMembershipDAO.GetExpiredList(ctx, size)
}

// GetFreeExpiredList 获取过期的Free会员列表
func (s *UserMembershipService) GetFreeExpiredList(ctx context.Context, size int) ([]model.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetFreeExpiredList")
	defer span.Finish()

	return s.userMembershipDAO.GetExpiredListByMemberType(ctx, int(pb.MembershipType_MEMBERSHIP_TYPE_FREE), size)
}

// GetProExpiredList 获取过期的Pro会员列表
func (s *UserMembershipService) GetProExpiredList(ctx context.Context, size int) ([]model.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.GetProExpiredList")
	defer span.Finish()

	return s.userMembershipDAO.GetExpiredListByMemberType(ctx, int(pb.MembershipType_MEMBERSHIP_TYPE_PRO), size)
}

// ExpiredAccountToFree 过期会员降级为Free会员
func (s *UserMembershipService) ExpiredAccountToFree(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserMembershipService.ExpiredAccountToFree")
	defer span.Finish()

	userMembership, err := s.GetById(ctx, id)
	if err != nil {
		return err
	}
	if userMembership == nil {
		return errors.Biz("user membership not found")
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		// 降级为Free会员
		userMembership.Type = int32(pb.MembershipType_MEMBERSHIP_TYPE_FREE)
		err1 := s.userMembershipDAO.Modify(txCtx, userMembership)
		if err1 != nil {
			return err1
		}
		// 会员积分和附加积分过期清零
		err2 := s.creditService.CreditAccountExpiredAllCredit(txCtx, userMembership.Id)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
