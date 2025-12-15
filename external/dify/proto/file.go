package proto

import "mime/multipart"

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	File      *multipart.FileHeader `json:"-"`
	User      string                `json:"user"`
	MediaType string                `json:"-"` // 不序列化，仅用于内部处理
}

// NewFileUploadRequest 创建新的文件上传请求
func NewFileUploadRequest(user string) *FileUploadRequest {
	return &FileUploadRequest{
		User:      user,
		MediaType: "application/octet-stream", // 默认媒体类型
	}
}

// WithMediaType 设置媒体类型
func (r *FileUploadRequest) WithMediaType(mediaType string) *FileUploadRequest {
	r.MediaType = mediaType
	return r
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
	MimeType  string `json:"mime_type"`
	CreatedBy string `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
}
