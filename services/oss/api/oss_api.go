package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/oss"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/services/oss/service"
)

// OssAPI OSS API处理器
type OssAPI struct {
	ossService service.OssServiceInterface
	logger     logging.Logger
	tracer     opentracing.Tracer
}

// NewOSSAPI 创建OSS API处理器
func NewOSSAPI(logger logging.Logger, tracer opentracing.Tracer, ossService service.OssServiceInterface) *OssAPI {
	return &OssAPI{
		ossService: ossService,
		logger:     logger,
		tracer:     tracer,
	}
}

// 获取上传
func (api *OssAPI) GetS3UploadToken(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OssServiceAPI.GetS3UploadToken")
	defer span.Finish()
	req := &pb.GetUploadTokenRequest{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "param error")
		return
	}

	//设置默认桶类型为public, TODO:  根据bucketEnum设置桶类型
	uniqueId := idgen.GenerateUUID()

	uploadResponse, err := api.ossService.GetS3UploadToken(ctx, req.BucketEnum, uniqueId, req.FileName, req.GetFileSHA256(), int64(req.GetFileSize()), req.BizMetadata, utils.GetStringFromPtr(req.CallbackTopic, ""), req.KeyPolicy)
	if err != nil {
		response.ErrorNoData(c, err.Error())
		return
	}
	response.Success(c, "success", uploadResponse)
}

// GetDownloadTempUrl 获取解析后的文件的临时url
func (api *OssAPI) GetDownloadTempUrl(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OssServiceAPI.GetDownloadTeamUrl")
	defer span.Finish()

	//
	req := &pb.GetS3DownloadTempUrlReq{}
	if err := transport.BindProto(c, req); err != nil {
		response.ErrorNoData(c, "param error")
		return
	}
	url, err := api.ossService.GetDownloadTempUrlByFileSHA256AndType(ctx, req.GetDataType(), req.GetFileSHA256(), req.GetVersion())
	if err != nil {
		response.ErrorNoData(c, err.Error())
		return
	}

	resp := &pb.GetS3DownloadTempUrlResp{
		FileSHA256:      req.GetFileSHA256(),
		DownloadTempUrl: url,
	}

	response.Success(c, "success", resp)
}

// 测试文件上传到S3
func (api *OssAPI) UploadToS3(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "OssAPI.UploadToS3")
	defer span.Finish()

	//生成一个随机内容的text文件
	uniqueId := idgen.GenerateUUID()
	filename := "test.txt"
	//插入内容
	content := "月光透过画廊的落地窗，在光滑的大理石地板上勾勒出冰冷的银边。展厅中央，巨大的青铜雕塑《缠绕》在阴影中投下交错的轮廓。\n\n“你总能找到最安静的地方，”伊拉（Elara）轻声说，她的指尖温柔地划过李昂（Leo）的手臂。她的长发如黑色的瀑布，眼眸像午夜的星辰，明亮而炽热。她的一举一动都充满了生命力，仿佛一团随时会燃烧的火焰。\n李昂只是微笑，他的身形清瘦而挺拔，眼神却深邃如古井，总能将周围的一切尽收眼底。他的气质是内敛的，像一柄收在鞘中的利刃，只有在必要时才会显露锋芒。他握住伊拉的手，将她的手心贴在自己的胸口，“只有在安静的地方，才能听清心跳。”"
	reader := strings.NewReader(content)
	objectSize := int64(len(content))
	contentType := "text/plain"
	metadata := map[string]string{
		"uniqueId": uniqueId,
	}

	// 使用提供的对象键生成器或默认生成器
	now := time.Now()
	datePrefix := now.Format("2006-01-02")
	objectKey := fmt.Sprintf("%s/%s/%s", datePrefix, uniqueId, filename)
	//打印相关参数
	api.logger.Info("调用存储接口上传对象", "bucket", pb.OSSBucketEnum_PUBLIC.String(), "objectKey", objectKey, "size", objectSize, "contentType", contentType, "metadata", metadata)
	// 上传对象
	err := api.ossService.UploadObject(ctx, pb.OSSBucketEnum_PUBLIC, objectKey, reader, objectSize, contentType, metadata)
	if err != nil {
		api.logger.Error("调用存储接口上传对象失败", "bucket", pb.OSSBucketEnum_PUBLIC.String(), "objectKey", objectKey, "size", objectSize, "contentType", contentType, "metadata", metadata)
		response.ErrorNoData(c, err.Error())
		return
	}
	api.logger.Info("调用存储接口上传对象成功", "bucket", pb.OSSBucketEnum_PUBLIC.String(), "objectKey", objectKey, "size", objectSize, "contentType", contentType, "metadata", metadata)
	response.Success(c, "success", nil)
}
