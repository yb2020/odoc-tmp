package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

// 处理图片
func HandleImage(base64String string) (*multipart.FileHeader, error) {
	// 解码Base64数据
	imgData, err := HandleImageToBytes(base64String)
	if err != nil {
		return nil, fmt.Errorf("解码Base64失败: %w", err)
	}

	// 确定文件扩展名和MIME类型
	contentType := http.DetectContentType(imgData)
	var fileExt string
	switch contentType {
	case "image/jpeg":
		fileExt = ".jpg"
	case "image/png":
		fileExt = ".png"
	case "image/gif":
		fileExt = ".gif"
	case "image/webp":
		fileExt = ".webp"
	default:
		fileExt = ".bin" // 默认二进制文件扩展名
	}

	// 创建一个表单请求，用于获取 multipart.FileHeader
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 创建表单文件字段
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"file", "image"+fileExt))
	h.Set("Content-Type", contentType)

	// 创建表单文件部分
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("创建表单文件部分失败: %w", err)
	}

	// 写入图片数据
	if _, err := part.Write(imgData); err != nil {
		return nil, fmt.Errorf("写入图片数据失败: %w", err)
	}

	// 关闭writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭writer失败: %w", err)
	}

	// 创建一个HTTP请求
	req, err := http.NewRequest("POST", "/dummy", body)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 解析表单
	if err := req.ParseMultipartForm(32 << 20); err != nil { // 32MB 最大内存
		return nil, fmt.Errorf("解析表单失败: %w", err)
	}

	// 获取文件
	_, fileHeader, err := req.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("获取表单文件失败: %w", err)
	}

	return fileHeader, nil
}

func HandleImageToBytes(base64String string) ([]byte, error) {
	// 检查并移除 Base64 前缀
	if idx := strings.Index(base64String, ";base64,"); idx > 0 {
		// 移除前缀，只保留Base64编码部分
		base64String = base64String[idx+8:]
	}
	// 解码Base64数据
	return base64.StdEncoding.DecodeString(base64String)
}
