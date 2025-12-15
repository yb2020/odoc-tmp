package dto

import "encoding/json"

// FullTranslateUploadResult 全文翻译上传结果
type FullTranslateUploadResult struct {
	Token    string                        `json:"token"`
	Progress int                           `json:"progress"`
	Code     int                           `json:"code"`
	Message  string                        `json:"message"`
	Data     FullTranslateUploadResultData `json:"data"`
	Version  string                        `json:"version"`
	Status   int                           `json:"status"`
}

// FullTranslateUploadResultData 全文翻译上传结果数据
type FullTranslateUploadResultData struct {
	Status         string  `json:"status"`
	ProcessingTime float64 `json:"processing_time"`
	TaskId         string  `json:"task_id"`
}

// FullTranslateProgressResult 全文翻译进度结果
type FullTranslateProgressResult struct {
	Token     string                          `json:"token"`
	Progress  int                             `json:"progress"`
	Message   string                          `json:"message"`
	Alignment json.RawMessage                 `json:"alignment,omitempty"` // 保留原始 JSON 数据
	FixId     int64                           `json:"fixId"`
	Data      FullTranslateProgressResultData `json:"data"`
}

// FullTranslateProgressResultData 全文翻译进度结果数据
type FullTranslateProgressResultData struct {
	ObjectKey  string `json:"objectKey"`
	BucketName string `json:"bucketName"`
}

// CheckFileLangResult 检查文件语言结果
type CheckFileLangResult struct {
	Language string `json:"language"`
}
