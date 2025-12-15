package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/ratelimit"
)

func init() {
	// 注册基于时间的插件工厂
	ratelimit.RegisterPluginFactory("time_based", NewTimeBasedPlugin)
}

// TimeRange 时间范围
type TimeRange struct {
	Start string // 格式: "HH:MM"
	End   string // 格式: "HH:MM"
}

// TimeBasedPlugin 基于时间的插件
type TimeBasedPlugin struct {
	timeRanges []TimeRange
	logger     logging.Logger
}

// NewTimeBasedPlugin 创建基于时间的插件
func NewTimeBasedPlugin(params map[string]interface{}, logger logging.Logger) (ratelimit.RateLimiterPlugin, error) {
	// 从参数中获取时间范围
	timeRangesParam, ok := params["time_ranges"]
	if !ok {
		return nil, fmt.Errorf("缺少必要参数: time_ranges")
	}

	// 转换参数类型
	timeRangesSlice, ok := timeRangesParam.([]interface{})
	if !ok {
		return nil, fmt.Errorf("time_ranges 参数类型错误，应为时间范围数组")
	}

	// 构建时间范围
	timeRanges := make([]TimeRange, 0, len(timeRangesSlice))
	for _, item := range timeRangesSlice {
		rangeMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("time_range 参数类型错误，应为包含 start 和 end 的对象")
		}

		start, ok := rangeMap["start"].(string)
		if !ok {
			return nil, fmt.Errorf("time_range.start 参数类型错误，应为字符串")
		}

		end, ok := rangeMap["end"].(string)
		if !ok {
			return nil, fmt.Errorf("time_range.end 参数类型错误，应为字符串")
		}

		timeRanges = append(timeRanges, TimeRange{
			Start: start,
			End:   end,
		})
	}

	return &TimeBasedPlugin{
		timeRanges: timeRanges,
		logger:     logger,
	}, nil
}

// Name 插件名称
func (p *TimeBasedPlugin) Name() string {
	return "time_based_plugin"
}

// DoFilter 过滤方法
func (p *TimeBasedPlugin) DoFilter(ctx context.Context, key string, tokens int64) (bool, error) {
	// 获取当前时间
	now := time.Now()
	currentTime := now.Format("15:04")

	// 检查当前时间是否在任一时间范围内
	for _, timeRange := range p.timeRanges {
		if isTimeInRange(currentTime, timeRange.Start, timeRange.End) {
			p.logger.Info("current time in specified range, bypassing rate limit", 
				"current", currentTime, "range", fmt.Sprintf("%s-%s", timeRange.Start, timeRange.End))
			return true, nil
		}
	}

	return false, nil
}

// isTimeInRange 判断时间是否在范围内
func isTimeInRange(current, start, end string) bool {
	// 如果结束时间小于开始时间，表示跨天
	if end < start {
		return current >= start || current <= end
	}
	// 正常情况
	return current >= start && current <= end
}
