package service

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	docPb "github.com/yb2020/odoc-proto/gen/go/doc"
	pb "github.com/yb2020/odoc-proto/gen/go/paper"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/bean"
	"github.com/yb2020/odoc/services/paper/constants"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperService 论文服务实现
type PaperService struct {
	paperDAO *dao.PaperDAO
	logger   logging.Logger
	tracer   opentracing.Tracer
}

// NewPaperService 创建新的论文服务
func NewPaperService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperDAO *dao.PaperDAO,
) *PaperService {
	return &PaperService{
		logger:   logger,
		tracer:   tracer,
		paperDAO: paperDAO,
	}
}

// CreatePaper 创建论文
func (s *PaperService) SavePaper(ctx context.Context, paper *model.Paper) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.CreatePaper")
	defer span.Finish()
	if err := s.paperDAO.Save(ctx, paper); err != nil {
		s.logger.Error("创建论文失败", "error", err)
		return errors.Biz("paper.errors.create_failed")
	}
	return nil
}

// 修改论文
func (s *PaperService) ModifyPaper(ctx context.Context, paper *model.Paper) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.ModifyPaper")
	defer span.Finish()
	if err := s.paperDAO.Modify(ctx, paper); err != nil {
		s.logger.Error("修改论文失败", "error", err)
		return errors.Biz("paper.errors.modify_failed")
	}
	return nil
}

// GetPaperById 根据ID获取论文
func (s *PaperService) GetPaperById(ctx context.Context, id string) (*model.Paper, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.GetPaperById")
	defer span.Finish()

	paper, err := s.paperDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文失败", "id", id, "error", err)
		return nil, errors.Biz("paper.errors.get_failed")
	}

	if paper == nil {
		return nil, errors.Biz("paper.errors.not_found")
	}

	return paper, nil
}

// GetPaperByPaperId 根据论文ID获取论文
func (s *PaperService) GetPaperByPaperId(ctx context.Context, paperId string) (*model.Paper, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.GetPaperByPaperId")
	defer span.Finish()

	paper, err := s.paperDAO.FindByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("paper.errors.get_failed")
	}

	if paper == nil {
		return nil, errors.Biz("paper.errors.not_found")
	}

	return paper, nil
}

// GetPapersByOwnerId 根据拥有者ID获取论文列表
func (s *PaperService) GetPapersByOwnerId(ctx context.Context, ownerId string) ([]*model.Paper, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.GetPapersByOwnerId")
	defer span.Finish()

	papers, err := s.paperDAO.FindByOwnerId(ctx, ownerId)
	if err != nil {
		s.logger.Error("获取论文列表失败", "owner_id", ownerId, "error", err)
		return nil, errors.Biz("paper.errors.list_failed")
	}

	result := make([]*model.Paper, 0, len(papers))
	for i := range papers {
		result = append(result, &papers[i])
	}

	return result, nil
}

// UpdatePaper 更新论文
func (s *PaperService) UpdatePaper(ctx context.Context, paper *model.Paper) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.UpdatePaper")
	defer span.Finish()

	// 获取论文
	existingPaper, err := s.paperDAO.FindById(ctx, paper.Id)
	if err != nil {
		s.logger.Error("获取论文失败", "id", paper.Id, "error", err)
		return false, errors.Biz("paper.errors.get_failed")
	}

	if existingPaper == nil {
		return false, errors.Biz("paper.errors.not_found")
	}

	// 更新论文
	if err := s.paperDAO.UpdateById(ctx, paper); err != nil {
		s.logger.Error("更新论文失败", "id", paper.Id, "error", err)
		return false, errors.Biz("paper.errors.update_failed")
	}

	return true, nil
}

// DeletePaper 删除论文
func (s *PaperService) DeletePaper(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.DeletePaper")
	defer span.Finish()

	// 删除论文
	if err := s.paperDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文失败", "id", id, "error", err)
		return false, errors.Biz("paper.errors.delete_failed")
	}

	return true, nil
}

// ListPapers 列出论文
func (s *PaperService) ListPapers(ctx context.Context, page, size int) ([]*model.Paper, int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.ListPapers")
	defer span.Finish()

	// 获取论文总数
	total, err := s.paperDAO.Count(ctx)
	if err != nil {
		s.logger.Error("获取论文总数失败", "error", err)
		return nil, 0, errors.Biz("paper.errors.count_failed")
	}

	if total == 0 {
		return []*model.Paper{}, 0, nil
	}

	// 获取论文列表
	papers, err := s.paperDAO.List(ctx, size, (page-1)*size)
	if err != nil {
		s.logger.Error("获取论文列表失败", "error", err)
		return nil, 0, errors.Biz("paper.errors.list_failed")
	}

	result := make([]*model.Paper, 0, len(papers))
	for i := range papers {
		result = append(result, &papers[i])
	}

	return result, total, nil
}

// GetPaperBaseInfoById 根据ID获取论文基本信息 MOCK:TODO
func (s *PaperService) GetPaperBaseInfoById(ctx context.Context, id string) (*bean.PaperBaseInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.GetPaperBaseInfo")
	defer span.Finish()

	paper, err := s.paperDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文基本信息失败", "id", id, "error", err)
		return nil, errors.Biz("paper.errors.get_failed")
	}

	if paper == nil {
		return nil, errors.Biz("paper.errors.not_found")
	}

	return &bean.PaperBaseInfo{
		PaperId:    paper.Id,
		PaperTitle: "paper.Title",
	}, nil
}

// 论文详情页-获取论文信息
func (s *PaperService) GetPaperDetailInfo(ctx context.Context, paperId string) (*pb.PaperDetailInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.GetPaperDetailInfo")
	defer span.Finish()

	response := &pb.PaperDetailInfo{
		HasPublicPdf:        true,
		HasPrivatePdf:       false,
		ReadStatus:          pb.ReadStatusEnum_READABLE_PUBLIC,
		ShowPaperDetail:     true,
		PaperCommentedCount: 0, // 评论数
		NewPaper:            true,
	}

	paper, err := s.paperDAO.FindByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("paper.errors.get_failed")
	}
	if paper == nil {
		return response, nil
	}
	//设置paper信息
	response.PaperId = paper.Id
	response.Title = paper.Title
	response.PublishDate = paper.PublishDate
	response.OriginalAbstract = paper.Abstract

	authorList := []*pb.BasePaperAuthorInfo{}
	if paper.Authors != "" {
		//json反序列化
		authors := []*docPb.AuthorInfo{}
		err := json.Unmarshal([]byte(paper.Authors), &authors)
		if err != nil {
			s.logger.Error("parse author list failed", "error", err.Error())
			return nil, err
		}
		for _, author := range authors {
			authorList = append(authorList, &pb.BasePaperAuthorInfo{
				Name: author.Literal,
			})
		}
	}
	response.AuthorList = authorList
	//TODO: 标签数据缺失
	//设置笔记数
	// response.NoteCount = paper.NoteCount
	//这里只可能是1，因为这里的论文全部都是私有的，后续如果更改了，需要重新设计
	response.CollectCount = 1
	//设置问题数
	// response.QuestionAnswerCount = 0
	//需要补充
	// response.PdfId =
	// response.Doi =
	// response.Venues =
	// response.PrimaryVenue =
	// response.VenueTags =
	//论文精读
	// response.PaperResource =
	//摘要翻译   暂无
	//补充机器翻译 暂无
	//同行评审数量  暂无

	return response, nil
}

// isPrivatePaper 判断论文是否为私有
func (s *PaperService) IsPrivatePaper(ctx context.Context, paperId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.isPrivatePaper")
	defer span.Finish()

	paper, err := s.paperDAO.FindByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文失败", "paper_id", paperId, "error", err)
		return false, errors.Biz("paper.errors.get_failed")
	}

	if paper == nil {
		return false, nil
	}

	return constants.PaperStatus(paper.Status) == constants.Private, nil
}

// GetPaperVersions 获取论文版本列表
func (s *PaperService) GetPaperVersions(ctx context.Context, paperId string) (*pb.GetPaperVersionsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperService.GetPaperVersions")
	defer span.Finish()

	paper, err := s.paperDAO.FindById(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("paper.errors.get_failed")
	}
	if paper == nil {
		return nil, errors.Biz("param error ! paper not found")
	}
	privatePaperVersionInfo := &pb.PaperVersionInfo{
		Type:        pb.PaperVersionType_PAPER_PDF,
		Name:        "用户上传",
		JumpUrl:     "",
		CurVersion:  true,
		LastVersion: true,
		DatePrefix:  paper.PublishDate,
	}

	versionResponse := &pb.GetPaperVersionsResponse{
		PrivateVersions: []*pb.PaperVersionInfo{privatePaperVersionInfo},
	}
	return versionResponse, nil
}
