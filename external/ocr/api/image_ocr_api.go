package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"strings"
	"time"

	parsepb "github.com/yb2020/odoc/proto/gen/go/parsed"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/http_client"
	idgen "github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
)

// ImageOCRApiService OCR服务
type ImageOCRApiService struct {
	client         http_client.HttpClient
	logger         logging.Logger
	extractTextURL string
	config         *config.Config
	httpClient     http_client.HttpClient
}

// NewImageOCRApiService 创建OCR服务
func NewImageOCRApiService(
	logger logging.Logger,
	extractTextURL string,
	config *config.Config,
	httpClient http_client.HttpClient,
) *ImageOCRApiService {
	return &ImageOCRApiService{
		client:         http_client.NewHttpClient(logger),
		extractTextURL: extractTextURL,
		logger:         logger,
		config:         config,
		httpClient:     httpClient,
	}
}

// UploadImageResponse 上传img响应结构
// 实际返回的是一个文件路径的数组
type UploadImageResponse []string

type TempDataResult struct {
	Msg     string `json:"msg"`
	EventId string `json:"event_id"`
	Output  struct {
		Data []json.RawMessage `json:"data"`
	} `json:"output"`
	Success bool   `json:"success"`
	Title   string `json:"title"`
}

// ExtractTextFromBase64 从base64编码的图片中提取文本
func (s *ImageOCRApiService) ExtractTextFromBase64(ctx context.Context, base64ImageData string) (string, error) {
	// 创建请求体
	requestBody := map[string]interface{}{
		"images": []string{base64ImageData},
	}
	// 序列/ 使用 http_client 发送请求
	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	responseBody, err := s.client.Post(s.extractTextURL, requestBody, headers)
	if err != nil {
		return "", err
	}

	// 解析响应
	var response struct {
		Results [][]struct {
			Text string `json:"text"`
		} `json:"results"`
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}

	// 提取文本
	var result string
	if len(response.Results) > 0 && len(response.Results[0]) > 0 {
		for _, item := range response.Results[0] {
			result += item.Text
		}
	}

	return result, nil
}

func (s *ImageOCRApiService) ExtractImgToText(ctx context.Context, imgBytes []byte) (string, error) {
	// 生成随机上传ID
	uploadId := idgen.GenerateUUID()
	// 1. 上传PDF   返回一个文件路径
	tempFilePath, err := s.uploadImage(ctx, imgBytes, idgen.GenerateUUID()+".png", uploadId)
	if err != nil {
		return "", err
	}

	// 2. 循环检查进度直到完成
	isCompleted := false
	maxRetries := 60 // 最大重试次数，防止无限循环
	retryCount := 0
	var uploadProgressResponse *parsepb.UploadProgressResponse
	for !isCompleted && retryCount < maxRetries {
		isCompleted, uploadProgressResponse, err = s.uploadProgress(ctx, uploadId)
		if err != nil {
			retryCount++
			time.Sleep(5 * time.Second)
			continue
		}
		if !isCompleted {
			// 等待一段时间再检查
			time.Sleep(1 * time.Second)
			retryCount++
		}
	}

	if !isCompleted {
		s.logger.Error("PDF parse timeout or failed")
		return "", errors.Biz("PDF parse timeout or failed")
	}
	//
	fileName := uploadProgressResponse.OrigName
	//文件大小
	fileSize := uploadProgressResponse.ChunkSize
	// 生成会话哈希值
	sessionHash := idgen.GenerateUUID()
	// 3. 加入解析队列
	resultID, err := s.join(ctx, tempFilePath, uploadId, sessionHash, fileName, int64(fileSize))
	if err != nil {
		return "", err
	}
	// 4. 获取数据
	resultData, err := s.getData(ctx, resultID, sessionHash)
	if err != nil {
		return "", err
	}
	// 5. 解析resultData   这个data中的第二条数据是我们想要的数据
	output := resultData.Output
	if output == nil {
		return "", errors.Biz("output is nil")
	}
	//这里正常第二条数据就是我们想要的结果，但是我们不确定是否会随着版本进行变化，所以这里判断字符串中是否包含哪些属性作为条件判断
	for _, rawData := range output.Data {
		// 保留您现有的判断条件
		if strings.Contains(rawData, "\"data\"") && strings.Contains(rawData, "\"headers\"") {
			// 1. 定义一个临时的结构体来接收解析后的数据
			var ocrBlock struct {
				Data [][]interface{} `json:"data"`
				// Headers 字段我们在这里不需要，所以可以不定义
			}
			// 2. 将 rawData (它是一个JSON字符串) 反序列化到我们的临时结构体中
			err := json.Unmarshal([]byte(rawData), &ocrBlock)
			if err != nil {
				// 如果解析失败，说明这个字符串虽然包含了 "data" 和 "headers"
				// 但它不是一个合法的JSON，或者结构不匹配。记录日志并继续循环。
				// s.logger.Warn("Found a matching raw string, but failed to unmarshal it", "error", err, "rawData", rawData)
				continue
			}
			// 3. 提取并拼接文本
			var textBuilder strings.Builder
			for i, row := range ocrBlock.Data {
				// 每个 row 是一个子数组, e.g., [0, "Batch normalization...", 0.92436]
				// 确保 row 至少有两个元素，并且第二个元素是字符串
				if len(row) >= 2 {
					if text, ok := row[1].(string); ok {
						textBuilder.WriteString(text)
						// (可选) 添加换行符
						if i < len(ocrBlock.Data)-1 {
							textBuilder.WriteString("\n")
						}
					}
				}
			}
			finalText := textBuilder.String()
			if finalText == "" {
				// 虽然找到了数据块，但没有提取出任何文本，继续寻找下一个可能的数据块
				// s.logger.Warn("Found OCR data block, but no text could be extracted.", "rawData", rawData)
				continue
			}
			// 4. 成功提取并拼接了文本，直接返回结果
			return finalText, nil
		}
	}

	return "", errors.Biz("parse output data failed , parse result is empty")
}

func (s *ImageOCRApiService) uploadImage(ctx context.Context, imgContent []byte, fileName string, uploadId string) (string, error) {
	// 从配置文件获取URL
	uploadConfigURL := s.config.Translate.OCR.ExtractTextURL + s.config.Translate.OCR.UploadImageURL
	// 构建上传URL
	uploadURL := fmt.Sprintf("%s?upload_id=%s", uploadConfigURL, uploadId)

	// 设置请求头
	headers := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "multipart/form-data",
	}
	// 设置文件内容，注意表单字段名为'files'
	fileContents := map[string][]byte{
		"files": imgContent,
	}
	// 设置文件名
	fileNames := map[string]string{
		"files": fileName,
	}
	// 发送上传请求
	timeout := time.Duration(s.config.Translate.OCR.Timeout) * time.Minute
	responseData, err := s.httpClient.PostMultipartFormWithFileInput(uploadURL, nil, fileContents, fileNames, headers, timeout)
	if err != nil {
		s.logger.Error("Mineru upload failed", "error", err)
		return "", err
	}
	// 解析响应获取文件路径
	var response UploadImageResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		s.logger.Error("parse upload response failed", "error", err, "responseData", string(responseData))
		return "", errors.BizWrap("parse upload response failed", err)
	}
	// 检查响应是否为空
	if len(response) == 0 {
		s.logger.Error("upload pdf failed, response is empty array")
		return "", errors.Biz("upload pdf failed, response is empty array")
	}
	// 返回第一个文件路径作为最终的文件路径
	filePath := response[0]
	return filePath, nil
}

// 调用mineru接口  检查上传进度
// 返回值为true表示上传完成
func (s *ImageOCRApiService) uploadProgress(ctx context.Context, uploadId string) (bool, *parsepb.UploadProgressResponse, error) {

	OcrProgressURL := s.config.Translate.OCR.ExtractTextURL + s.config.Translate.OCR.ProgressURL
	progressURL := fmt.Sprintf("%s?upload_id=%s", OcrProgressURL, uploadId)
	// 设置请求头
	headers := map[string]string{
		"Accept": "text/event-stream", // 事件流格式
	}
	// 发送进度检查请求
	timeout := time.Duration(s.config.Translate.OCR.Timeout) * time.Minute
	responseData, err := s.httpClient.GetWithTimeout(progressURL, headers, timeout)
	if err != nil {
		s.logger.Error("check parse progress failed", "error", err, "uploadId", uploadId)
		return false, nil, errors.BizWrap("check parse progress failed", err)
	}
	// 响应是事件流格式，需要按行解析
	lines := strings.Split(string(responseData), "\n")
	// 默认不完成
	isCompleted := false
	var response parsepb.UploadProgressResponse
	// 处理每一行数据
	for _, line := range lines {
		// 跳过空行
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 去除前缀"data: "
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		} else {
			continue // 不是数据行，跳过
		}
		// 解析JSON   这里的响应结构和mineru中的一致，所以这里直接使用minerU的协议就好
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			continue // 跳过解析失败的行
		}
		// 检查是否完成
		if response.Msg == "done" {
			isCompleted = true
			break
		}
	}
	return isCompleted, &response, nil
}

func (s *ImageOCRApiService) join(ctx context.Context, filePath string, uploadId string, sessionHash string, fileName string, fileSize int64) (string, error) {

	// 获取配置的joinURL
	joinURL := s.config.Translate.OCR.ExtractTextURL + s.config.Translate.OCR.JoinURL

	// 构造文件URL
	fileURL := fmt.Sprintf("%s%s", s.config.Translate.OCR.ExtractTextURL+s.config.Translate.OCR.FileURL, filePath)

	//使用模型数组
	models := []string{
		"use_det",
		"use_cls",
		"use_rec",
	}
	// 根据API要求构造请求体
	requestBody := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"path":      filePath,
				"url":       fileURL,
				"orig_name": fileName,
				"size":      fileSize,
				"mime_type": "image/png",
				"meta": map[string]string{
					"_type": "gradio.FileData",
				},
			},
			s.config.Translate.OCR.TextScore,
			s.config.Translate.OCR.BoxThresh,
			s.config.Translate.OCR.UnclipRatio,
			s.config.Translate.OCR.MaxSideLen,
			s.config.Translate.OCR.UseDet.EngineType,
			s.config.Translate.OCR.UseDet.LangDet,
			s.config.Translate.OCR.UseDet.ModelType,
			s.config.Translate.OCR.UseDet.Version,
			s.config.Translate.OCR.UseCls.EngineType,
			s.config.Translate.OCR.UseCls.LangCls,
			s.config.Translate.OCR.UseCls.ModelType,
			s.config.Translate.OCR.UseCls.Version,
			s.config.Translate.OCR.UseRec.EngineType,
			s.config.Translate.OCR.UseRec.LangRec,
			s.config.Translate.OCR.UseRec.ModelType,
			s.config.Translate.OCR.UseRec.Version,
			s.config.Translate.OCR.ReturnWordBox,
			models,
		},
		"event_data":   nil,
		"fn_index":     0,
		"trigger_id":   37,
		"session_hash": sessionHash,
	}

	// 设置请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 发送请求
	timeout := time.Duration(s.config.Translate.OCR.Timeout) * time.Minute
	responseData, err := s.httpClient.PostWithTimeout(joinURL, requestBody, headers, timeout)
	if err != nil {
		s.logger.Error("merge parse result failed", "error", err, "filePath", filePath)
		return "", errors.BizWrap("merge parse result failed", err)
	}

	// 解析响应获取事件ID
	var response parsepb.JoinResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		s.logger.Error("parse merge response failed", "error", err, "responseData", string(responseData))
		return "", errors.BizWrap("parse merge response failed", err)
	}

	// 检查事件ID是否存在
	if response.EventId == "" {
		s.logger.Error("merge parse result failed", "response", string(responseData))
		return "", errors.Biz("merge parse result failed: no event id")
	}
	return response.EventId, nil
}

// GetData 获取解析数据
func (s *ImageOCRApiService) getData(ctx context.Context, eventID string, sessionHash string) (*parsepb.ProcessCompletedEventMessage, error) {
	// 获取配置的获取数据URL
	DataURL := s.config.Translate.OCR.ExtractTextURL + s.config.Translate.OCR.DataURL

	// 使用固定的获取数据URL，并添加session_hash参数
	dataURL := fmt.Sprintf("%s?session_hash=%s", DataURL, sessionHash)

	// 设置请求头
	headers := map[string]string{
		"Accept": "text/event-stream", // 事件流格式
	}

	// 发送获取数据请求
	timeout := time.Duration(s.config.Translate.OCR.Timeout) * time.Minute
	responseData, err := s.httpClient.GetWithTimeout(dataURL, headers, timeout)
	if err != nil {
		s.logger.Error("get parse data failed", "error", err, "eventID", eventID)
		return nil, errors.BizWrap("get parse data failed", err)
	}

	// 响应是事件流格式，需要按行解析
	lines := strings.Split(string(responseData), "\n")

	// 存储process_completed事件消息
	var result *parsepb.ProcessCompletedEventMessage
	// 处理每一行数据
	for _, line := range lines {
		// 跳过空行
		if line == "" {
			continue
		}
		// 去除前缀"data: "
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		} else {
			continue // 不是数据行，跳过
		}
		// 先解析消息类型
		var msgObj struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal([]byte(line), &msgObj); err != nil {
			continue // 跳过解析失败的行
		}
		// 这里的消息类型有四种 estimation,process_starts,process_completed,close_stream
		// 在处理完 process_completed 消息后
		if msgObj.Msg == "process_completed" {
			// 先创建一个临时结构体来解析复杂的 JSON 结构
			var tempResult TempDataResult
			if err := json.Unmarshal([]byte(line), &tempResult); err != nil {
				s.logger.Error("解析 process_completed 消息失败", "error", err)
				continue
			}
			// 创建最终的 ProcessCompletedEventMessage 结构
			result = &parsepb.ProcessCompletedEventMessage{
				Msg:     tempResult.Msg,
				EventId: tempResult.EventId,
				Success: tempResult.Success,
				Title:   tempResult.Title,
			}

			// 处理 Output 字段
			if len(tempResult.Output.Data) > 0 {
				// 创建字符串数组
				stringData := make([]string, 0, len(tempResult.Output.Data))
				// 将每个 json.RawMessage 转换为字符串
				for _, rawMsg := range tempResult.Output.Data {
					// 将 json.RawMessage 转换为字符串
					stringData = append(stringData, string(rawMsg))
				}
				// 设置到 result 中
				result.Output = &parsepb.EventOutput{
					Data: stringData,
				}
			}
		}
		// 检查是否是完成消息
		if msgObj.Msg == "close_stream" {
			break
		}
	}
	return result, nil
}

func (s *ImageOCRApiService) imgToBytes(img image.Image) ([]byte, error) {
	buffer := bytes.Buffer{}
	if err := png.Encode(&buffer, img); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
