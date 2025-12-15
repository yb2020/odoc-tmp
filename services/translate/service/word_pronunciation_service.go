package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc-proto/gen/go/translate"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/translate/dao"
	"github.com/yb2020/odoc/services/translate/model"
)

// 常量定义
const (
	TranslateSourceYoudao = "youdao"
	TranslateSourceBaidu  = "baidu"
)

// FileFormat 文件格式
type FileFormat string

const (
	FormatMpeg FileFormat = "mpeg"
	FormatMp3  FileFormat = "mp3"
)

// WordPronunciationService 单词发音服务
type WordPronunciationService struct {
	wordPronunciationDAO dao.WordPronunciationDAO
	config               *config.Config
	logger               logging.Logger
	httpClient           http_client.HttpClient
	jsEngine             *http_client.JSEngine
	tracer               opentracing.Tracer
}

// NewWordPronunciationService 创建单词发音服务
func NewWordPronunciationService(
	logger logging.Logger,
	config *config.Config,
	tracer opentracing.Tracer,
	wordPronunciationDAO dao.WordPronunciationDAO,
) *WordPronunciationService {
	// 创建JavaScript引擎
	jsEngine := http_client.NewJSEngine(logger)

	// 加载百度翻译sign生成脚本
	// 尝试多个可能的路径来定位脚本文件
	possiblePaths := []string{
		filepath.Join("resources", "js", "baidu_fanyi_get_sign.js"),                   // 相对于工作目录
		filepath.Join("..", "..", "resources", "js", "baidu_fanyi_get_sign.js"),       // 相对于服务目录
		filepath.Join("..", "..", "..", "resources", "js", "baidu_fanyi_get_sign.js"), // 更深层级
	}

	var scriptPath string
	var scriptErr error

	// 尝试每个可能的路径
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			scriptPath = path
			break
		}
	}

	// 如果没有找到脚本文件，记录错误
	if scriptPath == "" {
		logger.Error("找不到百度翻译脚本文件", "paths", possiblePaths)
		scriptErr = fmt.Errorf("找不到百度翻译脚本文件")
	} else {
		// 加载找到的脚本文件
		scriptErr = jsEngine.LoadScriptFromFile(scriptPath)
		if scriptErr != nil {
			logger.Error("加载百度翻译JavaScript脚本失败", "error", scriptErr, "path", scriptPath)
		} else {
			logger.Info("成功加载JavaScript脚本", "path", scriptPath)
		}
	}

	return &WordPronunciationService{
		wordPronunciationDAO: wordPronunciationDAO,
		config:               config,
		logger:               logger,
		httpClient:           http_client.NewHttpClient(logger, http_client.WithTimeout(30*time.Second)),
		jsEngine:             jsEngine,
		tracer:               tracer,
	}
}

// GetYdWordTranslation 从有道词典获取单词翻译
func (s *WordPronunciationService) GetYdWordTranslation(ctx context.Context, content string) (*model.WordPronunciation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WordPronunciationService.GetYdWordTranslation")
	defer span.Finish()

	// 先从数据库查找
	dbWordTranslates, err := s.wordPronunciationDAO.FindByTargetContentAndSource(ctx, content, TranslateSourceYoudao)
	if err != nil {
		s.logger.Error("查询数据库失败", "error", err)
		return nil, err
	}

	if len(dbWordTranslates) > 0 {
		wordTranslate := dbWordTranslates[0]
		if len(wordTranslate.TargetResp) == 0 {
			// 如果没有翻译结果，删除记录
			err = s.wordPronunciationDAO.DeleteById(ctx, wordTranslate.Id)
			if err != nil {
				s.logger.Error("删除无效记录失败", "error", err)
			}
		} else {
			return &wordTranslate, nil
		}
	}

	// 创建新的翻译记录
	wordTranslate := &model.WordPronunciation{
		TargetContent: content,
		AmericaFormat: string(FormatMpeg),
		BritishFormat: string(FormatMpeg),
		Source:        TranslateSourceYoudao,
		TargetResp:    []model.WordExp{},
	}

	// 爬取有道词典
	err = s.crawlYoudaoTranslation(wordTranslate)
	if err != nil {
		s.logger.Error("爬取有道词典失败", "error", err, "content", content)
		return wordTranslate, err
	}
	// 如果有翻译结果，保存到数据库
	if len(wordTranslate.TargetResp) > 0 {
		err = s.wordPronunciationDAO.Save(ctx, wordTranslate)
		if err != nil {
			s.logger.Error("保存翻译结果失败", "error", err)
		}
	}

	return wordTranslate, nil
}

// crawlYoudaoTranslation 爬取有道词典翻译
func (s *WordPronunciationService) crawlYoudaoTranslation(wordTranslate *model.WordPronunciation) error {
	// 构建请求URL
	url := fmt.Sprintf(s.config.Translate.Pronunciation.Crawl.YD.Result.URL, wordTranslate.TargetContent)

	// 使用 http_client 发送 GET 请求
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	}

	responseData, err := s.httpClient.Get(url, headers)
	if err != nil {
		return errors.Wrap(err, "发送请求失败")
	}

	s.logger.Info("爬取有道词典结果", "content", string(responseData))

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(responseData))
	if err != nil {
		return errors.Wrap(err, "解析HTML失败")
	}

	// 提取音标
	doc.Find(".per-phone").Each(func(i int, s *goquery.Selection) {
		if s.Find(":contains('英')").Length() > 0 {
			wordTranslate.BritishSymbol = s.Find(".phonetic").Text()
		} else if s.Find(":contains('美')").Length() > 0 {
			wordTranslate.AmericaSymbol = s.Find(".phonetic").Text()
		}
	})

	// 提取释义
	doc.Find(".simple.dict-module").Each(func(i int, s *goquery.Selection) {
		s.Find(".word-exp").Each(func(j int, sel *goquery.Selection) {
			wordExp := model.WordExp{}

			// 提取词性
			posElements := sel.Find(".pos")
			if posElements.Length() > 0 {
				wordExp.Part = posElements.Text()
			}

			// 提取释义
			transElements := sel.Find(".trans")
			if transElements.Length() > 0 {
				wordExp.TargetContent = []string{transElements.Text()}
			}

			wordTranslate.TargetResp = append(wordTranslate.TargetResp, wordExp)
		})
	})

	// 获取发音
	britishPronounceURL := fmt.Sprintf(s.config.Translate.Pronunciation.Crawl.YD.Pronounce.URL, wordTranslate.TargetContent, 1)
	britishPronounceBase64, err := s.getPronounceBase64(britishPronounceURL)
	if err != nil {
		s.logger.Error("获取英式发音失败", "error", err)
	} else {
		wordTranslate.BritishPronunciation = britishPronounceBase64
	}

	americaPronounceURL := fmt.Sprintf(s.config.Translate.Pronunciation.Crawl.YD.Pronounce.URL, wordTranslate.TargetContent, 2)
	americaPronounceBase64, err := s.getPronounceBase64(americaPronounceURL)
	if err != nil {
		s.logger.Error("获取美式发音失败", "error", err)
	} else {
		wordTranslate.AmericaPronunciation = americaPronounceBase64
	}

	return nil
}

// getPronounceBase64 获取发音的Base64编码
func (s *WordPronunciationService) getPronounceBase64(pronounceURL string) (string, error) {
	// 使用 http_client 发送 GET 请求
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	}

	responseData, err := s.httpClient.Get(pronounceURL, headers)
	if err != nil {
		return "", errors.Wrap(err, "发送请求失败")
	}

	// 转换为Base64
	return base64.StdEncoding.EncodeToString(responseData), nil
}

// GetBdWordTranslation 从百度词典获取单词翻译
func (s *WordPronunciationService) GetBdWordTranslation(ctx context.Context, content string) (*model.WordPronunciation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WordPronunciationService.GetBdWordTranslation")
	defer span.Finish()

	// 先从数据库查找
	dbWordTranslates, err := s.wordPronunciationDAO.FindByTargetContentAndSource(ctx, content, TranslateSourceBaidu)
	if err != nil {
		s.logger.Error("查询数据库失败", "error", err)
		return nil, err
	}

	if len(dbWordTranslates) > 0 {
		wordTranslate := dbWordTranslates[0]
		if len(wordTranslate.TargetResp) == 0 {
			// 如果没有翻译结果，删除记录
			err = s.wordPronunciationDAO.DeleteById(ctx, wordTranslate.Id)
			if err != nil {
				s.logger.Error("删除无效记录失败", "error", err)
			}
		} else {
			return &wordTranslate, nil
		}
	}

	// 创建新的翻译记录
	wordTranslate := &model.WordPronunciation{
		TargetContent: content,
		AmericaFormat: string(FormatMp3),
		BritishFormat: string(FormatMp3),
		Source:        TranslateSourceBaidu,
		TargetResp:    []model.WordExp{},
	}

	// 爬取百度翻译
	err = s.crawlBaiduTranslation(ctx, wordTranslate)
	if err != nil {
		s.logger.Error("爬取百度翻译失败", "error", err)
	}

	// 如果有翻译结果，保存到数据库
	if len(wordTranslate.TargetResp) > 0 {
		err = s.wordPronunciationDAO.Save(ctx, wordTranslate)
		if err != nil {
			s.logger.Error("保存翻译结果失败", "error", err)
		}
	}

	return wordTranslate, nil
}

// crawlBaiduTranslation 爬取百度词典翻译
func (s *WordPronunciationService) crawlBaiduTranslation(ctx context.Context, wordTranslate *model.WordPronunciation) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WordPronunciationService.crawlBaiduTranslation")
	defer span.Finish()

	// 获取请求参数
	formData, headers, err := s.getBdTranslateResultRequest(wordTranslate.TargetContent)
	if err != nil {
		return errors.Wrap(err, "获取请求参数失败")
	}

	// 使用 http_client 发送 POST 请求
	responseData, err := s.httpClient.Post(
		s.config.Translate.Pronunciation.Crawl.BD.URL,
		strings.NewReader(formData),
		headers,
	)
	if err != nil {
		return errors.Wrap(err, "发送请求失败")
	}

	s.logger.Info("爬取百度词典结果", "content", string(responseData))

	// 解析JSON
	var result map[string]interface{}
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		return errors.Wrap(err, "解析JSON失败")
	}

	// 提取翻译结果
	dictResult, ok := result["dict_result"].(map[string]interface{})
	if !ok {
		return errors.System(errors.ErrorTypeBiz, "未找到dict_result字段", nil)
	}

	simpleMeans, ok := dictResult["simple_means"].(map[string]interface{})
	if !ok {
		return errors.System(errors.ErrorTypeBiz, "未找到simple_means字段", nil)
	}

	symbols, ok := simpleMeans["symbols"].([]interface{})
	if !ok || len(symbols) == 0 {
		return errors.System(errors.ErrorTypeBiz, "未找到symbols字段或为空", nil)
	}

	symbolsInfo, ok := symbols[0].(map[string]interface{})
	if !ok {
		return errors.System(errors.ErrorTypeBiz, "symbols字段格式错误", nil)
	}

	// 提取音标
	if phEn, ok := symbolsInfo["ph_en"].(string); ok {
		wordTranslate.BritishSymbol = phEn
	}
	if phAm, ok := symbolsInfo["ph_am"].(string); ok {
		wordTranslate.AmericaSymbol = phAm
	}

	// 提取释义
	parts, ok := symbolsInfo["parts"].([]interface{})
	if ok && len(parts) > 0 {
		for _, part := range parts {
			partObj, ok := part.(map[string]interface{})
			if !ok {
				continue
			}

			wordExp := model.WordExp{}
			if partStr, ok := partObj["part"].(string); ok {
				wordExp.Part = partStr
			}

			means, ok := partObj["means"].([]interface{})
			if !ok || len(means) == 0 {
				continue
			}

			// 将所有释义合并为一个字符串，用分号分隔
			var targetContent []string
			meansList := make([]string, 0, len(means))
			for _, mean := range means {
				if meanStr, ok := mean.(string); ok {
					meansList = append(meansList, meanStr)
				}
			}
			targetContent = append(targetContent, strings.Join(meansList, ";"))
			wordExp.TargetContent = targetContent
			wordTranslate.TargetResp = append(wordTranslate.TargetResp, wordExp)
		}
	}

	// 获取发音
	britishPronounceURL := fmt.Sprintf("https://dict.baidu.com/speech?lan=en&text=%s&type=1", wordTranslate.TargetContent)
	britishPronounceBase64, err := s.getPronounceBase64(britishPronounceURL)
	if err != nil {
		s.logger.Error("获取英式发音失败", "error", err)
	} else {
		wordTranslate.BritishPronunciation = britishPronounceBase64
	}

	americaPronounceURL := fmt.Sprintf("https://dict.baidu.com/speech?lan=en&text=%s&type=2", wordTranslate.TargetContent)
	americaPronounceBase64, err := s.getPronounceBase64(americaPronounceURL)
	if err != nil {
		s.logger.Error("获取美式发音失败", "error", err)
	} else {
		wordTranslate.AmericaPronunciation = americaPronounceBase64
	}

	return nil
}

// GetBdSign 获取百度翻译签名
func (s *WordPronunciationService) GetBdSign(content string) (string, error) {
	result, err := s.jsEngine.CallFunction("e", content)
	if err != nil {
		return "", errors.Wrap(err, "调用JS函数失败")
	}
	return fmt.Sprintf("%v", result), nil
}

// getBdTranslateResultRequest 获取百度翻译请求参数
func (s *WordPronunciationService) getBdTranslateResultRequest(content string) (string, map[string]string, error) {
	sign, err := s.GetBdSign(content)
	if err != nil {
		return "", nil, errors.Wrap(err, "获取签名失败")
	}

	// 构建请求体
	data := fmt.Sprintf("from=en&to=zh&query=%s&simple_means_flag=3&sign=%s&token=832518bea069e0a79f2cb1808aeb6d16&domain=common", content, sign)

	// 准备请求头
	headers := map[string]string{
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		"Accept":           "*/*",
		"Accept-Language":  "zh-CN,zh;q=0.9",
		"Connection":       "keep-alive",
		"Origin":           "https://fanyi.baidu.com",
		"Referer":          "https://fanyi.baidu.com/translate?aldtype=16047&query=strongest&keyfrom=baidu&smartresult=dict&lang=auto2zh",
		"X-Requested-With": "XMLHttpRequest",
	}

	return data, headers, nil
}

// GetTranslateResp 根据渠道获取翻译响应
func (s *WordPronunciationService) GetTranslateResp(ctx context.Context, channel string, content string, pdfId string) (*translate.TranslateResponse, error) {
	var wordTranslate *model.WordPronunciation
	var err error

	// 根据渠道选择翻译来源
	switch channel {
	case TranslateSourceBaidu:
		wordTranslate, err = s.GetBdWordTranslation(ctx, content)
	case TranslateSourceYoudao:
		wordTranslate, err = s.GetYdWordTranslation(ctx, content)
	default:
		return nil, errors.System(errors.ErrorTypeBiz, "not support translate channel", nil)
	}

	if err != nil {
		return nil, err
	}

	if wordTranslate == nil {
		return nil, nil
	}

	// 生成请求ID
	id := idgen.GenerateUUID()

	// 将所有翻译结果合并为一个字符串，用逗号分隔
	var targetJoined string
	for _, exp := range wordTranslate.TargetResp {
		targetJoined += strings.Join(exp.TargetContent, ",") + ","
	}
	// 去掉最后一个逗号
	if len(targetJoined) > 0 {
		targetJoined = targetJoined[:len(targetJoined)-1]
	}

	// 构建响应
	requestIdStr := fmt.Sprintf("%d", id)

	// 转换 model.WordExp 到 translate.TargetResp
	targetRespProto := make([]*translate.TargetResp, 0, len(wordTranslate.TargetResp))
	for _, exp := range wordTranslate.TargetResp {
		targetRespProto = append(targetRespProto, &translate.TargetResp{
			Part:          exp.Part,
			TargetContent: exp.TargetContent,
		})
	}

	resp := &translate.TranslateResponse{
		TargetContent:        []string{wordTranslate.TargetContent},
		RequestId:            &requestIdStr,
		BritishSymbol:        &wordTranslate.BritishSymbol,
		AmericaSymbol:        &wordTranslate.AmericaSymbol,
		BritishFormat:        &wordTranslate.BritishFormat,
		BritishPronunciation: &wordTranslate.BritishPronunciation,
		AmericaFormat:        &wordTranslate.AmericaFormat,
		AmericaPronunciation: &wordTranslate.AmericaPronunciation,
		TargetResp:           targetRespProto,
	}

	return resp, nil
}
