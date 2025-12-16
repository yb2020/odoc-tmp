package service

import (
	"context"
	std_errors "errors"
	"math"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	pb "github.com/yb2020/odoc/proto/gen/go/nav"
	"github.com/yb2020/odoc/services/nav/dao"
	"github.com/yb2020/odoc/services/nav/model"
	"gorm.io/gorm"
)

const (
	// SortOrderGap 定义了两个项目排序顺序之间的默认间隔。
	SortOrderGap int32 = 65536 // 2^16
	// MinSortOrderGap 定义了插入一个项目而无需重新平衡所需的最小间隙。
	MinSortOrderGap int32 = 2
)

// WebsiteService 学术网站服务实现
type WebsiteService struct {
	websiteDAO *dao.WebsiteDAO
	logger     logging.Logger
	tracer     opentracing.Tracer
	navConfig  *config.NavConfig
}

// NewWebsiteService 创建一个新的网站服务。
func NewWebsiteService(websiteDAO *dao.WebsiteDAO, logger logging.Logger, tracer opentracing.Tracer, navConfig *config.NavConfig) *WebsiteService {
	return &WebsiteService{
		websiteDAO: websiteDAO,
		logger:     logger,
		tracer:     tracer,
		navConfig:  navConfig,
	}
}

func (s *WebsiteService) CreateWebsite(ctx context.Context, userId string, source int32, req *pb.CreateWebsiteRequest) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.CreateWebsite")
	defer span.Finish()

	// 获取当前最大排序号，以便将新网站附加到末尾。
	maxSortOrder, err := s.websiteDAO.GetMaxSortOrder(ctx, userId)
	if err != nil {
		s.logger.Error("Failed to get max sort order for user", "error", err, "userId", userId)
		return "0", err
	}

	// 检查潜在的排序号溢出。
	if int64(maxSortOrder)+int64(SortOrderGap) > math.MaxInt32 {
		// 如果检测到溢出，首先为用户重新平衡所有网站。
		s.logger.Info("SortOrder overflow detected, triggering rebalance", "userId", userId)
		_, err = s.rebalanceWebsites(ctx, userId)
		if err != nil {
			s.logger.Error("Failed to rebalance websites on overflow", "error", err, "userId", userId)
			return "0", err
		}
		// 重新平衡后，获取新的最大排序号。
		maxSortOrder, err = s.websiteDAO.GetMaxSortOrder(ctx, userId)
		if err != nil {
			s.logger.Error("Failed to get max sort order after rebalance", "error", err, "userId", userId)
			return "0", err
		}
	}

	website := &model.Website{
		UserId:    userId,
		Source:    source,
		IconUrl:   req.IconUrl,
		Name:      req.Name,
		Url:       req.Url,
		OpenType:  int32(req.OpenType),
		SortOrder: maxSortOrder + SortOrderGap, // 设置排序号为最后一个。
	}
	website.Id = idgen.GenerateUUID()

	err = s.websiteDAO.Save(ctx, website)
	if err != nil {
		return "0", err
	}

	return website.Id, nil
}

func (s *WebsiteService) DeleteWebsite(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.DeleteWebsite")
	defer span.Finish()

	websiteExist, err := s.GetById(ctx, id)
	if err != nil {
		return err
	}
	if websiteExist == nil {
		return errors.Biz("website not found")
	}

	return s.websiteDAO.DeleteById(ctx, id)
}

func (s *WebsiteService) GetById(ctx context.Context, id string) (*model.Website, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.GetById")
	defer span.Finish()

	return s.websiteDAO.FindExistById(ctx, id)
}

func (s *WebsiteService) GetWebsitePbById(ctx context.Context, id string) (*pb.Website, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.GetWebsitePbById")
	defer span.Finish()

	website, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if website == nil {
		return nil, errors.Biz("website not found")
	}

	return &pb.Website{
		Id:       website.Id,
		UserId:   website.UserId,
		Source:   pb.WebsiteSource(website.Source),
		IconUrl:  website.IconUrl,
		Name:     website.Name,
		Url:      website.Url,
		OpenType: pb.WebsiteOpenType(website.OpenType),
	}, nil
}

func (s *WebsiteService) UpdateWebsite(ctx context.Context, req *pb.UpdateWebsiteRequest) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.UpdateWebsite")
	defer span.Finish()

	website, err := s.GetById(ctx, req.Id)
	if err != nil {
		return "0", err
	}
	if website == nil {
		return "0", errors.Biz("website not found")
	}

	website.IconUrl = req.IconUrl
	website.Name = req.Name
	website.Url = req.Url
	website.OpenType = int32(req.OpenType)

	return website.Id, s.websiteDAO.ModifyExcludeNull(ctx, website)
}

// ModifyWebsiteSortOrder 修改网站排序
func (s *WebsiteService) ModifyWebsiteSortOrder(ctx context.Context, id string, sortOrder int32) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.ModifyWebsiteSortOrder")
	defer span.Finish()

	website, err := s.GetById(ctx, id)
	if err != nil {
		return err
	}
	if website == nil {
		return errors.Biz("website not found")
	}

	website.SortOrder = sortOrder

	return s.websiteDAO.ModifyExcludeNull(ctx, website)
}

// GetUserWebsiteListBySortOrder 获取用户学术网站列表(按SortOrder升序排序)
func (s *WebsiteService) GetUserWebsiteListBySortOrder(ctx context.Context, userId string) ([]model.Website, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.GetUserWebsiteListBySortOrder")
	defer span.Finish()

	websiteList, err := s.websiteDAO.GetUserWebsiteListBySortOrder(ctx, userId)
	if err != nil {
		return nil, err
	}

	return websiteList, nil
}

func (s *WebsiteService) GetUserWebsiteList(ctx context.Context, userId string) ([]*pb.Website, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.GetUserWebsiteList")
	defer span.Finish()

	websiteModels, err := s.websiteDAO.GetUserWebsiteListBySortOrder(ctx, userId)
	if err != nil {
		return nil, err
	}

	websiteList := make([]*pb.Website, 0, len(websiteModels))
	for _, m := range websiteModels {
		websiteList = append(websiteList, &pb.Website{
			Id:        m.Id,
			Source:    pb.WebsiteSource(m.Source),
			IconUrl:   m.IconUrl,
			Name:      m.Name,
			Url:       m.Url,
			OpenType:  pb.WebsiteOpenType(m.OpenType),
			SortOrder: m.SortOrder,
		})
	}

	return websiteList, nil
}

// rebalanceWebsites 为用户的所有网站重新分配排序号，使其均匀分布。
func (s *WebsiteService) rebalanceWebsites(ctx context.Context, userId string) ([]*pb.WebsiteSortUpdate, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.rebalanceWebsites")
	defer span.Finish()

	s.logger.Info("Rebalance triggered, re-assigning sort orders for user", "userId", userId)

	allWebsites, err := s.websiteDAO.GetUserWebsiteListBySortOrder(ctx, userId)
	if err != nil {
		s.logger.Error("Failed to get user website list for rebalance", "error", err, "userId", userId)
		return nil, err
	}

	updates := make([]*pb.WebsiteSortUpdate, 0, len(allWebsites))
	batchUpdateParams := make(map[string]int32, len(allWebsites))
	for i, website := range allWebsites {
		newSortOrder := int32(i+1) * SortOrderGap
		if website.SortOrder != newSortOrder {
			batchUpdateParams[website.Id] = newSortOrder
			updates = append(updates, &pb.WebsiteSortUpdate{
				Id:        website.Id,
				SortOrder: newSortOrder,
			})
		}
	}

	if len(batchUpdateParams) > 0 {
		err = s.websiteDAO.BatchUpdateSortOrder(ctx, batchUpdateParams)
		if err != nil {
			s.logger.Error("Failed to batch update sort order during rebalance", "error", err, "userId", userId)
			return nil, err
		}
	}

	return updates, nil
}

func (s *WebsiteService) ReorderWebsites(ctx context.Context, userId string, req *pb.ReorderWebsitesRequest) (bool, []*pb.WebsiteSortUpdate, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.ReorderWebsites")
	defer span.Finish()

	var rebalanced bool
	var updatesForApi []*pb.WebsiteSortUpdate

	// 1. 权限检查：验证网站是否属于该用户。
	movedWebsite, err := s.websiteDAO.FindExistById(ctx, req.Id)
	if err != nil {
		if std_errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, std_errors.New("website not found")
		}
		return false, nil, err
	}
	if movedWebsite.UserId != userId {
		return false, nil, std_errors.New("permission denied")
	}

	// 2. 获取邻居的排序号
	var prevSortOrder int32
	if req.PrevId != "0" {
		prevWebsite, err := s.websiteDAO.FindExistById(ctx, req.PrevId)
		if err != nil {
			if std_errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil, std_errors.New("previous website not found")
			}
			return false, nil, err
		}
		prevSortOrder = prevWebsite.SortOrder
	} else {
		prevSortOrder = 0 // 移动到列表的开头。
	}

	var nextSortOrder int32
	if req.NextId != "0" {
		nextWebsite, err := s.websiteDAO.FindExistById(ctx, req.NextId)
		if err != nil {
			if std_errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil, std_errors.New("next website not found")
			}
			return false, nil, err
		}
		nextSortOrder = nextWebsite.SortOrder
	} else {
		// 移动到列表的末尾。
		nextSortOrder = prevSortOrder + SortOrderGap*2
	}

	// 3. 检查是否需要重新平衡。
	if req.PrevId != "0" && req.NextId != "0" && nextSortOrder-prevSortOrder < MinSortOrderGap {
		// 需要重新平衡，获取用户的所有网站并重新分配排序号。
		rebalanced = true
		updates, err := s.rebalanceWebsites(ctx, userId)
		if err != nil {
			return false, nil, err
		}
		updatesForApi = updates

		// 重新平衡后，我们需要找到移动项的新邻居
		// 以计算其最终的排序号。
		newSortOrdersMap := make(map[string]int32, len(updatesForApi))
		for _, u := range updatesForApi {
			newSortOrdersMap[u.Id] = u.SortOrder
		}

		if req.PrevId != "0" {
			prevSortOrder = newSortOrdersMap[req.PrevId]
		} else {
			prevSortOrder = 0
		}

		if req.NextId != "0" {
			nextSortOrder = newSortOrdersMap[req.NextId]
		} else {
			// 如果没有下一个项目，说明它被移到了末尾。
			nextSortOrder = prevSortOrder + SortOrderGap*2
		}

		newSortOrder := (prevSortOrder + nextSortOrder) / 2
		// 移动项的排序号也需要在数据库中更新。
		if err := s.websiteDAO.UpdateSortOrder(ctx, req.Id, newSortOrder); err != nil {
			return false, nil, err
		}
		// 将移动项的最终更新添加到将要返回的列表中。
		updatesForApi = append(updatesForApi, &pb.WebsiteSortUpdate{Id: req.Id, SortOrder: newSortOrder})

	} else {
		// 无需重新平衡，只需计算新的排序号并更新。
		newSortOrder := prevSortOrder + (nextSortOrder-prevSortOrder)/2
		if err := s.websiteDAO.UpdateSortOrder(ctx, req.Id, newSortOrder); err != nil {
			return false, nil, err
		}
		updatesForApi = append(updatesForApi, &pb.WebsiteSortUpdate{Id: req.Id, SortOrder: newSortOrder})
	}

	// 4. 准备最终响应。
	return rebalanced, updatesForApi, nil
}

// CheckUserInitSystemWebsite 检查用户是否初始化了系统网站
func (s *WebsiteService) CheckUserInitSystemWebsite(ctx context.Context, userId string) (bool, error) {
	return s.websiteDAO.CheckUserInitSystemWebsite(ctx, userId)
}

// InitUserSystemWebsite 初始化用户系统网站
func (s *WebsiteService) InitUserSystemWebsite(ctx context.Context, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "WebsiteService.InitUserSystemWebsite")
	defer span.Finish()

	initData := s.navConfig.Website.InitData
	for _, data := range initData {
		_, err := s.CreateWebsite(ctx, userId, int32(pb.WebsiteSource_WebsiteSource_System), &pb.CreateWebsiteRequest{
			Name:     data.Name,
			Url:      data.URL,
			IconUrl:  data.IconURL,
			OpenType: pb.WebsiteOpenType_WebsiteOpenType_NewTab,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
