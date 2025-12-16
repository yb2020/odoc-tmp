package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/internal/biz"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/eventbus"
	"github.com/yb2020/odoc/pkg/logging"
	pb "github.com/yb2020/odoc/proto/gen/go/membership"
	orderPb "github.com/yb2020/odoc/proto/gen/go/order"
	"github.com/yb2020/odoc/services/membership/dto"
	"github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/membership/model"
	userEvent "github.com/yb2020/odoc/services/user/event"
)

// MembershipService 会员服务实现
type MembershipService struct {
	logger logging.Logger
	tracer opentracing.Tracer

	msConfigService       *ConfigService
	userMembershipService interfaces.IUserMembershipService
	orderService          interfaces.IOrderService
	creditService         interfaces.ICreditService
	creditPaymentService  interfaces.ICreditPaymentService
}

func NewMembershipService(logger logging.Logger, tracer opentracing.Tracer, msConfigService *ConfigService, userMembershipService interfaces.IUserMembershipService, orderService interfaces.IOrderService, creditService interfaces.ICreditService, creditPaymentService interfaces.ICreditPaymentService) *MembershipService {
	return &MembershipService{
		logger:                logger,
		tracer:                tracer,
		msConfigService:       msConfigService,
		userMembershipService: userMembershipService,
		orderService:          orderService,
		creditService:         creditService,
		creditPaymentService:  creditPaymentService,
	}
}

// GetBaseInfo 获取用户会员基础信息
func (s *MembershipService) GetBaseInfo(ctx context.Context, userId string) (*pb.UserMembershipBaseInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.GetBaseInfo")
	defer span.Finish()
	_, err := s.CheckAccountAndNew(ctx, userId)
	if err != nil {
		return nil, err
	}
	return s.userMembershipService.GetBaseInfo(ctx, userId)
}

// GetInfo 获取用户会员信息
func (s *MembershipService) GetInfo(ctx context.Context, userId string) (*pb.UserMembershipInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.GetInfo")
	defer span.Finish()
	_, err := s.CheckAccountAndNew(ctx, userId)
	if err != nil {
		return nil, err
	}
	return s.userMembershipService.GetInfo(ctx, userId)
}

// CheckAccountAndNew 检查用户会员账号并创建
func (s *MembershipService) CheckAccountAndNew(ctx context.Context, userId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.CheckAccountAndNew")
	defer span.Finish()

	// 1.检查用户会员账号
	userMembership, err := s.userMembershipService.GetByUserId(ctx, userId)
	if err != nil {
		return "0", err
	}
	if userMembership == nil {
		// 用户会员账号不存在,创建用户会员账号
		return s.NewMembershipAccount(ctx, userId)
	}
	return userMembership.Id, nil
}

// NewMembershipAccount 创建用户会员账号并自动订阅Free会员
func (s *MembershipService) NewMembershipAccount(ctx context.Context, userId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.NewMembershipAccount")
	defer span.Finish()

	// 1.创建用户会员账号
	msId, err := s.userMembershipService.NewAccount(ctx, userId)
	if err != nil {
		return "0", errors.BizWithStatus(biz.Membership_Status_UserAccountNewAccountFailed, "new user account failed")
	}

	// 2.自动订阅Free会员订单
	err = s.AutoSubscribeMemberFree(ctx, userId)
	if err != nil {
		return "0", errors.Biz("subscribe free member failed")
	}

	return msId, nil
}

// DeleteMembershipAccount 删除用户会员账号
func (s *MembershipService) DeleteMembershipAccount(ctx context.Context, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.DeleteMembershipAccount")
	defer span.Finish()

	return s.userMembershipService.DeleteAccount(ctx, userId)
}

// AutoSubscribeMemberFree 自动订阅Free会员，无需支付自动完成订单状态
func (s *MembershipService) AutoSubscribeMemberFree(ctx context.Context, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.SubscribeMemberFree")
	defer span.Finish()

	// 1.订阅Free会员订单
	orderId, err := s.orderService.Subscribe(ctx, userId, orderPb.OrderType_ORDER_TYPE_SUB_FREE, 1)
	if err != nil {
		return err
	}

	// 2.自动完成订单状态
	err = s.orderService.DoOrderPaySuccessHandler(ctx, orderId, "0", "")
	if err != nil {
		return err
	}

	return nil
}

// abstractCreditFuns 抽象调用功能接口
// checkPermissionAndPayConfigFun: 不能为nil 校验用户权限和支付配置，返回值：hasPermission: 是否有权限，isNeedCreditPay: 是否需要积分支付，newCreditPayOrder: 支付积分，error: 错误
// funType: 调用操作类型
// invokeFun: 调用操作
// autoConfirm: 是否自动确认积分服务
func (s *MembershipService) abstractCreditFuns(ctx context.Context, checkPermissionAndPayConfigFun func(xctx context.Context) (bool, bool, *dto.NewCreditPayOrder, error), funType pb.CreditServiceType, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.abstractCreditFuns")
	defer span.Finish()

	userId, _ := userContext.GetUserID(ctx)
	s.logger.Info("CallCreditFuns", "userId", userId, "funType", funType)

	userMembership, err := s.userMembershipService.GetByUserId(ctx, userId)
	if err != nil {
		return err
	}

	if userMembership == nil {
		// 用户会员账号不存在,创建用户会员账号
		_, err = s.NewMembershipAccount(ctx, userId)
		if err != nil {
			return err
		}
	}

	// 校验用户权限和获取功能扣减积分等信息
	if checkPermissionAndPayConfigFun == nil {
		return errors.BizWithStatus(biz.Membership_Status_CreditServiceCheckPermissionAndPayConfigFunNoNull, "checkPermissionAndPayConfigFun can not be nil")
	}
	hasPermission, isNeedCreditPay, payCredit, err := checkPermissionAndPayConfigFun(ctx)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.BizWithStatus(biz.Membership_Status_CreditServicePermissionDenied, "function permission denied")
	}

	if isNeedCreditPay {
		//1.检测积分/附加积分是否足够支付
		isEnough, err := s.creditService.CheckCreditEnough(ctx, userId, userMembership.Id, payCredit.Credit)
		if err != nil {
			return err
		}
		if !isEnough {

			//检测是权限使用附加积分和检测附加积分是否足够支付
			useAddOnCreditPerm, _ := s.userMembershipService.HasPermUseAddOnCredit(ctx, userId)
			if !useAddOnCreditPerm {
				return errors.BizWithStatus(biz.Membership_Status_CreditNotEnough, "credit not enough")
			}
			isAddOnEnough, err := s.creditService.CheckAddOnCreditEnough(ctx, userId, userMembership.Id, payCredit.Credit)
			if err != nil {
				return err
			}
			if !isAddOnEnough {
				return errors.BizWithStatus(biz.Membership_Status_CreditAddOnNotEnough, "add on credit not enough")
			}
			//使用附加积分支付
			payCredit.CreditType = pb.CreditType(pb.CreditType_CREDIT_TYPE_ADD_ON_CREDIT)
		}

		//2. 创建Credit支付订单
		recordId, err := s.creditPaymentService.NewPaymentOrder(ctx, userMembership.Id, userId, payCredit)
		if err != nil {
			return err
		}
		s.logger.Info("CallCreditFuns", "recordId", recordId)
		if recordId == "0" {
			return errors.BizWithStatus(biz.Membership_Status_CreditPaymentRecordNotFound, "credit payment record not found")
		}

		//3.会员账户支付支付订单
		err = s.creditService.Pay(ctx, userId, userMembership.Id, recordId)
		if err != nil {
			return err
		}

		//4.执行回调
		if invokeFun != nil {
			err := invokeFun(ctx, recordId)
			if err != nil {
				// 执行回调失败，回退积分服务
				s.creditService.Retrieve(ctx, recordId)
				return err
			}
		}

		//5.自动确认支付订单
		if autoConfirm {
			s.logger.Info("msg", "autoconfirm payment order", "recordId", recordId)
			err = s.creditService.Confirm(ctx, recordId)
			if err != nil {
				return err
			}
		} else {
			//ignore
			s.logger.Info("msg", "confirm payment order by user self", "recordId", recordId)
		}
	} else {
		//2.执行回调
		if invokeFun != nil {
			err := invokeFun(ctx, "0")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CreditFunDocsUpload 调用积分接口上传接口功能
// fileSize: 文件大小(byte)
// filePageCount: 文件页数
// storageCapacity: 存储容量(byte)
// autoConfirm: 是否自动确认
func (s *MembershipService) CreditFunDocsUpload(ctx context.Context, fileSize int64, filePageCount int32, storageCapacity int64, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.CreditFunDocsUpload")
	defer span.Finish()

	err := s.abstractCreditFuns(ctx, func(xctx context.Context) (bool, bool, *dto.NewCreditPayOrder, error) {
		// 判断用户文档上传权限
		memberConfig, err := s.userMembershipService.GetUserConfig(xctx)
		if err != nil {
			return false, false, nil, err
		}
		//判断上传到单个文件
		//单位换算 1MB = 1048576 byte （1024*1024）
		maxSize := memberConfig.Docs.DocUploadMaxSize * 1048576
		if maxSize < fileSize {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Docs_Upload_OverFileSize, "upload file size too large")
		} else if memberConfig.Docs.DocUploadMaxPageCount < filePageCount {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Docs_Upload_OverPageSize, "upload file page count too large")
		}

		//判断用户存储总量
		maxStorageCapacity := memberConfig.Docs.MaxStorageCapacity * 1048576
		if maxStorageCapacity < (storageCapacity + fileSize) {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Docs_Upload_OverMaxStorage, "storage capacity not enough")
		}
		return true, false, nil, nil
	}, pb.CreditServiceType_CREDIT_SERVICE_TYPE_DOCS_UPLOAD, invokeFun, autoConfirm)
	return err
}

// CreditFunAi 调用积分接口AI辅读功能
// funType: AI辅读类型
// modelKey: AI模型key
// invokeFun: 调用AI辅读接口
// autoConfirm: 是否自动确认
func (s *MembershipService) CreditFunAi(ctx context.Context, funType pb.CreditServiceType, modelKey string, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.CreditFunAiCopilot")
	defer span.Finish()

	err := s.abstractCreditFuns(ctx, func(xctx context.Context) (bool, bool, *dto.NewCreditPayOrder, error) {
		// 判断用户AI辅读权限
		memberConfig, err := s.userMembershipService.GetUserConfig(xctx)
		if err != nil {
			return false, false, nil, err
		}
		if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_AI_COPILOT {
			if !memberConfig.AI.Copilot.IsEnable {
				return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Ai_Copilot_NotEnabled, "ai copilot not enabled")
			}
			aiModelId := -1
			models := memberConfig.AI.Copilot.Models
			for i, model := range models {
				if model.Key == modelKey {
					aiModelId = i
					break
				}
			}
			if aiModelId == -1 {
				return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Ai_Copilot_ModelNotFound, "ai model not found")
			}

			aiModel := memberConfig.AI.Copilot.Models[aiModelId]

			if !aiModel.IsEnable {
				return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Ai_Copilot_ModelNotEnabled, "ai model not enabled")
			}

			isNeedCreditPay := true
			creditCost := aiModel.CreditCost
			if aiModel.IsFree {
				isNeedCreditPay = false
			}

			// 构建支付信息
			payCredit := &dto.NewCreditPayOrder{
				CreditType:  pb.CreditType_CREDIT_TYPE_CREDIT,
				ServiceType: funType,
				Credit:      creditCost,
				//AddOnCredit: 0,
			}
			return true, isNeedCreditPay, payCredit, nil
		}
		return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditServiceTypeUnknown, "credit service type unknown")

	}, funType, invokeFun, autoConfirm)
	return err
}

// CreditFunTranslate 调用翻译接口功能
// funType: 翻译类型
// filePageCount: 文件页数
// invokeFun: 调用翻译接口
// autoConfirm: 是否自动确认
func (s *MembershipService) CreditFunTranslate(ctx context.Context, funType pb.CreditServiceType, filePageCount int32, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.CreditFunTranslate")
	defer span.Finish()

	err := s.abstractCreditFuns(ctx, func(xctx context.Context) (bool, bool, *dto.NewCreditPayOrder, error) {
		// 判断用户翻译权限
		memberConfig, err := s.userMembershipService.GetUserConfig(xctx)
		if err != nil {
			return false, false, nil, err
		}

		creditCost := int64(0)
		// creditPayType := dto.CreditPayType_Service_Cost
		if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_OCR {
			creditCost = memberConfig.Translate.OcrCreditCost
			// creditPayType = dto.CreditPayType_Translate_Ocr
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_WORD {
			creditCost = memberConfig.Translate.WordTranslateCreditCost
			// creditPayType = dto.CreditPayType_Translate_Word
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_FULLTEXT {
			creditCost = memberConfig.Translate.FullTextTranslateCreditCost
			// creditPayType = dto.CreditPayType_Translate_FullText
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_AI {
			creditCost = memberConfig.Translate.AiTranslationCreditCost
			// creditPayType = dto.CreditPayType_Translate_Ai
		}

		if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_OCR && !memberConfig.Translate.IsOcr {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Translate_OcrNotEnabled, "translate not enabled")
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_WORD && !memberConfig.Translate.IsWordTranslate {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Translate_WordNotEnabled, "translate not enabled")
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_FULLTEXT && !memberConfig.Translate.IsFullTextTranslate {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Translate_FullTextNotEnabled, "translate not enabled")
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_AI && !memberConfig.Translate.IsAiTranslation {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Translate_AiNotEnabled, "translate not enabled")
		}

		// 判断全文翻译页数是否超过限制
		if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_FULLTEXT && memberConfig.Translate.FullTextTranslateMaxPageCount < int64(filePageCount) {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Translate_FullText_OverPageSize, "translate full text page count too large")
		}

		// 构建支付信息
		payCredit := &dto.NewCreditPayOrder{
			CreditType:  pb.CreditType_CREDIT_TYPE_CREDIT,
			ServiceType: funType,
			Credit:      creditCost,
			// AddOnCredit: 0,
		}

		return true, true, payCredit, nil
	}, funType, invokeFun, autoConfirm)
	return err
}

// CreditFunNote 调用笔记接口功能
// funType: 笔记类型
// invokeFun: 调用笔记接口
// autoConfirm: 是否自动确认
func (s *MembershipService) CreditFunNote(ctx context.Context, funType pb.CreditServiceType, invokeFun func(xctx context.Context, sessionId string) error, autoConfirm bool) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.CreditFunNote")
	defer span.Finish()

	err := s.abstractCreditFuns(ctx, func(xctx context.Context) (bool, bool, *dto.NewCreditPayOrder, error) {
		// 判断用户笔记权限
		memberConfig, err := s.userMembershipService.GetUserConfig(xctx)
		if err != nil {
			return false, false, nil, err
		}
		if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_SUMMARY && !memberConfig.Note.IsNoteSummary {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Note_SummaryNotEnabled, "note summary not enabled")
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_WORD && !memberConfig.Note.IsNoteWord {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Note_WordNotEnabled, "note word not enabled")
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_EXTRACT && !memberConfig.Note.IsNoteExtract {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Note_ExtractNotEnabled, "note extract not enabled")
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_MANAGE && !memberConfig.Note.IsNoteManage {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Note_ManageNotEnabled, "note manage not enabled")
		} else if funType == pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_PDF_DOWNLOAD && !memberConfig.Note.IsNotePdfDownload {
			return false, false, nil, errors.BizWithStatus(biz.Membership_Status_CreditService_Note_PdfDownloadNotEnabled, "note pdf download not enabled")
		}
		return true, false, nil, nil
	}, funType, invokeFun, autoConfirm)
	return err
}

// ConfirmCreditFun 二段提交成功确认接口
// sessionId: 会话ID
func (s *MembershipService) ConfirmCreditFun(ctx context.Context, sessionId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.ConfirmCreditFun")
	defer span.Finish()
	// s.logger.Info("msg", "ConfirmCreditFun", "sessionId", sessionId)

	return s.creditService.Confirm(ctx, sessionId)
}

// RetrieveCreditFun 二段提交失败确认接口
// sessionId: 会话ID
func (s *MembershipService) RetrieveCreditFun(ctx context.Context, sessionId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.RetrieveCreditFun")
	defer span.Finish()
	// s.logger.Info("msg", "RetrieveCreditFun, do something for failed", "sessionId", sessionId)

	return s.creditService.Retrieve(ctx, sessionId)
}

// RefundCreditFun 发起退款接口，用于积分接口功能执行失败后，发起退款，暂不实现功能
// sessionId: 会话ID
// func (s *MembershipService) RefundCreditFun(ctx context.Context, sessionId int64) error {
// 	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.RefundCreditFun")
// 	defer span.Finish()

// 	s.logger.Info("msg", "RetrieveCreditFun, now no implement", "sessionId", sessionId)

// 	return nil
// }

// 处理用户通知事件（由user模块发布消息--》会员模块订阅处理）
func (s *MembershipService) HandleUserNotifyEventHandler(ctx context.Context, eventType eventbus.EventType, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.HandleUserNotifyEventHandler")
	defer span.Finish()

	switch eventType {
	case userEvent.UserRegisterEvent:
		_, err := s.NewMembershipAccount(ctx, userId)
		return err
	case userEvent.UserDeletedEvent:
		err := s.DeleteMembershipAccount(ctx, userId)
		return err
	}

	return nil
}

// GetAccountExpiredList 获取过期的用户会员列表
func (s *MembershipService) GetAccountExpiredList(ctx context.Context, size int) ([]model.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.GetAccountExpiredList")
	defer span.Finish()

	return s.userMembershipService.GetExpiredList(ctx, size)
}

// GetFreeAccountExpiredList 获取过期的Free用户会员列表
func (s *MembershipService) GetFreeAccountExpiredList(ctx context.Context, size int) ([]model.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.GetFreeAccountExpiredList")
	defer span.Finish()

	return s.userMembershipService.GetFreeExpiredList(ctx, size)
}

// GetProAccountExpiredList 获取过期的Pro用户会员列表
func (s *MembershipService) GetProAccountExpiredList(ctx context.Context, size int) ([]model.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.GetProAccountExpiredList")
	defer span.Finish()

	return s.userMembershipService.GetProExpiredList(ctx, size)
}

// HandleAccountExpired 处理过期的用户会员, id为会员id
func (s *MembershipService) HandleAccountExpired(ctx context.Context, id string, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipService.HandleAccountExpired")
	defer span.Finish()

	isExpired, userType, err := s.userMembershipService.CheckExpired(ctx, userId)
	if err != nil {
		return err
	}
	s.logger.Info("msg", "account info", "id", id, "isExpired", isExpired, "userType", userType)

	if !isExpired {
		return errors.Biz("account not expired")
	}

	if userType == int32(pb.MembershipType_MEMBERSHIP_TYPE_FREE) {
		// 1.自动订阅Free会员订单
		return s.AutoSubscribeMemberFree(ctx, userId)
	} else if userType == int32(pb.MembershipType_MEMBERSHIP_TYPE_PRO) {
		// 1.过期会员降级为Free会员
		err := s.userMembershipService.ExpiredAccountToFree(ctx, id)
		if err != nil {
			return err
		}
		// 2.自动订阅Free会员订单
		return s.AutoSubscribeMemberFree(ctx, userId)
	}

	return nil
}
