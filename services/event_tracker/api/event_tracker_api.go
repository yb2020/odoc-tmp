package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/proto/gen/go/tracker"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	transport "github.com/yb2020/odoc/pkg/transport"
)

const (
	secretKey     = "rp1012162"
	algorithmName = "HmacSHA256"
)

type EventTrackerAPI struct {
	logger logging.Logger
	tracer opentracing.Tracer
}

func NewEventTrackerAPI(logger logging.Logger,
	tracer opentracing.Tracer,
) *EventTrackerAPI {
	return &EventTrackerAPI{
		logger: logger,
		tracer: tracer,
	}
}

// validateSignature 验证签名
func (api *EventTrackerAPI) validateSignature(rawString, signature string) bool {
	if rawString == "" || signature == "" {
		return false
	}

	// 创建 HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(rawString))

	// 计算签名
	computedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 比较签名
	if computedSignature != signature {
		api.logger.Error("msg", "签名不匹配", "rawString", rawString, "signature", signature, "computedSignature", computedSignature)
		return false
	}

	return true
}

func (api *EventTrackerAPI) TrackEvent(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "EventTrackerAPI.TrackEvent")
	defer span.Finish()
	openSwitch := true
	if openSwitch {
		return
	}

	api.logger.Info("msg", "TrackEvent", "ctx", ctx)

	// 使用 Proto 绑定器解析请求体
	reqParams := &tracker.EventTrackerEncodeRequest{}
	if err := transport.BindProto(c, reqParams); err != nil {
		api.logger.Warn("msg", "解析事件追踪请求失败", "error", err)
		c.Error(err)
		return
	}

	// 获取原始请求体
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		api.logger.Error("msg", "读取请求体失败", "error", err.Error())
		c.Error(err)
		return
	}
	// 重新设置请求体，因为 ReadAll 会消耗掉 body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// 拼接签名字符串，使用原始 JSON 请求体
	rawSignStr := fmt.Sprintf("rd=%d&timestamp=%d\n%s", reqParams.Rd, reqParams.Timestamp, string(rawBody))

	// 验证签名
	if !api.validateSignature(rawSignStr, reqParams.Signature) {
		api.logger.Error("msg", "签名验证失败")
		c.Error(errors.New("sign error"))
		return
	}

	// 解析请求体到 proto 对象，用于后续处理
	reqBody := &tracker.EventTrackerBodyRequest{}
	if err := transport.BindProto(c, reqBody); err != nil {
		api.logger.Error("msg", "解析请求体失败", "error", err.Error())
		c.Error(err)
		return
	}

	// api.logger.Debug("msg", "签名验证成功", "rawSignStr", rawSignStr, "signature", reqParams.Signature, "computedSignature", reqParams.Signature)

	// 返回成功响应
	response.SuccessNoData(c, "success")
}
