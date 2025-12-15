package service

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/oss"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/oss"
	"github.com/yb2020/odoc/services/oss/dao"
	"github.com/yb2020/odoc/services/oss/model"
	paperService "github.com/yb2020/odoc/services/paper/service"
)

// ObjectKeyGenerator 是一个生成对象键的函数类型
type ObjectKeyGenerator func(uniqueId, filename string) string

// OssServiceInterface OSS服务接口
type OssServiceInterface interface {
	// SetPaperPdfParsedService 设置PaperPdfParsedService
	SetPaperPdfParsedService(paperPdfParsedService *paperService.PaperPdfParsedService) error

	// SetObjectKeyGenerator 设置自定义的对象键生成器
	SetObjectKeyGenerator(generator ObjectKeyGenerator)

	// GetS3UploadToken 获取Minio上传令牌
	GetS3UploadToken(ctx context.Context, bucketType pb.OSSBucketEnum, uniqueId, filename, fileSHA256 string, fileSize int64, metadata map[string]string, topic string, keyPolicy pb.OSSKeyPolicyEnum) (*pb.GetUploadTokenResponse, error)

	// GetDownloadTempUrlByFileSHA256AndType 获取解析后的文件的临时下载地址
	GetDownloadTempUrlByFileSHA256AndType(ctx context.Context, parsedDataType pb.ParsedDataEnum, fileSHA256 string, version string) (string, error)

	// GetS3UploadTokenWithCustomObjectKey 获取Minio上传令牌并指定对象键生成器
	GetS3UploadTokenWithCustomObjectKey(ctx context.Context, bucketType pb.OSSBucketEnum, uniqueId, filename, fileSHA256 string, fileSize int64, metadata map[string]string, topic string, objectKeyGen ObjectKeyGenerator) (*pb.GetUploadTokenResponse, error)

	// UpdateFileStatus 更新文件记录状态
	UpdateFileStatus(ctx context.Context, id string, status string, fileSize int64) error

	// GetRecordByID 根据 ID 获取文件记录
	GetRecordByID(ctx context.Context, id string) (*model.OssRecord, error)

	// GetRecordByBucketNameAndObjectKey 根据对象键获取文件记录
	GetRecordByBucketNameAndObjectKey(ctx context.Context, bucketName, objectKey string) (*model.OssRecord, error)

	// GetDefaultFileTemporaryURL 获取文件的临时访问地址（默认10分钟）
	GetDefaultFileTemporaryURL(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string) (string, error)

	// GetFileTemporaryURL 获取文件的临时访问地址
	GetFileTemporaryURL(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string, expiresTime int) (string, error)

	// GetPermanentURL 获取对象的永久URL（仅适用于公开桶）
	GetPermanentURL(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string) (string, error)

	// UploadObject 直接上传对象到存储服务
	UploadObject(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string, reader io.Reader, objectSize int64, contentType string, metadata map[string]string) error

	// UploadObjectAndSaveRecord 上传对象并保存记录
	UploadObjectAndSaveRecord(ctx context.Context, bucketType pb.OSSBucketEnum, fileName string, fileSHA256 string, reader io.Reader, objectSize int64, objectKeyGen ObjectKeyGenerator, topic string, constomMetadata map[string]string) (*model.OssRecord, error)

	// DownloadObject 从存储服务下载对象
	DownloadObject(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string) (io.ReadCloser, error)

	// DownloadObjectAsBytes 下载对象并返回字节数组
	DownloadObjectAsBytes(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string) ([]byte, error)

	// GetSuccessRecordByFileSHA256 根据文件SHA256获取成功记录
	GetSuccessRecordByFileSHA256(ctx context.Context, fileSHA256 string) (*model.OssRecord, error)

	// SaveOssRecord 保存OSS记录
	SaveOssRecord(ctx context.Context, record *model.OssRecord) error
}

// OssDeps OSS 服务依赖
type OssDeps struct {
	OssDao  *dao.OssDAO
	Config  *config.Config
	Logger  logging.Logger
	Tracer  opentracing.Tracer
	Storage oss.StorageInterface
}
