package api

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/services/pay/service"
)

// PaymentAPI 支付相关的API处理器
type PaymentAPI struct {
	paymentService *service.PaymentService
	logger         logging.Logger
}

// NewPaymentAPI 创建一个新的支付API处理器
func NewPaymentAPI(paymentService *service.PaymentService, logger logging.Logger) *PaymentAPI {
	return &PaymentAPI{
		paymentService: paymentService,
		logger:         logger,
	}
}

// InitiatePaymentRequest 发起支付请求参数
type InitiatePaymentRequest struct {
	UserId          string            `json:"user_id" binding:"required"`
	OrderId         string            `json:"order_id" binding:"required"`
	Amount          int64             `json:"amount" binding:"required,gt=0"`
	Currency        string            `json:"currency" binding:"required"`
	Description     string            `json:"description"`
	Channel         string            `json:"channel" binding:"required"`
	PaymentMethodId string            `json:"payment_method_id" binding:"required"`
	Metadata        map[string]string `json:"metadata"`
	ReturnURL       string            `json:"return_url"`
}

// PrePay 发起支付
// @Summary 发起支付
// @Description 创建新的支付交易
// @Tags 支付
// @Accept json
// @Produce json
// @Param request body InitiatePaymentRequest true "支付请求参数"
// @Success 200 {object} response.Response "支付创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /payments/prepay [post]
func (api *PaymentAPI) PrePay(c *gin.Context) {
	var req InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.logger.Error("msg", "请求参数绑定失败", "error", err.Error())
		response.Error(c, "请求参数错误: "+err.Error(), nil)
		return
	}

	// 转换为服务层参数
	params := &service.CreatePaymentParams{
		UserId:          req.UserId,
		OrderId:         req.OrderId,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Description:     req.Description,
		Channel:         req.Channel,
		PaymentMethodId: req.PaymentMethodId,
		Metadata:        req.Metadata,
		ReturnURL:       req.ReturnURL,
	}

	// 调用服务层创建支付
	paymentResult, err := api.paymentService.CreatePayment(c, params)
	if err != nil {
		api.logger.Error("msg", "创建支付失败", "error", err.Error())
		response.Error(c, "创建支付失败: "+err.Error(), nil)
		return
	}

	api.logger.Info("msg", "创建支付成功", "paymentResult", paymentResult)

	response.SuccessNoData(c, "Success")
}

// HandleStripeWebhook 处理Stripe Webhook回调
// @Summary 处理Stripe Webhook回调
// @Description 接收并处理来自Stripe的Webhook事件
// @Tags 支付
// @Accept json
// @Produce json
// @Success 200 {string} string "Webhook处理成功"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "服务器内部错误"
// @Router /payments/stripe/webhook [post]
func (api *PaymentAPI) HandleStripeWebhook(c *gin.Context) {
	// 读取请求体
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		api.logger.Error("msg", "读取Webhook请求体失败", "error", err.Error())
		response.Error(c, "读取请求体失败", nil)
		return
	}

	// 获取Stripe签名
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		api.logger.Error("msg", "缺少Stripe-Signature请求头")
		response.Error(c, "缺少Stripe-Signature请求头", nil)
		return
	}

	// 处理Webhook
	err = api.paymentService.HandleWebhook(c, "STRIPE", body, signature)
	if err != nil {
		api.logger.Error("msg", "处理Stripe Webhook失败", "error", err.Error())
		response.Error(c, "处理Webhook失败: "+err.Error(), nil)
		return
	}

	response.SuccessNoData(c, "Webhook处理成功")
}

// GetPaymentStatus 获取支付状态
// @Summary 获取支付状态
// @Description 根据支付ID获取支付状态
// @Tags 支付
// @Produce json
// @Param payment_id path int true "支付ID"
// @Success 200 {object} response.Response "支付状态"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "支付记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /payments/{payment_id}/status [get]
func (api *PaymentAPI) GetPaymentStatus(c *gin.Context) {
	paymentId := c.Param("payment_id")
	if paymentId == "" {
		api.logger.Error("msg", "支付ID格式错误", "paymentId", paymentId)
		response.Error(c, "支付ID格式错误", nil)
		return
	}

	_, err := api.paymentService.GetPaymentStatus(c, paymentId)
	if err != nil {
		api.logger.Error("msg", "获取支付状态失败", "paymentId", paymentId, "error", err.Error())
		response.Error(c, "获取支付状态失败: "+err.Error(), nil)
		return
	}

	response.SuccessNoData(c, "Success")
}

// CreateRefundRequest 创建退款请求参数
type CreateRefundRequest struct {
	Amount int64  `json:"amount" binding:"required,gt=0"`
	Reason string `json:"reason"`
}

// CreateRefund 创建退款
// @Summary 创建退款
// @Description 对指定支付创建退款
// @Tags 支付
// @Accept json
// @Produce json
// @Param payment_id path int true "支付ID"
// @Param request body CreateRefundRequest true "退款请求参数"
// @Success 200 {object} response.Response "退款创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "支付记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /payments/{payment_id}/refund [post]
func (api *PaymentAPI) CreateRefund(c *gin.Context) {
	paymentId := c.Param("payment_id")
	if paymentId == "" {
		api.logger.Error("msg", "支付ID格式错误", "paymentId", paymentId)
		response.Error(c, "支付ID格式错误", nil)
		return
	}

	var req CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.logger.Error("msg", "请求参数绑定失败", "error", err.Error())
		response.Error(c, "请求参数错误: "+err.Error(), nil)
		return
	}

	err := api.paymentService.CreateRefund(c, paymentId, req.Amount, req.Reason)
	if err != nil {
		api.logger.Error("msg", "创建退款失败", "paymentId", paymentId, "error", err.Error())
		response.Error(c, "创建退款失败: "+err.Error(), nil)
		return
	}

	response.SuccessNoData(c, "Success")
}

// GetUserPayments 获取用户支付记录
// @Summary 获取用户支付记录
// @Description 获取指定用户的支付记录
// @Tags 支付
// @Produce json
// @Param user_id query string true "用户ID"
// @Param status query string false "支付状态"
// @Success 200 {object} response.Response "用户支付记录"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /payments/user [get]
func (api *PaymentAPI) GetUserPayments(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		response.Error(c, "用户ID不能为空", nil)
		return
	}

	status := c.Query("status")

	_, err := api.paymentService.GetPaymentsByUserId(c, userId, status)
	if err != nil {
		api.logger.Error("msg", "获取用户支付记录失败", "userId", userId, "error", err.Error())
		response.Error(c, "获取用户支付记录失败: "+err.Error(), nil)
		return
	}

	response.SuccessNoData(c, "Success")
}

// // RegisterRoutes 注册API路由
// func (api *PaymentAPI) RegisterRoutes(router *gin.RouterGroup) {
// 	paymentGroup := router.Group("/payments")
// 	{
// 		paymentGroup.POST("/initiate", api.InitiatePayment)
// 		paymentGroup.POST("/stripe/webhook", api.HandleStripeWebhook)
// 		paymentGroup.GET("/:payment_id/status", api.GetPaymentStatus)
// 		paymentGroup.POST("/:payment_id/refund", api.CreateRefund)
// 		paymentGroup.GET("/user", api.GetUserPayments)
// 	}
// }
