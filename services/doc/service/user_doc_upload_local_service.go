package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/opentracing/opentracing-go"
	docpb "github.com/yb2020/odoc-proto/gen/go/doc"
	osspb "github.com/yb2020/odoc-proto/gen/go/oss"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/cache"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/squid_proxy"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services/doc/model"
	ossConstant "github.com/yb2020/odoc/services/oss/constant"
	ossModel "github.com/yb2020/odoc/services/oss/model"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperModel "github.com/yb2020/odoc/services/paper/model"
	pdfModel "github.com/yb2020/odoc/services/pdf/model"
)

// UserDocUploadLocalService 用户文档本地上传服务
type UserDocUploadLocalService struct {
	logger         logging.Logger
	tracer         opentracing.Tracer
	cache          cache.Cache
	config         *config.Config
	userDocService *UserDocService
	ossService     ossService.OssServiceInterface
}

// NewUserDocUploadLocalService 创建用户文档本地上传服务
func NewUserDocUploadLocalService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	cache cache.Cache,
	config *config.Config,
	userDocService *UserDocService,
	ossService ossService.OssServiceInterface,
) *UserDocUploadLocalService {
	return &UserDocUploadLocalService{
		logger:         logger,
		tracer:         tracer,
		cache:          cache,
		config:         config,
		userDocService: userDocService,
		ossService:     ossService,
	}
}

// UploadLocalFile 处理本地文件上传
func (s *UserDocUploadLocalService) UploadLocalFile(ctx context.Context, userId string, fileName, localFilePath, folderId string) (*docpb.GetUploadTokenResp, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocUploadLocalService.UploadLocalFile")
	defer span.Finish()

	// 1. 验证本地文件是否存在
	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Biz("local file not found: " + localFilePath)
		}
		return nil, errors.Biz("failed to access local file: " + err.Error())
	}

	// 2. 计算文件SHA256和PDF页码
	fileSHA256, pageCount, err := s.calculateFileInfo(localFilePath)
	if err != nil {
		return nil, errors.Biz("failed to calculate file info: " + err.Error())
	}

	// 3. 检查是否已存在相同文件
	existingUserDoc, existingPdf, _, err := s.userDocService.GetUserUploadBaseDataBySHA256AndUserId(ctx, userId, fileSHA256)
	if err == nil && existingUserDoc != nil {
		// 文件已存在，直接返回
		return &docpb.GetUploadTokenResp{
			NeedUpload: false,
			NeedParsed: false,
			OssInfo: &osspb.OSSS3Info{
				BucketName: existingPdf.OssBucketName,
				ObjectKey:  existingPdf.OssObjectKey,
				FileName:   fileName,
				FileSize:   fmt.Sprintf("%d", existingPdf.Size),
				FileSHA256: existingPdf.FileSHA256,
			},
		}, nil
	}

	// 4. 复制文件到本地存储目录
	objectKey, err := s.copyFileToLocalStorage(localFilePath, fileSHA256, fileName)
	if err != nil {
		return nil, err
	}

	// 5. 创建OSS记录（状态为成功）
	ossRecord, err := s.createOssRecord(ctx, userId, fileSHA256, objectKey, fileName, fileInfo.Size())
	if err != nil {
		return nil, err
	}

	// 6. 创建基础数据记录（UserDoc, Paper, PaperPdf）
	userDoc, paper, paperPdf, err := s.createBaseRecords(ctx, userId, fileSHA256, objectKey, fileName, fileInfo.Size(), pageCount, folderId, ossRecord)
	if err != nil {
		return nil, err
	}

	s.logger.Info("local file upload completed",
		"userId", userId,
		"fileSHA256", fileSHA256,
		"userDocId", userDoc.Id,
		"paperId", paper.PaperId,
		"pdfId", paperPdf.Id)

	return &docpb.GetUploadTokenResp{
		NeedUpload: false,
		NeedParsed: false,
		OssInfo: &osspb.OSSS3Info{
			BucketName: ossRecord.BucketName,
			ObjectKey:  ossRecord.ObjectKey,
			FileName:   ossRecord.FileName,
			FileSize:   fmt.Sprintf("%d", ossRecord.FileSize),
			FileSHA256: ossRecord.FileSHA256,
		},
	}, nil
}

// calculateFileInfo 计算文件的SHA256值和PDF页码
func (s *UserDocUploadLocalService) calculateFileInfo(filePath string) (string, int, error) {
	// 读取文件内容
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", 0, err
	}

	// 计算SHA256
	hash := sha256.Sum256(fileContent)
	fileSHA256 := hex.EncodeToString(hash[:])

	// 计算PDF页码
	reader := bytes.NewReader(fileContent)
	pageCount, err := squid_proxy.GetPdfPageCount(reader)
	if err != nil {
		s.logger.Warn("failed to get pdf page count, using default 0", "error", err.Error(), "filePath", filePath)
		pageCount = 0
	}

	return fileSHA256, pageCount, nil
}

// copyFileToLocalStorage 复制文件到本地存储目录
func (s *UserDocUploadLocalService) copyFileToLocalStorage(srcPath, fileSHA256, fileName string) (string, error) {

	basePath := utils.ResolveRelativePath(s.config.OSS.Local.BasePath)
	// 创建目标目录: basePath/fileSHA256/source/
	targetDir := filepath.Join(basePath, fileSHA256, "source")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", errors.Biz("failed to create storage directory: " + err.Error())
	}

	// 生成唯一文件名
	ext := filepath.Ext(fileName)
	uuidFileName := idgen.GenerateUUID() + ext
	targetPath := filepath.Join(targetDir, uuidFileName)

	// 复制文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", errors.Biz("failed to open source file: " + err.Error())
	}
	defer srcFile.Close()

	dstFile, err := os.Create(targetPath)
	if err != nil {
		return "", errors.Biz("failed to create target file: " + err.Error())
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return "", errors.Biz("failed to copy file: " + err.Error())
	}

	// 返回objectKey格式: fileSHA256/source/uuidFileName
	objectKey := fmt.Sprintf("%s/source/%s", fileSHA256, uuidFileName)
	return objectKey, nil
}

// createOssRecord 创建OSS记录
func (s *UserDocUploadLocalService) createOssRecord(ctx context.Context, userId, fileSHA256, objectKey, fileName string, fileSize int64) (*ossModel.OssRecord, error) {
	// 获取PDF bucket配置
	record := &ossModel.OssRecord{
		BucketName:  "Local",
		ObjectKey:   objectKey,
		FileName:    fileName,
		FileSize:    fileSize,
		FileSHA256:  fileSHA256,
		ContentType: "application/pdf",
		Status:      ossConstant.FileStatusSuccess, // 本地上传直接设置为成功
		Visibility:  ossConstant.VisibilityPrivate,
		IsTemp:      false,
	}
	record.Id = idgen.GenerateUUID()
	record.CreatorId = userId

	// 保存记录
	if err := s.ossService.SaveOssRecord(ctx, record); err != nil {
		return nil, errors.Biz("failed to save oss record: " + err.Error())
	}

	return record, nil
}

// createBaseRecords 创建基础数据记录
func (s *UserDocUploadLocalService) createBaseRecords(ctx context.Context, userId, fileSHA256, objectKey, fileName string, fileSize int64, pageCount int, folderId string, ossRecord *ossModel.OssRecord) (*model.UserDoc, *paperModel.Paper, *pdfModel.PaperPdf, error) {
	// 生成ID
	paperId := idgen.GenerateUUID()
	pdfId := idgen.GenerateUUID()
	userDocId := idgen.GenerateUUID()

	// 从文件名提取标题（去掉扩展名）
	docName := fileName
	ext := filepath.Ext(fileName)
	if ext != "" {
		docName = fileName[:len(fileName)-len(ext)]
	}

	// 创建Paper记录
	paper := &paperModel.Paper{
		PaperId:     paperId,
		OwnerId:     userId,
		Title:       docName,
		ParseStatus: int(docpb.UserDocParsedStatusEnum_BASE_DATA_GENERATED),
	}
	paper.Id = paperId
	paper.CreatorId = userId

	// 创建PaperPdf记录
	paperPdf := &pdfModel.PaperPdf{
		PaperId:       paperId,
		FileSHA256:    fileSHA256,
		Size:          fileSize,
		PageCount:     pageCount,
		OssBucketName: "Local",
		OssObjectKey:  objectKey,
	}
	paperPdf.Id = pdfId
	paperPdf.CreatorId = userId

	// 创建UserDoc记录
	userDoc := &model.UserDoc{
		UserId:      userId,
		PaperId:     paperId,
		PdfId:       pdfId,
		DocName:     docName,
		PaperTitle:  docName,
		ParseStatus: int(docpb.UserDocParsedStatusEnum_BASE_DATA_GENERATED),
	}
	userDoc.Id = userDocId
	userDoc.CreatorId = userId

	// 使用事务保存所有记录
	if err := s.userDocService.SaveUploadRecords(ctx, paperPdf, paper, userDoc, folderId); err != nil {
		return nil, nil, nil, errors.Biz("failed to save upload records: " + err.Error())
	}

	return userDoc, paper, paperPdf, nil
}
