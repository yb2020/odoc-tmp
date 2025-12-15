package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/oss"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	basemodel "github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/pkg/oss"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services/oss/constant"
	"github.com/yb2020/odoc/services/oss/dao"
	"github.com/yb2020/odoc/services/oss/model"
	"github.com/yb2020/odoc/services/oss/util"
	paperService "github.com/yb2020/odoc/services/paper/service"
	parseConstant "github.com/yb2020/odoc/services/parse/constant"
)

// LocalOssService 本地OSS服务实现
type LocalOssService struct {
	ossDao  *dao.OssDAO
	config  *config.Config
	logger  logging.Logger
	tracer  opentracing.Tracer
	storage oss.StorageInterface
	// 默认的对象键生成器
	defaultObjectKeyGenerator ObjectKeyGenerator
	paperPdfParsedService     *paperService.PaperPdfParsedService
}

func (s *LocalOssService) SetPaperPdfParsedService(paperPdfParsedService *paperService.PaperPdfParsedService) error {
	if paperPdfParsedService == nil {
		return errors.Biz("paperPdfParsedService cannot be nil")
	}
	s.paperPdfParsedService = paperPdfParsedService
	return nil
}

// NewLocalOssService 创建本地OSS服务
func NewLocalOssService(ossDao *dao.OssDAO, config *config.Config, logger logging.Logger, tracer opentracing.Tracer, storage oss.StorageInterface) *LocalOssService {
	service := &LocalOssService{
		ossDao:  ossDao,
		config:  config,
		logger:  logger,
		tracer:  tracer,
		storage: storage,
		defaultObjectKeyGenerator: func(uniqueId, filename string) string {
			ext := filepath.Ext(filename)
			uuidFileName := idgen.GenerateUUID() + ext
			now := time.Now()
			datePrefix := now.Format("2006-01-02")
			return fmt.Sprintf("%s/%s/%s", datePrefix, uniqueId, uuidFileName)
		},
	}
	return service
}

// SetObjectKeyGenerator 设置自定义的对象键生成器
func (s *LocalOssService) SetObjectKeyGenerator(generator ObjectKeyGenerator) {
	if generator != nil {
		s.defaultObjectKeyGenerator = generator
	}
}

// GetS3UploadToken 获取Minio上传令牌
func (s *LocalOssService) GetS3UploadToken(ctx context.Context, bucketType pb.OSSBucketEnum, uniqueId, filename, fileSHA256 string, fileSize int64, metadata map[string]string, topic string, keyPolicy pb.OSSKeyPolicyEnum) (*pb.GetUploadTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.GetS3UploadToken")
	defer span.Finish()

	//根据keyPolicy设置key生成策略
	var objectKeyGen ObjectKeyGenerator
	switch keyPolicy {
	case pb.OSSKeyPolicyEnum_UPLOAD_PUBLIC_SHORT:
		objectKeyGen = func(uniqueId, fileName string) string {
			base62ID, _ := idgen.EncodeStringUUIDToBase62(idgen.GenerateUUID())
			shortDir := idgen.GetShardDirectory(base62ID)
			ext := filepath.Ext(filename)
			shortFileName := base62ID + ext
			return fmt.Sprintf("%s/%s", shortDir, shortFileName)
		}
	case pb.OSSKeyPolicyEnum_UPLOAD_PDF:
		objectKeyGen = func(uniqueId, fileName string) string {
			return fmt.Sprintf("%s/%s/%s", fileSHA256, parseConstant.SourcePdfCatalog, fileName)
		}
	default:
		objectKeyGen = s.defaultObjectKeyGenerator
	}
	uploadResponse, err := s.getUploadURLWithCustomObjectKey(ctx, bucketType, uniqueId, filename, fileSHA256, fileSize, metadata, topic, objectKeyGen)
	if err != nil {
		return nil, errors.Biz("oss.GetS3UploadToken error, get upload url error!")
	}

	return s.prepareS3UploadTokenResponse(uploadResponse)
}

// 获取解析后的文件的临时下载地址
func (s *LocalOssService) GetDownloadTempUrlByFileSHA256AndType(ctx context.Context, parsedDataType pb.ParsedDataEnum, fileSHA256 string, version string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.GetDownloadTeamUrlByFileSHA256AndType")
	defer span.Finish()

	if s.paperPdfParsedService == nil {
		return "", errors.Biz("oss.GetDownloadTeamUrlByFileSHA256AndType error, paperPdfParsedService is nil!")
	}

	// 根据fileSHA256查询文件的paperPdfParsedRecord
	paperPdfParsedRecord, err := s.paperPdfParsedService.GetBySourcePdfFileSHA256AndTypeAndVersion(ctx, fileSHA256, parsedDataType.String(), version)
	if err != nil {
		return "", errors.Biz("oss.GetDownloadTeamUrlByFileSHA256AndType error, get paper pdf parsed record error!")
	}

	if paperPdfParsedRecord == nil {
		return "", errors.Biz("oss.GetDownloadTeamUrlByFileSHA256AndType error, paper pdf parsed record is nil!")
	}
	// 获取桶类型
	bucketType := constant.BucketTypeToEnum(s.config, paperPdfParsedRecord.BucketName)
	//获取文件的临时访问地址
	url, err := s.GetDefaultFileTemporaryURL(ctx, bucketType, paperPdfParsedRecord.ObjectKey)
	if err != nil {
		return "", errors.Biz("oss.GetDownloadTeamUrlByFileSHA256AndType error, get default file temporary url error!")
	}
	return url, nil
}

// GetS3UploadTokenWithCustomObjectKey 获取Minio上传令牌并指定对象键生成器
func (s *LocalOssService) GetS3UploadTokenWithCustomObjectKey(ctx context.Context, bucketType pb.OSSBucketEnum, uniqueId, filename, fileSHA256 string, fileSize int64, metadata map[string]string, topic string, objectKeyGen ObjectKeyGenerator) (*pb.GetUploadTokenResponse, error) {
	s.logger.Info("GetS3UploadTokenWithCustomObjectKey", "bucketType", bucketType.String(), "bucketEnum", int32(bucketType))
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.GetS3UploadTokenWithCustomObjectKey")
	defer span.Finish()

	if objectKeyGen == nil {
		return nil, errors.Biz("oss.GetS3UploadTokenWithCustomObjectKey error, objectKeyGen is nil!")
	}

	uploadResponse, err := s.getUploadURLWithCustomObjectKey(ctx, bucketType, uniqueId, filename, fileSHA256, fileSize, metadata, topic, objectKeyGen)
	if err != nil {
		s.logger.Error("oss.GetS3UploadTokenWithCustomObjectKey error, get upload url error!", "error", err.Error())
		return nil, errors.Biz("oss.GetS3UploadTokenWithCustomObjectKey error, get upload url error!")
	}

	return s.prepareS3UploadTokenResponse(uploadResponse)
}

// prepareS3UploadTokenResponse 准备Minio上传令牌响应
func (s *LocalOssService) prepareS3UploadTokenResponse(uploadResponse *oss.PreSignedUploadResponse) (*pb.GetUploadTokenResponse, error) {
	// 创建响应对象
	resp := &pb.GetUploadTokenResponse{
		NeedUpload: uploadResponse.NeedUpload,
	}
	if uploadResponse.PublicUrl != "" {
		resp.PublicUrl = uploadResponse.PublicUrl
	}

	// 当需要上传时设置 UploadInfo
	if uploadResponse.NeedUpload && uploadResponse.UploadInfo != nil {
		resp.UploadInfo = &pb.UploadS3Info{
			Url:     uploadResponse.UploadInfo.URL,
			Method:  uploadResponse.UploadInfo.Method,
			Headers: uploadResponse.UploadInfo.Headers,
		}
	} else if !uploadResponse.NeedUpload && uploadResponse.OssInfo != nil {
		// 当不需要上传时设置 OssInfo
		resp.OssInfo = &pb.OSSS3Info{
			BucketName: uploadResponse.OssInfo.BucketName,
			ObjectKey:  uploadResponse.OssInfo.ObjectKey,
			FileName:   uploadResponse.OssInfo.FileName,
			FileSize:   fmt.Sprintf("%d", uploadResponse.OssInfo.FileSize),
			FileSHA256: uploadResponse.OssInfo.FileSHA256,
		}
	}
	return resp, nil
}

// prepareUploadURL 准备上传URL的私有方法，处理公共逻辑
func (s *LocalOssService) getUploadURLWithCustomObjectKey(ctx context.Context, bucketEnum pb.OSSBucketEnum, uniqueId, filename, fileSHA256 string, fileSize int64, metadata map[string]string, topic string, objectKeyGen ObjectKeyGenerator) (*oss.PreSignedUploadResponse, error) {
	//根据文件的fileSHA256查询文件是否已经存在上传成功的记录
	ossUploadSuccessRecord, err := s.ossDao.GetSuccessRecordByFileSHA256(ctx, fileSHA256)
	if err != nil {
		return nil, errors.BizWrap("根据 fileSHA256 获取成功文件记录失败", err)
	}
	if ossUploadSuccessRecord != nil {
		// 文件已存在，返回不需要上传的信息
		return &oss.PreSignedUploadResponse{
			NeedUpload: false,
			OssInfo: &oss.OssInfo{
				BucketName: ossUploadSuccessRecord.BucketName,
				ObjectKey:  ossUploadSuccessRecord.ObjectKey,
				FileName:   ossUploadSuccessRecord.FileName,
				FileSize:   ossUploadSuccessRecord.FileSize,
				FileSHA256: ossUploadSuccessRecord.FileSHA256,
			},
		}, nil
	}
	bucketType := constant.EnumToBucketType(s.config, bucketEnum)
	// 获取桶配置
	bucketConfig, err := s.storage.GetBucketConfig(bucketType)
	s.logger.Info("oss.GetUploadURL", "bucketConfig.Name", bucketConfig.Name, "bucketConfig.Public", bucketConfig.Public, "bucketConfig.Versioning", bucketConfig.Versioning, "bucketConfig.IsTemp", bucketConfig.IsTemp, "bucketConfig.Expiration", bucketConfig.Expiration, "bucketType", bucketEnum.String(), "bucketEnum", int32(bucketEnum))
	if err != nil {
		s.logger.Error("oss.GetUploadURL error, get bucket config error!", "error", err.Error())
		return nil, errors.Biz("oss.GetUploadURL error, get bucket config error!")
	}
	//根据文件扩展名设置 content-type
	contentType := getLocalContentType(filename)

	// 使用提供的对象键生成器或默认生成器
	var objectKey string
	if objectKeyGen != nil {
		objectKey = objectKeyGen(uniqueId, filename)
	} else {
		objectKey = s.defaultObjectKeyGenerator(uniqueId, filename)
	}

	id := idgen.GenerateUUID()

	// 1. 准备业务元数据 - 存储到数据库记录中
	bizMetadataStr := ""
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err == nil {
			bizMetadataStr = string(metadataBytes)
		} else {
			s.logger.Error("failed to marshal metadata for database", "error", err.Error())
		}
	}

	// 2. 准备MinIO元数据 - 只保留记录ID，简化元数据处理
	minioMetadata := map[string]string{
		"RecordId": fmt.Sprintf("%d", id),
	}

	// 3. 创建OSS记录，将业务元数据和回调主题存储到数据库中
	record := &model.OssRecord{
		BaseModel: basemodel.BaseModel{
			Id: id,
		},
		BucketName:    bucketConfig.Name,
		ObjectKey:     objectKey,
		FileName:      filename,
		ContentType:   contentType,
		FileSHA256:    fileSHA256,
		Status:        constant.FileStatusPending,
		Visibility:    constant.VisibilityPrivate,
		IsTemp:        false,
		CallbackTopic: topic,          // 回调主题直接存储在数据库记录中
		BizMetadata:   bizMetadataStr, // 业务元数据JSON字符串
	}

	// 保存记录到数据库
	if err := s.ossDao.Save(ctx, record); err != nil {
		// 保存记录失败，阻止URL生成，直接返回错误
		s.logger.Error("oss.GetUploadURL error, save ossUploadRecord error!", "error", err.Error())
		return nil, errors.Biz("oss.GetUploadURL error, save ossUploadRecord error!")
	}

	// 传递给 storage 的 GeneratePreSignedUpload 方法的应该是 bucket 的逻辑名称，而不是 bucketConfig.Name
	// storage 内部会通过逻辑名称去获取真实的 bucket 名称
	uploadResponse, err := s.storage.GeneratePreSignedUpload(ctx, bucketType, objectKey, contentType, fileSize, minioMetadata)
	if err != nil {
		return nil, errors.BizWrap("oss.GetUploadURL error, generate pre-signed upload URL error!", err)
	}
	//如果是公共桶，则返回公共桶的地址
	if bucketConfig.Public {
		permanentURL, err := s.storage.GetPermanentURL(ctx, bucketType, objectKey)
		if err == nil {
			uploadResponse.PublicUrl = permanentURL
		}
	}
	return uploadResponse, nil
}

// getLocalContentType 根据文件名获取 content-type，如果无法识别则返回默认类型
func getLocalContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if mimeType, ok := util.ExtensionMimeTypeMap[ext]; ok {
		return mimeType
	}
	return util.FileTypes.DEFAULT.MimeType
}

// UpdateFileStatus 更新文件记录状态
func (s *LocalOssService) UpdateFileStatus(ctx context.Context, id string, status string, fileSize int64) error {
	s.logger.Info("msg", "更新文件记录状态", "id", id, "status", status, "fileSize", fileSize)

	if err := s.ossDao.UpdateFileStatus(ctx, id, status, fileSize); err != nil {
		return errors.BizWrap("更新文件记录状态失败", err)
	}

	return nil
}

// GetRecordByID 根据 ID 获取文件记录
func (s *LocalOssService) GetRecordByID(ctx context.Context, id string) (*model.OssRecord, error) {
	record, err := s.ossDao.GetRecordByID(ctx, id)
	if err != nil {
		return nil, errors.BizWrap("根据 ID 获取文件记录失败", err)
	}
	if record == nil {
		return nil, errors.Biz("oss.error.record_not_found")
	}
	return record, nil
}

// GetRecordByObjectKey 根据对象键获取文件记录
func (s *LocalOssService) GetRecordByBucketNameAndObjectKey(ctx context.Context, bucketName, objectKey string) (*model.OssRecord, error) {
	record, err := s.ossDao.GetByBucketNameAndObjectKey(ctx, bucketName, objectKey)
	if err != nil {
		return nil, errors.BizWrap("根据对象键获取文件记录失败", err)
	}
	if record == nil {
		return nil, errors.Biz("oss.error.record_not_found")
	}
	return record, nil
}

// GetDefaultFileTemporaryURL 获取文件的临时访问地址
// 根据 OSS 记录或对象键名生成临时访问 URL
// bucketType: 存储桶类型
// objectKey: 对象键名
// URL 有效期（分钟），默认为 10 分钟
func (s *LocalOssService) GetDefaultFileTemporaryURL(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string) (string, error) {
	s.logger.Info("msg", "获取文件临时访问地址",
		"bucketType", bucketType,
		"objectKey", objectKey)

	// 设置默认过期时间
	expiresMinutes := 10 // 默认 10 分钟

	// 检查对象键名是否为空
	if objectKey == "" {
		return "", errors.Biz("oss.error.object_key_empty")
	}

	// 获取桶名称
	bucketName := constant.EnumToBucketType(s.config, bucketType)
	_, err := s.storage.GetBucketConfig(bucketName)
	if err != nil {
		return "", errors.Biz("oss.error.bucket_not_exists")
	}

	// 生成临时 URL
	expires := time.Duration(expiresMinutes) * time.Minute
	url, err := s.storage.GetObjectURL(ctx, bucketName, objectKey, expires)
	if err != nil {
		return "", errors.BizWrap("生成临时访问地址失败", err)
	}

	return url, nil
}

// GetFileTemporaryURL 获取文件的临时访问地址
// 根据 OSS 记录或对象键名生成临时访问 URL
// bucketType: 存储桶类型
// objectKey: 对象键名
// expiresTime: URL 有效期（秒），默认为 60 秒
func (s *LocalOssService) GetFileTemporaryURL(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string, expiresTime int) (string, error) {
	s.logger.Info("msg", "获取文件临时访问地址",
		"bucketType", bucketType,
		"objectKey", objectKey,
		"expiresTime", expiresTime)

	// 根据 objectKey 和 bucketType 查询 OSS 记录，获取 fileSHA256
	record, err := s.ossDao.GetByBucketNameAndObjectKey(ctx, "Local", objectKey)
	if err != nil {
		return "", errors.BizWrap("根据对象键获取文件记录失败", err)
	}
	if record == nil {
		return "", errors.Biz("oss.error.record_not_found")
	}

	// 本地存储路径格式: basePath/fileSHA256/source/uuidFileName
	basePath := utils.ResolveRelativePath(s.config.OSS.Local.BasePath)

	// 从 objectKey 中提取文件名部分（最后一个路径段）
	filePath := filepath.Join(basePath, objectKey)

	// 如果是相对路径，转换为绝对路径
	if !filepath.IsAbs(filePath) {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return "", errors.BizWrap("转换绝对路径失败", err)
		}
		filePath = absPath
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.Biz("oss.error.file_not_found")
	}

	return filePath, nil
}

// GetPermanentURL 获取对象的永久URL（仅适用于公开桶）
// 注意：只有公开桶中的对象才能使用此方法
// bucketType: 存储桶类型
// objectKey: 对象键名
func (s *LocalOssService) GetPermanentURL(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.GetPermanentURL")
	defer span.Finish()

	// 获取桶名称
	bucketName := constant.EnumToBucketType(s.config, bucketType)
	_, err := s.storage.GetBucketConfig(bucketName)
	if err != nil {
		return "", errors.BizWrap("获取存储桶配置失败", err)
	}
	// 调用存储接口获取永久URL
	url, err := s.storage.GetPermanentURL(ctx, bucketName, objectKey)
	if err != nil {
		return "", errors.BizWrap("获取永久URL失败", err)
	}

	return url, nil
}

// UploadObject 直接上传对象到存储服务
// bucketType: 存储桶类型
// objectKey: 对象键名
// reader: 文件内容读取器
// objectSize: 对象大小（字节）
// contentType: 内容类型
// metadata: 元数据
func (s *LocalOssService) UploadObject(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string, reader io.Reader, objectSize int64, contentType string, metadata map[string]string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.UploadObject")
	defer span.Finish()

	// 获取桶名称
	bucketName := constant.EnumToBucketType(s.config, bucketType)

	// 准备选项
	opts := make(map[string]string)
	userMetadata := make(map[string]string)
	if contentType != "" {
		opts["Content-Type"] = contentType
	}

	// 添加元数据
	if metadata != nil {
		for k, v := range metadata {
			userMetadata["X-Amz-Meta-"+k] = v
		}
	}
	//打印调用方法参数
	s.logger.Info("调用存储接口上传对象", "bucket", bucketName, "objectKey", objectKey, "size", objectSize, "contentType", contentType, "metadata", metadata)
	err := s.storage.Upload(ctx, bucketName, objectKey, reader, objectSize, opts, userMetadata)
	if err != nil {
		s.logger.Error("上传对象失败", "error", err.Error(), "bucket", bucketName, "objectKey", objectKey)
		return errors.BizWrap("上传对象失败", err)
	}

	s.logger.Info("上传对象成功", "bucket", bucketName, "objectKey", objectKey, "size", objectSize)
	return nil
}

// UploadObjectAndSaveRecord 上传对象并保存记录
// bucketType: 存储桶类型
// uniqueId: 唯一ID，用于构建对象名称
// fileName: 原始文件名
// fileSHA256: 文件SHA256值
// reader: 文件内容读取器
// objectSize: 对象大小（字节）
// objectKeyGen: 可选的对象键生成器，如果为 nil 则使用默认生成器
func (s *LocalOssService) UploadObjectAndSaveRecord(
	ctx context.Context,
	bucketType pb.OSSBucketEnum,
	fileName string,
	fileSHA256 string,
	reader io.Reader,
	objectSize int64,
	objectKeyGen ObjectKeyGenerator,
	topic string,
	constomMetadata map[string]string,
) (*model.OssRecord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.UploadObjectAndSaveRecord")
	defer span.Finish()

	//根据文件扩展名设置 content-type
	contentType := getLocalContentType(fileName)

	// 获取桶配置
	bucketConfig, err := s.storage.GetBucketConfig(constant.EnumToBucketType(s.config, bucketType))
	if err != nil {
		return nil, errors.BizWrap("获取存储桶配置失败", err)
	}
	if objectKeyGen == nil {
		return nil, errors.Biz("未提供对象键生成器")
	}
	// 使用提供的对象键生成器或默认生成器
	objectKey := objectKeyGen(fileSHA256, fileName)

	// 1. 准备业务元数据 - 存储到数据库记录中
	bizMetadataStr := ""
	if constomMetadata != nil {
		metadataBytes, err := json.Marshal(constomMetadata)
		if err == nil {
			bizMetadataStr = string(metadataBytes)
		} else {
			s.logger.Error("failed to marshal metadata for database", "error", err.Error())
		}
	}
	// 创建OSS记录
	id := idgen.GenerateUUID()
	record := &model.OssRecord{
		BaseModel: basemodel.BaseModel{
			Id: id,
		},
		BucketName:    bucketConfig.Name,
		ObjectKey:     objectKey,
		FileName:      fileName,
		ContentType:   contentType,
		FileSHA256:    fileSHA256,
		Status:        constant.FileStatusPending,
		Visibility:    constant.VisibilityPrivate,
		IsTemp:        false,
		CallbackTopic: topic,          // 回调主题直接存储在数据库记录中
		BizMetadata:   bizMetadataStr, // 业务元数据JSON字符串
	}

	// 保存记录到数据库
	if err := s.ossDao.Save(ctx, record); err != nil {
		return nil, errors.BizWrap("保存OSS记录失败", err)
	}
	metadata := map[string]string{
		"RecordId": fmt.Sprintf("%d", id),
	}
	// 上传对象
	if err := s.UploadObject(ctx, bucketType, objectKey, reader, objectSize, contentType, metadata); err != nil {
		return nil, err
	}
	return record, nil
}

// DownloadObject 从存储服务下载对象   适用于文件过大的场景，可以自由控制文件的读取方式
// bucketType: 存储桶类型
// objectKey: 对象键名
// 返回文件内容的读取器和错误
func (s *LocalOssService) DownloadObject(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string) (io.ReadCloser, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.DownloadObject")
	defer span.Finish()

	// 获取桶名称
	bucketName := constant.EnumToBucketType(s.config, bucketType)
	// 从存储服务下载对象
	reader, err := s.storage.Download(ctx, bucketName, objectKey)
	if err != nil {
		s.logger.Error("从存储服务下载对象失败", "error", err.Error(), "bucket", bucketName, "objectKey", objectKey)
		return nil, errors.Biz("oss.errors.download_failed")
	}

	// 将 io.Reader 转换为 io.ReadCloser
	readCloser, ok := reader.(io.ReadCloser)
	if !ok {
		// 如果不是 io.ReadCloser，则需要包装
		return io.NopCloser(reader), nil
	}

	return readCloser, nil
}

// DownloadObjectAsBytes 下载对象并返回字节数组  适用于文件较小的场景
func (s *LocalOssService) DownloadObjectAsBytes(ctx context.Context, bucketType pb.OSSBucketEnum, objectKey string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.DownloadObjectAsBytes")
	defer span.Finish()

	if objectKey == "" {
		return nil, errors.Biz("objectKey or bucketType is empty")
	}
	reader, err := s.DownloadObject(ctx, bucketType, objectKey)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (s *LocalOssService) GetSuccessRecordByFileSHA256(ctx context.Context, fileSHA256 string) (*model.OssRecord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.DownloadObjectAsBytes")
	defer span.Finish()

	//根据文件的fileSHA256查询文件是否已经存在上传成功的记录
	ossUploadSuccessRecord, err := s.ossDao.GetSuccessRecordByFileSHA256(ctx, fileSHA256)
	if err != nil {
		return nil, errors.BizWrap("根据 fileSHA256 获取成功文件记录失败", err)
	}
	return ossUploadSuccessRecord, nil
}

// SaveOssRecord 保存OSS记录
func (s *LocalOssService) SaveOssRecord(ctx context.Context, record *model.OssRecord) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LocalOssService.SaveOssRecord")
	defer span.Finish()

	if err := s.ossDao.Save(ctx, record); err != nil {
		s.logger.Error("save oss record failed", "error", err.Error())
		return errors.BizWrap("save oss record failed", err)
	}
	return nil
}
