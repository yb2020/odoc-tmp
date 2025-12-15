package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pay/dao"
	"github.com/yb2020/odoc/services/pay/model"
	"github.com/yb2020/odoc/services/pay/provider"
)

// PaymentService 支付服务，协调支付记录DAO和支付提供商
type PaymentService struct {
	paymentRecordDAO       *dao.PaymentRecordDAO
	paymentProviderFactory *provider.PaymentProviderFactory
	logger                 logging.Logger
}

// NewPaymentService 创建一个新的支付服务
func NewPaymentService(
	paymentRecordDAO *dao.PaymentRecordDAO,
	paymentProviderFactory *provider.PaymentProviderFactory,
	logger logging.Logger,
) *PaymentService {
	return &PaymentService{
		paymentRecordDAO:       paymentRecordDAO,
		paymentProviderFactory: paymentProviderFactory,
		logger:                 logger,
	}
}

// CreatePaymentParams 创建支付的参数
type CreatePaymentParams struct {
	UserId          string            // 用户ID
	OrderId         string            // 关联的业务订单ID
	Amount          int64             // 金额，以最小货币单位表示（例如：分）
	Currency        string            // ISO 4217 货币代码 (例如 "CNY", "USD")
	Description     string            // 支付描述
	Channel         string            // 支付渠道
	PaymentMethodId string            // 支付方式ID
	Metadata        map[string]string // 元数据
	ReturnURL       string            // 支付完成后的返回URL
}

// CreatePaymentResult 创建支付的结果
type CreatePaymentResult struct {
	PaymentId      string // 支付记录ID
	ProviderTxId   string // 支付渠道的交易ID
	Status         string // 支付状态
	ClientSecret   string // 客户端密钥（如Stripe的client_secret）
	RedirectURL    string // 重定向URL（对于需要重定向的支付方式）
	RequiresAction bool   // 是否需要额外的客户操作
	ErrorCode      string // 错误码（如果有）
	ErrorMessage   string // 错误信息（如果有）
}

// CreatePayment 创建支付
func (s *PaymentService) CreatePayment(ctx context.Context, params *CreatePaymentParams) (*CreatePaymentResult, error) {
	// 获取支付提供商
	paymentProvider, err := s.paymentProviderFactory.GetProvider(params.Channel)
	if err != nil {
		s.logger.Error("msg", "获取支付提供商失败", "channel", params.Channel, "error", err.Error())
		return nil, errors.Biz("获取支付提供商失败")
	}

	// 创建支付记录
	paymentRecord := &model.PaymentRecord{
		UserId:      params.UserId,
		OrderId:     params.OrderId,
		Amount:      params.Amount,
		Currency:    params.Currency,
		Status:      model.PaymentStatusPending,
		Channel:     params.Channel,
		Description: params.Description,
	}

	// 序列化元数据
	if params.Metadata != nil {
		metadataJSON, err := json.Marshal(params.Metadata)
		if err != nil {
			s.logger.Error("msg", "序列化元数据失败", "error", err.Error())
			return nil, errors.Biz("序列化元数据失败")
		}
		paymentRecord.Metadata = string(metadataJSON)
	}

	// 保存支付记录
	err = s.paymentRecordDAO.Save(ctx, paymentRecord)
	if err != nil {
		s.logger.Error("msg", "创建支付记录失败", "error", err.Error())
		return nil, errors.Biz("创建支付记录失败")
	}

	// 构建支付提供商参数
	chargeParams := &provider.ChargeParams{
		Amount:          params.Amount,
		Currency:        params.Currency,
		Description:     params.Description,
		OrderId:         params.OrderId,
		UserId:          params.UserId,
		PaymentMethodId: params.PaymentMethodId,
		Metadata:        params.Metadata,
		ReturnURL:       params.ReturnURL,
	}

	// 调用支付提供商创建支付
	chargeResult, err := paymentProvider.CreateCharge(ctx, chargeParams)
	if err != nil {
		// 更新支付记录为失败状态
		s.updatePaymentRecordStatus(ctx, paymentRecord.Id, model.PaymentStatusFailed, chargeResult)
		s.logger.Error("msg", "调用支付提供商创建支付失败", "error", err.Error())
		return nil, errors.Biz("调用支付提供商创建支付失败")
	}

	// 更新支付记录
	s.updatePaymentRecordStatus(ctx, paymentRecord.Id, chargeResult.Status, chargeResult)

	// 构建返回结果
	result := &CreatePaymentResult{
		PaymentId:      paymentRecord.Id,
		ProviderTxId:   chargeResult.ProviderTxId,
		Status:         chargeResult.Status,
		ClientSecret:   chargeResult.ClientSecret,
		RedirectURL:    chargeResult.RedirectURL,
		RequiresAction: chargeResult.RequiresAction,
		ErrorCode:      chargeResult.ErrorCode,
		ErrorMessage:   chargeResult.ErrorMessage,
	}

	return result, nil
}

// HandleWebhook 处理支付提供商的Webhook回调
func (s *PaymentService) HandleWebhook(ctx context.Context, channel string, requestData []byte, signature string) error {
	// 获取支付提供商
	paymentProvider, err := s.paymentProviderFactory.GetProvider(channel)
	if err != nil {
		s.logger.Error("msg", "获取支付提供商失败", "channel", channel, "error", err.Error())
		return errors.Biz("获取支付提供商失败")
	}

	// 解析Webhook数据
	webhookEvent, err := paymentProvider.HandleWebhook(ctx, requestData, signature)
	if err != nil {
		s.logger.Error("msg", "解析Webhook数据失败", "error", err.Error())
		return errors.Biz("解析Webhook数据失败")
	}

	// 查找对应的支付记录
	paymentRecord, err := s.paymentRecordDAO.GetByProviderTxId(ctx, webhookEvent.ProviderTxId)
	if err != nil {
		s.logger.Error("msg", "查询支付记录失败", "error", err.Error())
		return errors.Biz("查询支付记录失败")
	}

	// 更新支付记录状态
	if webhookEvent.Status == model.PaymentStatusSucceeded {
		// 支付成功，更新支付记录
		updateFields := map[string]interface{}{
			"status":         model.PaymentStatusSucceeded,
			"provider_tx_id": webhookEvent.ProviderTxId,
			"paid_at":        time.Now(),
		}
		err = s.paymentRecordDAO.UpdateFields(ctx, paymentRecord.Id, updateFields)
		if err != nil {
			s.logger.Error("msg", "更新支付记录失败", "error", err.Error())
			return errors.Biz("更新支付记录失败")
		}

		// TODO: 触发支付成功后的业务逻辑，如更新订单状态、发送通知等
		// 这里可以通过事件总线或直接调用其他服务来完成

	} else if webhookEvent.Status == model.PaymentStatusFailed || webhookEvent.Status == model.PaymentStatusCanceled {
		// 支付失败或取消，更新支付记录
		updateFields := map[string]interface{}{
			"status":         webhookEvent.Status,
			"provider_tx_id": webhookEvent.ProviderTxId,
		}
		err = s.paymentRecordDAO.UpdateFields(ctx, paymentRecord.Id, updateFields)
		if err != nil {
			s.logger.Error("msg", "更新支付记录失败", "error", err.Error())
			return errors.Biz("更新支付记录失败")
		}

		// TODO: 触发支付失败后的业务逻辑，如释放库存、通知用户等
	}

	return nil
}

// CreateRefund 创建退款
func (s *PaymentService) CreateRefund(ctx context.Context, paymentId string, amount int64, reason string) error {
	// 查找支付记录
	paymentRecord, err := s.paymentRecordDAO.FindById(ctx, paymentId)
	if err != nil {
		s.logger.Error("msg", "查询支付记录失败", "error", err.Error())
		return errors.Biz("查询支付记录失败")
	}

	// 检查支付状态
	if paymentRecord.Status != model.PaymentStatusSucceeded {
		return errors.Biz("只有成功的支付才能退款")
	}

	// 检查退款金额
	if amount <= 0 || amount > paymentRecord.Amount {
		return errors.Biz("无效的退款金额")
	}

	// 获取支付提供商
	paymentProvider, err := s.paymentProviderFactory.GetProvider(paymentRecord.Channel)
	if err != nil {
		s.logger.Error("msg", "获取支付提供商失败", "channel", paymentRecord.Channel, "error", err.Error())
		return errors.Biz("获取支付提供商失败")
	}

	// 调用支付提供商创建退款
	refundResult, err := paymentProvider.CreateRefund(ctx, paymentRecord.ProviderTxId, amount, reason)
	if err != nil {
		s.logger.Error("msg", "创建退款失败", "error", err.Error())
		return errors.Biz("创建退款失败")
	}

	// 更新支付记录状态
	updateFields := map[string]interface{}{
		"status": model.PaymentStatusRefunded,
	}
	err = s.paymentRecordDAO.UpdateFields(ctx, paymentId, updateFields)
	if err != nil {
		s.logger.Error("msg", "更新支付记录失败", "error", err.Error())
		return errors.Biz("更新支付记录失败")
	}

	// TODO: 创建退款记录
	// TODO: 触发退款成功后的业务逻辑，如更新订单状态、发送通知等

	s.logger.Info("msg", "退款成功", "paymentId", paymentId, "refundId", refundResult.RefundId)
	return nil
}

// GetPaymentStatus 获取支付状态
func (s *PaymentService) GetPaymentStatus(ctx context.Context, paymentId string) (string, error) {
	// 查找支付记录
	paymentRecord, err := s.paymentRecordDAO.FindById(ctx, paymentId)
	if err != nil {
		s.logger.Error("msg", "查询支付记录失败", "error", err.Error())
		return "", errors.Biz("查询支付记录失败")
	}

	// 如果支付状态是终态，直接返回
	if isTerminalStatus(paymentRecord.Status) {
		return paymentRecord.Status, nil
	}

	// 如果支付状态不是终态，查询支付提供商获取最新状态
	if paymentRecord.ProviderTxId != "" {
		paymentProvider, err := s.paymentProviderFactory.GetProvider(paymentRecord.Channel)
		if err != nil {
			s.logger.Error("msg", "获取支付提供商失败", "channel", paymentRecord.Channel, "error", err.Error())
			return paymentRecord.Status, nil // 返回数据库中的状态
		}

		// 查询支付提供商获取最新状态
		latestStatus, err := paymentProvider.GetChargeStatus(ctx, paymentRecord.ProviderTxId)
		if err != nil {
			s.logger.Error("msg", "查询支付状态失败", "error", err.Error())
			return paymentRecord.Status, nil // 返回数据库中的状态
		}

		// 如果状态有变化，更新支付记录
		if latestStatus != paymentRecord.Status {
			updateFields := map[string]interface{}{
				"status": latestStatus,
			}

			// 如果支付成功，更新支付时间
			if latestStatus == model.PaymentStatusSucceeded {
				updateFields["paid_at"] = time.Now()
			}

			err = s.paymentRecordDAO.UpdateFields(ctx, paymentId, updateFields)
			if err != nil {
				s.logger.Error("msg", "更新支付记录失败", "error", err.Error())
			}

			return latestStatus, nil
		}
	}

	return paymentRecord.Status, nil
}

// GetPaymentsByUserId 获取用户的支付记录
func (s *PaymentService) GetPaymentsByUserId(ctx context.Context, userId string, status string) ([]*model.PaymentRecord, error) {
	records, err := s.paymentRecordDAO.GetByUserIdAndStatus(ctx, userId, status)
	if err != nil {
		return nil, err
	}

	// 转换为指针切片
	result := make([]*model.PaymentRecord, len(records))
	for i := range records {
		result[i] = &records[i]
	}

	return result, nil
}

// updatePaymentRecordStatus 更新支付记录状态
func (s *PaymentService) updatePaymentRecordStatus(ctx context.Context, paymentId string, status string, chargeResult *provider.ChargeResult) {
	updateFields := map[string]interface{}{
		"status":         status,
		"provider_tx_id": chargeResult.ProviderTxId,
	}

	// 如果有错误信息，更新错误信息
	if chargeResult.ErrorCode != "" {
		updateFields["provider_error_code"] = chargeResult.ErrorCode
	}
	if chargeResult.ErrorMessage != "" {
		updateFields["provider_error_message"] = chargeResult.ErrorMessage
	}

	// 如果支付成功，更新支付时间
	if status == model.PaymentStatusSucceeded {
		updateFields["paid_at"] = time.Now()
	}

	// 更新支付方式类型
	if chargeResult.PaymentMethodType != "" {
		updateFields["payment_method_type"] = chargeResult.PaymentMethodType
	}

	err := s.paymentRecordDAO.UpdateFields(ctx, paymentId, updateFields)
	if err != nil {
		s.logger.Error("msg", "更新支付记录失败", "paymentId", paymentId, "error", err.Error())
	}
}

// isTerminalStatus 判断支付状态是否是终态
func isTerminalStatus(status string) bool {
	return status == model.PaymentStatusSucceeded ||
		status == model.PaymentStatusFailed ||
		status == model.PaymentStatusCanceled ||
		status == model.PaymentStatusRefunded
}
