package util

import (
	"reflect"
	"strings"
)

// FileType 文件类型定义
type FileType struct {
	Extension   string // 文件扩展名（包含点）
	MimeType    string // MIME类型
	Description string // 文件类型描述
}

// FileTypes 文件类型映射
var FileTypes = struct {
	// 图片类型
	JPEG     FileType
	PNG      FileType
	GIF      FileType
	BMP      FileType
	WEBP     FileType
	SVG      FileType
	ICO      FileType
	TIFF     FileType
	
	// 文档类型
	PDF      FileType
	DOC      FileType
	DOCX     FileType
	XLS      FileType
	XLSX     FileType
	PPT      FileType
	PPTX     FileType
	TXT      FileType
	RTF      FileType
	ODT      FileType
	ODS      FileType
	ODP      FileType
	
	// 音频类型
	MP3      FileType
	WAV      FileType
	OGG      FileType
	AAC      FileType
	M4A      FileType
	
	// 视频类型
	MP4      FileType
	AVI      FileType
	MOV      FileType
	WMV      FileType
	FLV      FileType
	MKV      FileType
	WEBM     FileType
	
	// 压缩文件
	ZIP      FileType
	RAR      FileType
	SEVENZ   FileType
	TAR      FileType
	GZ       FileType
	
	// 代码和标记语言
	HTML     FileType
	CSS      FileType
	JS       FileType
	JSON     FileType
	XML      FileType
	YAML     FileType
	MD       FileType
	
	// 其他类型
	DEFAULT  FileType
}{
	// 图片类型
	JPEG:     FileType{".jpg,.jpeg", "image/jpeg", "JPEG图像"},
	PNG:      FileType{".png", "image/png", "PNG图像"},
	GIF:      FileType{".gif", "image/gif", "GIF图像"},
	BMP:      FileType{".bmp", "image/bmp", "BMP图像"},
	WEBP:     FileType{".webp", "image/webp", "WebP图像"},
	SVG:      FileType{".svg", "image/svg+xml", "SVG矢量图"},
	ICO:      FileType{".ico", "image/x-icon", "ICO图标"},
	TIFF:     FileType{".tif,.tiff", "image/tiff", "TIFF图像"},
	
	// 文档类型
	PDF:      FileType{".pdf", "application/pdf", "PDF文档"},
	DOC:      FileType{".doc", "application/vnd.ms-word", "Word文档"},
	DOCX:     FileType{".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "Word文档(OOXML)"},
	XLS:      FileType{".xls", "application/vnd.ms-excel", "Excel表格"},
	XLSX:     FileType{".xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "Excel表格(OOXML)"},
	PPT:      FileType{".ppt", "application/vnd.ms-powerpoint", "PowerPoint演示文稿"},
	PPTX:     FileType{".pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation", "PowerPoint演示文稿(OOXML)"},
	TXT:      FileType{".txt", "text/plain; charset=utf-8", "文本文件"},
	RTF:      FileType{".rtf", "application/rtf", "富文本格式"},
	ODT:      FileType{".odt", "application/vnd.oasis.opendocument.text", "OpenDocument文本文档"},
	ODS:      FileType{".ods", "application/vnd.oasis.opendocument.spreadsheet", "OpenDocument电子表格"},
	ODP:      FileType{".odp", "application/vnd.oasis.opendocument.presentation", "OpenDocument演示文稿"},
	
	// 音频类型
	MP3:      FileType{".mp3", "audio/mpeg", "MP3音频"},
	WAV:      FileType{".wav", "audio/wav", "WAV音频"},
	OGG:      FileType{".ogg", "audio/ogg", "OGG音频"},
	AAC:      FileType{".aac", "audio/aac", "AAC音频"},
	M4A:      FileType{".m4a", "audio/mp4", "M4A音频"},
	
	// 视频类型
	MP4:      FileType{".mp4", "video/mp4", "MP4视频"},
	AVI:      FileType{".avi", "video/x-msvideo", "AVI视频"},
	MOV:      FileType{".mov", "video/quicktime", "QuickTime视频"},
	WMV:      FileType{".wmv", "video/x-ms-wmv", "Windows Media视频"},
	FLV:      FileType{".flv", "video/x-flv", "Flash视频"},
	MKV:      FileType{".mkv", "video/x-matroska", "Matroska视频"},
	WEBM:     FileType{".webm", "video/webm", "WebM视频"},
	
	// 压缩文件
	ZIP:      FileType{".zip", "application/zip", "ZIP压缩文件"},
	RAR:      FileType{".rar", "application/vnd.rar", "RAR压缩文件"},
	SEVENZ:   FileType{".7z", "application/x-7z-compressed", "7Z压缩文件"},
	TAR:      FileType{".tar", "application/x-tar", "TAR归档文件"},
	GZ:       FileType{".gz", "application/gzip", "GZIP压缩文件"},
	
	// 代码和标记语言
	HTML:     FileType{".html,.htm", "text/html; charset=utf-8", "HTML文档"},
	CSS:      FileType{".css", "text/css; charset=utf-8", "CSS样式表"},
	JS:       FileType{".js", "application/javascript; charset=utf-8", "JavaScript代码"},
	JSON:     FileType{".json", "application/json; charset=utf-8", "JSON数据"},
	XML:      FileType{".xml", "application/xml; charset=utf-8", "XML数据"},
	YAML:     FileType{".yml,.yaml", "application/yaml; charset=utf-8", "YAML数据"},
	MD:       FileType{".md,.markdown", "text/markdown; charset=utf-8", "Markdown文档"},
	
	// 其他类型
	DEFAULT:  FileType{"", "application/octet-stream", "二进制数据"},
}

// ExtensionMimeTypeMap 扩展名到MIME类型的映射
var ExtensionMimeTypeMap = initExtensionMimeTypeMap()

// initExtensionMimeTypeMap 初始化扩展名到MIME类型的映射
func initExtensionMimeTypeMap() map[string]string {
	mapping := make(map[string]string)
	
	// 遍历所有FileTypes，为每个扩展名创建映射
	v := reflect.ValueOf(FileTypes)
	for i := 0; i < v.NumField(); i++ {
		fileType := v.Field(i).Interface().(FileType)
		if fileType.Extension == "" {
			continue
		}
		
		// 处理多个扩展名的情况
		for _, ext := range strings.Split(fileType.Extension, ",") {
			mapping[strings.TrimSpace(ext)] = fileType.MimeType
		}
	}
	
	return mapping
}
