package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	userContext "github.com/yb2020/odoc/pkg/context"

	// "github.com/yb2020/odoc/pkg/errors" // 暂时未使用

	pb "github.com/yb2020/odoc-proto/gen/go/note"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/pkg/transport"
	docService "github.com/yb2020/odoc/services/doc/service"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	paperService "github.com/yb2020/odoc/services/paper/service"
	pdfInterface "github.com/yb2020/odoc/services/pdf/interfaces"
	userService "github.com/yb2020/odoc/services/user/service"
)

// PaperNoteAPI 论文笔记API处理器
type PaperNoteAPI struct {
	service noteInterface.IPaperNoteService

	pdfService     pdfInterface.IPaperPdfService
	userDocService *docService.UserDocService
	userService    *userService.UserService
	paperService   *paperService.PaperService
	logger         logging.Logger
	tracer         opentracing.Tracer
}

// NewPaperNoteAPI 创建论文笔记API处理器
func NewPaperNoteAPI(service noteInterface.IPaperNoteService, pdfService pdfInterface.IPaperPdfService, userDocService *docService.UserDocService, userService *userService.UserService, paperService *paperService.PaperService, logger logging.Logger, tracer opentracing.Tracer) *PaperNoteAPI {
	return &PaperNoteAPI{
		service:        service,
		pdfService:     pdfService,
		userDocService: userDocService,
		userService:    userService,
		paperService:   paperService,
		logger:         logger,
		tracer:         tracer,
	}
}

// SetPaperPdfService 设置论文PDF服务，用于解决循环依赖问题
// func (s *PaperNoteAPI) SetPaperPdfService(pdfService *pdfInterface.IPaperPdfService) error {
// 	if pdfService == nil {
// 		return errors.New("pdfService cannot be nil")
// 	}
// 	s.pdfService = pdfService
// 	return nil
// }

//------ 旧接口翻译扩展 ------//

//------ 旧接口翻译扩展 ------//

// @api_path: /api/note/paperNote/getPaperNoteBaseInfoById
// @method: POST
// @content-type: application/json
// @summary: 获取论文笔记基础信息
func (api *PaperNoteAPI) GetPaperNoteBaseInfoById(c *gin.Context) {

	req := &pb.GetPaperNoteBaseInfoByIdReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	paperNoteBaseInfo, err := api.service.GetPaperNoteBaseInfoById(c.Request.Context(), req.NoteId)
	if err != nil {
		api.logger.Error("msg", "获取论文笔记基础信息失败", "error", err.Error())
		response.ErrorNoData(c, "获取论文笔记基础信息失败")
		return
	}

	response.Success(c, "获取论文笔记基础信息成功", paperNoteBaseInfo)
}

// SetPaperPdfService 设置论文PDF服务
func (api *PaperNoteAPI) SetPaperPdfService(pdfService pdfInterface.IPaperPdfService) error {
	if pdfService == nil {
		return errors.New("pdfService cannot be nil")
	}
	api.pdfService = pdfService
	return nil
}

// @api_path: /api/note/paperNote/getOwnerPaperNoteBaseInfo
// @method: POST
// @content-type: application/json
// @summary: 获取自己的笔记基础信息
func (api *PaperNoteAPI) GetOwnerPaperNoteBaseInfo(c *gin.Context) {
	// 获取用户ID
	userId, _ := userContext.GetUserID(c.Request.Context())
	api.logger.Info("msg", "获取论文笔记基础信息By ID", "userID", userId)

	req := &pb.GetOwnerPaperNoteBaseInfoReq{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	paperNoteBaseInfo, err := api.service.GetOwnerPaperNoteBaseInfo(c.Request.Context(), userId, req.PdfId)
	if err != nil {
		api.logger.Error("msg", "获取论文笔记基础信息失败", "error", err.Error())
		response.ErrorNoData(c, err.Error())
		return
	}
	if paperNoteBaseInfo == nil {
		response.ErrorNoData(c, "无法获取笔记信息且无法新建笔记")
		return
	}

	response.Success(c, "success", paperNoteBaseInfo)
}

// @api_path: /api/note/paperNote/downloadNotePdf
// @method: GET
// @content-type: application/json
// @summary: 下载笔记PDF
func (api *PaperNoteAPI) DownloadNotePdf(c *gin.Context) {
	req := &pb.DownloadNotePdfRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	if req.GetNoteId() == "0" {
		response.ErrorNoData(c, "missing noteId")
		return
	}

	pdfBytes, err := api.service.GetDownloadNoteMarkPdf(c.Request.Context(), req.GetNoteId())
	if err != nil {
		api.logger.Error("msg", "下载笔记PDF失败", "error", err.Error(), "noteId", req.GetNoteId())
		response.ErrorNoData(c, "下载笔记PDF失败")
		return
	}

	// 设置HTTP响应头 "Content-Type"，告诉浏览器我们发送的是一个PDF文件。
	c.Header("Content-Type", "application/pdf")
	// 设置HTTP响应头 "Content-Disposition" 为 "attachment"，这会触发浏览器的文件下载行为，
	// "filename" 用于指定下载文件的默认名称。
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%d.pdf\"", req.GetNoteId()))
	// 将PDF的实际字节数据写入HTTP响应体，并设置HTTP状态码为 200 (http.StatusOK)，表示请求成功。
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// @api_path: /api/note/paperNote/downloadNoteMarkdown
// @method: GET
// @content-type: application/json
// @summary: 下载笔记Markdown
func (api *PaperNoteAPI) DownloadNoteMarkdown(c *gin.Context) {
	req := &pb.DownloadNoteMarkdownRequest{}
	if err := transport.BindProto(c, req); err != nil {
		api.logger.Error("msg", "解析请求参数失败", "error", err.Error())
		response.ErrorNoData(c, "解析请求参数失败")
		return
	}

	if req.GetNoteId() == "0" {
		response.ErrorNoData(c, "missing noteId")
		return
	}

	pdfBytes, err := api.service.GetDownloadNoteMarkMarkdown(c.Request.Context(), req.GetNoteId())
	if err != nil {
		api.logger.Error("msg", "下载笔记PDF失败", "error", err.Error(), "noteId", req.GetNoteId())
		response.ErrorNoData(c, "下载笔记PDF失败")
		return
	}

	// 设置HTTP响应头 "Content-Type"，告诉浏览器我们发送的是一个PDF文件。
	c.Header("Content-Type", "application/zip")
	// 设置HTTP响应头 "Content-Disposition" 为 "attachment"，这会触发浏览器的文件下载行为，
	// "filename" 用于指定下载文件的默认名称。
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%d.zip\"", req.GetNoteId()))
	// 将PDF的实际字节数据写入HTTP响应体，并设置HTTP状态码为 200 (http.StatusOK)，表示请求成功。
	c.Data(http.StatusOK, "application/zip", pdfBytes)
}
