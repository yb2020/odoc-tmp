package service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/proto/gen/go/translate"
	"github.com/yb2020/odoc/services/translate/dao"
	"github.com/yb2020/odoc/services/translate/model"
)

// GlossaryService 默认的术语库服务实现
type GlossaryService struct {
	glossaryDAO dao.GlossaryDAO
	config      *config.Config
	logger      logging.Logger
	tracer      opentracing.Tracer
}

// NewGlossaryService 创建一个新的术语库服务
func NewGlossaryService(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	glossaryDAO dao.GlossaryDAO,
) *GlossaryService {
	return &GlossaryService{
		glossaryDAO: glossaryDAO,
		config:      config,
		logger:      logger,
		tracer:      tracer,
	}
}

// 检查文本是否合法
func (s *GlossaryService) checkText(texts ...string) error {
	for _, text := range texts {
		if text == "" {
			return errors.Biz("translate.glossary.errors.empty_text")
		}
		if len(text) > s.config.Translate.Glossary.MaxLength {
			return errors.Biz("translate.glossary.errors.text_too_long")
		}
	}
	return nil
}

// 检查是否有重复的术语
func (s *GlossaryService) checkDuplicate(ctx context.Context, userID string, matchCase bool, originalText string, ignoreID *string) error {
	glossaries, err := s.glossaryDAO.GetGlossariesByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("msg", "获取用户术语列表失败", "error", err.Error())
		return errors.BizWrap("translate.glossary.errors.invalid_data", err)
	}

	if len(glossaries) > s.config.Translate.Glossary.MaxSize {
		return errors.Biz("translate.glossary.errors.max_size_exceeded")
	}

	for _, g := range glossaries {
		// 如果是更新操作，跳过自身
		if ignoreID != nil && g.Id == *ignoreID {
			continue
		}

		// 根据是否区分大小写进行比较
		var match bool
		if matchCase {
			match = g.OriginalText == originalText
		} else {
			match = strings.EqualFold(g.OriginalText, originalText)
		}

		if match {
			return errors.Biz("translate.glossary.errors.entry_already_exists")
		}
	}

	return nil
}

// AddGlossary 添加术语条目
func (s *GlossaryService) AddGlossary(ctx context.Context, userId string, req *translate.AddGlossaryReq) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.AddGlossary")
	defer span.Finish()

	// 检查文本
	if err := s.checkText(req.OriginalText, req.TranslationText); err != nil {
		return "0", err
	}

	// 去除前后空格
	originalText := strings.TrimSpace(req.OriginalText)

	// 检查重复
	if err := s.checkDuplicate(ctx, userId, req.MatchCase, originalText, nil); err != nil {
		return "0", err
	}

	// 创建术语条目
	glossary := &model.Glossary{
		UserId:          userId,
		OriginalText:    originalText,
		TranslationText: req.TranslationText,
		MatchCase:       req.MatchCase,
		Ignored:         req.Ignored,
		IsPublic:        false, // 默认不公开
	}

	// 设置创建时间和修改者
	now := time.Now()
	glossary.CreatedAt = now
	glossary.UpdatedAt = now
	glossary.CreatorId = userId
	glossary.ModifierId = userId

	// 保存到数据库
	if err := s.glossaryDAO.Save(ctx, glossary); err != nil {
		s.logger.Error("msg", "创建术语条目失败", "error", err.Error())
		return "0", errors.BizWrap("translate.glossary.errors.entry_create_failed", err)
	}

	return glossary.Id, nil
}

// UpdateGlossary 更新术语条目
func (s *GlossaryService) UpdateGlossary(ctx context.Context, userId string, req *translate.UpdateGlossaryReq) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.UpdateGlossary")
	defer span.Finish()

	// 检查参数
	if req == nil {
		return errors.Biz("translate.glossary.errors.invalid_request")
	}

	// 检查文本是否合法
	if err := s.checkText(req.OriginalText, req.TranslationText); err != nil {
		return err
	}

	// 去除前后空格
	originalText := strings.TrimSpace(req.OriginalText)

	// 查询术语条目
	glossary, err := s.glossaryDAO.FindExistById(ctx, req.Id)
	if err != nil {
		s.logger.Error("msg", "查询术语条目失败", "error", err.Error())
		return errors.BizWrap("translate.glossary.errors.entry_not_found", err)
	}

	// 检查是否存在
	if glossary == nil {
		return errors.Biz("translate.glossary.errors.entry_not_found")
	}

	// 检查权限
	if glossary.UserId != userId {
		return errors.Biz("translate.glossary.errors.permission_denied")
	}

	// 检查重复
	id := req.Id
	if err := s.checkDuplicate(ctx, userId, req.MatchCase, originalText, &id); err != nil {
		return err
	}

	// 更新术语条目
	glossary.OriginalText = originalText
	glossary.TranslationText = req.TranslationText
	glossary.MatchCase = req.MatchCase
	glossary.Ignored = req.Ignored
	glossary.UpdatedAt = time.Now()
	glossary.ModifierId = userId

	// 保存到数据库
	if err := s.glossaryDAO.ModifyExcludeNull(ctx, glossary); err != nil {
		s.logger.Error("msg", "更新术语条目失败", "error", err.Error())
		return errors.BizWrap("translate.glossary.errors.entry_update_failed", err)
	}

	return nil
}

// DeleteGlossary 删除术语条目
func (s *GlossaryService) DeleteGlossary(ctx context.Context, userId string, id string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.DeleteGlossary")
	defer span.Finish()

	// 查询术语条目
	glossary, err := s.glossaryDAO.FindExistById(ctx, id)
	if err != nil {
		s.logger.Error("msg", "查询术语条目失败", "error", err.Error())
		return errors.BizWrap("translate.glossary.errors.entry_not_found", err)
	}

	// 检查是否存在
	if glossary == nil {
		return errors.Biz("translate.glossary.errors.entry_not_found")
	}

	// 检查权限
	if glossary.UserId != userId {
		return errors.Biz("translate.glossary.errors.permission_denied")
	}

	// 删除术语条目
	if err := s.glossaryDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("msg", "删除术语条目失败", "error", err.Error())
		return errors.BizWrap("translate.glossary.errors.entry_delete_failed", err)
	}

	return nil
}

// GetGlossaryList 获取术语条目列表
func (s *GlossaryService) GetGlossaryList(ctx context.Context, userID string, req *translate.GetGlossaryListReq) (*translate.GetGlossaryListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.GetGlossaryList")
	defer span.Finish()

	// 获取术语条目列表
	glossaries, err := s.glossaryDAO.GetGlossariesByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("msg", "获取术语条目列表失败", "error", err.Error())
		return nil, errors.BizWrap("translate.glossary.errors.invalid_data", err)
	}

	// 根据创建时间排序（默认降序）
	sort.Slice(glossaries, func(i, j int) bool {
		return glossaries[i].CreatedAt.After(glossaries[j].CreatedAt)
	})

	// 根据searchText进行过滤
	if req.SearchText != nil {
		glossaries = s.filterGlossariesBySearchText(glossaries, *req.SearchText)
	}

	// 计算总数
	total := int32(len(glossaries))

	// 分页
	start := int(req.CurrentPage * req.PageSize)
	if start >= len(glossaries) {
		start = 0
		req.CurrentPage = 0
	}
	end := start + int(req.PageSize)
	if end > len(glossaries) {
		end = len(glossaries)
	}

	// 分页后的数据
	pagedGlossaries := glossaries[start:end]

	// 转换为响应对象
	items := make([]*translate.GlossaryItem, 0, len(pagedGlossaries))
	for _, g := range pagedGlossaries {
		item := &translate.GlossaryItem{
			Id:              g.Id,
			OriginalText:    g.OriginalText,
			TranslationText: g.TranslationText,
			MatchCase:       g.MatchCase,
			Ignored:         g.Ignored,
		}
		items = append(items, item)
	}

	// 构建响应
	resp := &translate.GetGlossaryListResponse{
		Total: uint32(total),
		Items: items,
	}

	return resp, nil
}

func (s *GlossaryService) filterGlossariesBySearchText(glossaries []model.Glossary, searchText string) []model.Glossary {
	var filtered []model.Glossary
	for _, g := range glossaries {
		if strings.Contains(g.OriginalText, searchText) {
			filtered = append(filtered, g)
		}
	}
	return filtered
}

// ReplaceOriginalText 替换原文中的术语
func (s *GlossaryService) ReplaceOriginalText(ctx context.Context, userId string, content string) (*model.GlossaryTranslateModel, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.ReplaceOriginalText")
	defer span.Finish()

	// 创建返回对象
	glossaryBean := &model.GlossaryTranslateModel{
		UserSelectedText: content,
		RelInfos:         make([]model.RelInfo, 0),
	}

	// 获取用户的所有术语条目
	glossaries, err := s.glossaryDAO.GetGlossariesByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取用户术语条目列表失败", "userId", userId, "error", err.Error())
		return glossaryBean, fmt.Errorf("获取用户术语条目列表失败: %w", err)
	}

	// 按照匹配规则排序：1.区分大小写优先 2.原文长度降序 3.创建时间降序
	sort.Slice(glossaries, func(i, j int) bool {
		if glossaries[i].MatchCase != glossaries[j].MatchCase {
			return glossaries[i].MatchCase // 区分大小写的排在前面
		}
		if len(glossaries[i].OriginalText) != len(glossaries[j].OriginalText) {
			return len(glossaries[i].OriginalText) > len(glossaries[j].OriginalText) // 原文长的排在前面
		}
		return glossaries[i].CreatedAt.After(glossaries[j].CreatedAt) // 创建时间晚的排在前面
	})

	// 构建结果
	result := content

	// 遍历排序后的术语条目
	for _, glossary := range glossaries {
		// 构建特殊替换文本
		idStr := fmt.Sprintf("%d", glossary.Id)
		specialText := "rp" + idStr[len(idStr)-4:] + "sp"

		// 构建正则表达式
		var regex string
		var translationText string

		if glossary.Ignored {
			translationText = glossary.OriginalText
		} else {
			translationText = glossary.TranslationText
		}

		if glossary.MatchCase {
			// 区分大小写，精确匹配
			// regex = fmt.Sprintf("(?<![\\p{L}-])%s(?![\\p{L}-])", regexp.QuoteMeta(glossary.OriginalText))
			// 替换
			// re := regexp.MustCompile(regex)
			// result = re.ReplaceAllString(result, specialText)
			regex = fmt.Sprintf(`(^|[^\p{L}-])(%s)($|[^\p{L}-])`, regexp.QuoteMeta(glossary.OriginalText))
			re := regexp.MustCompile(regex)
			result = re.ReplaceAllString(result, "${1}"+specialText+"${3}")
		} else {
			// 不区分大小写
			// regex = fmt.Sprintf("(?<![\\p{L}-])%s(?![\\p{L}-])", regexp.QuoteMeta(strings.ToLower(glossary.OriginalText)))
			// // 替换
			// re := regexp.MustCompile(regex)
			// result = re.ReplaceAllString(strings.ToLower(result), specialText)
			regex = fmt.Sprintf(`(?i)(^|[^\p{L}-])(%s)($|[^\p{L}-])`, regexp.QuoteMeta(glossary.OriginalText))
			re := regexp.MustCompile(regex)
			result = re.ReplaceAllString(result, "${1}"+specialText+"${3}")
		}

		// 添加关联信息
		relInfo := model.RelInfo{
			Id:           glossary.Id,
			Key:          specialText,
			Translation:  translationText,
			OriginalText: glossary.OriginalText,
			MatchCase:    glossary.MatchCase,
		}
		glossaryBean.RelInfos = append(glossaryBean.RelInfos, relInfo)
	}

	// 设置替换后的文本
	glossaryBean.UserSelectedTextAfterReplace = result

	return glossaryBean, nil
}

// DealTranslationText 处理翻译结果中的术语
func (s *GlossaryService) DealTranslationText(ctx context.Context, translateResp *translate.TranslateResponse, glossaryBean *model.GlossaryTranslateModel) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.DealTranslationText")
	defer span.Finish()

	// 检查参数
	if glossaryBean == nil || len(glossaryBean.RelInfos) == 0 || len(translateResp.TargetContent) == 0 {
		return nil
	}

	// 获取翻译结果
	rtc := strings.TrimSpace(translateResp.TargetContent[0])

	// 查找匹配的术语
	var matchedRelInfo *model.RelInfo
	for i := range glossaryBean.RelInfos {
		relInfo := &glossaryBean.RelInfos[i]
		if relInfo.MatchCase {
			if relInfo.OriginalText == rtc {
				matchedRelInfo = relInfo
				break
			}
		} else {
			if strings.EqualFold(relInfo.OriginalText, rtc) {
				matchedRelInfo = relInfo
				break
			}
		}
	}

	// 如果找到匹配的术语，替换翻译结果
	if matchedRelInfo != nil {
		// 设置目标内容为原始选择的文本
		translateResp.TargetContent = []string{glossaryBean.UserSelectedText}

		// 创建目标响应
		targetResp := &translate.TargetResp{
			TargetContent: []string{matchedRelInfo.Translation},
		}

		// 设置目标响应
		translateResp.TargetResp = []*translate.TargetResp{targetResp}
	}

	return nil
}

// DealTranslationTextTargetContent 处理目标内容中的术语
func (s *GlossaryService) DealTranslationTextTargetContent(ctx context.Context, targetContent string, glossaryBean *model.GlossaryTranslateModel) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.DealTranslationTextTargetContent")
	defer span.Finish()

	// 检查参数
	if glossaryBean == nil || len(glossaryBean.RelInfos) == 0 {
		return nil
	}

	// 去除前后空格
	targetContent = strings.TrimSpace(targetContent)

	// 遍历所有术语，替换特殊标记
	for _, relInfo := range glossaryBean.RelInfos {
		if !strings.Contains(strings.ToLower(targetContent), strings.ToLower(relInfo.Key)) {
			continue
		}

		// 不区分大小写替换
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(relInfo.Key))
		targetContent = re.ReplaceAllString(targetContent, relInfo.Translation)
	}

	// 设置替换后的翻译
	glossaryBean.OriginalTranslationAfterReplace = targetContent

	return nil
}

// GetGlossariesByUserID 获取用户的所有术语条目
func (s *GlossaryService) GetGlossariesByUserID(ctx context.Context, userId string) ([]model.Glossary, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.GetGlossariesByUserID")
	defer span.Finish()

	// 获取用户的术语条目列表
	glossaries, err := s.glossaryDAO.GetGlossariesByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取用户术语条目列表失败", "error", err.Error())
		return nil, errors.BizWrap("translate.glossary.errors.invalid_data", err)
	}

	return glossaries, nil
}

// GetGlossary 获取术语条目
func (s *GlossaryService) GetGlossary(ctx context.Context, userId string, id string) (*translate.GlossaryItem, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GlossaryService.GetGlossary")
	defer span.Finish()

	// 查询术语条目
	glossary, err := s.glossaryDAO.FindExistById(ctx, id)
	if err != nil {
		s.logger.Error("msg", "查询术语条目失败", "error", err.Error())
		return nil, errors.BizWrap("translate.glossary.errors.entry_not_found", err)
	}

	// 检查是否存在
	if glossary == nil {
		return nil, errors.Biz("translate.glossary.errors.entry_not_found")
	}

	// 检查权限
	if glossary.UserId != userId {
		return nil, errors.Biz("translate.glossary.errors.permission_denied")
	}

	// 转换为响应对象
	return &translate.GlossaryItem{
		Id:              glossary.Id,
		OriginalText:    glossary.OriginalText,
		TranslationText: glossary.TranslationText,
		MatchCase:       glossary.MatchCase,
		Ignored:         glossary.Ignored,
	}, nil
}
