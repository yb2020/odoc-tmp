package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/membership"
	orderPb "github.com/yb2020/odoc-proto/gen/go/order"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/membership/dto"
)

// ConfigService 会员配置服务实现
type ConfigService struct {
	logger       logging.Logger
	tracer       opentracing.Tracer
	config       *config.Config
	memberConfig config.MembershipConfig
}

func NewConfigService(logger logging.Logger, tracer opentracing.Tracer, config *config.Config) *ConfigService {
	return &ConfigService{
		logger:       logger,
		tracer:       tracer,
		config:       config,
		memberConfig: config.Membership,
	}
}

// GetConfig 获取会员配置
func (s *ConfigService) GetConfig(ctx context.Context) config.MembershipConfig {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipConfigService.GetConfig")
	defer span.Finish()

	return s.memberConfig
}

// GetMemberFreeConfig 获取Free版本会员配置
func (s *ConfigService) GetMemberFreeConfig(ctx context.Context) config.MembershipTypeConfig {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipConfigService.GetMemberFreeConfig")
	defer span.Finish()

	return s.memberConfig.Free
}

// GetMemberProConfig 获取Pro版本会员配置
func (s *ConfigService) GetMemberProConfig(ctx context.Context) config.MembershipTypeConfig {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipConfigService.GetMemberProConfig")
	defer span.Finish()

	return s.memberConfig.Professional
}

// GetSubBaseInfo 根据订阅类型subType获取订阅基础信息
func (s *ConfigService) GetSubBaseInfo(ctx context.Context, orderType int32) *dto.MembershipSubBaseInfo {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipConfigService.GetSubBaseInfo")
	defer span.Finish()

	switch orderPb.OrderType(orderType) {
	case orderPb.OrderType_ORDER_TYPE_SUB_FREE:

		freeConfig := s.GetMemberFreeConfig(ctx)
		return &dto.MembershipSubBaseInfo{
			Type:          orderPb.OrderType(freeConfig.Base.SubInfo.Type),
			Name:          freeConfig.Name,
			IsFree:        freeConfig.IsFree,
			Credit:        freeConfig.Base.SubInfo.Credit,
			AddOnCredit:   freeConfig.Base.SubInfo.AddOnCredit,
			Duration:      freeConfig.Base.SubInfo.Duration,
			Price:         freeConfig.Base.SubInfo.Price,
			Currency:      freeConfig.Base.SubInfo.Currency,
			StripePayMode: freeConfig.Base.SubInfo.StripePayMode,
			StripePriceId: freeConfig.Base.SubInfo.StripePriceId,
		}

	case orderPb.OrderType_ORDER_TYPE_SUB_PRO:

		proConfig := s.GetMemberProConfig(ctx)
		return &dto.MembershipSubBaseInfo{
			Type:          orderPb.OrderType(proConfig.Base.SubInfo.Type),
			Name:          proConfig.Name,
			IsFree:        proConfig.IsFree,
			Credit:        proConfig.Base.SubInfo.Credit,
			AddOnCredit:   proConfig.Base.SubInfo.AddOnCredit,
			Duration:      proConfig.Base.SubInfo.Duration,
			Price:         proConfig.Base.SubInfo.Price,
			Currency:      proConfig.Base.SubInfo.Currency,
			StripePayMode: proConfig.Base.SubInfo.StripePayMode,
			StripePriceId: proConfig.Base.SubInfo.StripePriceId,
		}

	case orderPb.OrderType_ORDER_TYPE_SUB_PRO_ADD_ON_CREDIT:

		proConfig := s.GetMemberProConfig(ctx)
		return &dto.MembershipSubBaseInfo{
			Type:          orderPb.OrderType(proConfig.Base.SubAddOnCreditInfo.Type),
			Name:          proConfig.Base.SubAddOnCreditInfo.Name,
			IsFree:        proConfig.IsFree,
			Credit:        proConfig.Base.SubAddOnCreditInfo.Credit,
			AddOnCredit:   proConfig.Base.SubAddOnCreditInfo.AddOnCredit,
			Price:         proConfig.Base.SubAddOnCreditInfo.Price,
			Currency:      proConfig.Base.SubAddOnCreditInfo.Currency,
			Duration:      proConfig.Base.SubAddOnCreditInfo.Duration,
			StripePayMode: proConfig.Base.SubAddOnCreditInfo.StripePayMode,
			StripePriceId: proConfig.Base.SubAddOnCreditInfo.StripePriceId,
		}

	default:
		return nil
	}
}

// GetMemberConfigByType 根据会员类型获取会员配置
func (s *ConfigService) GetMemberConfigByType(ctx context.Context, memberType pb.MembershipType) *config.MembershipTypeConfig {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipConfigService.GetMemberConfigByType")
	defer span.Finish()

	if memberType == pb.MembershipType_MEMBERSHIP_TYPE_FREE {
		return &s.memberConfig.Free
	} else if memberType == pb.MembershipType_MEMBERSHIP_TYPE_PRO {
		return &s.memberConfig.Professional
	}

	return nil
}

// GetSubPlanInfo 根据订阅类型获取订阅计划信息
func (s *ConfigService) GetSubPlanInfo(ctx context.Context, memberType pb.MembershipType) *pb.MembershipSubPlanInfo {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MembershipConfigService.GetSubPlanInfo")
	defer span.Finish()

	memberConfig := s.GetMemberConfigByType(ctx, memberType)
	if memberConfig == nil {
		return nil
	}

	// Mapping config to pb.MembershipSubPlanInfo

	pbBase := &pb.BasePermission{
		SubInfo: &pb.SubInfo{
			Type:           uint32(memberConfig.Base.SubInfo.Type),
			Name:           memberConfig.Base.SubInfo.Name,
			Price:          uint64(memberConfig.Base.SubInfo.Price),
			OriginalPrice:  uint64(memberConfig.Base.SubInfo.OriginalPrice),
			Currency:       memberConfig.Base.SubInfo.Currency,
			Credit:         uint64(memberConfig.Base.SubInfo.Credit),
			OriginalCredit: uint64(memberConfig.Base.SubInfo.OriginalCredit),
			AddOnCredit:    uint64(memberConfig.Base.SubInfo.AddOnCredit),
			Duration:       uint32(memberConfig.Base.SubInfo.Duration),
			StripePayMode:  memberConfig.Base.SubInfo.StripePayMode,
			StripePriceId:  memberConfig.Base.SubInfo.StripePriceId,
		},
		IsEnableAddOnCredit:           memberConfig.Base.IsEnableAddOnCredit,
		IsEnableSubAddOnCredit:        memberConfig.Base.IsEnableSubAddOnCredit,
		MaxAddOnCreditSubCountOfMonth: uint32(memberConfig.Base.MaxAddOnCreditSubCountOfMonth),
		SubAddOnCreditInfo: &pb.SubInfo{
			Type:           uint32(memberConfig.Base.SubAddOnCreditInfo.Type),
			Name:           memberConfig.Base.SubAddOnCreditInfo.Name,
			Price:          uint64(memberConfig.Base.SubAddOnCreditInfo.Price),
			OriginalPrice:  uint64(memberConfig.Base.SubAddOnCreditInfo.OriginalPrice),
			Currency:       memberConfig.Base.SubAddOnCreditInfo.Currency,
			Credit:         uint64(memberConfig.Base.SubAddOnCreditInfo.Credit),
			OriginalCredit: uint64(memberConfig.Base.SubAddOnCreditInfo.OriginalCredit),
			AddOnCredit:    uint64(memberConfig.Base.SubAddOnCreditInfo.AddOnCredit),
			Duration:       uint32(0),
			StripePayMode:  memberConfig.Base.SubAddOnCreditInfo.StripePayMode,
			StripePriceId:  memberConfig.Base.SubAddOnCreditInfo.StripePriceId,
		},
	}

	pbDocs := &pb.DocsPermission{
		MaxStorageCapacity:            uint64(memberConfig.Docs.MaxStorageCapacity),
		MaxStorageCapacityOriginal:    uint64(memberConfig.Docs.MaxStorageCapacityOriginal),
		DocUploadMaxSize:              uint64(memberConfig.Docs.DocUploadMaxSize),
		DocUploadMaxSizeOriginal:      uint64(memberConfig.Docs.DocUploadMaxSizeOriginal),
		DocUploadMaxPageCount:         uint32(memberConfig.Docs.DocUploadMaxPageCount),
		DocUploadMaxPageCountOriginal: uint32(memberConfig.Docs.DocUploadMaxPageCountOriginal),
	}

	pbNote := &pb.NotePermission{
		IsNoteSummary:     memberConfig.Note.IsNoteSummary,
		IsNoteWord:        memberConfig.Note.IsNoteWord,
		IsNoteExtract:     memberConfig.Note.IsNoteExtract,
		IsNoteManage:      memberConfig.Note.IsNoteManage,
		IsNotePdfDownload: memberConfig.Note.IsNotePdfDownload,
	}

	pbTranslate := &pb.TranslatePermission{
		IsOcr:                               memberConfig.Translate.IsOcr,
		OcrCreditCost:                       uint64(memberConfig.Translate.OcrCreditCost),
		IsWordTranslate:                     memberConfig.Translate.IsWordTranslate,
		WordTranslateCreditCost:             uint64(memberConfig.Translate.WordTranslateCreditCost),
		IsFullTextTranslate:                 memberConfig.Translate.IsFullTextTranslate,
		FullTextTranslateCreditCost:         uint64(memberConfig.Translate.FullTextTranslateCreditCost),
		FullTextTranslateCreditCostOriginal: uint64(memberConfig.Translate.FullTextTranslateCreditCostOriginal),
		FullTextTranslateMaxPageCount:       uint64(memberConfig.Translate.FullTextTranslateMaxPageCount),
		IsAiTranslation:                     memberConfig.Translate.IsAiTranslation,
		AiTranslationCreditCost:             uint64(memberConfig.Translate.AiTranslationCreditCost),
	}

	var pbAICopilotModels []*pb.CopilotModel
	for _, model := range memberConfig.AI.Copilot.Models {
		pbAICopilotModels = append(pbAICopilotModels, &pb.CopilotModel{
			Key:        model.Key,
			Name:       model.Name,
			IsEnable:   model.IsEnable,
			IsFree:     model.IsFree,
			CreditCost: uint64(model.CreditCost),
		})
	}
	pbAI := &pb.AIPermission{
		Copilot: &pb.CopilotPermission{
			IsEnable: memberConfig.AI.Copilot.IsEnable,
			Models:   pbAICopilotModels,
		},
	}

	return &pb.MembershipSubPlanInfo{
		Type:        uint32(memberConfig.Type),
		Name:        memberConfig.Name,
		Description: memberConfig.Description,
		IsFree:      memberConfig.IsFree,
		Base:        pbBase,
		Docs:        pbDocs,
		Note:        pbNote,
		Ai:          pbAI,
		Translate:   pbTranslate,
	}
}
