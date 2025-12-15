package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"os"
	"path/filepath"

	"github.com/opentracing/opentracing-go"
	"github.com/signintech/gopdf"
	docpb "github.com/yb2020/odoc-proto/gen/go/doc"
	pb "github.com/yb2020/odoc-proto/gen/go/note"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	userDocService "github.com/yb2020/odoc/services/doc/interfaces"
	docModel "github.com/yb2020/odoc/services/doc/model"
	"github.com/yb2020/odoc/services/note/dao"
	noteInterface "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/note/model"
	paperService "github.com/yb2020/odoc/services/paper/service"
	pdfInterfaces "github.com/yb2020/odoc/services/pdf/interfaces"
	pdfModel "github.com/yb2020/odoc/services/pdf/model"
	userService "github.com/yb2020/odoc/services/user/service"
)

const (
	pageTopMargin    = 40.0
	pageBottomMargin = 40.0
	lineHeight       = 15.0
	fontSize         = 14.0
)

// PaperNoteService 论文笔记服务实现
type PaperNoteService struct {
	paperNoteDAO    *dao.PaperNoteDAO
	logger          logging.Logger
	tracer          opentracing.Tracer
	pdfService      pdfInterfaces.IPaperPdfService
	userDocService  userDocService.IUserDocService
	userService     *userService.UserService
	paperService    *paperService.PaperService
	noteWordService noteInterface.INoteWordService
}

// NewPaperNoteService 创建新的论文笔记服务
func NewPaperNoteService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperNoteDAO *dao.PaperNoteDAO,
	pdfService pdfInterfaces.IPaperPdfService,
	userDocService userDocService.IUserDocService,
	userService *userService.UserService,
	paperService *paperService.PaperService,
	wordService noteInterface.INoteWordService,
) *PaperNoteService {
	return &PaperNoteService{
		paperNoteDAO:    paperNoteDAO,
		logger:          logger,
		tracer:          tracer,
		pdfService:      pdfService,
		userDocService:  userDocService,
		userService:     userService,
		paperService:    paperService,
		noteWordService: wordService,
	}
}

// // SetPaperPdfService 设置论文PDF服务，用于解决循环依赖问题
// func (s *PaperNoteService) SetPaperPdfService(pdfService *pdfService.PaperPdfService) error {
// 	if pdfService == nil {
// 		return errors.Biz("pdfService cannot be nil")
// 	}
// 	s.pdfService = pdfService
// 	return nil
// }

func (s *PaperNoteService) SetNoteWordService(noteWordService noteInterface.INoteWordService) error {
	if noteWordService == nil {
		return errors.Biz("noteWordService cannot be nil")
	}
	s.noteWordService = noteWordService
	return nil
}

// CreatePaperNote 创建论文笔记
func (s *PaperNoteService) CreatePaperNote(ctx context.Context, note *model.PaperNote) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.CreatePaperNote")
	defer span.Finish()

	if err := s.paperNoteDAO.Create(ctx, note); err != nil {
		s.logger.Error("创建论文笔记失败", "error", err)
		return "0", errors.Biz("note.paper_note.errors.create_failed")
	}

	return note.Id, nil
}

// UpdatePaperNote 更新论文笔记
func (s *PaperNoteService) UpdatePaperNote(ctx context.Context, paperNote *model.PaperNote) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.UpdatePaperNote")
	defer span.Finish()

	// 获取论文笔记
	note, err := s.paperNoteDAO.FindById(ctx, paperNote.Id)
	if err != nil {
		s.logger.Error("获取论文笔记失败", "error", err)
		return false, errors.Biz("note.paper_note.errors.get_failed")
	}

	if note == nil {
		return false, errors.Biz("note.paper_note.errors.not_found")
	}

	// 更新论文笔记
	//note = paperNote

	// if err := s.paperNoteDAO.UpdatePaperNote(ctx, note); err != nil {
	// 	s.logger.Error("更新论文笔记失败", "error", err)
	// 	return false, errors.Biz("note.paper_note.errors.update_failed")
	// }
	if err := s.paperNoteDAO.ModifyExcludeNull(ctx, paperNote); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePaperNoteById 删除论文笔记
func (s *PaperNoteService) DeletePaperNoteById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.DeletePaperNoteById")
	defer span.Finish()

	// 删除论文笔记
	if err := s.paperNoteDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文笔记失败", "error", err)
		return false, errors.Biz("note.paper_note.errors.delete_failed")
	}

	return true, nil
}

// SetPaperPdfService 设置论文PDF服务
func (s *PaperNoteService) SetPaperPdfService(pdfService pdfInterfaces.IPaperPdfService) error {
	if pdfService == nil {
		return errors.Biz("pdfService cannot be nil")
	}
	s.pdfService = pdfService
	return nil
}

// GetPaperNoteById 根据ID获取论文笔记
func (s *PaperNoteService) GetPaperNoteById(ctx context.Context, id string) (*model.PaperNote, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetPaperNoteById")
	defer span.Finish()

	// 获取论文笔记
	note, err := s.paperNoteDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文笔记失败", "error", err)
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}

	return note, nil
}

// GetByIds 根据IDS获取论文笔记列表
func (s *PaperNoteService) GetByIds(ctx context.Context, ids []string) ([]model.PaperNote, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetPaperNoteById")
	defer span.Finish()

	// 获取论文笔记
	notes, err := s.paperNoteDAO.FindByIds(ctx, ids)
	if err != nil {
		s.logger.Error("获取论文笔记失败", "error", err)
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}

	return notes, nil
}

func (s *PaperNoteService) GetPaperNoteIdByPdfIdAndUserId(ctx context.Context, pdfId string, userId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetPaperNoteIdByPdfIdAndUserId")
	defer span.Finish()

	// 获取论文笔记
	note, err := s.paperNoteDAO.FindByPdfIdAndUserId(ctx, pdfId, userId)
	if err != nil {
		s.logger.Error("获取论文笔记失败", "error", err)
		return "0", errors.Biz("note.paper_note.errors.get_failed")
	}

	if note == nil {
		return "0", errors.Biz("note.paper_note.errors.not_found")
	}

	// 转换为响应格式
	return note.Id, nil
}

func (s *PaperNoteService) GetPaperNoteByPdfIdAndUserId(ctx context.Context, pdfId string, userId string) (*model.PaperNote, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetPaperNoteByPdfIdAndUserId")
	defer span.Finish()

	// 获取论文笔记
	note, err := s.paperNoteDAO.FindByPdfIdAndUserId(ctx, pdfId, userId)
	if err != nil {
		s.logger.Error("获取论文笔记失败", "error", err)
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}

	// 转换为响应格式
	return note, nil
}

func (s *PaperNoteService) GetPaperNoteBaseInfoById(ctx context.Context, noteId string) (*pb.PaperNoteBaseInfoResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetPaperNoteBaseInfoById")
	defer span.Finish()

	// 获取用户ID
	userId, _ := userContext.GetUserID(ctx)

	//1. 根据noteId获取论文笔记
	paperNote, err := s.GetPaperNoteById(ctx, noteId)
	if err != nil {
		s.logger.Error("msg", "论文笔记不存在", "error", err.Error())
		return nil, errors.Biz("note.paper_note.errors.not_found")
	}
	s.logger.Info("msg", "查询论文笔记", "paperNote", paperNote)

	//2. 根据pdfId获取PaperPdf
	paperPdf, err := s.pdfService.GetById(ctx, paperNote.PdfId)
	if err != nil {
		s.logger.Error("msg", "PaperPdf不存在", "error", err.Error())
		return nil, errors.Biz("note.paper_pdf.errors.not_found")
	}
	s.logger.Info("msg", "获取PaperPdf成功", "paperPdf", paperPdf)

	//3. 如果当前用户是笔记创建者且笔记的论文ID与PDF关联的论文ID不匹配，则异步更新笔记的论文ID
	// 判断条件：PDF文件关联的论文ID不为空且不是0 && 当前笔记的论文ID与PDF关联的论文ID不匹配 && 用户已登录 &&当前用户是这个笔记的创建者
	if paperPdf.PaperId != "" && // PDF文件关联的论文ID不为空
		paperNote.PaperId != paperPdf.PaperId && // 当前笔记的论文ID与PDF关联的论文ID不匹配
		paperNote.CreatorId == userId { // 用户已登录且是这个笔记的创建者
		// 异步更新笔记的论文ID
		paperNote.PaperId = paperPdf.PaperId
		success, err := s.UpdatePaperNote(ctx, paperNote)
		if err != nil {
			s.logger.Error("msg", "更新论文ID失败",
				"error", err.Error(),
				"noteId", paperNote.Id)
		} else if success {
			s.logger.Info("msg", "更新论文ID成功",
				"noteId", paperNote.Id,
				"paperId", paperPdf.PaperId)
		}
	}

	paperNoteBaseInfo := &pb.PaperNoteBaseInfoResponse{}

	//4. 设置是否点赞
	paperNoteBaseInfo.IsPaperCollected = false
	paperNoteBaseInfo.IsCollected = false
	paperNoteBaseInfo.IsLike = false // 设置是否点赞为false，新版本已经取消点赞功能
	paperNoteBaseInfo.LikeCount = 0  // 设置点赞数为0，新版本已经取消点赞功能

	//5. 设置UserDoc
	var userDoc *docModel.UserDoc
	if paperNote.PaperId != "" {
		userDoc, err = s.userDocService.GetByUserIDAndPaperID(ctx, paperNote.CreatorId, paperNote.PaperId)
		if err != nil {
			s.logger.Error("msg", "获取UserDoc失败", "error", err.Error())
			return nil, errors.Biz("doc.user_doc.errors.get_failed")
		}
	} else if paperNote.PdfId != "" {
		userDoc, err = s.userDocService.GetByUserIdAndPdfId(ctx, paperNote.CreatorId, paperNote.PdfId)
		if err != nil {
			s.logger.Error("msg", "获取UserDoc失败", "error", err.Error())
			return nil, errors.Biz("doc.user_doc.errors.get_failed")
		}
	}
	s.logger.Info("msg", "获取UserDoc成功", "userDoc", userDoc)
	if userDoc != nil {
		paperNoteBaseInfo.DocName = userDoc.DocName
		paperNoteBaseInfo.UserDocId = userDoc.Id
		paperNoteBaseInfo.PaperTitle = userDoc.PaperTitle
	}

	//6. 设置更新时间
	paperNoteBaseInfo.ModifyDate = strconv.FormatInt(paperNote.UpdatedAt.UnixNano()/1e6, 10) //毫秒级时间戳

	if paperNote.PaperId != "" {
		paperNoteBaseInfo.PaperId = paperNote.PaperId
		paperNoteBaseInfo.IsPrivatePaper, _ = s.paperService.IsPrivatePaper(ctx, paperNote.PaperId)
	}

	paperNoteBaseInfo.NoteId = paperNote.Id

	//7. 设置Pdf信息（贡献者、pdf链接）
	user, _ := s.userService.GetUserByID(ctx, paperPdf.CreatorId)
	if user != nil {
		paperNoteBaseInfo.PdfContributor = user.Username
	} else {
		paperNoteBaseInfo.PdfContributor = "第三方"
	}
	paperNoteBaseInfo, err = s.setPdfInfo(ctx, paperNoteBaseInfo, paperNote.PdfId)
	if err != nil {
		s.logger.Error("msg", "设置Pdf信息失败", "error", err.Error())
		return nil, errors.Biz("note.paper_pdf.errors.set_failed")
	}

	//8. 设置Title MOCK:TODO
	// paperBaseInfo, err := s.paperService.GetPaperBaseInfoById(ctx, paperNote.PaperId)
	// if err != nil {
	// 	s.logger.Error("msg", "获取论文基本信息失败", "error", err.Error())
	// 	return nil, errors.Biz("note.paper.errors.get_failed")
	// }
	// if paperBaseInfo != nil {
	// 	paperNoteBaseInfo.PaperTitle = paperBaseInfo.PaperTitle
	// }

	//9.设置笔记作者信息，新版本忽略此信息设置, MOCK:TODO
	// paperNoteBaseInfo.UserInfo = &pb.UserInfoBean{
	// 	Id:               12345,
	// 	NickName:         "研究学者",
	// 	Self:             false,
	// 	ShowName:         "研究学者",
	// 	AvatarUrl:        "https://example.com/avatars/researcher.png",
	// 	Tags:             "AI,机器学习,自然语言处理",
	// 	Description:      "专注于人工智能和机器学习领域的研究",
	// 	AuthorId:         67890,
	// 	AuthorName:       "张教授",
	// 	IsAuthentication: true,
	// 	IsPaperAuthor:    true,
	// }
	user, userErr := s.userService.GetUserByID(ctx, userDoc.CreatorId)
	if userErr != nil {
		s.logger.Error("msg", "获取用户信息失败", "error", userErr.Error())
		return nil, errors.Biz("user.user.errors.get_failed")
	}
	paperNoteBaseInfo.UserInfo = &pb.UserInfoBean{
		Id:               user.Id,
		NickName:         user.Nickname,
		Self:             false,
		ShowName:         user.Nickname,
		AvatarUrl:        "",
		Tags:             "",
		Description:      "",
		AuthorId:         "67890",
		AuthorName:       user.Username,
		IsAuthentication: true,
		IsPaperAuthor:    true,
	}

	paperNoteBaseInfo.IsGptWhite = true
	//10. 设置showAnnotation MOCK:TODO
	paperNoteBaseInfo.ShowAnnotation = true

	//paperNoteBaseInfo.PdfUrl = "https://pdf.cdn.readpaper.com/spd/2294053032.pdf"
	// pdfUrl, err := s.pdfService.GetPdfUrlById(ctx, paperNote.PdfId, 6*360)
	// if err != nil {
	// 	s.logger.Error("获取PDF文件地址失败", "error", err)
	// 	return nil, errors.Biz("pdf.paper_pdf.errors.get_url_failed")
	// }
	// paperNoteBaseInfo.PdfUrl = *pdfUrl

	// 模拟返回数据
	// paperNoteBaseInfo := &pb.PaperNoteBaseInfoResponse{
	// 	PdfId:          12345,
	// 	NoteId:         req.NoteId, // 使用请求中的笔记ID
	// 	PaperId:        67890,
	// 	PdfContributor: "张三",
	// 	PaperNoteCount: 5,
	// 	IsCollected:    true,
	// 	IsUserUpload:   false,
	// 	IsPrivatePaper: false,
	// 	ShowAnnotation: true,
	// 	PaperTitle:     "人工智能在医疗领域的应用研究",
	// 	SourceMark:     "IEEE",
	// 	CrawlUrl:       "https://example.com/papers/12345",
	// 	UploadUserId:   10086,
	// 	LicenceType:    "CC BY-NC-SA 4.0",
	// 	DocName:        "AI_Healthcare_Research.pdf",
	// 	NoteSummary:    "本文探讨了人工智能技术在医疗诊断、药物研发和健康管理等方面的应用，并分析了当前面临的挑战和未来发展趋势。",
	// 	PdfUrl:         "https://example.com/download/12345.pdf",
	// 	ModifyDate:     1711419600, // 2024-03-26 时间戳
	// 	// UserInfo 字段已被注释
	// 	IsPaperCollected:           true,
	// 	IsLike:                     false,
	// 	LikeCount:                  42,
	// 	IsGptWhite:                 true,
	// 	GptGrayTip:                 "AI 辅助阅读已启用",
	// 	UserDocId:                  54321,
	// 	AccessToAiAssistantReading: true,
	// }

	return paperNoteBaseInfo, nil
}

func (s *PaperNoteService) setPdfInfo(ctx context.Context, paperNoteBaseInfo *pb.PaperNoteBaseInfoResponse, pdfId string) (*pb.PaperNoteBaseInfoResponse, error) {
	paperPdf, _ := s.pdfService.GetById(ctx, pdfId)
	if paperPdf != nil {
		//查到pdf链接
		paperNoteBaseInfo.PdfId = paperPdf.Id
		pdfUrl, err := s.pdfService.GetPdfUrlById(ctx, paperPdf.Id, 6*360)
		if err != nil {
			s.logger.Error("获取PDF文件地址失败", "error", err)
			return nil, errors.Biz("pdf.paper_pdf.errors.get_url_failed")
		}
		paperNoteBaseInfo.PdfUrl = *pdfUrl

		//!CommonConstant.SYSTEM_ID.equals(dbPaperPdf.getCreatorId())
		if paperPdf.CreatorId != "0" {
			paperNoteBaseInfo.IsUserUpload = true
			paperNoteBaseInfo.SourceMark = paperNoteBaseInfo.PdfContributor
			paperNoteBaseInfo.UploadUserId = paperPdf.CreatorId
		} else {
			paperNoteBaseInfo.IsUserUpload = false
			//这个属性已经被删除
			// parsedURL, err := url.Parse(paperPdf.SourceUrl)
			// if err == nil {
			// 	paperNoteBaseInfo.SourceMark = parsedURL.Host // 获取主机名
			// }
			// paperNoteBaseInfo.CrawlUrl = paperPdf.SourceUrl
		}

	} else {
		paperNoteBaseInfo.PdfContributor = "第三方"
	}
	return paperNoteBaseInfo, nil
}

func (s *PaperNoteService) GetOwnerPaperNoteBaseInfo(ctx context.Context, userId string, pdfId string) (*pb.PaperNoteBaseInfoResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetOwnerPaperNoteBaseInfo")
	defer span.Finish()

	// 1. 判断Pdf 是否私密，如果私密则使用上传者的noteId
	paperPdf, err := s.pdfService.GetById(ctx, pdfId)
	if err != nil {
		s.logger.Error("msg", "获取PDF失败", "error", err.Error())
		return nil, errors.Biz("note.paper_pdf.errors.get_failed")
	}
	if paperPdf == nil {
		return nil, errors.Biz("can no find paperpdf")

	}
	s.logger.Debug("msg", "获取PDF成功", "pdf", paperPdf)
	// 2. TODO:查找pdfId和笔记权限
	//if (Objects.isNull(pdfId) || pdfAuthService.permissionDenied(pdfId, userId)) {
	// 	throw new BusinessException("您暂无权限在该pdf上做笔记");
	// }
	isDenied, err := s.pdfService.AuthPermissionDenied(ctx, pdfId, userId)
	if err != nil {
		s.logger.Error("msg", "Pdf权限认证失败", "error", err.Error())
	}
	if isDenied {
		return nil, errors.Biz("您暂无权限在该pdf上做笔记")
	}

	// 3. 查询笔记noteId或者创建笔记
	note, err := s.GetPaperNoteByPdfIdAndUserId(ctx, pdfId, userId)
	if err != nil {
		s.logger.Error("msg", "获取笔记失败", "error", err.Error())
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}
	if note == nil {
		// 创建笔记
		note = &model.PaperNote{
			PdfId:   pdfId,
			PaperId: paperPdf.PaperId,
		}
		note.Id = idgen.GenerateUUID()
		note.CreatorId = userId
		note.AnnotationPdfId = pdfId

		note.Id, err = s.CreatePaperNote(ctx, note)
		if err != nil {
			s.logger.Error("msg", "创建笔记失败", "error", err.Error())
			return nil, errors.Biz("note.paper_note.errors.create_failed")
		}

		//创建笔记后，需要将笔记NoteId更新到UserDoc表中
		err = s.userDocService.UpdateUserDocNoteIdByPdfId(ctx, pdfId, note.Id)
		if err != nil {
			s.logger.Error("msg", "更新UserDoc失败", "error", err.Error())
			return nil, errors.Biz("note.paper_note.errors.update_failed")
		}
	}

	// 4. 获取笔记基础信息
	paperNoteBaseInfo, err := s.GetPaperNoteBaseInfoById(ctx, note.Id)
	if err != nil {
		s.logger.Error("msg", "获取笔记基础信息失败", "error", err.Error())
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}
	//这个属性已经被删除
	// paperNoteBaseInfo.LicenceType = paperPdf.LicenceType
	return paperNoteBaseInfo, nil
}

// 根据用户ID获取笔记列表
func (s *PaperNoteService) SelectByUserIdLimit(ctx context.Context, userId string, limit int) ([]model.PaperNote, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.SelectByUserIdLimit")
	defer span.Finish()
	// 调用DAO层方法获取用户笔记列表，使用更新后的ctx
	return s.paperNoteDAO.SelectByUserIdLimit(ctx, userId, limit)
}

// 根据用户ID获取笔记数量
func (s *PaperNoteService) GetUserPaperNoteCount(ctx context.Context, userId string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetUserPaperNoteCount")
	defer span.Finish()
	// 调用DAO层方法获取用户笔记数量，使用更新后的ctx
	return s.paperNoteDAO.CountByUserId(ctx, userId)
}

// GetByPaperIdAndUserIdOrderByNoteCount 根据论文ID和用户ID获取笔记，按笔记数量和修改时间排序
func (s *PaperNoteService) GetByPaperIdAndUserIdOrderByNoteCount(ctx context.Context, paperId string, userId string) (*model.PaperNote, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetByPaperIdAndUserIdOrderByNoteCount")
	defer span.Finish()

	// 调用DAO层方法获取笔记
	return s.paperNoteDAO.GetByPaperIdAndUserIdOrderByNoteCount(ctx, paperId, userId)
}

// GetPdfByNoteId 根据NotId获取Pdf信息
func (s *PaperNoteService) GetPdfByNoteId(ctx context.Context, noteId string) (*pdfModel.PaperPdf, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetDocByNoteId")
	defer span.Finish()

	note, err := s.GetPaperNoteById(ctx, noteId)
	if err != nil {
		s.logger.Error("获取论文笔记失败", "error", err)
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}

	if note == nil {
		return nil, errors.Biz("note.paper_note.errors.not_found")
	}

	pdf, err := s.pdfService.GetById(ctx, note.PdfId)
	if err != nil {
		s.logger.Error("获取pdf失败", "error", err)
		return nil, errors.Biz("note.paper_pdf.errors.get_failed")
	}

	return pdf, nil
}

// GetDocByNoteId 根据NotId获取文献信息
func (s *PaperNoteService) GetDocByNoteId(ctx context.Context, noteId string) (*docModel.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetDocByNoteId")
	defer span.Finish()

	pdf, err := s.GetPdfByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取pdf失败", "error", err)
		return nil, errors.Biz("note.paper_pdf.errors.get_failed")
	}

	if pdf == nil {
		return nil, errors.Biz("note.paper_pdf.errors.not_found")
	}

	return s.pdfService.GetDocById(ctx, pdf.Id)
}

// GetNoteManageFolderInfosByDocIds 根据文献Ids获取文献信息
func (s *PaperNoteService) GetNoteManageFolderInfosByDocIds(ctx context.Context, userId string, docIds []string) ([]*pb.NoteManageFolderInfo, int32, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetNoteManageFolderInfosByDocIds")
	defer span.Finish()

	folderInfos, total, err := s.userDocService.GetFolderInfosByUserIdAndIds(ctx, userId, docIds)
	if err != nil {
		s.logger.Error("获取笔记列表失败", "error", err, "userId", userId)
		return nil, 0, errors.Biz("note.paper_note.errors.list_failed")
	}

	s.logger.Info("msg", "folderInfos info", "folderInfos", folderInfos)

	// resp := &pb.GetSummaryListResponse{
	// 	Total: uint32(total),
	// }

	var noteFolderInfos []*pb.NoteManageFolderInfo

	// 2. 文件夹及子文件夹和文件转换
	if len(folderInfos) > 0 {
		noteFolderInfos = make([]*pb.NoteManageFolderInfo, 0, len(folderInfos))

		// 定义递归函数处理文件夹转换
		var convertFolder func(folder *docpb.UserDocFolderInfo) *pb.NoteManageFolderInfo
		convertFolder = func(folder *docpb.UserDocFolderInfo) *pb.NoteManageFolderInfo {
			if folder == nil {
				return nil
			}

			// 创建文件夹信息对象
			folderInfo := &pb.NoteManageFolderInfo{
				Name:     folder.Name,
				Count:    folder.DocCount,
				FolderId: folder.FolderId,
			}

			// 递归处理子文件夹 - 支持任意多层级的文件夹结构
			if len(folder.ChildrenFolders) > 0 {
				folderInfo.ChildrenFolders = make([]*pb.NoteManageFolderInfo, 0, len(folder.ChildrenFolders))
				for _, child := range folder.ChildrenFolders {
					childFolder := convertFolder(child) // 递归调用
					if childFolder != nil {
						folderInfo.NoteWordCount += childFolder.NoteWordCount
						folderInfo.NoteAnnotateCount += childFolder.NoteAnnotateCount
						folderInfo.ChildrenFolders = append(folderInfo.ChildrenFolders, childFolder)
					}
				}
			}

			// 处理当前文件夹下的文档
			if len(folder.DocInfos) > 0 {
				folderInfo.DocInfos = make([]*pb.NoteManageDocInfo, 0, len(folder.DocInfos))
				for _, docInfo := range folder.DocInfos {
					if docInfo != nil {
						docInfoPB := &pb.NoteManageDocInfo{
							DocName:    docInfo.DocName,
							DocId:      docInfo.DocId,
							NoteId:     docInfo.NoteId,
							ModifyDate: &docInfo.ModifyDate,
						}
						wordCount, _ := s.noteWordService.GetCountByNoteId(ctx, docInfo.NoteId)
						docInfoPB.NoteWordCount = uint32(wordCount)
						folderInfo.NoteWordCount += docInfoPB.NoteWordCount
						noteAnnotateCount, _ := s.pdfService.GetCountPdfMarksByNoteId(ctx, docInfo.NoteId)
						docInfoPB.NoteAnnotateCount = uint32(noteAnnotateCount)
						folderInfo.NoteAnnotateCount += docInfoPB.NoteAnnotateCount
						folderInfo.DocInfos = append(folderInfo.DocInfos, docInfoPB)
					}
				}
			}

			return folderInfo
		}

		// 处理顶级文件夹
		for _, folder := range folderInfos {
			if folder != nil {
				folderInfo := convertFolder(folder)
				if folderInfo != nil {
					noteFolderInfos = append(noteFolderInfos, folderInfo)
				}
			}
		}
		s.logger.Info("文件夹转换完成", "foldersCount", len(noteFolderInfos))
	}

	return noteFolderInfos, total, nil
}

// GetUnclassifiedDocInfosByDocIds 根据文献Ids获取未分类的文献列表
func (s *PaperNoteService) GetUnclassifiedDocInfosByDocIds(ctx context.Context, userId string, docIds []string) ([]*pb.NoteManageDocInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetUnclassifiedDocInfosByDocIds")
	defer span.Finish()

	//业务逻辑: 根据用户userId查询未分类文献列表，然后过滤掉不含docIds的数据
	docInfos, err := s.userDocService.GetDocTreeRootNodeOfDocsByUserId(ctx, userId, true)
	if err != nil {
		s.logger.Error("msg", "获取用户未分类文献列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Biz("")
	}

	// docInfos 过滤掉不含docIds的数据
	// 创建docIds的映射，用于快速查找
	docIdsMap := make(map[string]bool)
	for _, id := range docIds {
		docIdsMap[id] = true
	}

	// 过滤docInfos，只保留docIds中包含的文档
	filteredDocInfos := make([]*docpb.SimpleUserDocInfo, 0)
	for _, docInfo := range docInfos {
		if docInfo != nil && docIdsMap[docInfo.DocId] {
			filteredDocInfos = append(filteredDocInfos, docInfo)
		}
	}
	docInfos = filteredDocInfos

	// 声明noteDocInfos变量
	var noteDocInfos []*pb.NoteManageDocInfo
	// 3. 未分类文献转换
	if len(docInfos) > 0 {
		noteDocInfos = make([]*pb.NoteManageDocInfo, 0, len(docInfos))
		for _, docInfo := range docInfos {
			if docInfo != nil {
				docInfoPB := &pb.NoteManageDocInfo{
					DocName:    docInfo.DocName,
					DocId:      docInfo.DocId,
					NoteId:     docInfo.NoteId,
					ModifyDate: &docInfo.ModifyDate,
				}
				wordCount, _ := s.noteWordService.GetCountByNoteId(ctx, docInfo.NoteId)
				docInfoPB.NoteWordCount = uint32(wordCount)
				noteAnnotateCount, _ := s.pdfService.GetCountPdfMarksByNoteId(ctx, docInfo.NoteId)
				docInfoPB.NoteAnnotateCount = uint32(noteAnnotateCount)
				noteDocInfos = append(noteDocInfos, docInfoPB)
			}
		}
		// s.logger.Info("未分类文档转换完成", "unclassifiedCount", len(noteDocInfos))
	}

	return noteDocInfos, nil
}

// GetNoteDocTreeNodeByFolderId 根据文件夹ID获取文件夹
func (s *PaperNoteService) GetNoteDocTreeNodeByFolderId(ctx context.Context, userId string, folderId string) (*pb.NoteManageFolderInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetNoteDocTreeNodeByFolderId")
	defer span.Finish()

	folderInfo, err := s.userDocService.GetDocTreeNodeByFolderId(ctx, userId, folderId, true)
	if err != nil {
		s.logger.Error("获取用户文件夹信息失败", "error", err)
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}

	if folderInfo != nil {
		// 定义递归函数处理文件夹转换
		var convertFolder func(folder *docpb.UserDocFolderInfo) *pb.NoteManageFolderInfo
		convertFolder = func(folder *docpb.UserDocFolderInfo) *pb.NoteManageFolderInfo {
			if folder == nil {
				return nil
			}

			// 创建文件夹信息对象
			folderInfo := &pb.NoteManageFolderInfo{
				Name:     folder.Name,
				Count:    folder.DocCount,
				FolderId: folder.FolderId,
			}

			// 递归处理子文件夹 - 支持任意多层级的文件夹结构
			if len(folder.ChildrenFolders) > 0 {
				folderInfo.ChildrenFolders = make([]*pb.NoteManageFolderInfo, 0, len(folder.ChildrenFolders))
				for _, child := range folder.ChildrenFolders {
					childFolder := convertFolder(child) // 递归调用
					if childFolder != nil {
						folderInfo.NoteWordCount += childFolder.NoteWordCount
						folderInfo.NoteAnnotateCount += childFolder.NoteAnnotateCount
						folderInfo.ChildrenFolders = append(folderInfo.ChildrenFolders, childFolder)
					}
				}
			}

			// 处理当前文件夹下的文档
			if len(folder.DocInfos) > 0 {
				folderInfo.DocInfos = make([]*pb.NoteManageDocInfo, 0, len(folder.DocInfos))
				for _, docInfo := range folder.DocInfos {
					if docInfo != nil {
						docInfoPB := &pb.NoteManageDocInfo{
							DocName:    docInfo.DocName,
							DocId:      docInfo.DocId,
							NoteId:     docInfo.NoteId,
							ModifyDate: &docInfo.ModifyDate,
						}
						wordCount, _ := s.noteWordService.GetCountByNoteId(ctx, docInfo.NoteId)
						docInfoPB.NoteWordCount = uint32(wordCount)
						folderInfo.NoteWordCount += docInfoPB.NoteWordCount
						noteAnnotateCount, _ := s.pdfService.GetCountPdfMarksByNoteId(ctx, docInfo.NoteId)
						docInfoPB.NoteAnnotateCount = uint32(noteAnnotateCount)
						folderInfo.NoteAnnotateCount += docInfoPB.NoteAnnotateCount
						folderInfo.DocInfos = append(folderInfo.DocInfos, docInfoPB)
					}
				}
			}

			return folderInfo
		}

		// 调用递归函数处理顶级文件夹
		return convertFolder(folderInfo), nil
	}

	return nil, nil
}

// GetAllNoteDocInfosByFolderId 根据文件夹ID获取文件夹及子文件夹下的所有文献，folderId=0时获取用户的所有文献；返回文献列表
func (s *PaperNoteService) GetAllNoteDocInfosByFolderId(ctx context.Context, userId string, folderId string) ([]*pb.NoteManageDocInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetAllNoteDocInfosByFolderId")
	defer span.Finish()
	docs, err := s.userDocService.GetAllDocsByFolderId(ctx, userId, folderId, true)
	if err != nil {
		s.logger.Error("获取用户文献列表失败", "error", err)
		return nil, errors.Biz("note.paper_note.errors.list_failed")
	}

	// docs 转换为 docsInfos
	docsInfos := make([]*pb.NoteManageDocInfo, 0, len(docs))
	for _, doc := range docs {
		docInfo := &pb.NoteManageDocInfo{
			DocName:    doc.DocName,
			DocId:      doc.DocId,
			NoteId:     doc.NoteId,
			ModifyDate: &doc.ModifyDate,
		}
		wordCount, _ := s.noteWordService.GetCountByNoteId(ctx, docInfo.NoteId)
		docInfo.NoteWordCount = uint32(wordCount)
		noteAnnotateCount, _ := s.pdfService.GetCountPdfMarksByNoteId(ctx, docInfo.NoteId)
		docInfo.NoteAnnotateCount = uint32(noteAnnotateCount)
		docsInfos = append(docsInfos, docInfo)
	}

	return docsInfos, nil
}

// GetAllNoteDocInfosPageByFolderId 根据文件夹ID获取文件夹及子文件夹下的所有文献分页，folderId=0时获取用户的所有文献分页；返回文献列表和总记录数
func (s *PaperNoteService) GetAllNoteDocInfosPageByFolderId(ctx context.Context, userId string, folderId string, page int32, pageSize int32) ([]*pb.NoteManageDocInfo, int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetAllNoteDocInfosPageByFolderId")
	defer span.Finish()

	docsInfos, err := s.GetAllNoteDocInfosByFolderId(ctx, userId, folderId)
	if err != nil {
		s.logger.Error("获取用户文献列表失败", "error", err)
		return nil, 0, errors.Biz("note.paper_note.errors.list_failed")
	}

	// 计算总数量
	total := int64(len(docsInfos))

	// 验证分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 计算起始和结束索引
	startIndex := (page - 1) * pageSize
	endIndex := page * pageSize

	// 检查是否超出范围
	if int32(total) <= startIndex {
		// 如果起始索引超出了总数量，返回空数组
		return []*pb.NoteManageDocInfo{}, total, nil
	}

	// 确保结束索引不超出数组边界
	if int32(total) < endIndex {
		endIndex = int32(total)
	}

	// 切片获取当前页的数据
	docsInfos = docsInfos[startIndex:endIndex]

	return docsInfos, int64(len(docsInfos)), nil
}

// GetDownloadNoteMarkPdf 生成包含笔记和标注的PDF文档
func (s *PaperNoteService) GetDownloadNoteMarkPdf(ctx context.Context, noteId string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetDownloadNoteMarkPdf")
	defer span.Finish()

	// 1. 根据noteId查询笔记信息
	note, err := s.GetPaperNoteById(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记信息失败", "error", err)
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}

	if note == nil {
		s.logger.Error("笔记不存在", "noteId", noteId)
		return nil, errors.Biz("note.paper_note.errors.not_found")
	}

	// 2. 获取标题
	doc, err := s.userDocService.GetByPdId(ctx, note.PdfId)
	if err != nil {
		s.logger.Error("获取文献失败", "error", err)
		return nil, errors.Biz("note.paper_pdf.errors.get_failed")
	}
	if doc == nil {
		s.logger.Error("文献不存在", "pdfId", note.PdfId)
		return nil, errors.Biz("note.paper_pdf.errors.not_found")
	}
	title := doc.DocName
	s.logger.Debug("获取文献成功", "title", title)

	// 3. 获取PdfMark标注
	marks, err := s.pdfService.GetPdfMarksByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}
	s.logger.Debug("获取PDF标记成功", "marks", marks)

	// 4. 生成PDF
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// 4.1 动态获取项目根目录，以构造字体文件的绝对路径
	// 尝试多个可能的路径来定位脚本文件
	possiblePaths := []string{
		filepath.Join("resources", "fonts", "SourceHanSansSC-Regular.ttf"),                   // 相对于工作目录
		filepath.Join("..", "..", "resources", "fonts", "SourceHanSansSC-Regular.ttf"),       // 相对于服务目录
		filepath.Join("..", "..", "..", "resources", "fonts", "SourceHanSansSC-Regular.ttf"), // 更深层级
	}

	var fontPath string

	// 尝试每个可能的路径
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			fontPath = path
			break
		}
	}

	// 如果没有找到脚本文件，记录错误
	if fontPath == "" {
		s.logger.Error("找不到字体文件", "paths", possiblePaths)
		return nil, fmt.Errorf("找不到字体文件")
	}

	s.logger.Info("加载字体文件的绝对路径", "path", fontPath)

	err = pdf.AddTTFFont("chinese", fontPath)
	if err != nil {
		s.logger.Error("加载字体失败", "error", err)
		return nil, errors.Wrap(err, "加载字体失败, 请检查字体文件路径")
	}

	err = pdf.SetFont("chinese", "", 14)
	if err != nil {
		s.logger.Error("设置字体失败", "error", err)
		return nil, errors.Wrap(err, "设置字体失败")
	}

	// 4.2写入标题
	pdf.SetX(10)
	pdf.SetY(10)
	err = pdf.Cell(nil, "【note】"+title)
	if err != nil {
		s.logger.Error("向PDF添加标题失败", "error", err)
	}
	pdf.SetY(pageTopMargin)

	// 4.3写入标注
	for i, mark := range marks {
		if i > 0 {
			// 4.3.1 添加分割线
			s.checkAndAddPage(&pdf, lineHeight)
			pdf.Line(10, pdf.GetY(), 585, pdf.GetY())
			pdf.Br(lineHeight)
		}

		// 4.3.2 写入页码
		if mark.Page != math.MaxInt32 {
			s.checkAndAddPage(&pdf, lineHeight+fontSize)
			err = pdf.Cell(nil, fmt.Sprintf("【第%d页】", mark.Page))
			if err != nil {
				s.logger.Error("向PDF添加页码失败", "error", err)
			}
			pdf.Br(lineHeight)
		}

		// 4.3.3 写入标签
		tags, err := s.pdfService.GetPdfMarkTagsByMarkId(ctx, mark.Id)
		if err != nil {
			s.logger.Error("获取PDF标记的标签失败", "error", err, "markId", mark.Id)
		}
		if len(tags) > 0 {
			s.checkAndAddPage(&pdf, lineHeight+fontSize)
			var tagNames []string
			for _, t := range tags {
				tagNames = append(tagNames, "#"+t.TagName)
			}
			err = pdf.Cell(nil, "【标签】"+strings.Join(tagNames, ","))
			if err != nil {
				s.logger.Error("向PDF添加标签失败", "error", err)
			}
			pdf.Br(lineHeight)
		}

		s.logger.Info("正在处理的标注数据", "mark_id", mark.Id, "pic_url", mark.PicUrl)

		// 4.3.4 原文摘要
		if mark.PicUrl != "" {
			s.checkAndAddPage(&pdf, lineHeight+fontSize)
			err = pdf.Cell(nil, "【原文摘要】")
			if err != nil {
				s.logger.Error("向PDF添加原文摘要头失败", "error", err)
			}
			pdf.Br(lineHeight)
			// 从URL获取图片
			s.logger.Info("开始下载图片", "url", mark.PicUrl)
			resp, err := http.Get(mark.PicUrl)
			if err != nil {
				s.logger.Error("下载图片失败", "error", err, "url", mark.PicUrl)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					err := fmt.Errorf("下载图片HTTP状态码非200: %d", resp.StatusCode)
					s.logger.Error(err.Error(), "url", mark.PicUrl)
					return nil, err
				}
				imgBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					s.logger.Error("读取图片内容失败", "error", err)
				} else {
					s.logger.Info("图片下载成功", "size_bytes", len(imgBytes))
					// 使用 image 标准库来解析图片尺寸
					imgReader := bytes.NewReader(imgBytes)
					imgConfig, _, err := image.DecodeConfig(imgReader)
					if err != nil {
						s.logger.Error("解析图片尺寸失败", "error", err)
					} else {
						imgH, err := gopdf.ImageHolderByBytes(imgBytes)
						if err != nil {
							s.logger.Error("创建图片持有者失败", "error", err)
						} else {
							// 根据固定宽度500，按比例计算高度
							maxWidth := 500.0
							height := float64(imgConfig.Height) * (maxWidth / float64(imgConfig.Width))
							rect := &gopdf.Rect{W: maxWidth, H: height}
							s.checkAndAddPage(&pdf, height+lineHeight)
							currentY := pdf.GetY()
							err = pdf.ImageByHolder(imgH, 10, currentY, rect)
							s.logger.Info("图片处理完成", "rect_height", rect.H)
							if err != nil {
								s.logger.Error("向PDF添加图片失败", "error", err)
							} else {
								pdf.SetY(currentY + rect.H)
								s.logger.Info("PDF Y坐标已更新", "new_y", pdf.GetY())
							}
							pdf.Br(lineHeight) // 图片后增加一些间距
						}
					}
				}
			}
		} else if mark.KeyContent != "" {
			s.checkAndAddPage(&pdf, lineHeight+fontSize)
			err = pdf.Cell(nil, "【原文摘要】")
			if err != nil {
				s.logger.Error("向PDF添加原文摘要头失败", "error", err)
			}
			pdf.Br(lineHeight)
			lines, err := pdf.SplitText(mark.KeyContent, 575)
			if err != nil {
				s.logger.Error("分割关键内容文本失败", "error", err)
			} else {
				for _, line := range lines {
					s.checkAndAddPage(&pdf, lineHeight)
					err = pdf.Cell(nil, line)
					if err != nil {
						s.logger.Error("向PDF添加关键内容失败", "error", err)
						break
					}
					pdf.Br(lineHeight)
				}
			}
			pdf.Br(5)
		}

		// 4.3.5写入批注笔记
		if mark.Idea != "" {
			err = pdf.Cell(nil, "【批注笔记】")
			if err != nil {
				s.logger.Error("向PDF添加批注笔记头失败", "error", err)
			}
			pdf.Br(15)
			lines, err := pdf.SplitText(mark.Idea, 575)
			if err != nil {
				s.logger.Error("分割想法文本失败", "error", err)
			} else {
				for _, line := range lines {
					err = pdf.Cell(nil, line)
					if err != nil {
						s.logger.Error("向PDF添加想法失败", "error", err)
						break
					}
					pdf.Br(15)
				}
			}
			pdf.Br(5) // 在整个块后增加一点额外的间距
		}
	}

	var buf bytes.Buffer
	err = pdf.Write(&buf)
	if err != nil {
		s.logger.Error("将PDF写入缓冲区失败", "error", err)
		return nil, errors.Wrap(err, "生成PDF失败")
	}

	return buf.Bytes(), nil
}

// checkAndAddPage 检查是否需要添加新页面
func (s *PaperNoteService) checkAndAddPage(pdf *gopdf.GoPdf, neededHeight float64) {
	if pdf.GetY()+neededHeight > gopdf.PageSizeA4.H-pageBottomMargin {
		pdf.AddPage()
		err := pdf.SetFont("chinese", "", 14)
		if err != nil {
			s.logger.Error("设置字体失败", "error", err)
		}
		pdf.SetY(pageTopMargin)
	}
}

// GetDownloadNoteMarkMarkdown 生成下载笔记标签Markdown
//  4. 生成Markdown逻辑
//     4.1. 遍历所有标注，构建Markdown字符串。
//     4.2. 对于标注中的图片，先用其原始URL占位。
//     4.3. 解析Markdown文本中的图片URL，将其替换为ZIP包内的相对路径（例如 `images/image_name.png`）。
//     4.4. 将Markdown文本和图片文件打包成一个ZIP压缩文件。
//     4.5. 异步下载这些图片，并将它们添加到ZIP压缩包的 `images` 目录下。
//     4.6. 将生成并处理好的Markdown文件写入ZIP包的根目录。
func (s *PaperNoteService) GetDownloadNoteMarkMarkdown(ctx context.Context, noteId string) ([]byte, error) {
	// 1. 根据noteId查询笔记信息
	note, err := s.GetPaperNoteById(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记信息失败", "error", err)
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}

	if note == nil {
		s.logger.Error("笔记不存在", "noteId", noteId)
		return nil, errors.Biz("note.paper_note.errors.not_found")
	}

	// 2. 获取标题
	doc, err := s.userDocService.GetByPdId(ctx, note.PdfId)
	if err != nil {
		s.logger.Error("获取文献失败", "error", err)
		return nil, errors.Biz("note.paper_pdf.errors.get_failed")
	}
	if doc == nil {
		s.logger.Error("文献不存在", "pdfId", note.PdfId)
		return nil, errors.Biz("note.paper_pdf.errors.not_found")
	}
	title := doc.DocName
	s.logger.Debug("获取文献成功", "title", title)

	// 3. 获取PdfMark标注
	marks, err := s.pdfService.GetPdfMarksByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}
	s.logger.Debug("获取PDF标记成功", "marks", marks)

	// 4. 生成Markdown内容
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("**【note】%s**\n\n", title))

	for i, mark := range marks {
		if i > 0 {
			builder.WriteString("---\n\n")
		}
		if mark.Page != math.MaxInt32 {
			builder.WriteString(fmt.Sprintf("**【第%d页】**\n", mark.Page))
		}

		// 获取标签
		tags, err := s.pdfService.GetPdfMarkTagsByMarkId(ctx, mark.Id)
		if err != nil {
			// 获取标签失败不影响整体流程，仅记录日志
			s.logger.Error("获取PDF标记的标签失败", "error", err, "markId", mark.Id)
		} else if len(tags) > 0 {
			tagNames := make([]string, len(tags))
			for i, tag := range tags {
				tagNames[i] = "#" + tag.TagName
			}
			builder.WriteString(fmt.Sprintf("**【标签】**%s\n", strings.Join(tagNames, ",")))
		}

		if mark.PicUrl != "" {
			builder.WriteString(fmt.Sprintf("**【原文摘要】**\n![image.png](%s)\n", mark.PicUrl))
		} else if mark.KeyContent != "" {
			builder.WriteString(fmt.Sprintf("**【原文摘要】**\n>%s\n<br>\n\n", mark.KeyContent))
		}

		if mark.Idea != "" {
			builder.WriteString(fmt.Sprintf("**【批注笔记】**\n%s\n", mark.Idea))
		}
		builder.WriteString("\n")
	}

	markdownText := builder.String()
	var markDownLinks []MarkDownLink
	markdownText, err = s.replaceUrlsWithRelativePaths(markdownText, &markDownLinks)
	if err != nil {
		// Log the error but return the markdown anyway, as some links might have been processed.
		s.logger.Error("errors occurred while replacing urls with relative paths", "error", err)
	}

	// 4.4, 4.5, 4.6: Package markdown and images into a zip file.
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Add markdown file to zip
	mdName := fmt.Sprintf("【note】%s.md", title)
	mdEntry, err := zipWriter.Create(mdName)
	if err != nil {
		s.logger.Error("failed to create markdown entry in zip", "error", err)
		return nil, errors.Biz("note.download.errors.zip_failed")
	}
	_, err = mdEntry.Write([]byte(markdownText))
	if err != nil {
		s.logger.Error("failed to write markdown to zip", "error", err)
		return nil, errors.Biz("note.download.errors.zip_failed")
	}

	// Asynchronously download images and add to zip
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, link := range markDownLinks {
		wg.Add(1)
		go func(l MarkDownLink) {
			defer wg.Done()
			err := s.downloadAndAddToZip(zipWriter, l, &mu)
			if err != nil {
				// Log error but continue, so the user gets a partial zip file
				s.logger.Error("failed to download or add image to zip", "url", l.Url, "error", err)
			}
		}(link)
	}
	wg.Wait()

	if err := zipWriter.Close(); err != nil {
		s.logger.Error("failed to close zip writer", "error", err)
		return nil, errors.Biz("note.download.errors.zip_failed")
	}

	return buf.Bytes(), nil
}

// MarkDownLink holds the original URL and the relative path for a linked resource.
type MarkDownLink struct {
	Url          string
	RelativePath string
}

// markdownImageRegex is a compiled regular expression to find markdown image links.
var markdownImageRegex = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)

// replaceUrlsWithRelativePaths finds all image URLs in markdown text, replaces them with
// a relative path, and returns a list of the mappings.
func (s *PaperNoteService) replaceUrlsWithRelativePaths(markdownText string, markDownLinks *[]MarkDownLink) (string, error) {
	var firstErr error
	newText := markdownImageRegex.ReplaceAllStringFunc(markdownText, func(match string) string {
		submatches := markdownImageRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match // Should not happen, but as a safeguard
		}
		originalUrl := submatches[1]

		// Skip empty or non-http URLs
		if originalUrl == "" || !strings.HasPrefix(originalUrl, "http") {
			return match
		}

		relativePath, err := s.getRelativePath(originalUrl)
		if err != nil {
			s.logger.Error("failed to get relative path for URL", "url", originalUrl, "error", err)
			if firstErr == nil {
				firstErr = err
			}
			return match // Return original match if path generation fails
		}

		*markDownLinks = append(*markDownLinks, MarkDownLink{Url: originalUrl, RelativePath: relativePath})
		return strings.Replace(match, originalUrl, relativePath, 1)
	})

	return newText, firstErr
}

// getRelativePath creates a relative file path under an 'images' directory from a full URL.
func (s *PaperNoteService) getRelativePath(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	fileName := filepath.Base(u.Path)
	if filepath.Ext(fileName) == "" {
		fileName += ".png" // Assume .png if no extension is present, matching Java logic
	}

	return filepath.ToSlash(filepath.Join("images", fileName)), nil
}

func (s *PaperNoteService) downloadAndAddToZip(zipWriter *zip.Writer, link MarkDownLink, mu *sync.Mutex) error {
	// Get the data from the URL first, before locking.
	resp, err := http.Get(link.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Lock before writing to the zip archive to prevent race conditions.
	mu.Lock()
	defer mu.Unlock()

	// Create a new entry in the zip file for the image.
	zipEntry, err := zipWriter.Create(link.RelativePath)
	if err != nil {
		return err
	}

	// Copy the image data to the zip entry.
	_, err = io.Copy(zipEntry, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *PaperNoteService) GetAllNoteByUserId(ctx context.Context, userId string) ([]model.PaperNote, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetAllNoteByUserId")
	defer span.Finish()

	// 获取用户笔记列表
	notes, err := s.paperNoteDAO.GetAllNoteByUserId(ctx, userId)
	if err != nil {
		s.logger.Error("get paper note failed", "error", err)
		return nil, errors.Biz("note.paper_note.errors.get_failed")
	}

	return notes, nil
}
