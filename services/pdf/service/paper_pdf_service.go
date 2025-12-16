package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/utils"
	"github.com/yb2020/odoc/proto/gen/go/common"
	pb "github.com/yb2020/odoc/proto/gen/go/pdf"
	user_pb "github.com/yb2020/odoc/proto/gen/go/user"
	userDocService "github.com/yb2020/odoc/services/doc/interfaces"
	docModel "github.com/yb2020/odoc/services/doc/model"

	noteInterfaces "github.com/yb2020/odoc/services/note/interfaces"
	paperNoteService "github.com/yb2020/odoc/services/note/service"
	ossConstant "github.com/yb2020/odoc/services/oss/constant"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperService "github.com/yb2020/odoc/services/paper/service"
	"github.com/yb2020/odoc/services/pdf/dao"
	pdfInterfaces "github.com/yb2020/odoc/services/pdf/interfaces"
	"github.com/yb2020/odoc/services/pdf/model"
	userService "github.com/yb2020/odoc/services/user/service"
)

// 确保 PaperPdfService 实现了 PaperPdfServiceInterface 接口
var _ pdfInterfaces.IPaperPdfService = (*PaperPdfService)(nil)

// PaperPdfService 论文PDF服务实现
type PaperPdfService struct {
	paperPDFDAO            *dao.PaperPDFDAO
	logger                 logging.Logger
	ossService             ossService.OssServiceInterface
	tracer                 opentracing.Tracer
	paperAccessService     *paperService.PaperAccessService
	paperService           *paperService.PaperService
	userService            *userService.UserService
	userDocService         userDocService.IUserDocService
	paperNoteService       noteInterfaces.IPaperNoteService
	paperNoteAccessService *paperNoteService.PaperNoteAccessService
	config                 *config.Config
	pdfMarkService         pdfInterfaces.IPdfMarkService
}

// NewPaperPDFService 创建新的论文PDF服务
func NewPaperPdfService(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperPDFDAO *dao.PaperPDFDAO,
	ossService ossService.OssServiceInterface,
	paperAccessService *paperService.PaperAccessService,
	paperService *paperService.PaperService,
	userService *userService.UserService,
	userDocService userDocService.IUserDocService,
	paperNoteService noteInterfaces.IPaperNoteService,
	paperNoteAccessService *paperNoteService.PaperNoteAccessService,
	pdfMarkService pdfInterfaces.IPdfMarkService,
) *PaperPdfService {
	return &PaperPdfService{
		config:                 config,
		paperPDFDAO:            paperPDFDAO,
		logger:                 logger,
		tracer:                 tracer,
		ossService:             ossService,
		paperAccessService:     paperAccessService,
		paperService:           paperService,
		userService:            userService,
		userDocService:         userDocService,
		paperNoteService:       paperNoteService,
		paperNoteAccessService: paperNoteAccessService,
		pdfMarkService:         pdfMarkService,
	}
}

// SetPdfMarkService a setter method to inject the PdfMarkService dependency.
// This is used to break the circular dependency between PaperPdfService and PdfMarkService.
func (s *PaperPdfService) SetPdfMarkService(pdfMarkService pdfInterfaces.IPdfMarkService) {
	s.pdfMarkService = pdfMarkService
}

// SavePaperPDF 保存论文PDF
func (s *PaperPdfService) SavePaperPDF(ctx context.Context, pdf *model.PaperPdf) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.SavePaperPDF")
	defer span.Finish()
	// 创建论文PDF
	return s.paperPDFDAO.Save(ctx, pdf)
}

// 修改论文PDF
func (s *PaperPdfService) ModifyPaperPDF(ctx context.Context, pdf *model.PaperPdf) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.ModifyPaperPDF")
	defer span.Finish()
	// 修改论文PDF
	return s.paperPDFDAO.UpdatePaperPDF(ctx, pdf)
}

// ListPaperPDFs 列出论文PDF
func (s *PaperPdfService) ListPaperPDFs(ctx context.Context, req *pb.ListPaperPDFsRequest) (*pb.ListPaperPDFsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.ListPaperPDFs")
	defer span.Finish()

	// 列出论文PDF
	pdfs, err := s.paperPDFDAO.ListPaperPDFs(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Error("列出论文PDF失败", "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.list_failed")
	}

	// 转换为响应格式
	responsePDFs := make([]*pb.PaperPDF, len(pdfs))
	for i, pdf := range pdfs {
		responsePDFs[i] = s.convertModelToPB(&pdf)
	}

	return &pb.ListPaperPDFsResponse{
		Pdfs: responsePDFs,
	}, nil
}

// CountPaperPDFs 获取论文PDF总数
func (s *PaperPdfService) CountPaperPDFs(ctx context.Context, req *pb.CountPaperPDFsRequest) (*pb.CountPaperPDFsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.CountPaperPDFs")
	defer span.Finish()

	// 获取论文PDF总数
	count, err := s.paperPDFDAO.CountPaperPDFs(ctx)
	if err != nil {
		s.logger.Error("获取论文PDF总数失败", "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.count_failed")
	}

	return &pb.CountPaperPDFsResponse{
		Count: uint64(count),
	}, nil
}

// GetById 根据ID获取论文PDF
func (s *PaperPdfService) GetById(ctx context.Context, id string) (*model.PaperPdf, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.GetById")
	defer span.Finish()

	// 获取论文PDF
	pdf, err := s.paperPDFDAO.GetPaperPDFByID(ctx, id)
	if err != nil {
		s.logger.Error("获取论文PDF失败", "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.get_failed")
	}

	return pdf, nil
}

// GetByIds 根据IDS获取论文PDF列表
func (s *PaperPdfService) GetByIds(ctx context.Context, ids []string) ([]model.PaperPdf, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.GetById")
	defer span.Finish()

	// 获取论文PDF
	pdfs, err := s.paperPDFDAO.FindByIds(ctx, ids)
	if err != nil {
		s.logger.Error("获取论文PDF失败", "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.get_failed")
	}

	return pdfs, nil
}

// GetByPaperId 根据论文ID获取PDF
func (s *PaperPdfService) GetByPaperId(ctx context.Context, paperId string) (*model.PaperPdf, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.GetByPaperId")
	defer span.Finish()

	// 获取论文PDF
	pdf, err := s.paperPDFDAO.GetPaperPDFByPaperID(ctx, paperId)
	if err != nil {
		s.logger.Error("根据论文ID获取PDF失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.get_by_paper_id_failed")
	}

	if pdf == nil {
		return nil, errors.Biz("pdf.paper_pdf.errors.not_found")
	}

	return pdf, nil
}

// GetByFileSHA256 根据FileSHA256获取PDF
func (s *PaperPdfService) GetByFileSHA256(ctx context.Context, fileSHA256 string) (*model.PaperPdf, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.GetByFileSHA256")
	defer span.Finish()

	// 获取论文PDF
	pdf, err := s.paperPDFDAO.GetPaperPDFByFileSHA256(ctx, fileSHA256)
	if err != nil {
		s.logger.Error("根据FileSHA256获取PDF失败", "file_sha256", fileSHA256, "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.get_by_file_sha256_failed")
	}

	if pdf == nil {
		return nil, errors.Biz("pdf.paper_pdf.errors.not_found")
	}

	return pdf, nil
}

func (s *PaperPdfService) GetDefaultPdfByPaperId(ctx context.Context, paperId string) (*model.PaperPdf, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.GetDefaultPdfByPaperId")
	defer span.Finish()

	// 获取默认PDF
	pdf, err := s.paperPDFDAO.GetDefaultPdfByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("根据论文ID获取默认PDF失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.get_default_by_paper_id_failed")
	}

	if pdf == nil {
		return nil, errors.Biz("pdf.paper_pdf.errors.not_found")
	}

	return pdf, nil
}

// GetPdfUrlById 获取PDF文件地址, duration: URL有效期(秒)
func (s *PaperPdfService) GetPdfUrlById(ctx context.Context, pdfId string, duration int) (*string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPDFService.GetPdfUrlById")
	defer span.Finish()

	// 获取PDF
	pdf, err := s.paperPDFDAO.GetPaperPDFByID(ctx, pdfId)
	if err != nil {
		s.logger.Error("获取论文PDF失败", "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.get_failed")
	}

	if pdf == nil {
		return nil, errors.Biz("pdf.paper_pdf.errors.not_found")
	}
	// 获取PDF文件地址
	fileUrl, err := s.ossService.GetFileTemporaryURL(ctx, ossConstant.BucketTypeToEnum(s.config, pdf.OssBucketName), pdf.OssObjectKey, utils.GetIntPtrValue(&duration, 60))
	if err != nil {
		s.logger.Error("获取PDF文件地址失败", "error", err)
		return nil, errors.Biz("pdf.paper_pdf.errors.get_url_failed")
	}
	return &fileUrl, nil
}

// Pdf权限认证是否许可
func (s *PaperPdfService) AuthPermissionDenied(ctx context.Context, pdfId string, userId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.AuthPermissionDenied")
	defer span.Finish()

	pdf, err := s.paperPDFDAO.GetPaperPDFByID(ctx, pdfId)
	if err != nil {
		s.logger.Error("获取论文PDF失败", "error", err)
		return false, errors.Biz("pdf.paper_pdf.errors.get_failed")
	}
	if pdf == nil {
		return false, errors.Biz("pdf.paper_pdf.errors.not_found")
	}

	if pdf.CreatorId == "" || pdf.CreatorId == userId {
		return false, nil
	}

	// TODO:判断用户是否有权限打开某论文
	// accessPermision, _ := s.paperAccessService.IsUserHasPaperAccess(ctx, userId, pdf.Id)
	// if accessPermision {
	// 	return false, nil
	// }

	return true, nil
}

func (s *PaperPdfService) GetPdfStatusInfo(ctx context.Context, pdfId string, noteId string) (*pb.GetPdfStatusInfoResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.getPdfStatusInfo")
	defer span.Finish()

	paperPdf, err := s.GetById(ctx, pdfId)
	if err != nil {
		s.logger.Error("msg", "获取PDF状态信息失败", "error", err.Error())
		return nil, errors.Biz("pdf.paper_pdf.errors.get_failed")
	}
	if paperPdf == nil {
		return nil, errors.Biz("pdf.paper_pdf.errors.not_found")
	}

	pdfStatusInfo := &pb.GetPdfStatusInfoResponse{}
	pdfUrl, _ := s.GetPdfUrlById(ctx, paperPdf.Id, 60)
	pdfStatusInfo.PdfUrl = *pdfUrl

	//设置论文相关信息
	if paperPdf.PaperId != "" {
		//MOCK:TODO
		basePaperInfo, err := s.paperService.GetPaperBaseInfoById(ctx, paperPdf.PaperId)
		if err != nil {
			s.logger.Error("msg", "获取论文基本信息失败", "error", err.Error())
			return nil, errors.Biz("pdf.paper_pdf.errors.get_failed")
		}
		if basePaperInfo != nil {
			pdfStatusInfo.PaperId = paperPdf.PaperId
			pdfStatusInfo.PaperTitle = basePaperInfo.PaperTitle

			defaultPdf, err := s.GetDefaultPdfByPaperId(ctx, paperPdf.PaperId)
			if err != nil {
				s.logger.Error("msg", "获取默认PDF失败", "error", err.Error())
				return nil, errors.Biz("pdf.paper_pdf.errors.get_failed")
			}
			if defaultPdf != nil {
				pdfStatusInfo.AuthPdfId = defaultPdf.Id
			}
		}
	}

	//如果pdf为私密,返回该pdf上传者信息
	if paperPdf.CreatorId != "" {
		//获取作者信息 TODO
		authInfo, err := s.userService.GetAuthInfo(ctx, paperPdf.CreatorId)
		if err != nil {
			s.logger.Error("msg", "获取作者信息失败", "error", err.Error())
			return nil, errors.Biz("pdf.paper_pdf.errors.get_failed")
		}
		if authInfo != nil {
			// 将bean.AuthorBean转换为proto的user.AuthorBean
			pbAuthorBean := &user_pb.AuthorBean{
				Id:               authInfo.Id,
				NickName:         authInfo.NickName,
				Description:      authInfo.Description,
				Usn:              authInfo.Usn,
				UsnCanModify:     authInfo.UsnCanModify,
				Mobile:           authInfo.Mobile,
				Email:            authInfo.Email,
				IsWxPublicBind:   authInfo.IsWxPublicBind,
				Self:             authInfo.Self,
				ShowName:         authInfo.ShowName,
				AvatarUrl:        authInfo.AvatarUrl,
				Tags:             authInfo.Tags,
				AuthorId:         authInfo.AuthorId,
				AuthorName:       authInfo.AuthorName,
				IsAuthentication: authInfo.IsAuthentication,
				IsCert:           authInfo.IsCert,
				IsPaperAuthor:    authInfo.IsPaperAuthor,
				Profession:       authInfo.Profession,
				ResearchField:    authInfo.ResearchField,
				SchoolCompany:    authInfo.SchoolCompany,
			}

			// 如果有封禁信息，也进行转换
			if authInfo.BanInfo != nil {
				pbAuthorBean.BanInfo = &user_pb.UserBanInfo{
					BanFlag:    authInfo.BanInfo.BanFlag,
					BanReason:  authInfo.BanInfo.BanReason,
					BanRemark:  authInfo.BanInfo.BanRemark,
					BanEndTime: authInfo.BanInfo.BanEndTime,
				}
			}

			pdfStatusInfo.PdfOwnerInfo = pbAuthorBean
		}
	}

	// 获取用户ID
	userId, _ := userContext.GetUserID(ctx)
	// 4.df用户态
	if userId != "" {
		// 如果用户是创建者，设置为所有者状态
		if userId == paperPdf.CreatorId {
			pdfStatusInfo.PdfUserStatus = pb.UserStatusEnum_OWNER
		} else {
			// 其他登录用户为访客状态
			pdfStatusInfo.PdfUserStatus = pb.UserStatusEnum_GUEST
		}
	} else {
		// 未登录用户为游客状态
		pdfStatusInfo.PdfUserStatus = pb.UserStatusEnum_TOURIST
	}

	//设置userDoc
	//5. 设置UserDoc
	var userDoc *docModel.UserDoc
	if paperPdf.PaperId != "" {
		userDoc, err = s.userDocService.GetByUserIDAndPaperID(ctx, userId, paperPdf.PaperId)
		if err != nil {
			s.logger.Error("msg", "获取UserDoc失败", "error", err.Error())
			return nil, errors.Biz("pdf.paper_pdf.errors.get_user_doc_failed")
		}
	} else {
		userDoc, err = s.userDocService.GetByUserIdAndPdfId(ctx, userId, paperPdf.Id)
		if err != nil {
			s.logger.Error("msg", "获取UserDoc失败", "error", err.Error())
			return nil, errors.Biz("pdf.paper_pdf.errors.get_user_doc_failed")
		}
	}
	if userDoc != nil {
		pdfStatusInfo.DocName = userDoc.DocName
		pdfStatusInfo.PaperTitle = userDoc.PaperTitle
	}

	// 6.设置笔记信息
	if noteId != "" {
		paperNote, err := s.paperNoteService.GetPaperNoteById(ctx, noteId)
		if err != nil {
			s.logger.Error("msg", "获取笔记失败", "error", err.Error())
			return nil, errors.Biz("pdf.paper_pdf.errors.get_note_failed")
		}
		s.logger.Info("msg", "获取笔记成功", "note", paperNote)
		paperNoteAccess, err := s.paperNoteAccessService.GetPaperNoteAccessByNoteId(ctx, noteId)
		if err != nil {
			s.logger.Error("msg", "获取笔记访问记录失败", "error", err.Error())
			return nil, errors.Biz("pdf.paper_pdf.errors.get_note_access_failed")
		}
		//设置笔记是否公开
		if paperNoteAccess != nil {
			pdfStatusInfo.NoteOpenAccessFlag = paperNoteAccess.OpenStatus
		} else {
			pdfStatusInfo.NoteOpenAccessFlag = false
		}

		//设置笔记用户态
		if paperNote != nil {
			if userId == "" {
				pdfStatusInfo.NoteUserStatus = pb.UserStatusEnum_TOURIST
			} else if userId == paperNote.CreatorId {
				pdfStatusInfo.NoteUserStatus = pb.UserStatusEnum_OWNER
			} else {
				pdfStatusInfo.NoteUserStatus = pb.UserStatusEnum_GUEST
			}
			//查看他人笔记，未入库的PDF不显示标题 TODO

		}
	}

	// 7.设置pdf访问权限
	if userId == "" {
		pdfStatusInfo.HasPdfAccessFlag = pdfStatusInfo.AuthPdfId != "0"
	} else {
		if paperPdf.PaperId == "" {
			//无paperId
			pdfStatusInfo.HasPdfAccessFlag = func() bool {
				r, _ := s.userDocService.GetByUserIdAndPdfId(ctx, userId, paperPdf.Id)
				return r != nil
			}()
		} else {
			//有paperId
			if pdfStatusInfo.AuthPdfId != "0" {
				pdfStatusInfo.HasPdfAccessFlag = true
			} else {
				//用户在该paperid下有一篇文献是否有权限 MOCK:TODO
				pdfStatusInfo.HasPdfAccessFlag = false
			}
		}
	}

	// 8.处理特殊权限 TODO
	// //8.1 管理员能看见任何人的笔记内容
	// //8.2 笔记如果公开，任何人有PDF访问权限
	if pdfStatusInfo.NoteOpenAccessFlag {
		pdfStatusInfo.HasPdfAccessFlag = true
	}
	// //8.3 nacos配置bug开关，任何人有PDF访问权限

	return pdfStatusInfo, nil
}

// 转换模型到PB
func (s *PaperPdfService) convertModelToPB(pdf *model.PaperPdf) *pb.PaperPDF {
	return &pb.PaperPDF{
		Id:         pdf.Id,
		AppId:      pdf.AppId,
		PaperId:    pdf.PaperId,
		Creator:    pdf.Creator,
		CreatorId:  pdf.CreatorId,
		Modifier:   pdf.Modifier,
		ModifierId: pdf.ModifierId,
		FileSHA256: pdf.FileSHA256,
		ParseCount: uint32(pdf.ParseCount),
		Size:       float32(pdf.Size),
		PageCount:  uint32(pdf.PageCount),
		// Language:          pdf.Language,
		// OssBucketName:     pdf.OssBucketName,
		// OssObjectKey:      pdf.OssObjectKey,
		CreatedAt: uint64(pdf.CreatedAt.Unix()),
		UpdatedAt: uint64(pdf.UpdatedAt.Unix()),
	}
}

// GetVisitablePdfIdByPaperIdAndUserId 获取用户可做笔记的pdfId
// 对应Java版本的getVisitablePdfIdByPaperIdAndUserId方法
func (s *PaperPdfService) GetVisitablePdfIdByPaperIdAndUserId(ctx context.Context, paperId string, userId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.GetVisitablePdfIdByPaperIdAndUserId")
	defer span.Finish()

	// 首先尝试从paperNoteService获取用户对该论文的笔记
	paperNote, err := s.paperNoteService.GetByPaperIdAndUserIdOrderByNoteCount(ctx, paperId, userId)
	if err != nil {
		s.logger.Error("获取用户笔记失败", "error", err.Error())
		// 继续处理，不返回错误
	}
	// 如果找到了笔记并且笔记有pdfId，直接返回
	if paperNote != nil && paperNote.PdfId != "" {
		return paperNote.PdfId, nil
	}
	// 否则尝试获取用户可访问的PDF
	paperPdf, err := s.getVisitablePdfByPaperIdAndUserId(ctx, paperId, userId)
	if err != nil {
		s.logger.Error("获取用户可访问的PDF失败", "error", err.Error())
		return "0", err
	}
	// 如果找到了可访问的PDF，返回其ID
	if paperPdf != nil {
		return paperPdf.Id, nil
	}
	// 如果都没找到，返回0
	return "0", nil
}

// getVisitablePdfByPaperIdAndUserId 获取用户可访问的PDF
// 对应Java版本的getVisitablePdfByPaperIdAndUserId方法
func (s *PaperPdfService) getVisitablePdfByPaperIdAndUserId(ctx context.Context, paperId string, userId string) (*model.PaperPdf, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.getVisitablePdfByPaperIdAndUserId")
	defer span.Finish()

	// 调用DAO层方法获取用户可访问的PDF
	return s.paperPDFDAO.GetVisitablePdfIdByPaperIdAndUserId(ctx, paperId, userId)
}

// GetByUserIdAndFileSHA256 根据用户ID和FileSHA256获取PDF
func (s *PaperPdfService) GetByUserIdAndFileSHA256(ctx context.Context, userId string, fileSHA256 string) (*model.PaperPdf, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.GetByUserIdAndFileSHA256")
	defer span.Finish()

	// 调用DAO层方法获取PDF记录
	return s.paperPDFDAO.GetByUserIdAndFileSHA256(ctx, userId, fileSHA256)
}

// GetDocById 根据PdfId获取文献信息
func (s *PaperPdfService) GetDocById(ctx context.Context, id string) (*docModel.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteService.GetDocById")
	defer span.Finish()

	pdf, err := s.GetById(ctx, id)
	if err != nil {
		s.logger.Error("获取pdf失败", "error", err)
		return nil, errors.Biz("note.paper_pdf.errors.get_failed")
	}

	if pdf == nil {
		return nil, errors.Biz("note.paper_pdf.errors.not_found")
	}

	return s.userDocService.GetByPdId(ctx, pdf.Id)
}

// GetMarkTagInfosByFolderId 根据folderId获取用户markTagInfos
func (s *PaperPdfService) GetMarkTagInfosByFolderId(ctx context.Context, userId string, folderId string) ([]*common.PdfMarkTagInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.GetMarkTagInfosByFolderId")
	defer span.Finish()

	return s.pdfMarkService.GetMarkTagInfosByFolderId(ctx, userId, folderId)
}

func (s *PaperPdfService) GetPdfFileTotalSize(ctx context.Context, userId string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.GetPdfFileTotalSize")
	defer span.Finish()

	pdfs, err := s.paperPDFDAO.GetUserPaperPdfList(ctx, userId)
	if err != nil {
		s.logger.Error("获取用户pdf失败", "error", err)
		return 0, errors.Biz("pdf.paper_pdf.errors.get_failed")
	}
	totalSize := int64(0)
	for _, pdf := range pdfs {
		totalSize += pdf.Size
	}
	return totalSize, nil
}

func (s *PaperPdfService) GetCountPdfMarksByNoteId(ctx context.Context, noteId string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.GetCountPdfMarksByNoteId")
	defer span.Finish()

	// 获取PDF标记
	count, err := s.pdfMarkService.GetCountPdfMarksByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return 0, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}

	return count, nil
}

// GetPdfMarksByNoteId 根据笔记ID获取PDF标记
func (s *PaperPdfService) GetPdfMarksByNoteId(ctx context.Context, noteId string) ([]model.PdfMark, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.GetPdfMarksByNoteId")
	defer span.Finish()

	// 获取PDF标记
	marks, err := s.pdfMarkService.GetPdfMarksByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}

	return marks, nil
}

// GetPdfMarkTagsByMarkId 根据markId获取PDF标记标签信息
func (s *PaperPdfService) GetPdfMarkTagsByMarkId(ctx context.Context, markId string) ([]*common.AnnotateTag, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfService.GetPdfMarkTagsByMarkId")
	defer span.Finish()

	// 获取PDF标记标签
	tags, err := s.pdfMarkService.GetAnnotateTagsByMarkId(ctx, markId)
	if err != nil {
		s.logger.Error("获取PDF标记标签失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	return tags, nil
}
