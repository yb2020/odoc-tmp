package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	pb "github.com/yb2020/odoc/proto/gen/go/pay"
	msInterfaces "github.com/yb2020/odoc/services/membership/interfaces"
	"github.com/yb2020/odoc/services/pay/service"
)

// StripeCheckoutAPI 封装了 Stripe Checkout 相关的 API 处理器
type StripeCheckoutAPI struct {
	checkoutService *service.StripeCheckoutService
	orderService    msInterfaces.IOrderService
	logger          logging.Logger
	tracer          opentracing.Tracer
}

// NewStripeCheckoutAPI 创建一个新的 CheckoutAPI 实例
func NewStripeCheckoutAPI(cs *service.StripeCheckoutService, orderService msInterfaces.IOrderService, logger logging.Logger, tracer opentracing.Tracer) *StripeCheckoutAPI {
	return &StripeCheckoutAPI{
		checkoutService: cs,
		orderService:    orderService,
		logger:          logger,
		tracer:          tracer,
	}
}

// @api_path: /api/pay/stripe/get-publishable-key
// @method: GET
// @summary: 获取 Stripe Publishable Key
func (api *StripeCheckoutAPI) GetStripePublishableKey(c *gin.Context) {
	publicKey := api.checkoutService.PublishableKey()
	if publicKey == "" {
		response.ErrorNoData(c, "public key not found")
		return
	}

	res := &pb.GetStripePublishableKeyResp{
		PublishableKey: publicKey,
	}

	response.Success(c, "success", res)
}

// @api_path: /api/pay/stripe-checkout/create-session
// @method: POST
// @summary: 创建 Stripe Checkout Session
func (api *StripeCheckoutAPI) CreateCheckoutSession(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "StripeCheckoutAPI.CreateCheckoutSessionReq")
	defer span.Finish()

	//userId, _ := userContext.GetUserID(ctx)
	req := &pb.CreateCheckoutSessionReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "Failed to parse request parameters", "error", err.Error())
		response.ErrorNoData(c, "Failed to parse request parameters")
		return
	}

	// TODO: 订单信息校验
	order, err := api.orderService.GetById(ctx, req.OrderId)
	if err != nil {
		api.logger.Error("msg", "Failed to get order by id", "error", err.Error())
		response.ErrorNoData(c, "Failed to get order by id")
		return
	}
	if order == nil {
		response.ErrorNoData(c, "order not found")
		return
	}

	checkoutSessionParams := service.CreateCheckoutSessionParams{
		PayMode:  order.StripePayMode,
		PriceId:  order.StripePriceId,
		Name:     order.SubName,
		Amount:   order.PayAmount,
		Quantity: int64(order.NumberCount),
		Currency: order.Currency,
		OrderId:  req.OrderId,
		UserId:   order.UserId,
	}

	result, err := api.checkoutService.CreateCheckoutSession(ctx, checkoutSessionParams)
	if err != nil {
		api.logger.Error("createCheckoutSessionHandler: 创建Stripe Checkout Session失败: %v", err)
		response.ErrorNoData(c, "创建支付会话失败: "+err.Error())
		return
	}

	CreateCheckoutSessionResp := &pb.CreateCheckoutSessionResp{
		SessionId: result.SessionId,
	}
	response.Success(c, "success", CreateCheckoutSessionResp)
}

// @api_path: /api/pay/stripe-checkout/cancel-subscription-at-period-end
// @method: POST
// @summary: 取消订阅，在周期结束时生效
func (api *StripeCheckoutAPI) CancelSubscriptionAtPeriodEnd(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "StripeCheckoutAPI.CancelSubscriptionAtPeriodEnd")
	defer span.Finish()
	var req struct {
		SubId string `json:"subId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorNoData(c, "invalid request: "+err.Error())
		return
	}
	subId := req.SubId
	if subId == "" {
		response.ErrorNoData(c, "subId is empty")
		return
	}

	_, err := api.checkoutService.CancelSubscriptionAtPeriodEnd(ctx, subId)
	if err != nil {
		api.logger.Error("cancelSubscriptionAtPeriodEndHandler: 取消订阅失败: %v", err)
		response.ErrorNoData(c, "取消订阅失败: "+err.Error())
		return
	}
	response.SuccessNoData(c, "success")
}

// @api_path: /api/pay/stripe-checkout/cancel-subscription-immediately
// @method: POST
// @summary: 立即取消订阅
func (api *StripeCheckoutAPI) CancelSubscriptionImmediately(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "StripeCheckoutAPI.CancelSubscriptionImmediately")
	defer span.Finish()
	var req struct {
		SubId string `json:"subId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorNoData(c, "invalid request: "+err.Error())
		return
	}
	subId := req.SubId
	if subId == "" {
		response.ErrorNoData(c, "subId is empty")
		return
	}

	_, err := api.checkoutService.CancelSubscriptionImmediately(ctx, subId)
	if err != nil {
		api.logger.Error("CancelSubscriptionImmediately: 取消订阅失败: %v", err)
		response.ErrorNoData(c, "取消订阅失败: "+err.Error())
		return
	}
	response.SuccessNoData(c, "success")
}

// @api_path: /services/pay/stripe-checkout/webhook
// @method: POST
// @summary: 处理 Stripe Checkout Webhook 事件
func (a *StripeCheckoutAPI) HandleCheckoutWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536) // 64KB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("handleCheckoutWebhook: 读取请求体失败: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "读取请求失败"})
		return
	}

	signatureHeader := c.GetHeader("Stripe-Signature")
	err = a.checkoutService.HandleCheckoutWebhook(c.Request.Context(), payload, signatureHeader)
	if err != nil {
		a.logger.Error("handleCheckoutWebhook: 处理Webhook失败: %v", err)
		// 根据错误类型返回不同的状态码，例如签名验证失败返回 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook 处理失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
