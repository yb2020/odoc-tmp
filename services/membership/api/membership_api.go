package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/membership"
	"github.com/yb2020/odoc/config"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/membership/service"
	userService "github.com/yb2020/odoc/services/user/service"
)

type MembershipApi struct {
	logger            logging.Logger
	tracer            opentracing.Tracer
	config            *config.Config
	membershipService interfaces.IMembershipService
	msConfigService   *service.ConfigService
	userService       *userService.UserService
}

func NewMembershipApi(logger logging.Logger,
	tracer opentracing.Tracer,
	config *config.Config,
	membershipService interfaces.IMembershipService,
	msConfigService *service.ConfigService,
	userService *userService.UserService,
) *MembershipApi {
	return &MembershipApi{
		logger:            logger,
		tracer:            tracer,
		config:            config,
		membershipService: membershipService,
		msConfigService:   msConfigService,
		userService:       userService,
	}
}

func (api *MembershipApi) Test(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.Test")
	defer span.Finish()

	// membershipConfig := api.membershipService.GetMemberConfig(ctx)
	// api.logger.Info("membershipConfig", "membershipConfig", membershipConfig)

	userId, _ := userContext.GetUserID(ctx)
	api.logger.Info("userId", "userId", userId)

	// // 1.创建用户会员账号
	// userMembership, err := api.membershipService.NewMembershipAccount(ctx, userId)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }
	// api.logger.Info("userMembership", "userMembership", userMembership)

	response.SuccessNoData(c, "Success")
}

func (api *MembershipApi) TestSubFree(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.TestSubFree")
	defer span.Finish()

	userId, _ := userContext.GetUserID(ctx)
	api.logger.Info("userId", "userId", userId)
	// 1.自动订阅Free会员
	err := api.membershipService.AutoSubscribeMemberFree(ctx, userId)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessNoData(c, "Success")
}

func (api *MembershipApi) TestCallCreditFunsDocsUpload(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.Test")
	defer span.Finish()

	// 4.调用积分接口上传接口功能
	err := api.membershipService.CreditFunDocsUpload(ctx, 1048576, 10, 1048576, func(xctx context.Context, sessionId string) error {
		api.logger.Info("msg", "do something", "sessionId", sessionId)
		return nil
	}, true)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessNoData(c, "Success")
}

func (api *MembershipApi) TestCallCreditFunsAi(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.Test")
	defer span.Finish()

	sessionId1 := "0"
	// 1.调用积分接口AI Copilot接口功能
	err := api.membershipService.CreditFunAi(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_AI_COPILOT, "gpt-4o-mini", func(xctx context.Context, sessionId string) error {
		api.logger.Info("msg", "do something", "sessionId", sessionId)
		sessionId1 = sessionId
		// return errors.Biz("test error")
		return nil
	}, false)
	if err != nil {
		c.Error(err)
		return
	}

	api.logger.Info("sessionId1", "sessionId1", sessionId1)

	//2.确认积分服务成功
	err = api.membershipService.ConfirmCreditFun(ctx, sessionId1)
	if err != nil {
		c.Error(err)
		return
	}

	//3.回退积分服务失败
	err = api.membershipService.RetrieveCreditFun(ctx, sessionId1)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessNoData(c, "Success")
}

func (api *MembershipApi) TestCallCreditFunsTranslate(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.Test")
	defer span.Finish()

	// 1.调用积分接口OCR翻译接口功能
	// err := api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_OCR, 0, func(xctx context.Context, sessionId int64) error {
	// 	api.logger.Info("msg", "do something", "sessionId", sessionId)
	// 	return nil
	// }, true)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	// // 2.调用积分接口单词翻译接口功能
	// err = api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_WORD, 0, func(xctx context.Context, sessionId int64) error {
	// 	api.logger.Info("msg", "do something", "sessionId", sessionId)
	// 	return nil
	// }, true)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	sessionId2 := "0"
	// 3.调用积分接口全文翻译接口功能
	err := api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_FULLTEXT, 0, func(xctx context.Context, sessionId string) error {
		api.logger.Info("msg", "do something", "sessionId", sessionId)
		sessionId2 = sessionId
		return nil
	}, false)
	if err != nil {
		c.Error(err)
		return
	}
	api.membershipService.RetrieveCreditFun(ctx, sessionId2)

	// 4.调用积分接口AI翻译接口功能
	// err = api.membershipService.CreditFunTranslate(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_TRANSLATE_AI, 0, func(xctx context.Context, sessionId int64) error {
	// 	api.logger.Info("msg", "do something", "sessionId", sessionId)
	// 	return nil
	// }, true)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	response.SuccessNoData(c, "Success")
}

func (api *MembershipApi) TestCallCreditFunsNote(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.Test")
	defer span.Finish()

	// 1.调用笔记总结接口功能
	err := api.membershipService.CreditFunNote(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_SUMMARY, func(xctx context.Context, sessionId string) error {
		api.logger.Info("msg", "do something", "sessionId", sessionId)
		return nil
	}, true)
	if err != nil {
		c.Error(err)
		return
	}

	// 2.调用笔记单词接口功能
	err = api.membershipService.CreditFunNote(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_WORD, func(xctx context.Context, sessionId string) error {
		api.logger.Info("msg", "do something", "sessionId", sessionId)
		return nil
	}, true)
	if err != nil {
		c.Error(err)
		return
	}

	// 3.调用笔记提取接口功能
	err = api.membershipService.CreditFunNote(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_EXTRACT, func(xctx context.Context, sessionId string) error {
		api.logger.Info("msg", "do something", "sessionId", sessionId)
		return nil
	}, true)
	if err != nil {
		c.Error(err)
		return
	}

	// 4.调用笔记管理接口功能
	err = api.membershipService.CreditFunNote(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_MANAGE, func(xctx context.Context, sessionId string) error {
		api.logger.Info("msg", "do something", "sessionId", sessionId)
		return nil
	}, true)
	if err != nil {
		c.Error(err)
		return
	}

	// 5.调用笔记PDF下载接口功能
	err = api.membershipService.CreditFunNote(ctx, pb.CreditServiceType_CREDIT_SERVICE_TYPE_NOTE_PDF_DOWNLOAD, func(xctx context.Context, sessionId string) error {
		api.logger.Info("msg", "do something", "sessionId", sessionId)
		return nil
	}, true)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessNoData(c, "Success")
}

//==========正式API======================================//

// @api /api/public/membership/get-subplan-infos
// @method GET
// @apiDescription 获取订阅信息列表
func (api *MembershipApi) GetSubPlanInfos(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.GetSubPlanInfos")
	defer span.Finish()

	subInfoList := make([]*pb.MembershipSubPlanInfo, 0)
	subInfoFree := api.msConfigService.GetSubPlanInfo(ctx, pb.MembershipType_MEMBERSHIP_TYPE_FREE)
	if subInfoFree != nil {
		subInfoList = append(subInfoList, subInfoFree)
	}

	subInfoPro := api.msConfigService.GetSubPlanInfo(ctx, pb.MembershipType_MEMBERSHIP_TYPE_PRO)
	if subInfoPro != nil {
		subInfoList = append(subInfoList, subInfoPro)
	}

	response.Success(c, "Success", &pb.GetSubPlanInfosResponse{
		SubPlanInfos: subInfoList,
	})
}

// @api /api/membership/get-base-info
// @method GET
// @apiDescription 获取用户会员基础信息
func (api *MembershipApi) GetBaseInfo(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.GetBaseInfo")
	defer span.Finish()

	userId, _ := userContext.GetUserID(ctx)
	userMembershipBaseInfo, err := api.membershipService.GetBaseInfo(ctx, userId)
	if err != nil {
		c.Error(err)
		return
	}

	rep := &pb.GetBaseInfoResponse{
		BaseInfo: userMembershipBaseInfo,
	}
	response.Success(c, "Success", rep)
}

// @api /api/membership/get-info
// @method GET
// @apiDescription 获取用户会员信息
func (api *MembershipApi) GetInfo(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.GetInfo")
	defer span.Finish()

	userId, _ := userContext.GetUserID(ctx)
	userMembershipInfo, err := api.membershipService.GetInfo(ctx, userId)
	if err != nil {
		c.Error(err)
		return
	}

	rep := &pb.GetInfoResponse{
		Info: userMembershipInfo,
	}
	response.Success(c, "Success", rep)
}

// @api /api/membership/user/profile
// @method GET
// @apiDescription 获取用户会员信息
func (api *MembershipApi) GetMembershipAndUserInfo(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.GetUserInfo")
	defer span.Finish()

	userId, _ := userContext.GetUserID(ctx)
	userMembershipInfo, err := api.membershipService.GetInfo(ctx, userId)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := api.userService.GetUserByID(ctx, userId)
	if err != nil {
		c.Error(err)
		return
	}

	rep := &pb.GetMembershipProfileResponse{
		Info: userMembershipInfo,
		User: user.ToProto(),
	}
	response.Success(c, "Success", rep)
}

// @api /api/public/membership/get-login-page-config
// @method GET
// @apiDescription 获取登录页面配置
func (api *MembershipApi) GetLoginPageConfig(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.GetLoginPageConfig")
	defer span.Finish()

	rep := &pb.GetLoginPageConfigResponse{
		IsGoogleLoginEnabled:           api.config.LoginPage.IsGoogleLoginEnabled,
		IsUsernamePasswordLoginEnabled: api.config.LoginPage.IsUsernamePasswordLoginEnabled,
		IsRegisterEnabled:              api.config.LoginPage.IsRegisterEnabled,
		IsForgetPasswordEnabled:        api.config.LoginPage.IsForgetPasswordEnabled,
	}
	response.Success(c, "Success", rep)
}

// @api /api/membership/get-user-credit
// @method GET
// @apiDescription 获取用户积分
func (api *MembershipApi) GetUserCredit(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "MembershipApi.GetUserCredit")
	defer span.Finish()

	userId, _ := userContext.GetUserID(ctx)
	userBaseInfo, err := api.membershipService.GetBaseInfo(ctx, userId)
	if err != nil {
		c.Error(err)
		return
	}
	api.logger.Info("msg", "get user credit", "userCredit", userBaseInfo.Credit)

	rep := &pb.GetUserCreditResponse{
		Credit:      userBaseInfo.Credit,
		AddOnCredit: userBaseInfo.AddOnCredit,
	}
	response.Success(c, "Success", rep)
}
