package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/cache"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/utils"
	docpb "github.com/yb2020/odoc/proto/gen/go/doc"
	osspb "github.com/yb2020/odoc/proto/gen/go/oss"
	"github.com/yb2020/odoc/services/doc/constant"
	"github.com/yb2020/odoc/services/doc/model"
	membershipInterface "github.com/yb2020/odoc/services/membership/interfaces"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperModel "github.com/yb2020/odoc/services/paper/model"
	paperService "github.com/yb2020/odoc/services/paper/service"
	parseConstant "github.com/yb2020/odoc/services/parse/constant"
	pdfService "github.com/yb2020/odoc/services/pdf/interfaces"
	pdfModel "github.com/yb2020/odoc/services/pdf/model"
)

// UserDocUploadService 用户文档上传服务
type UserDocUploadService struct {
	logger                logging.Logger
	tracer                opentracing.Tracer
	cache                 cache.Cache
	config                *config.Config
	userDocService        *UserDocService
	paperPdfService       pdfService.IPaperPdfService
	paperPdfParsedService *paperService.PaperPdfParsedService
	ossService            ossService.OssServiceInterface
	membershipService     membershipInterface.IMembershipService
}

// NewUserDocUploadService 创建用户文档上传服务
func NewUserDocUploadService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	cache cache.Cache,
	config *config.Config,
	userDocService *UserDocService,
	paperPdfService pdfService.IPaperPdfService,
	paperPdfParsedService *paperService.PaperPdfParsedService,
	ossService ossService.OssServiceInterface,
	membershipService membershipInterface.IMembershipService,
) *UserDocUploadService {
	return &UserDocUploadService{
		logger:                logger,
		tracer:                tracer,
		cache:                 cache,
		config:                config,
		userDocService:        userDocService,
		paperPdfService:       paperPdfService,
		paperPdfParsedService: paperPdfParsedService,
		ossService:            ossService,
		membershipService:     membershipService,
	}
}

// GetPdfUploadToken 获取PDF上传令牌
func (s *UserDocUploadService) GetPdfUploadToken(ctx context.Context, userId string, req *osspb.GetUploadTokenRequest) (*docpb.GetUploadTokenResp, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocUploadService.GetPdfUploadToken")
	defer span.Finish()

	fileSHA256 := req.GetFileSHA256()
	fileSize := req.GetFileSize()
	filePage := req.GetFilePage()

	// 1. 获取已使用的存储容量
	useStorageCapacity, err := s.paperPdfService.GetPdfFileTotalSize(ctx, userId)
	if err != nil {
		return nil, errors.Biz("get storage capacity failed")
	}

	// 2. 判断当前用户是否存在已经解析好或者没解析好的相同的文件数据
	userDoc, paperPdf, paper, err := s.userDocService.GetUserUploadBaseDataBySHA256AndUserId(ctx, userId, fileSHA256)
	if err != nil {
		// 3. 验证会员权益
		err = s.membershipService.CreditFunDocsUpload(ctx, int64(fileSize), int32(filePage), useStorageCapacity, func(xctx context.Context, sessionId string) error {
			s.logger.Info("msg", "do something", "sessionId", sessionId)
			return nil
		}, true)
		if err != nil {
			return nil, err
		}
	}

	// 4. 构建上传元数据
	fileName := req.GetFileName()
	folderId := req.GetFolderId()
	uuidFileName, metadata := s.buildPdfUploadMetadata(fileName, fileSHA256, folderId)

	// 5. 获取S3上传令牌
	uploadResponse, err := s.requestS3UploadToken(ctx, userId, fileSHA256, int64(fileSize), uuidFileName, metadata)
	if err != nil {
		return nil, err
	}

	// 6. 构建响应并缓存token
	resp, err := s.buildUploadTokenResponse(ctx, userId, fileSHA256, uploadResponse)
	if err != nil {
		return nil, err
	}

	// 7. 如果需要上传，直接返回
	if uploadResponse.NeedUpload {
		return resp, nil
	}

	// 8. 处理已存在文件的解析逻辑
	needParsed, err := s.handleExistingFileParse(ctx, userId, fileSHA256, metadata, resp.OssInfo, userDoc, paper, paperPdf)
	if err != nil {
		return nil, err
	}
	resp.NeedParsed = needParsed

	return resp, nil
}

// buildPdfUploadMetadata 构建PDF上传元数据
func (s *UserDocUploadService) buildPdfUploadMetadata(fileName, fileSHA256, folderId string) (string, map[string]string) {
	ext := filepath.Ext(fileName)
	uuidFileName := idgen.GenerateUUID() + ext

	metadata := map[string]string{
		"fileSHA256": fileSHA256,
		"fileName":   uuidFileName,
	}
	if folderId != "0" {
		metadata["folderId"] = folderId
	}
	return uuidFileName, metadata
}

// requestS3UploadToken 请求S3上传令牌
func (s *UserDocUploadService) requestS3UploadToken(ctx context.Context, userId string, fileSHA256 string, fileSize int64, uuidFileName string, metadata map[string]string) (*osspb.GetUploadTokenResponse, error) {
	topicName := s.config.RocketMQ.Topic.UploadCallback.Name

	customObjectKeyGenerator := func(userId, fileName string) string {
		return fmt.Sprintf("%s/%s/%s", fileSHA256, parseConstant.SourcePdfCatalog, fileName)
	}

	s.logger.Info("GetPdfUploadToken", "bucketType", osspb.OSSBucketEnum_PDF.String(), "bucketEnum", int32(osspb.OSSBucketEnum_PDF))
	s.logger.Info("GetPdfUploadToken", "uniqueId", userId, "uuidFileName", uuidFileName, "fileSHA256", fileSHA256)

	uploadResponse, err := s.ossService.GetS3UploadTokenWithCustomObjectKey(ctx, osspb.OSSBucketEnum_PDF, userId, uuidFileName, fileSHA256, fileSize, metadata, topicName, customObjectKeyGenerator)
	if err != nil {
		return nil, errors.Biz("get upload token failed")
	}
	return uploadResponse, nil
}

// buildUploadTokenResponse 构建上传令牌响应并缓存
func (s *UserDocUploadService) buildUploadTokenResponse(ctx context.Context, userId string, fileSHA256 string, uploadResponse *osspb.GetUploadTokenResponse) (*docpb.GetUploadTokenResp, error) {
	resp := &docpb.GetUploadTokenResp{
		NeedUpload: uploadResponse.NeedUpload,
		UploadInfo: uploadResponse.UploadInfo,
		OssInfo:    uploadResponse.OssInfo,
		NeedParsed: true,
	}

	// 获取一个公共token并写入redis
	token, err := utils.GenerateServiceTokenDefaultWithDefaults(userId)
	if err != nil {
		return nil, errors.Biz("generate token failed")
	}
	resp.Token = token

	// 缓存解析状态
	parseCacheKey := fmt.Sprintf("%s%s", constant.ParsePDFStatusTokenKeyPrefix, token)
	cacheValue := docpb.UserDocParsedStatusObj{
		HasSha256File: true,
		Sha256Value:   fileSHA256,
	}
	s.cache.SetNotBizPrefix(ctx, parseCacheKey, &cacheValue, 30*time.Minute)

	return resp, nil
}

// handleExistingFileParse 处理已存在文件的解析逻辑
func (s *UserDocUploadService) handleExistingFileParse(ctx context.Context, userIdStr, fileSHA256 string, metadata map[string]string, ossInfo *osspb.OSSS3Info, userDoc *model.UserDoc, paper *paperModel.Paper, paperPdf *pdfModel.PaperPdf) (bool, error) {
	// 查询解析记录
	hasExist, err := s.paperPdfParsedService.HasExistBySourcePdfFileSHA256AndVersion(ctx, fileSHA256, parseConstant.ParseVersion)
	if err != nil {
		return false, errors.Biz("get parse record failed")
	}

	if hasExist {
		return false, nil
	}

	// 没有解析记录，需要触发解析
	if userDoc != nil {
		// 用户已有文档记录，发送解析消息
		if err := s.sendParsePdfHeaderMessage(ctx, userIdStr, userDoc, paper, paperPdf); err != nil {
			return false, err
		}
	}

	return true, nil
}

// sendParsePdfHeaderMessage 发送PDF解析头部消息
func (s *UserDocUploadService) sendParsePdfHeaderMessage(ctx context.Context, userIdStr string, userDoc *model.UserDoc, paper *paperModel.Paper, paperPdf *pdfModel.PaperPdf) error {
	return nil
}
