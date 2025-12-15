package translate

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/logging"
)

// GoogleFreeTranslateClient Google Free翻译客户端
type GoogleFreeTranslateClient struct {
	config config.Config
	client http_client.HttpClient
	logger logging.Logger
}

// GoogleFreeTranslateRequest Google Free翻译请求
type GoogleFreeTranslateRequest struct {
	Token          string `json:"token"`
	Text           string `json:"text"`
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage"`
}

// GoogleFreeTranslateResponse Google Free翻译响应
type GoogleFreeTranslateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Text string `json:"text"`
	} `json:"data"`
}

// NewGoogleFreeTranslateClient 创建Google Free翻译客户端
func NewGoogleFreeTranslateClient(config config.Config, client http_client.HttpClient, logger logging.Logger) *GoogleFreeTranslateClient {
	return &GoogleFreeTranslateClient{
		config: config,
		client: client,
		logger: logger,
	}
}

// Translate 翻译文本
func (c *GoogleFreeTranslateClient) Translate(content, sourceLanguage, targetLanguage string) (string, error) {
	response, err := c.FreeTranslate(content, sourceLanguage, targetLanguage)
	if err != nil {
		c.logger.Error("Google Free翻译失败", "error", err)
		response, err = c.SimpleFreeTranslate(content, sourceLanguage, targetLanguage)
		if err != nil {
			c.logger.Error("Google Simple Free翻译失败", "error", err)
			return "", errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("Google translate failed: %s", err.Error()), err)
		}
	}

	// 检查翻译结果
	if response.Code != 0 {
		return "", errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("Google translate failed: %s", response.Message), nil)
	}
	if response.Data.Text == "" {
		return "", errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("Google translate failed: %s", response.Message), nil)
	}

	// 转换为统一的翻译结果格式
	result := map[string]interface{}{
		"errorCode":   0,
		"translation": []string{response.Data.Text},
	}

	// 序列化结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		c.logger.Error("序列化翻译结果失败", "error", err)
		return "", errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("serialize result failed: %s", err.Error()), err)
	}

	return string(resultJSON), nil
}

// FreeTranslate 翻译文本
func (c *GoogleFreeTranslateClient) FreeTranslate(content, sourceLanguage, targetLanguage string) (*GoogleFreeTranslateResponse, error) {
	// 构建请求URL
	apiURL := fmt.Sprintf("%s%s", c.config.Translate.Text.Channel.Google.BaseURL, c.config.Translate.Text.Channel.Google.FreeUri)

	// 构建请求体
	request := GoogleFreeTranslateRequest{
		Token:          c.config.Translate.Text.Channel.Google.Token,
		Text:           content,
		SourceLanguage: sourceLanguage,
		TargetLanguage: targetLanguage,
	}

	// 构建请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 发送请求
	c.logger.Info("发送Google Free翻译请求", "content", content, "from", sourceLanguage, "to", targetLanguage)
	responseData, err := c.client.PostWithTimeout(apiURL, request, headers, time.Duration(5*time.Second))
	if err != nil {
		c.logger.Error("发送Google Free翻译请求失败", "error", err)
		return nil, errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("发送Google Free翻译请求失败: %s", err.Error()), err)
	}

	if len(responseData) == 0 {
		c.logger.Error("接收到空的Google Free翻译响应")
		return nil, errors.System(errors.ErrorTypeExternalInterface, "接收到空的Google Free翻译响应", nil)
	}

	c.logger.Info("接收到Google Free翻译原始结果：", "responseData", string(responseData))

	// 解析响应
	var response GoogleFreeTranslateResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		c.logger.Error("解析Google Free翻译响应失败", "error", err)
		return nil, errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("解析Google Free翻译响应失败: %s", err.Error()), err)
	}

	// 详细调试结构体解析结果
	c.logger.Info("结构体解析结果", "code", response.Code, "message", response.Message, "data_text_length", len(response.Data.Text))
	if response.Data.Text != "" {
		c.logger.Info("结构体解析的text内容", "text", response.Data.Text)
	} else {
		c.logger.Error("结构体解析后text字段为空")
	}

	// 检查响应状态码
	if response.Code != 0 {
		c.logger.Error("Google Free翻译返回错误", "code", response.Code, "message", response.Message)
		return nil, errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("Google Free翻译返回错误: %s", response.Message), nil)
	}

	c.logger.Info("发送Google Free翻译结果：", "response", response)
	return &response, nil
}

// SimpleFreeTranslate 简单翻译文本
func (c *GoogleFreeTranslateClient) SimpleFreeTranslate(content, sourceLanguage, targetLanguage string) (*GoogleFreeTranslateResponse, error) {
	// 构建请求URL
	apiURL := fmt.Sprintf("%s%s", c.config.Translate.Text.Channel.Google.BaseURL, c.config.Translate.Text.Channel.Google.SimpleFreeUri)

	// 构建请求体
	request := GoogleFreeTranslateRequest{
		Token:          c.config.Translate.Text.Channel.Google.Token,
		Text:           content,
		SourceLanguage: sourceLanguage,
		TargetLanguage: targetLanguage,
	}

	// 构建请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 发送请求
	c.logger.Info("发送Google Simple Free翻译请求", "content", content, "from", sourceLanguage, "to", targetLanguage)
	responseData, err := c.client.PostWithTimeout(apiURL, request, headers, time.Duration(5*time.Second))
	if err != nil {
		c.logger.Error("发送Google Simple Free翻译请求失败", "error", err)
		return nil, errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("发送Google Simple Free翻译请求失败: %s", err.Error()), nil)
	}

	if len(responseData) == 0 {
		c.logger.Error("接收到空的Google Simple Free翻译响应")
		return nil, errors.System(errors.ErrorTypeExternalInterface, "接收到空的Google Simple Free翻译响应", nil)
	}

	c.logger.Info("接收到Google Simple Free翻译原始结果：", "responseData", string(responseData))

	// 先用通用方式解析JSON来调试结构
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(responseData, &rawResponse); err != nil {
		c.logger.Error("解析原始JSON失败", "error", err)
		return nil, errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("解析原始JSON失败: %s", err.Error()), err)
	}

	// 详细调试原始解析结果
	c.logger.Info("原始JSON解析成功", "keys", fmt.Sprintf("%v", getMapKeys(rawResponse)))
	if data, ok := rawResponse["data"].(map[string]interface{}); ok {
		if text, exists := data["text"]; exists {
			textStr := fmt.Sprintf("%v", text)
			c.logger.Info("原始解析中的text字段", "text", textStr, "length", len(textStr), "type", fmt.Sprintf("%T", text))
		} else {
			c.logger.Error("原始解析中没有找到text字段")
		}
	} else {
		c.logger.Error("原始解析中data字段类型不正确", "data_type", fmt.Sprintf("%T", rawResponse["data"]))
	}

	// 解析响应
	var response GoogleFreeTranslateResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		c.logger.Error("解析Google Simple Free翻译响应失败", "error", err)
		return nil, errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("解析Google Simple Free翻译响应失败: %s", err.Error()), err)
	}

	// 详细调试结构体解析结果
	c.logger.Info("结构体解析结果", "code", response.Code, "message", response.Message, "data_text_length", len(response.Data.Text))
	if response.Data.Text != "" {
		c.logger.Info("结构体解析的text内容", "text", response.Data.Text)
	} else {
		c.logger.Error("结构体解析后text字段为空")
	}

	// 检查响应状态码
	if response.Code != 0 {
		c.logger.Error("Google Simple Free翻译返回错误", "code", response.Code, "message", response.Message)
		return nil, errors.System(errors.ErrorTypeExternalInterface, fmt.Sprintf("Google Simple Free翻译返回错误: %s", response.Message), nil)
	}

	c.logger.Info("发送Google Simple Free翻译结果：", "response", response)
	// 直接返回翻译结果
	return &response, nil
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
