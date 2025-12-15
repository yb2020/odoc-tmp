package translate

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/logging"
)

// YoudaoTranslateClient 有道翻译客户端
type YoudaoTranslateClient struct {
	config config.Config
	client http_client.HttpClient
	logger logging.Logger
}

// YoudaoTranslateResponse 有道翻译响应
type YoudaoTranslateResponse struct {
	ErrorCode   string      `json:"errorCode"`
	Query       string      `json:"query"`
	Translation []string    `json:"translation"`
	Basic       BasicInfo   `json:"basic,omitempty"`
	Web         []WebInfo   `json:"web,omitempty"`
	Dict        DictInfo    `json:"dict,omitempty"`
	Webdict     WebdictInfo `json:"webdict,omitempty"`
	L           string      `json:"l"`
	TSpeakURL   string      `json:"tSpeakUrl,omitempty"`
	SpeakURL    string      `json:"speakUrl,omitempty"`
}

// BasicInfo 基本翻译信息
type BasicInfo struct {
	Phonetic string   `json:"phonetic,omitempty"`
	UK       string   `json:"uk-phonetic,omitempty"`
	US       string   `json:"us-phonetic,omitempty"`
	UKSpeech string   `json:"uk-speech,omitempty"`
	USSpeech string   `json:"us-speech,omitempty"`
	Explains []string `json:"explains,omitempty"`
	Wfs      []Wf     `json:"wfs,omitempty"`
}

// Wf 词形变化
type Wf struct {
	Wf Wfs `json:"wf"`
}

// Wfs 词形
type Wfs struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// WebInfo 网络释义
type WebInfo struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

// DictInfo 词典信息
type DictInfo struct {
	URL string `json:"url"`
}

// WebdictInfo 网络词典信息
type WebdictInfo struct {
	URL string `json:"url"`
}

// NewYoudaoTranslateClient 创建有道翻译客户端
func NewYoudaoTranslateClient(config config.Config, client http_client.HttpClient, logger logging.Logger) *YoudaoTranslateClient {
	return &YoudaoTranslateClient{
		config: config,
		client: client,
		logger: logger,
	}
}

// Translate 翻译文本
func (c *YoudaoTranslateClient) Translate(content, sourceLanguage, targetLanguage string) (string, error) {
	// 如果源语言或目标语言为空，设置默认值, 先固定为en和zh-CHS
	// if sourceLanguage == "" {
	sourceLanguage = "en"
	// }
	// if targetLanguage == "" {
	targetLanguage = "zh-CHS"
	// }

	paramsMap := createRequestParams(sourceLanguage, targetLanguage, content, "")

	AddAuthParams(c.config.Translate.Text.Channel.Youdao.AppId, c.config.Translate.Text.Channel.Youdao.AppSecretKey, paramsMap)

	// 构建请求头
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// 将 map[string][]string 转换为 map[string]string
	formData := make(map[string]string)
	for key, values := range paramsMap {
		if len(values) > 0 {
			formData[key] = values[0]
		}
	}

	// 发送请求
	c.logger.Info("发送有道翻译请求", "content", content, "from", sourceLanguage, "to", targetLanguage)
	responseData, err := c.client.PostForm(c.config.Translate.Text.Channel.Youdao.ApiUrl, formData, headers)
	if err != nil {
		c.logger.Error("发送有道翻译请求失败", "error", err)
		return "", fmt.Errorf("发送有道翻译请求失败: %w", err)
	}

	// 解析响应
	var response YoudaoTranslateResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		c.logger.Error("解析有道翻译响应失败", "error", err)
		return "", fmt.Errorf("解析有道翻译响应失败: %w", err)
	}

	// 检查错误码
	if response.ErrorCode != "0" {
		c.logger.Error("有道翻译返回错误", "errorCode", response.ErrorCode)
		return "", fmt.Errorf("有道翻译返回错误: %s", response.ErrorCode)
	}

	// 转换为统一的翻译结果格式
	result := map[string]interface{}{
		"errorCode":   0,
		"translation": response.Translation,
	}

	// 如果有网络释义，添加到结果中
	if len(response.Web) > 0 {
		result["web"] = response.Web
	}

	// 序列化结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		c.logger.Error("序列化翻译结果失败", "error", err)
		return "", fmt.Errorf("序列化翻译结果失败: %w", err)
	}

	return string(resultJSON), nil
}

func createRequestParams(sourceLanguage, targetLanguage, content, vocabId string) map[string][]string {

	/*
		note: 将下列变量替换为需要请求的参数
		取值参考文档: https://ai.youdao.com/DOCSIRMA/html/%E8%87%AA%E7%84%B6%E8%AF%AD%E8%A8%80%E7%BF%BB%E8%AF%91/API%E6%96%87%E6%A1%A3/%E6%96%87%E6%9C%AC%E7%BF%BB%E8%AF%91%E6%9C%8D%E5%8A%A1/%E6%96%87%E6%9C%AC%E7%BF%BB%E8%AF%91%E6%9C%8D%E5%8A%A1-API%E6%96%87%E6%A1%A3.html
	*/
	return map[string][]string{
		"q":       {content},
		"from":    {sourceLanguage},
		"to":      {targetLanguage},
		"vocabId": {vocabId},
	}
}

// AddAuthParams 添加鉴权相关参数 -
// appKey : 应用ID
// salt : 随机值
// curtime : 当前时间戳(秒)
// signType : 签名版本
// sign : 请求签名
// @param appKey    您的应用ID
// @param appSecret 您的应用密钥
// @param paramsMap 请求参数表
func AddAuthParams(appKey string, appSecret string, params map[string][]string) {
	qs := params["q"]
	if qs == nil {
		qs = params["img"]
	}
	var q string
	for i := range qs {
		q += qs[i]
	}
	salt := getUuid()
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	sign := CalculateSign(appKey, appSecret, q, salt, curtime)
	params["appKey"] = []string{appKey}
	params["salt"] = []string{salt}
	params["curtime"] = []string{curtime}
	params["signType"] = []string{"v3"}
	params["sign"] = []string{sign}
}

func AddAuthParamsWithQ(appKey string, appSecret string, q string) map[string]interface{} {
	salt := getUuid()
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	sign := CalculateSign(appKey, appSecret, q, salt, curtime)
	return map[string]interface{}{
		"appKey":   appKey,
		"salt":     salt,
		"curtime":  curtime,
		"signType": "v3",
		"sign":     sign,
	}
}

// CalculateSign 计算v3鉴权签名 -
// 计算方式 : sign = sha256(appKey + input(q) + salt + curtime + appSecret)

// @param appKey    您的应用ID
// @param appSecret 您的应用密钥
// @param q         请求内容
// @param salt      随机值
// @param curtime   当前时间戳(秒)
// @return 鉴权签名sign
func CalculateSign(appKey string, appSecret string, q string, salt string, curtime string) string {
	strSrc := appKey + getInput(q) + salt + curtime + appSecret
	return encrypt(strSrc)
}

func encrypt(strSrc string) string {
	bt := []byte(strSrc)
	bts := sha256.Sum256(bt)
	return hex.EncodeToString(bts[:])
}

func getInput(q string) string {
	str := []rune(q)
	strLen := len(str)
	if strLen <= 20 {
		return q
	} else {
		return string(str[:10]) + strconv.Itoa(strLen) + string(str[strLen-10:])
	}
}

func getUuid() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
