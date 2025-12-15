package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	pb "github.com/yb2020/odoc-proto/gen/go/oss"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/services/oss/constant"
	serviceModel "github.com/yb2020/odoc/services/oss/model"
	"github.com/yb2020/odoc/services/oss/service"
)

// CallbackAPI OSS回调API
type CallbackAPI struct {
	config  *config.Config
	logger  logging.Logger
	tracer  opentracing.Tracer
	service *service.CallbackService
}

// NewCallbackAPI 创建OSS回调API
func NewCallbackAPI(config *config.Config, logger logging.Logger, tracer opentracing.Tracer, service *service.CallbackService) *CallbackAPI {
	return &CallbackAPI{
		config:  config,
		logger:  logger,
		tracer:  tracer,
		service: service,
	}
}

// HandleCallback 处理OSS回调
func (api *CallbackAPI) HandleCallback(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "CallbackAPI.HandleCallback")
	defer span.Finish()
	authToken := c.GetHeader(constant.HeaderAuthorization)
	// 如果有Bearer前缀，则去除前缀
	if len(authToken) > len(constant.AuthTokenPrefix) && authToken[:len(constant.AuthTokenPrefix)] == constant.AuthTokenPrefix {
		authToken = authToken[len(constant.AuthTokenPrefix):]
	}
	if authToken != api.config.OSS.S3.Webhook.AuthToken {
		response.ErrorNoData(c, "认证令牌无效")
		return
	}
	// 读取请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// 记录原始数据
	api.logger.Info("msg", "收到MinIO回调原始数据", "body", string(bodyBytes))
	// 使用 BindProto 解析通知数据
	var notification pb.MinioCallbackNotification
	if err := json.Unmarshal(bodyBytes, &notification); err != nil {
		api.logger.Warn("msg", "解析通知数据失败", "error", err.Error())
		c.Error(err)
		return
	}
	// 处理对象键中的URL编码
	for i, record := range notification.Records {
		if record.S3 != nil && record.S3.Object != nil {
			// URL解码对象键
			decodedKey, err := url.QueryUnescape(record.S3.Object.Key)
			if err == nil && decodedKey != record.S3.Object.Key {
				notification.Records[i].S3.Object.Key = decodedKey
			}
		}
	}
	// 处理回调通知
	if err := api.service.HandleCallback(ctx, &notification); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, "处理成功", nil)
}

// HandleS3Callback 处理S3回调
func (api *CallbackAPI) HandleS3Callback(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "CallbackAPI.HandleS3Callback")
	defer span.Finish()

	authToken := c.Query("token")
	if authToken != api.config.OSS.S3.Webhook.AuthToken {
		response.ErrorNoData(c, "认证令牌无效")
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	api.logger.Info("msg", "收到S3回调原始数据", "body", string(bodyBytes))

	// 先尝试判断是否为直接的S3事件消息（检查是否有Records字段）
	var rawMessage map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawMessage); err != nil {
		api.logger.Error("msg", "解析JSON数据失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 检查是否为直接的S3事件消息（有Records字段但没有Type字段）
	if _, hasRecords := rawMessage["Records"]; hasRecords && rawMessage["Type"] == nil {
		api.logger.Info("msg", "检测到直接的S3事件消息")

		// 解析为S3事件
		var s3Event pb.MinioCallbackNotification
		if err := json.Unmarshal(bodyBytes, &s3Event); err != nil {
			api.logger.Error("msg", "解析S3事件失败", "error", err.Error())
			c.Error(err)
			return
		}

		// 处理对象键中的URL编码
		for i, record := range s3Event.Records {
			if record.S3 != nil && record.S3.Object != nil {
				// URL解码对象键
				decodedKey, err := url.QueryUnescape(record.S3.Object.Key)
				if err == nil && decodedKey != record.S3.Object.Key {
					s3Event.Records[i].S3.Object.Key = decodedKey
				}
			}
		}

		// 处理回调通知
		if err := api.service.HandleCallback(ctx, &s3Event); err != nil {
			api.logger.Error("msg", "处理S3事件失败", "error", err.Error())
			c.Error(err)
			return
		}

		response.Success(c, "处理成功", nil)
		return
	}

	// 否则按SNS消息处理
	var snsMessage serviceModel.SNSMessage
	if err := json.Unmarshal(bodyBytes, &snsMessage); err != nil {
		api.logger.Warn("msg", "解析S3回调JSON数据失败", "error", err.Error())
		c.Error(err)
		return
	}

	// 根据消息类型分别处理
	switch snsMessage.Type {
	case "SubscriptionConfirmation":
		// 这是订阅确认请求，需要自动访问SubscribeURL来完成确认
		api.logger.Info("msg", "收到SNS订阅确认请求，正在自动确认", "url", snsMessage.SubscribeURL)

		resp, err := http.Get(snsMessage.SubscribeURL)
		if err != nil {
			api.logger.Error("msg", "自动确认SNS订阅失败", "url", snsMessage.SubscribeURL, "error", err.Error())
			// 即使自动确认失败，也告知上游收到请求，以便手动重试
			response.Success(c, "订阅确认请求已收到，但自动确认失败", nil)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			api.logger.Info("msg", "成功自动确认SNS订阅", "url", snsMessage.SubscribeURL, "status", resp.Status)
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			api.logger.Warn("msg", "自动确认SNS订阅返回非成功状态", "url", snsMessage.SubscribeURL, "status", resp.Status, "body", string(bodyBytes))
		}

		response.Success(c, "订阅确认请求已收到并自动确认", nil)
		return
	case "Notification":
		// 这是文件上传等通知
		if err := api.service.HandleS3Notification(ctx, snsMessage); err != nil {
			c.Error(err)
			return
		}
	default:
		api.logger.Warn("msg", "收到未知的SNS消息类型", "type", snsMessage.Type)
	}

	response.Success(c, "处理成功", nil)
}
