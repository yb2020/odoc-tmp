package paginate

import "time"

type PaginateOptions struct {
	// 精确匹配条件 (字段名 -> 值)
	Equals map[string]interface{}

	// 模糊匹配条件 (字段名 -> 值)
	Like map[string]interface{}

	// 范围条件 (字段名 -> 范围值)
	Ranges map[string]*RangeCondition

	// 自定义条件 (SQL片段 -> 参数值)
	Custom map[string][]interface{}

	// 排序方式 (字段名 -> 排序方向)
	OrderBy map[string]string
}

// RangeCondition 范围条件
type RangeCondition struct {
	Min interface{}
	Max interface{}
}

// NewOptions 创建一个新的分页查询选项
func NewOptions() *PaginateOptions {
	return &PaginateOptions{
		Equals: make(map[string]interface{}),
		Like:   make(map[string]interface{}),
		Ranges: make(map[string]*RangeCondition),
		Custom: make(map[string][]interface{}),
		OrderBy: map[string]string{
			"id": "DESC", // 默认按ID降序排序
		},
	}
}

// AddEquals 添加精确匹配条件
func (o *PaginateOptions) AddEquals(field string, value interface{}) *PaginateOptions {
	if value != nil && value != "" {
		o.Equals[field] = value
	}
	return o
}

// AddLike 添加模糊匹配条件
func (o *PaginateOptions) AddLike(field string, value string) *PaginateOptions {
	if value != "" {
		o.Like[field] = value
	}
	return o
}

// AddRange 添加范围条件
func (o *PaginateOptions) AddRange(field string, min, max interface{}) *PaginateOptions {
	if min != nil || max != nil {
		o.Ranges[field] = &RangeCondition{
			Min: min,
			Max: max,
		}
	}
	return o
}

// AddTimeRange 添加时间范围条件
func (o *PaginateOptions) AddTimeRange(field string, minTimestamp, maxTimestamp uint64) *PaginateOptions {
	var min, max interface{}

	if minTimestamp > 0 {
		min = time.Unix(int64(minTimestamp), 0)
	}

	if maxTimestamp > 0 {
		max = time.Unix(int64(maxTimestamp), 0)
	}

	return o.AddRange(field, min, max)
}

// AddOrder 添加排序条件
func (o *PaginateOptions) AddOrder(field, direction string) *PaginateOptions {
	if field != "" && direction != "" {
		o.OrderBy[field] = direction
	}
	return o
}
